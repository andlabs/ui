// 5 december 2014

static const handlerfunc keyDownHandlers[] = {
	NULL,
};

static const handlerfunc keyUpHandlers[] = {
	NULL,
};

static const handlerfunc charHandlers[] = {
	NULL,
};

static const handlerfunc mouseMoveHandlers[] = {
	NULL,
};

static const handlerfunc mouseLeaveHandlers[] = {
	NULL,
};

static const handlerfunc lbuttonDownHandlers[] = {
	NULL,
};

static const handlerfunc lbuttonUpHandlers[] = {
	NULL,
};

static const handlerfunc mouseWheelHandlers[] = {
	NULL,
};

// TODO WM_MOUSEHOVER, other mouse buttons

HANDLER(eventHandlers)
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
