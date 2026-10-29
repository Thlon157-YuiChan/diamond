package view

import (
	"os"

	"github.com/chihaya/diamond/internal/app"
	"github.com/chihaya/diamond/internal/model"
	"github.com/mappu/miqt/qt"
)

// MainWindow wraps the application's main window.
type MainWindow struct {
	w           *qt.QMainWindow
	app         *app.App
	routeEditor *RouteEditor
	diagramView *DiagramView
}

// NewMainWindow creates and configures the main application window.
func NewMainWindow(a *app.App) *MainWindow {
	ret := &MainWindow{app: a}
	ret.w = qt.NewQMainWindow2()
	ret.w.SetWindowTitle("Diamond - Railway Operation Diagram")
	ret.w.Resize(1200, 800)

	// Menu bar
	mnu := qt.NewQMenuBar2()

	fileMenu := mnu.AddMenuWithTitle("&File")
	newAct := fileMenu.AddAction("&New")
	newAct.SetShortcut(qt.NewQKeySequence2("Ctrl+N"))
	newAct.OnTriggered(func() {
		ret.app.Project = &model.Project{Version: model.CurrentSchemaVersion}
		ret.app.Project.Routes = append(ret.app.Project.Routes, model.Route{ID: 1, Name: "Route 1", Direction: model.DirectionUp, Stations: []model.Station{}})
		ret.app.Project.Diagrams = append(ret.app.Project.Diagrams, model.Diagram{ID: 1, Name: "平日", RouteID: 1, Trains: []model.Train{}})
		ret.app.RouteIndex = 0
		ret.app.DiagramIndex = 0
		ret.app.FilePath = ""
		ret.app.Dirty = false
		if ret.routeEditor != nil {
			ret.routeEditor.refreshRoutes()
			ret.routeEditor.refreshStations()
		}
		if ret.diagramView != nil {
			ret.diagramView.Refresh()
		}
	})
	openAct := fileMenu.AddAction("&Open...")
	openAct.SetShortcut(qt.NewQKeySequence2("Ctrl+O"))
	openAct.OnTriggered(func() { /* TODO: file dialog */ })
	saveAct := fileMenu.AddAction("&Save")
	saveAct.SetShortcut(qt.NewQKeySequence2("Ctrl+S"))
	saveAct.OnTriggered(func() {
		_ = ret.app.Save(ret.app.FilePath)
	})
	fileMenu.AddSeparator()
	exit := fileMenu.AddAction("E&xit")
	exit.SetShortcut(qt.NewQKeySequence2("Ctrl+Q"))
	exit.OnTriggered(func() { os.Exit(0) })

	viewMenu := mnu.AddMenuWithTitle("&View")
	themeAct := viewMenu.AddAction("Toggle &Theme")
	themeAct.OnTriggered(func() { ApplyTheme(ThemeDark) })

	helpMenu := mnu.AddMenuWithTitle("&Help")
	aboutQt := helpMenu.AddAction("About &Qt")
	aboutQt.OnTriggered(func() { qt.QApplication_AboutQt() })

	ret.w.SetMenuBar(mnu)

	// Central tabs
	tabs := qt.NewQTabWidget(ret.w.QWidget)
	re := NewRouteEditor(ret.app)
	dv := NewDiagramView(ret.app)
	ret.routeEditor = re
	ret.diagramView = dv
	re.onChanged = func() { dv.Refresh() }
	tabs.AddTab(re.Widget(), "Route")
	tabs.AddTab(dv.Widget(), "Diagram")
	tabs.OnCurrentChanged(func(index int) {
		if index == 1 {
			dv.Refresh()
		}
	})
	ret.w.SetCentralWidget(tabs.QWidget)

	// Status bar
	sb := ret.w.StatusBar()
	sb.ShowMessage2("Ready", 0)

	return ret
}

// Show displays the window.
func (m *MainWindow) Show() {
	m.w.Show()
	if m.diagramView != nil {
		m.diagramView.Refresh()
	}
}
