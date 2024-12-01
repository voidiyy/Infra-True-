package app

import (
	"Infra/gui/appTheme"
	"gioui.org/app"
)

type UI struct {
	Theme appTheme.Theme

	window  *app.Window
	sidebar *appSidebar
	header  *appHeader
}

func New(w *app.Window, appVersion string) *UI {
	u := &UI{
		window: w,
	}

	fontCollection :=
}
