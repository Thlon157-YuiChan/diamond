package main

import (
	"fmt"
	"os"

	"github.com/chihaya/diamond/internal/view"
	"github.com/mappu/miqt/qt"
)

func main() {
	qt.NewQApplication(os.Args)

	w := view.NewMainWindow()
	w.Show()

	qt.QApplication_Exec()
	fmt.Println("Diamond exited")
}
