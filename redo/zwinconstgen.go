// +build ignore

// 24 may 2014
package main

import (
	"fmt"
	"os"
	"strings"
	"go/token"
	"go/ast"
	"go/parser"
	"sort"
	"io/ioutil"
	"path/filepath"
	"text/template"
	"os/exec"
)

func getPackage(path string) (pkg *ast.Package) {
	fileset := token.NewFileSet()		// parser.ParseDir() actually writes to this; not sure why it doesn't return one instead
	filter := func(i os.FileInfo) bool {
		return strings.HasSuffix(i.Name(), "_windows.go")
	}
	pkgs, err := parser.ParseDir(fileset, path, filter, parser.AllErrors | parser.ParseComments)
	if err != nil {
		panic(err)
	}
	if len(pkgs) != 1 {
		panic("more than one package found")
	}
	for k, _ := range pkgs {		// get the sole key
		pkg = pkgs[k]
	}
	return pkg
}

type walker struct {
	desired	func(string) bool
}

var known = map[string]string{}
var unknown = map[string]struct{}{}

func (w *walker) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.Ident:			// constant or structure?
		if w.desired(n.Name) {
			if n.Obj != nil {
				delete(unknown, n.Name)
				kind := n.Obj.Kind.String()
				if known[n.Name] != "" && known[n.Name] != kind {
					panic(n.Name + "(" + kind + ") already known to be a " + known[n.Name])
				}
				known[n.Name] = kind
			} else if _, ok := known[n.Name]; !ok {		// only if not known
				unknown[n.Name] = struct{}{}
			}
		}
	}
	return w
}

func gatherNames(pkg *ast.Package) {
	desired := func(name string) bool {
		return strings.HasPrefix(name, "c_") ||		// constants
			strings.HasPrefix(name, "s_")			// structs
	}
	for _, f := range pkg.Files {
		for _, d := range f.Decls {
			ast.Walk(&walker{desired}, d)
		}
		for _, c := range f.Comments {
			for _, cc := range c.List {
				readComment(cc)
			}
		}
	}
}

var funcs []string
var dlls = map[string]struct{}{}

// parse comments of the form "// wfunc dll Name argtypes {ret[,noerr]|void}"
// TODO clean this up
func readComment(c *ast.Comment) {
	words := strings.Split(c.Text, " ")[1:]		// strip leading //
	if len(words) <= 0 || words[0] != "wfunc" {
		return
	}
	dll := words[1]
	dlls[dll] = struct{}{}
	name := words[2]
	args := make([]string, 0, len(words))
	for i := 3; i < len(words) - 1; i++ {
		args = append(args, words[i])
	}
	ret := words[len(words) - 1]

	funcs = append(funcs, fmt.Sprintf("var fv_%s = %s.NewProc(%q)", name, dll, name))

	r1 := "r1"
	err := "err"
	assign := ":="

	r := rune('a')
	argspec := ""
	for _, t := range args {
		argspec += fmt.Sprintf("%c %s, ", r, t)
		r++
	}
	retspec := ""
	if ret != "void" {
		r := strings.Split(ret, ",")
		retspec = "(" + r[0]
		if len(r) > 1 && r[1] == "noerr" {
			err = "_"
		} else {
			retspec += ", error"
		}
		retspec += ")"
	} else {
		r1 = "_"
		err = "_"
		assign = "="
	}
	funcs = append(funcs, fmt.Sprintf("func f_%s(%s) %s {", name, argspec, retspec))

	call := fmt.Sprintf("\t%s, _, %s %s fv_%s.Call(", r1, err, assign, name)
	r = rune('a')
	for _, t := range args {
		call += "uintptr("
		if t[0] == '*' {
			call += "unsafe.Pointer("
		}
		call += fmt.Sprintf("%c", r)
		if t[0] == '*' {
			call += ")"
		}
		call += "), "
		r++
	}
	call += ")"
	funcs = append(funcs, call)

	if ret != "void" {
		r := strings.Split(ret, ",")
		retspec = "return (" + r[0] + ")("
		if r[0][0] == '*' {
			retspec += "unsafe.Pointer("
		}
		retspec += "r1"
		if r[0][0] == '*' {
			retspec += ")"
		}
		retspec += ")"
		if len(r) > 1 && r[1] == "noerr" {
			// do nothing
		} else {
			retspec += ", err"
		}
		funcs = append(funcs, retspec)
	}

	funcs = append(funcs, "}")
}

// for backwards compatibiilty reasons, Windows defines GetWindowLongPtr()/SetWindowLongPtr() as a macro which expands to GetWindowLong()/SetWindowLong() on 32-bit systems
// we'll just simulate that here
var gwlpNames = map[string]string{
	"386":		"etWindowLongW",
	"amd64":		"etWindowLongPtrW",
}

// in reality these use LONG_PTR for the actual values; LONG_PTR is a signed value, but for our use case it doesn't really matter
func genGetSetWindowLongPtr(targetarch string) {
	name := gwlpNames[targetarch]

	funcs = append(funcs, fmt.Sprintf("var fv_GetWindowLongPtrW = user32.NewProc(%q)", "G" + name))
	funcs = append(funcs, "func f_GetWindowLongPtrW(hwnd uintptr, which uintptr) uintptr {")
	funcs = append(funcs, "\tres, _, _ := fv_GetWindowLongPtrW.Call(hwnd, which)")
	funcs = append(funcs, "\treturn res")
	funcs = append(funcs, "}")

	funcs = append(funcs, fmt.Sprintf("var fv_SetWindowLongPtrW = user32.NewProc(%q)", "S" + name))
	funcs = append(funcs, "func f_SetWindowLongPtrW(hwnd uintptr, which uintptr, value uintptr) {")
	funcs = append(funcs, "\tfv_SetWindowLongPtrW.Call(hwnd, which, value)")
	funcs = append(funcs, "}")
}

const outTemplate = `package main
import (
	"fmt"
	"bytes"
	"reflect"
	"go/format"
	"strings"
)
// #define UNICODE
// #define _UNICODE
// #define STRICT
// #define STRICT_TYPED_ITEMIDS
// /* get Windows version right; right now Windows XP */
// #define WINVER 0x0501
// #define _WIN32_WINNT 0x0501
// #define _WIN32_WINDOWS 0x0501		/* according to Microsoft's winperf.h */
// #define _WIN32_IE 0x0600				/* according to Microsoft's sdkddkver.h */
// #define NTDDI_VERSION 0x05010000	/* according to Microsoft's sdkddkver.h */
// #include <windows.h>
// #include <commctrl.h>
// #include <stdint.h>
{{range .Consts}}// uintptr_t {{.}} = (uintptr_t) ({{noprefix .}});
{{end}}import "C"		// notice the lack of newline in the template
// MinGW will generate handle pointers as pointers to some structure type under some conditions I don't fully understand; here's full overrides
var handleOverrides = []string{
	"HWND",
	"HINSTANCE",
	"HICON",
	"HCURSOR",
	"HBRUSH",
	"HMENU",
	// These are all pointers to functions; handle them identically to handles.
	"WNDPROC",
}
func winName(t reflect.Type) string {
	for _, s := range handleOverrides {
		if strings.Contains(t.Name(), s) {
			return "uintptr"
		}
	}
	switch t.Kind() {
	case reflect.UnsafePointer:
		return "uintptr"
	case reflect.Ptr:
		return "*" + winName(t.Elem())
	case reflect.Struct:
		// the t.Name() will be the cgo-mangled name; get the original name out
		parts := strings.Split(t.Name(), "_")
		part := parts[len(parts) - 1]
		// many Windows API types have struct tagXXX as their declarator
		// if you wonder why, see http://blogs.msdn.com/b/oldnewthing/archive/2008/03/26/8336829.aspx?Redirected=true
		if strings.HasPrefix(part, "tag") {
			part = part[3:]
		}
		return "s_" + part
	case reflect.Array:
		return fmt.Sprintf("[%d]%s", t.Len(), winName(t.Elem()))
	}
	return t.Kind().String()
}
func main() {
	buf := new(bytes.Buffer)
	fmt.Fprintln(buf, "package ui")
	fmt.Fprintln(buf, "import (")
	fmt.Fprintln(buf, "\t\"syscall\"")
	fmt.Fprintln(buf, "\t\"unsafe\"")
	fmt.Fprintln(buf, ")")

	// constants
{{range .Consts}}	fmt.Fprintf(buf, "const %s = %d\n", {{printf "%q" .}}, C.{{.}})
{{end}}

	// structures
	var t reflect.Type
{{range .Structs}}	t = reflect.TypeOf(C.{{noprefix .}}{})
	fmt.Fprintf(buf, "type %s struct {\n", {{printf "%q" .}})
	for i := 0; i < t.NumField(); i++ {
		fmt.Fprintf(buf, "\t%s %s\n", t.Field(i).Name, winName(t.Field(i).Type))
	}
	fmt.Fprintf(buf, "}\n")
{{end}}

	// let's generate names for window procedure types
	fmt.Fprintf(buf, "\n")
	fmt.Fprintf(buf, "type t_UINT %s\n", winName(reflect.TypeOf(C.UINT(0))))
	fmt.Fprintf(buf, "type t_WPARAM %s\n", winName(reflect.TypeOf(C.WPARAM(0))))
	fmt.Fprintf(buf, "type t_LPARAM %s\n", winName(reflect.TypeOf(C.LPARAM(0))))
	fmt.Fprintf(buf, "type t_LRESULT %s\n", winName(reflect.TypeOf(C.LRESULT(0))))
	// and one for GetMessageW()
	fmt.Fprintf(buf, "type t_BOOL %s\n", winName(reflect.TypeOf(C.BOOL(0))))

	// functions
{{range .Funcs}}	fmt.Fprintf(buf, "%s\n", {{printf "%q" .}})
{{end}}

	// DLLs
{{range .DLLs}}		fmt.Fprintf(buf, "var %s = syscall.NewLazyDLL(%q)\n", {{printf "%q" .}}, {{printf "%s.dll" . | printf "%q"}})
{{end}}

	// and finally done
	res, err := format.Source(buf.Bytes())
	if err != nil { panic(err.Error() + "\n" + string(buf.Bytes())) }
	fmt.Printf("%s", res)
}
`

type templateArgs struct {
	Consts	[]string
	Structs	[]string
	Funcs	[]string
	DLLs		[]string
}

var templateFuncs = template.FuncMap{
	"noprefix":	func(s string) string {
		return s[2:]
	},
}

func main() {
	if len(os.Args) < 3 {
		panic("usage: " + os.Args[0] + " path goarch [go-command-options...]")
	}
	pkgpath := os.Args[1]
	targetarch := os.Args[2]
	if _, ok := gwlpNames[targetarch]; !ok {
		panic("unknown target windows/" + targetarch)
	}
	goopts := os.Args[3:]		// valid if len(os.Args) == 3; in that case this will just be a slice of length zero

	pkg := getPackage(pkgpath)
	gatherNames(pkg)

	// if we still have some known, I didn't clean things up completely
	if len(known) > 0 {
		knowns := ""
		for ident, kind := range known {
			if kind != "var" && kind != "const" {
				continue
			}
			knowns += "\n" + ident + " (" + kind + ")"
		}
		panic("error: the following are still known!" + knowns)		// has a newline already
	}

	// keep sorted for git
	consts := make([]string, 0, len(unknown))
	structs := make([]string, 0, len(unknown))
	sorteddlls := make([]string, 0, len(dlls))
	for ident, _ := range unknown {
		if strings.HasPrefix(ident, "s_") {
			structs = append(structs, ident)
			continue
		}
		consts = append(consts, ident)
	}
	for dll, _ := range dlls {
		sorteddlls = append(sorteddlls, dll)
	}
	sort.Strings(consts)
	sort.Strings(structs)
	sort.Strings(sorteddlls)

	// and finally
	genGetSetWindowLongPtr(targetarch)

	// thanks to james4k in irc.freenode.net/#go-nuts
	tmpdir, err := ioutil.TempDir("", "windowsconstgen")
	if err != nil {
		panic(err)
	}
	genoutname := filepath.Join(tmpdir, "gen.go")
	f, err := os.Create(genoutname)
	if err != nil {
		panic(err)
	}

	t := template.Must(template.New("winconstgenout").Funcs(templateFuncs).Parse(outTemplate))
	err = t.Execute(f, &templateArgs{
		Consts:		consts,
		Structs:		structs,
		Funcs:		funcs,
		DLLs:		sorteddlls,
	})
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("go", "run")
	cmd.Args = append(cmd.Args, goopts...)		// valid if len(goopts) == 0; in that case this will just be a no-op
	cmd.Args = append(cmd.Args, genoutname)
	f, err = os.Create(filepath.Join(pkgpath, "zconstants_windows_" + targetarch + ".go"))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	cmd.Stdout = f
	cmd.Stderr = os.Stderr
	// we need to preserve the environment EXCEPT FOR the variables we're overriding
	// thanks to raggi and smw in irc.freenode.net/#go-nuts
	for _, ev := range os.Environ() {
		if strings.HasPrefix(ev, "GOOS=") ||
			strings.HasPrefix(ev, "GOARCH=") ||
			strings.HasPrefix(ev, "CGO_ENABLED=") {
			continue
		}
		cmd.Env = append(cmd.Env, ev)
	}
	cmd.Env = append(cmd.Env,
		"GOOS=windows",
		"GOARCH=" + targetarch,
		"CGO_ENABLED=1")		// needed as it's not set by default in cross-compiles
	err = cmd.Run()
	if err != nil {
		// TODO find a way to get the exit code
		os.Exit(1)
	}

	// TODO remove the temporary directory
}
