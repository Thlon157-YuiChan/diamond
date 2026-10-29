package view

import (
	"fmt"

	"github.com/chihaya/diamond/internal/app"
	"github.com/chihaya/diamond/internal/model"
	"github.com/mappu/miqt/qt"
)

const (
	stationHeight = 40
	hourWidth     = 120
	pixelsPerKm   = 50
	topMargin     = 20
	leftMargin    = 90
)

// DiagramView displays a railway operation diagram.
type DiagramView struct {
	w         *qt.QWidget
	scroll    *qt.QScrollArea
	canvas    *qt.QWidget
	app       *app.App
	pxPerTick float64
}

// NewDiagramView creates a diagram view bound to app state.
func NewDiagramView(a *app.App) *DiagramView {
	ret := &DiagramView{app: a, pxPerTick: float64(hourWidth) / float64(model.TicksPerHour)}
	ret.canvas = qt.NewQWidget2()
	ret.canvas.SetMinimumSize2(800, 400)
	ret.canvas.OnPaintEvent(ret.paint)
	ret.scroll = qt.NewQScrollArea2()
	ret.scroll.SetWidget(ret.canvas)
	ret.scroll.SetWidgetResizable(false)
	ret.w = ret.scroll.QAbstractScrollArea.QFrame.QWidget
	return ret
}

// Widget returns the scroll area widget.
func (d *DiagramView) Widget() *qt.QWidget { return d.w }

// Refresh triggers a repaint.
func (d *DiagramView) Refresh() {
	if d.canvas != nil {
		d.canvas.Update()
	}
}

func (d *DiagramView) paint(super func(*qt.QPaintEvent), ev *qt.QPaintEvent) {
	p := qt.NewQPainter2(d.canvas.QPaintDevice)
	defer p.Delete()

	route := d.app.CurrentRoute()
	diag := d.app.CurrentDiagram()

	p.FillRect5(0, 0, d.canvas.Width(), d.canvas.Height(), qt.NewQColor3(250, 250, 250))
	p.SetPen(qt.NewQColor3(0, 0, 0))

	if route == nil {
		return
	}

	numStations := len(route.Stations)
	canvasH := topMargin + stationHeight
	if route.EqualSpacing {
		canvasH = topMargin + numStations*stationHeight
	} else if numStations > 0 {
		canvasH = topMargin + int(route.Stations[numStations-1].PositionKm*pixelsPerKm) + stationHeight
	}
	canvasW := leftMargin + int(float64(model.TicksPerDay)*d.pxPerTick)
	if canvasW < 800 {
		canvasW = 800
	}
	if canvasH < 400 {
		canvasH = 400
	}
	d.canvas.SetMinimumSize2(canvasW, canvasH)

	for i := range route.Stations {
		y := stationCenterY(route, i)
		p.DrawText3(4, y+4, route.Stations[i].Name)
		// station baseline
		p.DrawLine2(leftMargin, y, canvasW, y)
	}

	for h := 0; h <= 24; h++ {
		x := leftMargin + int(float64(h*model.TicksPerHour)*d.pxPerTick)
		p.DrawLine2(x, topMargin, x, canvasH)
		label := fmt.Sprintf("%02d", h%24)
		p.DrawText3(x-6, topMargin-4, label)
	}

	if diag != nil {
		for _, tr := range diag.Trains {
			p.SetPen(qt.NewQColor3(0, 0, 180))
			for i := 1; i < len(tr.Stops); i++ {
				prev := tr.Stops[i-1]
				cur := tr.Stops[i]
				pi := stationIndex(route, prev.StationID)
				ci := stationIndex(route, cur.StationID)
				if pi < 0 || ci < 0 {
					continue
				}
				x1 := leftMargin + int(float64(prev.Departure)*d.pxPerTick)
				if x1 == leftMargin {
					x1 = leftMargin + int(float64(prev.Arrival)*d.pxPerTick)
				}
				x2 := leftMargin + int(float64(cur.Arrival)*d.pxPerTick)
				if x2 == leftMargin {
					x2 = leftMargin + int(float64(cur.Departure)*d.pxPerTick)
				}
				y1 := stationCenterY(route, pi)
				y2 := stationCenterY(route, ci)
				p.DrawLine2(x1, y1, x2, y2)
			}
		}
	}
}

func stationCenterY(route *model.Route, i int) int {
	if route.EqualSpacing {
		return topMargin + i*stationHeight + stationHeight/2
	}
	if i < 0 || i >= len(route.Stations) {
		return topMargin
	}
	return topMargin + int(route.Stations[i].PositionKm*pixelsPerKm)
}

func stationIndex(route *model.Route, id int64) int {
	for i, s := range route.Stations {
		if s.ID == id {
			return i
		}
	}
	return -1
}
