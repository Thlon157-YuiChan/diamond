package view

import (
	"strconv"

	"github.com/chihaya/diamond/internal/app"
	"github.com/chihaya/diamond/internal/model"
	"github.com/mappu/miqt/qt"
)

// RouteEditor provides UI for editing routes and stations.
type RouteEditor struct {
	w             *qt.QWidget
	app           *app.App
	onChanged     func()
	loading       bool
	routeList     *qt.QListWidget
	stationRouteCb *qt.QComboBox
	list          *qt.QListWidget
	routeName *qt.QLineEdit
	nameEd    *qt.QLineEdit
	posEd     *qt.QLineEdit
	platSp    *qt.QSpinBox
	equalCb   *qt.QCheckBox
	dirCb     *qt.QComboBox
}

// NewRouteEditor creates a route editor.
func NewRouteEditor(a *app.App) *RouteEditor {
	ret := &RouteEditor{app: a}
	ret.w = qt.NewQWidget2()
	layout := qt.NewQVBoxLayout(ret.w)

	tabs := qt.NewQTabWidget(ret.w)

	// --- Route settings tab ---
	routeTab := qt.NewQWidget2()
	routeLayout := qt.NewQVBoxLayout(routeTab)

	split := qt.NewQHBoxLayout2()
	ret.routeList = qt.NewQListWidget(routeTab)
	ret.routeList.OnCurrentRowChanged(func(row int) {
		if ret.loading {
			return
		}
		if row < 0 || row >= len(ret.app.Project.Routes) {
			return
		}
		ret.loading = true
		ret.app.RouteIndex = row
		ret.stationRouteCb.SetCurrentIndex(row)
		ret.loading = false
		ret.refreshRouteFields()
		ret.refreshStations()
	})
	split.AddWidget(ret.routeList.QWidget)

	rform := qt.NewQFormLayout2()
	ret.routeName = qt.NewQLineEdit2()
	ret.routeName.OnTextChanged(func(s string) { ret.applyRouteName() })
	rform.AddRow3("Route Name", ret.routeName.QWidget)
	ret.dirCb = qt.NewQComboBox2()
	ret.dirCb.AddItem("Up")
	ret.dirCb.AddItem("Down")
	ret.dirCb.OnCurrentIndexChanged(func(idx int) { ret.applyRoute() })
	rform.AddRow3("Direction", ret.dirCb.QWidget)
	ret.equalCb = qt.NewQCheckBox3("Equal spacing")
	ret.equalCb.OnStateChanged(func(v int) { ret.applyRoute() })
	rform.AddRow3("", ret.equalCb.QWidget)
	split.AddLayout(rform.QLayout)
	routeLayout.AddLayout(split.QLayout)

	routeBtnRow := qt.NewQHBoxLayout2()
	addRouteBtn := qt.NewQPushButton5("Add Route", routeTab)
	addRouteBtn.OnClicked(func() {
		idx := len(ret.app.Project.Routes) + 1
		ret.app.AddRoute(model.Route{ID: int64(idx), Name: "Route " + strconv.Itoa(idx), Direction: model.DirectionUp, Stations: []model.Station{}})
		ret.refreshRoutes()
		ret.routeList.SetCurrentRow(len(ret.app.Project.Routes) - 1)
		ret.app.Dirty = true
	})
	routeBtnRow.AddWidget(addRouteBtn.QWidget)

	delRouteBtn := qt.NewQPushButton5("Remove Route", routeTab)
	delRouteBtn.OnClicked(func() {
		routes := ret.app.Project.Routes
		row := ret.routeList.CurrentRow()
		if row < 0 || row >= len(routes) {
			return
		}
		ret.app.Project.Routes = append(routes[:row], routes[row+1:]...)
		if ret.app.RouteIndex >= len(ret.app.Project.Routes) {
			ret.app.RouteIndex = len(ret.app.Project.Routes) - 1
		}
		if ret.app.RouteIndex < 0 {
			ret.app.RouteIndex = 0
		}
		ret.refreshRoutes()
		if len(ret.app.Project.Routes) > 0 {
			newRow := row
			if newRow >= len(ret.app.Project.Routes) {
				newRow = len(ret.app.Project.Routes) - 1
			}
			ret.routeList.SetCurrentRow(newRow)
		}
		ret.app.Dirty = true
	})
	routeBtnRow.AddWidget(delRouteBtn.QWidget)
	routeLayout.AddLayout(routeBtnRow.QLayout)

	tabs.AddTab(routeTab, "Route Settings")

	// --- Station settings tab ---
	stationTab := qt.NewQWidget2()
	stationLayout := qt.NewQVBoxLayout(stationTab)

	// target route selector
	routeSelRow := qt.NewQHBoxLayout2()
	routeSelLabel := qt.NewQLabel3("Target Route:")
	routeSelRow.AddWidget(routeSelLabel.QWidget)
	ret.stationRouteCb = qt.NewQComboBox2()
	ret.stationRouteCb.OnCurrentIndexChanged(func(idx int) {
		if ret.loading {
			return
		}
		if idx < 0 || idx >= len(ret.app.Project.Routes) {
			return
		}
		ret.loading = true
		ret.app.RouteIndex = idx
		ret.routeList.SetCurrentRow(idx)
		ret.loading = false
		ret.refreshRouteFields()
		ret.refreshStations()
	})
	routeSelRow.AddWidget(ret.stationRouteCb.QWidget)
	stationLayout.AddLayout(routeSelRow.QLayout)

	split2 := qt.NewQHBoxLayout2()
	ret.list = qt.NewQListWidget(stationTab)
	ret.list.OnCurrentRowChanged(func(row int) { ret.loadStation(row) })
	split2.AddWidget(ret.list.QWidget)

	sform := qt.NewQFormLayout2()
	ret.nameEd = qt.NewQLineEdit2()
	ret.nameEd.OnTextChanged(func(s string) { ret.applyStation() })
	sform.AddRow3("Name", ret.nameEd.QWidget)
	ret.posEd = qt.NewQLineEdit2()
	ret.posEd.OnTextChanged(func(s string) { ret.applyStation() })
	sform.AddRow3("Position km", ret.posEd.QWidget)
	ret.platSp = qt.NewQSpinBox2()
	ret.platSp.SetMinimum(1)
	ret.platSp.SetMaximum(99)
	ret.platSp.OnValueChanged(func(v int) { ret.applyStation() })
	sform.AddRow3("Platforms", ret.platSp.QWidget)
	split2.AddLayout(sform.QLayout)
	stationLayout.AddLayout(split2.QLayout)

	btnRow := qt.NewQHBoxLayout2()
	addStBtn := qt.NewQPushButton5("Add Station", stationTab)
	addStBtn.OnClicked(func() {
		r := ret.app.CurrentRoute()
		if r == nil {
			return
		}
		r.Stations = append(r.Stations, model.Station{ID: int64(len(r.Stations) + 1), Name: "New Station", Platforms: 1})
		ret.refreshStations()
		ret.list.SetCurrentRow(len(r.Stations) - 1)
		ret.app.Dirty = true
	})
	btnRow.AddWidget(addStBtn.QWidget)

	delStBtn := qt.NewQPushButton5("Remove", stationTab)
	delStBtn.OnClicked(func() {
		r := ret.app.CurrentRoute()
		if r == nil {
			return
		}
		row := ret.list.CurrentRow()
		if row < 0 || row >= len(r.Stations) {
			return
		}
		r.Stations = append(r.Stations[:row], r.Stations[row+1:]...)
		ret.refreshStations()
		if len(r.Stations) > 0 {
			newRow := row
			if newRow >= len(r.Stations) {
				newRow = len(r.Stations) - 1
			}
			ret.list.SetCurrentRow(newRow)
		}
		ret.app.Dirty = true
	})
	btnRow.AddWidget(delStBtn.QWidget)
	stationLayout.AddLayout(btnRow.QLayout)

	tabs.AddTab(stationTab, "Station Settings")

	layout.AddWidget(tabs.QWidget)

	ret.refreshRoutes()
	ret.refreshRouteFields()
	ret.refreshStations()
	return ret
}

// Widget returns the editor widget.
func (r *RouteEditor) Widget() *qt.QWidget { return r.w }

func (r *RouteEditor) refreshRoutes() {
	r.loading = true
	r.routeList.Clear()
	r.stationRouteCb.Clear()
	for _, rt := range r.app.Project.Routes {
		r.routeList.AddItem(rt.Name)
		r.stationRouteCb.AddItem(rt.Name)
	}
	if len(r.app.Project.Routes) > 0 {
		if r.app.RouteIndex < 0 || r.app.RouteIndex >= len(r.app.Project.Routes) {
			r.app.RouteIndex = 0
		}
		r.routeList.SetCurrentRow(r.app.RouteIndex)
		r.stationRouteCb.SetCurrentIndex(r.app.RouteIndex)
	}
	r.loading = false
	r.refreshRouteFields()
	r.refreshStations()
}

func (r *RouteEditor) refreshRouteFields() {
	r.loading = true
	rt := r.app.CurrentRoute()
	if rt == nil {
		r.loading = false
		return
	}
	r.routeName.SetText(rt.Name)
	if rt.Direction == model.DirectionDown {
		r.dirCb.SetCurrentIndex(1)
	} else {
		r.dirCb.SetCurrentIndex(0)
	}
	r.equalCb.SetChecked(rt.EqualSpacing)
	r.loading = false
}

func (r *RouteEditor) refreshStations() {
	r.loading = true
	r.list.Clear()
	rt := r.app.CurrentRoute()
	if rt == nil {
		r.loading = false
		return
	}
	for _, st := range rt.Stations {
		r.list.AddItem(st.Name)
	}
	r.loading = false
	if r.onChanged != nil {
		r.onChanged()
	}
}

func (r *RouteEditor) loadStation(row int) {
	rt := r.app.CurrentRoute()
	if rt == nil || row < 0 || row >= len(rt.Stations) {
		return
	}
	r.loading = true
	st := rt.Stations[row]
	r.nameEd.SetText(st.Name)
	r.posEd.SetText(formatFloat(st.PositionKm))
	r.platSp.SetValue(st.Platforms)
	r.loading = false
}

func (r *RouteEditor) applyRouteName() {
	if r.loading {
		return
	}
	rt := r.app.CurrentRoute()
	if rt == nil {
		return
	}
	rt.Name = r.routeName.Text()
	row := r.routeList.CurrentRow()
	if row >= 0 && row < r.routeList.Count() {
		item := r.routeList.Item(row)
		if item != nil {
			item.SetText(rt.Name)
		}
	}
	// keep target route combo in sync
	if row >= 0 && row < r.stationRouteCb.Count() {
		r.stationRouteCb.SetItemText(row, rt.Name)
	}
	r.app.Dirty = true
}

func (r *RouteEditor) applyRoute() {
	if r.loading {
		return
	}
	rt := r.app.CurrentRoute()
	if rt == nil {
		return
	}
	rt.EqualSpacing = r.equalCb.IsChecked()
	if r.dirCb.CurrentIndex() == 1 {
		rt.Direction = model.DirectionDown
	} else {
		rt.Direction = model.DirectionUp
	}
	r.app.Dirty = true
	if r.onChanged != nil {
		r.onChanged()
	}
}

func (r *RouteEditor) applyStation() {
	if r.loading {
		return
	}
	rt := r.app.CurrentRoute()
	if rt == nil {
		return
	}
	row := r.list.CurrentRow()
	if row < 0 || row >= len(rt.Stations) {
		return
	}
	pos, _ := strconv.ParseFloat(r.posEd.Text(), 64)
	rt.Stations[row].Name = r.nameEd.Text()
	rt.Stations[row].PositionKm = pos
	rt.Stations[row].Platforms = r.platSp.Value()
	item := r.list.Item(row)
	if item != nil {
		item.SetText(rt.Stations[row].Name)
	}
	r.app.Dirty = true
	if r.onChanged != nil {
		r.onChanged()
	}
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}