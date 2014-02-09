// 7 february 2014
package main

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

func wndProc(hwnd HWND, msg uint32, wparam WPARAM, lparam LPARAM) LRESULT {
	switch msg {
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
		return DefWindowProc(hwnd, msg, wparam, lparam)
	}
	fatalf("major bug: forgot a return on wndProc for message %d", msg)
	panic("unreachable")
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
		HbrBackground:	HBRUSH(COLOR_WINDOW + 1),
	}
	_, err = RegisterClass(wc)
	if err != nil {
		fatalf("error registering window class: %v", err)
	}

	hwnd, err := CreateWindowEx(
		WS_EX_OVERLAPPEDWINDOW,
		className, "Main Window",
		WS_OVERLAPPEDWINDOW,
		CW_USEDEFAULT, CW_USEDEFAULT, 320, 240,
		NULL, NULL, hInstance, NULL)
	if err != nil {
		fatalf("error creating window: %v", err)
	}

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

