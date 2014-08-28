	case WM_MOUSELEAVE:
		tablePushed(t->gotable, -1, -1);			// in case button held as drag out
		// and let the list view do its thing
		return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
