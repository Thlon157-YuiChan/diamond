package view

import "github.com/mappu/miqt/qt"

// Theme represents light/dark theme.
type Theme int

const (
	ThemeLight Theme = iota
	ThemeDark
)

// ApplyTheme applies the given theme to the application.
func ApplyTheme(t Theme) {
	if t != ThemeDark {
		return
	}
	pal := qt.NewQPalette()
	pal.SetColor2(qt.QPalette__Window, qt.NewQColor3(45, 45, 48))
	pal.SetColor2(qt.QPalette__WindowText, qt.NewQColor3(220, 220, 220))
	pal.SetColor2(qt.QPalette__Base, qt.NewQColor3(30, 30, 30))
	pal.SetColor2(qt.QPalette__Text, qt.NewQColor3(220, 220, 220))
	pal.SetColor2(qt.QPalette__Button, qt.NewQColor3(60, 60, 60))
	pal.SetColor2(qt.QPalette__ButtonText, qt.NewQColor3(220, 220, 220))
	qt.QApplication_SetPalette(pal)
}
