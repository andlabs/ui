// 7 february 2014
package main

//+build skip

import (
	"fmt"
	"os"
	"runtime"
)

func fatalf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	_, err := MessageBox(NULL,
		"An internal error has occured:\n" + s,
		os.Args[0],
		MB_OK | MB_ICONERROR)
	if err == nil {
		os.Exit(1)
	}
	panic(fmt.Sprintf("error trying to warn user of internal error: %v\ninternal error:\n%s", err, s))
}

const (
	IDC_BUTTON = 100 + iota
	IDC_VARCOMBO
	IDC_FIXCOMBO
	IDC_EDIT
	IDC_LIST
	IDC_LABEL
	IDC_CHECK
)

var varCombo, fixCombo, edit, list HWND

func wndProc(hwnd HWND, msg uint32, wParam WPARAM, lParam LPARAM) LRESULT {
	switch msg {
	case WM_COMMAND:
		if wParam.LOWORD() == IDC_BUTTON {
			buttonclick := "neither clicked nor double clicked (somehow)"
			if wParam.HIWORD() == BN_CLICKED {
				buttonclick = "clicked"
			} else if wParam.HIWORD() == BN_DOUBLECLICKED {
				buttonclick = "double clicked"
			}

			varText, err := getText(varCombo)
			if err != nil {
				fatalf("error getting variable combo box text: %v", err)
			}

			fixTextWM, err := getText(fixCombo)
			if err != nil {
				fatalf("error getting fixed combo box text with WM_GETTEXT: %v", err)
			}

			fixTextIndex, err := SendMessage(fixCombo, CB_GETCURSEL, 0, 0)
			if err != nil {
				fatalf("error getting fixed combo box current selection: %v", err)
			}
			// TODO get text from index

			editText, err := getText(edit)
			if err != nil {
				fatalf("error getting edit field text: %v", err)
			}

			listIndex, err := SendMessage(list, LB_GETCURSEL, 0, 0)
			if err != nil {
				fatalf("error getting fixed list box current selection: %v", err)
			}
			// TODO get text from index

			checkState, err := IsDlgButtonChecked(hwnd, IDC_CHECK)
			if err != nil {
				fatalf("error getting checkbox check state: %v", err)
			}

			MessageBox(hwnd,
				fmt.Sprintf("button state: %s\n" +
					"variable combo box text: %s\n" +
					"fixed combo box text with WM_GETTEXT: %s\n" +
					"fixed combo box current index: %d\n" +
					"edit field text: %s\n" +
					"list box current index: %d\n" +
					"check box checked: %v\n",
					buttonclick, varText, fixTextWM, fixTextIndex, editText, listIndex, checkState == BST_CHECKED),
				"note",
				MB_OK)
		}
		return 0
	case WM_GETMINMAXINFO:
		mm := lParam.MINMAXINFO()
		mm.PtMinTrackSize.X = 320
		mm.PtMinTrackSize.Y = 240
		return 0
	case WM_SIZE:
		if wParam != SIZE_MINIMIZED {
			resize(hwnd)
		}
		return 0
	case WM_CLOSE:
		err := DestroyWindow(hwnd)
		if err != nil {
			fatalf("error destroying window: %v", err)
		}
		return 0
	case WM_DESTROY:
		err := PostQuitMessage(0)
		if err != nil {
			fatalf("error posting quit message: %v", err)
		}
		return 0
	default:
		return DefWindowProc(hwnd, msg, wParam, lParam)
	}
	fatalf("major bug: forgot a return on wndProc for message %d", msg)
	panic("unreachable")
}

func setFontAll(hwnd HWND, lParam LPARAM) (cont bool) {
	_, err := SendMessage(hwnd, WM_SETFONT, WPARAM(lParam), LPARAM(TRUE))
	if err != nil {
		fatalf("error setting window font: %v", err)
	}
	return true
}

func resize(hwnd HWND) {
	cr, err := GetClientRect(hwnd)
	if err != nil {
		fatalf("error getting window client rect: %v", err)
	}
	cr.Bottom -= 80		// Y position of listbox
	cr.Bottom -= 20		// amount of pixels to leave behind
	err = SetWindowPos(list,
		HWND_TOP,
		20, 80, 100, int(cr.Bottom),
		0)
	if err != nil {
		fatalf("error resizing listbox: %v", err)
	}
}

const className = "mainwin"

func main() {
	runtime.LockOSThread()

	hInstance, err := getWinMainhInstance()
	if err != nil {
		fatalf("error getting WinMain hInstance: %v", err)
	}
	nCmdShow, err := getWinMainnCmdShow()
	if err != nil {
		fatalf("error getting WinMain nCmdShow: %v", err)
	}
	font, err := getStandardWindowFont()
	if err != nil {
		fatalf("error getting standard window font: %v", err)
	}

	icon, err := LoadIcon_ResourceID(NULL, IDI_APPLICATION)
	if err != nil {
		fatalf("error getting window icon: %v", err)
	}
	cursor, err := LoadCursor_ResourceID(NULL, IDC_ARROW)
	if err != nil {
		fatalf("error getting window cursor: %v", err)
	}

	wc := &WNDCLASS{
		LpszClassName:	className,
		LpfnWndProc:		wndProc,
		HInstance:		hInstance,
		HIcon:			icon,
		HCursor:			cursor,
		HbrBackground:	HBRUSH(COLOR_BTNFACE + 1),
	}
	_, err = RegisterClass(wc)
	if err != nil {
		fatalf("error registering window class: %v", err)
	}

	hwnd, err := CreateWindowEx(
		0,
		className, "Main Window",
		WS_OVERLAPPEDWINDOW,
		CW_USEDEFAULT, CW_USEDEFAULT, 320, 240,
		NULL, NULL, hInstance, NULL)
	if err != nil {
		fatalf("error creating window: %v", err)
	}

	const controlStyle = WS_CHILD | WS_VISIBLE | WS_TABSTOP

	_, err = CreateWindowEx(
		0,
		"BUTTON", "Click Me",
		BS_PUSHBUTTON | controlStyle,
		20, 20, 100, 20,
		hwnd, HMENU(IDC_BUTTON), hInstance, NULL)
	if err != nil {
		fatalf("error creating button: %v", err)
	}

	varCombo, err = CreateWindowEx(
		0,
		"COMBOBOX", "",
		CBS_DROPDOWN | CBS_AUTOHSCROLL | controlStyle,
		140, 20, 100, 20,
		hwnd, HMENU(IDC_VARCOMBO), hInstance, NULL)
	if err != nil {
		fatalf("error creating variable combo box: %v", err)
	}
	vcItems := []string{"a", "b", "c", "d"}
	for _, v := range vcItems {
		_, err := SendMessage(varCombo, CB_ADDSTRING, 0,
			LPARAMFromString(v))
		if err != nil {
			fatalf("error adding %q to variable combo box: %v", v, err)
		}
	}

	fixCombo, err = CreateWindowEx(
		0,
		"COMBOBOX", "",
		CBS_DROPDOWNLIST | controlStyle,
		140, 50, 100, 20,
		hwnd, HMENU(IDC_FIXCOMBO), hInstance, NULL)
	if err != nil {
		fatalf("error creating fixed combo box: %v", err)
	}
	fcItems := []string{"e", "f", "g", "h"}
	for _, v := range fcItems {
		_, err := SendMessage(fixCombo, CB_ADDSTRING, 0,
			LPARAMFromString(v))
		if err != nil {
			fatalf("error adding %q to fixed combo box: %v", v, err)
		}
	}

	edit, err = CreateWindowEx(
		0,
		"EDIT", "",
		ES_AUTOHSCROLL | ES_NOHIDESEL | WS_BORDER | controlStyle,
		20, 50, 100, 20,
		hwnd, HMENU(IDC_EDIT), hInstance, NULL)
	if err != nil {
		fatalf("error creating edit field: %v", err)
	}

	list, err = CreateWindowEx(
		0,
		"LISTBOX", "",
		LBS_STANDARD | controlStyle,
		20, 80, 100, 100,
		hwnd, HMENU(IDC_FIXCOMBO), hInstance, NULL)
	if err != nil {
		fatalf("error creating list box: %v", err)
	}
	lItems := []string{"i", "j", "k", "l"}
	for _, v := range lItems {
		_, err := SendMessage(list, LB_ADDSTRING, 0,
			LPARAMFromString(v))
		if err != nil {
			fatalf("error adding %q to list box: %v", v, err)
		}
		// TODO check actual return value as THAT indicates an error
	}

	_, err = CreateWindowEx(
		0,
		"STATIC", "Label",
		SS_NOPREFIX | controlStyle,
		140, 80, 100, 20,
		hwnd, HMENU(IDC_FIXCOMBO), hInstance, NULL)
	if err != nil {
		fatalf("error creating label: %v", err)
	}

	_, err = CreateWindowEx(
		0,
		"BUTTON", "Checkbox",
		BS_AUTOCHECKBOX | controlStyle,
		140, 110, 100, 20,
		hwnd, HMENU(IDC_CHECK), hInstance, NULL)
	if err != nil {
		fatalf("error creating checkbox: %v", err)
	}

	setFontAll(hwnd, LPARAM(font))
	err = EnumChildWindows(hwnd, setFontAll, LPARAM(font))
	if err != nil {
		fatalf("error setting font on controls: %v", err)
	}
	resize(hwnd)

	_, err = ShowWindow(hwnd, nCmdShow)
	if err != nil {
		fatalf("error showing window: %v", err)
	}
	err = UpdateWindow(hwnd)
	if err != nil {
		fatalf("error updating window: %v", err)
	}

	for {
		msg, done, err := GetMessage(NULL, 0, 0)
		if err != nil {
			fatalf("error getting message: %v", err)
		}
		if done {
			break
		}
		_, err = TranslateMessage(msg)
		if err != nil {
			fatalf("error translating message: %v", err)
		}
		_, err = DispatchMessage(msg)
		if err != nil {
			fatalf("error dispatching message: %v", err)
		}
	}
}

