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
	openAct.OnTriggered(func() {
		path := qt.QFileDialog_GetOpenFileName4(ret.w.QWidget, "Open Project", "", "JSON files (*.json);;All files (*)")
		if path == "" {
			return
		}
		if err := ret.app.Load(path); err != nil {
			ret.w.StatusBar().ShowMessage2("Open failed", 0)
			return
		}
		if ret.routeEditor != nil {
			ret.routeEditor.refreshRoutes()
			ret.routeEditor.refreshRouteFields()
			ret.routeEditor.refreshStations()
		}
		if ret.diagramView != nil {
			ret.diagramView.Refresh()
		}
		ret.w.StatusBar().ShowMessage2("Loaded: "+path, 0)
	})
	saveAct := fileMenu.AddAction("&Save")
	saveAct.SetShortcut(qt.NewQKeySequence2("Ctrl+S"))
	saveAct.OnTriggered(func() {
		path := ret.app.FilePath
		if path == "" {
			path = qt.QFileDialog_GetSaveFileName4(ret.w.QWidget, "Save Project", "", "JSON files (*.json);;All files (*)")
			if path == "" {
				return
			}
		}
		if err := ret.app.Save(path); err != nil {
			ret.w.StatusBar().ShowMessage2("Save failed", 0)
			return
		}
		ret.w.StatusBar().ShowMessage2("Saved: "+path, 0)
	})
	fileMenu.AddSeparator()
	exit := fileMenu.AddAction("E&xit")
	exit.SetShortcut(qt.NewQKeySequence2("Ctrl+Q"))
	exit.OnTriggered(func() { os.Exit(0) })

	viewMenu := mnu.AddMenuWithTitle("&View")
	themeAct := viewMenu.AddAction("Toggle &Theme")
	themeAct.SetShortcut(qt.NewQKeySequence2("Ctrl+T"))
	themeAct.OnTriggered(func() {
		if CurrentTheme() == ThemeDark {
			ApplyTheme(ThemeLight)
		} else {
			ApplyTheme(ThemeDark)
		}
		if ret.diagramView != nil {
			ret.diagramView.Refresh()
		}
		ret.w.StatusBar().ShowMessage2("Theme: "+themeName(CurrentTheme()), 0)
	})

	helpMenu := mnu.AddMenuWithTitle("&Help")
	aboutQt := helpMenu.AddAction("About &Qt")
	aboutQt.OnTriggered(func() { qt.QApplication_AboutQt() })

	// Toolbar (visible controls)
	tb := qt.NewQToolBar3()
	tb.SetMovable(false)
	newTb := tb.AddAction("New")
	newTb.SetShortcut(qt.NewQKeySequence2("Ctrl+N"))
	newTb.OnTriggered(func() {
		ret.app.Project = &model.Project{Version: model.CurrentSchemaVersion}
		ret.app.Project.Routes = append(ret.app.Project.Routes, model.Route{ID: 1, Name: "Route 1", Direction: model.DirectionUp, Stations: []model.Station{}})
		ret.app.Project.Diagrams = append(ret.app.Project.Diagrams, model.Diagram{ID: 1, Name: "平日", RouteID: 1, Trains: []model.Train{}})
		ret.app.RouteIndex = 0
		ret.app.DiagramIndex = 0
		ret.app.FilePath = ""
		ret.app.Dirty = false
		if ret.routeEditor != nil {
			ret.routeEditor.refreshRoutes()
			ret.routeEditor.refreshRouteFields()
			ret.routeEditor.refreshStations()
		}
		if ret.diagramView != nil {
			ret.diagramView.Refresh()
		}
	})
	openTb := tb.AddAction("Open")
	openTb.SetShortcut(qt.NewQKeySequence2("Ctrl+O"))
	openTb.OnTriggered(func() {
		path := qt.QFileDialog_GetOpenFileName4(ret.w.QWidget, "Open Project", "", "JSON files (*.json);;All files (*)")
		if path == "" {
			return
		}
		if err := ret.app.Load(path); err != nil {
			ret.w.StatusBar().ShowMessage2("Open failed", 0)
			return
		}
		if ret.routeEditor != nil {
			ret.routeEditor.refreshRoutes()
			ret.routeEditor.refreshRouteFields()
			ret.routeEditor.refreshStations()
		}
		if ret.diagramView != nil {
			ret.diagramView.Refresh()
		}
		ret.w.StatusBar().ShowMessage2("Loaded: "+path, 0)
	})
	saveTb := tb.AddAction("Save")
	saveTb.SetShortcut(qt.NewQKeySequence2("Ctrl+S"))
	saveTb.OnTriggered(func() {
		path := ret.app.FilePath
		if path == "" {
			path = qt.QFileDialog_GetSaveFileName4(ret.w.QWidget, "Save Project", "", "JSON files (*.json);;All files (*)")
			if path == "" {
				return
			}
		}
		if err := ret.app.Save(path); err != nil {
			ret.w.StatusBar().ShowMessage2("Save failed", 0)
			return
		}
		ret.w.StatusBar().ShowMessage2("Saved: "+path, 0)
	})
	themeTb := tb.AddAction("Theme")
	themeTb.SetShortcut(qt.NewQKeySequence2("Ctrl+T"))
	themeTb.OnTriggered(func() {
		if CurrentTheme() == ThemeDark {
			ApplyTheme(ThemeLight)
		} else {
			ApplyTheme(ThemeDark)
		}
		if ret.diagramView != nil {
			ret.diagramView.Refresh()
		}
		ret.w.StatusBar().ShowMessage2("Theme: "+themeName(CurrentTheme()), 0)
	})
	ret.w.AddToolBarWithToolbar(tb)

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
