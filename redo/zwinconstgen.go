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
	"os/exec"
)

func getPackage(path string) (pkg *ast.Package) {
	fileset := token.NewFileSet()		// parser.ParseDir() actually writes to this; not sure why it doesn't return one instead
	filter := func(i os.FileInfo) bool {
		return strings.HasSuffix(i.Name(), "_windows.go")
	}
	pkgs, err := parser.ParseDir(fileset, path, filter, parser.AllErrors)
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
	case *ast.Comment:		// function?
		// TODO
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
	}
}

// for backwards compatibiilty reasons, Windows defines GetWindowLongPtr()/SetWindowLongPtr() as a macro which expands to GetWindowLong()/SetWindowLong() on 32-bit systems
// we'll just simulate that here
var gwlpNames = map[string]string{
	"386":		"etWindowLongW",
	"amd64":		"etWindowLongPtrW",
}

func writeLine(f *os.File, line string) {
	fmt.Fprintf(f, "%s\n", line)
}

const cgopreamble = `
#define UNICODE
#define _UNICODE
#define STRICT
#define STRICT_TYPED_ITEMIDS
/* get Windows version right; right now Windows XP */
#define WINVER 0x0501
#define _WIN32_WINNT 0x0501
#define _WIN32_WINDOWS 0x0501		/* according to Microsoft's winperf.h */
#define _WIN32_IE 0x0600			/* according to Microsoft's sdkddkver.h */
#define NTDDI_VERSION 0x05010000	/* according to Microsoft's sdkddkver.h */
#include <windows.h>
#include <commctrl.h>
#include <stdint.h>
`

func writeConstCast(f *os.File, c string) {
	cc := c[2:]		// strip leading c_
	fmt.Fprintf(f, "uintptr_t %s = (uintptr_t) (%s);\n", c, cc)
}

func writeConstPrint(f *os.File, c string) {
	fmt.Fprintf(f, "\tfmt.Fprintf(buf, \"const %%s = %%d\\n\", %q, C.%s)\n", c, c)
}

func writeStructPrint(f *os.File, s string) {
	cs := s[2:]		// strip leading s_
	fmt.Fprintf(f, "\tt = reflect.TypeOf(C.%s{})\n", cs)
	// TODO
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
	for ident, _ := range unknown {
		if strings.HasPrefix(ident, "s_") {
			structs = append(structs, ident)
			continue
		}
		consts = append(consts, ident)
	}
	sort.Strings(consts)
	sort.Strings(structs)

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

	writeLine(f, "package main")
	writeLine(f, "import \"fmt\"")
	writeLine(f, "import \"bytes\"")
	writeLine(f, "import \"runtime\"")
	writeLine(f, "import \"go/format\"")
	writeLine(f, "/*")
	writeLine(f, cgopreamble)
	for _, const := range consts {
		writeConstCast(f, const)
	}
	writeLine(f, "*/")
	writeLine(f, "import \"C\")
	writeLine(f, "func main() {")
	writeLine(f, "\tbuf := new(bytes.Buffer)")
	writeLine(f, "\tfmt.Fprintln(buf, \"package main\")")
	writeLine(f, "\tfmt.Fprintln(buf, \"func main() {\")")
	for _, c := range consts {
		writeConstPrint(f, c)
	}
	writeLine(f, "\tfmt.Fprintln(buf, \"var t reflect.Type\")")
	writeLine(f, "\tfmt.Fprintln(buf, \"var s string\")")
	for _, s := range structs {
		writeStructPrint(f, s)
	}
	writeLine(f, "\tfmt.Fprintln(buf, \"}\")")
	writeLine(f, "\tres, err := format.Source(buf.Bytes())")
	writeLine(f, "\tif err != nil { panic(err) }")
	wrtieLine(f, "\tfmt.Printf(\"%s\", res)")
	writeLine(f, "}")

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
