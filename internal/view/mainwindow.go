package view

import (
	"os"

	"github.com/mappu/miqt/qt"
)

// MainWindow wraps the application's main window.
type MainWindow struct {
	w *qt.QMainWindow
}

// NewMainWindow creates and configures the main application window.
func NewMainWindow() *MainWindow {
	ret := &MainWindow{}
	ret.w = qt.NewQMainWindow2()
	ret.w.SetWindowTitle("Diamond - Railway Operation Diagram")
	ret.w.Resize(1200, 800)

	// Menu bar
	mnu := qt.NewQMenuBar2()

	fileMenu := mnu.AddMenuWithTitle("&File")
	exit := fileMenu.AddAction("E&xit")
	exit.SetShortcut(qt.NewQKeySequence2("Ctrl+Q"))
	exit.OnTriggered(func() {
		os.Exit(0)
	})

	helpMenu := mnu.AddMenuWithTitle("&Help")
	aboutQt := helpMenu.AddAction("About &Qt")
	aboutQt.OnTriggered(func() {
		qt.QApplication_AboutQt()
	})

	ret.w.SetMenuBar(mnu)

	// Central widget placeholder
	central := qt.NewQWidget2()
	ret.w.SetCentralWidget(central)

	// Status bar
	sb := ret.w.StatusBar()
	sb.ShowMessage2("Ready", 0)

	return ret
}

// Show displays the window.
func (m *MainWindow) Show() {
	m.w.Show()
}

// Widget returns the underlying QMainWindow.
func (m *MainWindow) Widget() *qt.QMainWindow {
	return m.w
}
