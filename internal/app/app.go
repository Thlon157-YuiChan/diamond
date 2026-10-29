package app

import (
	"path/filepath"

	"github.com/chihaya/diamond/internal/io"
	"github.com/chihaya/diamond/internal/model"
)

// App holds application state.
type App struct {
	Project      *model.Project
	FilePath     string
	RouteIndex   int
	DiagramIndex int
	Dirty        bool
}

// New creates a new empty application state.
func New() *App {
	a := &App{Project: &model.Project{Version: model.CurrentSchemaVersion}}
	a.Project.Routes = append(a.Project.Routes, model.Route{ID: 1, Name: "Route 1", Direction: model.DirectionUp, Stations: []model.Station{}})
	a.Project.Diagrams = append(a.Project.Diagrams, model.Diagram{ID: 1, Name: "平日", RouteID: 1, Trains: []model.Train{}})
	return a
}

// Load loads a project from path.
func (a *App) Load(path string) error {
	p, err := io.LoadProject(path)
	if err != nil {
		return err
	}
	a.Project = p
	a.FilePath = path
	a.Dirty = false
	if a.RouteIndex >= len(p.Routes) {
		a.RouteIndex = 0
	}
	if a.DiagramIndex >= len(p.Diagrams) {
		a.DiagramIndex = 0
	}
	return nil
}

// Save saves the project to its path (or given path).
func (a *App) Save(path string) error {
	if path == "" {
		path = a.FilePath
	}
	if path == "" {
		return nil
	}
	if err := io.SaveProject(path, a.Project); err != nil {
		return err
	}
	a.FilePath = path
	a.Dirty = false
	return nil
}

// CurrentRoute returns the currently selected route (or nil).
func (a *App) CurrentRoute() *model.Route {
	if a.Project == nil || a.RouteIndex < 0 || a.RouteIndex >= len(a.Project.Routes) {
		return nil
	}
	return &a.Project.Routes[a.RouteIndex]
}

// CurrentDiagram returns the currently selected diagram (or nil).
func (a *App) CurrentDiagram() *model.Diagram {
	if a.Project == nil || a.DiagramIndex < 0 || a.DiagramIndex >= len(a.Project.Diagrams) {
		return nil
	}
	return &a.Project.Diagrams[a.DiagramIndex]
}

// AddRoute adds a new route and selects it.
func (a *App) AddRoute(r model.Route) {
	a.Project.Routes = append(a.Project.Routes, r)
	a.RouteIndex = len(a.Project.Routes) - 1
	a.Dirty = true
}

// AddDiagram adds a new diagram and selects it.
func (a *App) AddDiagram(d model.Diagram) {
	a.Project.Diagrams = append(a.Project.Diagrams, d)
	a.DiagramIndex = len(a.Project.Diagrams) - 1
	a.Dirty = true
}

// BaseName returns file base name.
func (a *App) BaseName() string {
	if a.FilePath == "" {
		return "untitled.json"
	}
	return filepath.Base(a.FilePath)
}
