package fonts

import (
	"embed"
	"fmt"
	"gioui.org/font"
)

//go:embed fonts/*
var fonts embed.FS

func Prepare() ([]font.FontFace, error) {

}

func getFont(path string) ([]byte, error) {
	data, err := fonts.ReadFile(fmt.Sprintf("fonts/%s", path))
	if err != nil {

	}

}
