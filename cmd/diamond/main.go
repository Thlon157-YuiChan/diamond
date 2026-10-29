package main

import (
	"os"

	"github.com/chihaya/diamond/internal/app"
	"github.com/chihaya/diamond/internal/view"
	"github.com/mappu/miqt/qt"
)

func main() {
	qt.NewQApplication(os.Args)

	a := app.New()
	w := view.NewMainWindow(a)
	w.Show()

	qt.QApplication_Exec()
}
