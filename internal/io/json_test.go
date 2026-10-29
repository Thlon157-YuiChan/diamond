package io

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chihaya/diamond/internal/model"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.json")
	p := &model.Project{
		Routes:   []model.Route{{ID: 1, Name: "R", Stations: []model.Station{{ID: 1, Name: "S", Platforms: 2}}}},
		Diagrams: []model.Diagram{{ID: 1, Name: "平日", RouteID: 1, Trains: []model.Train{{ID: 1, Number: "101", Stops: []model.Stop{{StationID: 1, Arrival: 0, Departure: 4}}}}}},
	}
	if err := SaveProject(path, p); err != nil {
		t.Fatal(err)
	}
	got, err := LoadProject(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Version != model.CurrentSchemaVersion {
		t.Errorf("version = %d", got.Version)
	}
	if len(got.Routes) != 1 || len(got.Diagrams[0].Trains) != 1 {
		t.Errorf("round trip mismatch")
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
}
