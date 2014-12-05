// 5 december 2014

static handlerfunc keyDownHandlers[] = {
	NULL,
};

static handlerfunc keyUpHandlers[] = {
	NULL,
};

static handlerfunc charHandlers[] = {
	NULL,
};

static handlerfunc mouseMoveHandlers[] = {
	NULL,
};

static handlerfunc mouseLeaveHandlers[] = {
	NULL,
};

static handlerfunc lbuttonDownHandlers[] = {
	NULL,
};

static handlerufnc lbuttonUpHandlers[] = {
	NULL,
};

static handlerfunc mouseWheelHandlers[] = {
	NULL,
};

// TODO WM_MOUSEHOVER, other mouse buttons

HANDLER(events)
{
	switch (uMsg) {
#define eventHandler(msg, array) \
	case msg: \
		return runHandlers(array, t, uMsg, wParam, lParam, lResult);
	eventHandler(WM_KEYDOWN, keyDownHandlers)
	eventHandler(WM_KEYUP, keyUpHandlers)
	eventHandler(WM_CHAR, charHandlers)
	eventHandler(WM_MOUSEMOVE, mouseMoveHandlers)
	eventHandler(WM_MOUSELEAVE, mouseLeaveHandlers)
	eventHandler(WM_LBUTTONDOWN, lbuttonDownHandlers)
	eventHandler(WM_LBUTTONUP, lbuttonUpHandlers)
	eventHandler(WM_MOUSEWHEEL, mouseWheelHandlers)
#undef eventHandler
	}
	return FALSE;
}
