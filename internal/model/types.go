package model

// Direction は列車の進行方向。
type Direction int

const (
	DirectionUp Direction = iota
	DirectionDown
)

// ConnectionType は前後列車の接続種別。
type ConnectionType int

const (
	ConnNone       ConnectionType = iota // 0: 既定（非接続）
	ConnSameTrain                        // 同一列車（同方向継続）
	ConnTypeChange                       // 種別変更
	ConnJoinSplit                        // 増解結
	ConnTurnback                         // 折り返し（方向転換）
	ConnDepotIn                          // 入区
	ConnDepotOut                         // 出区
)

type Station struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	PositionKm float64 `json:"position_km"`
	Platforms  int     `json:"platforms"`
}

type Route struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Direction    Direction `json:"direction"`
	EqualSpacing bool      `json:"equal_spacing"`
	Stations     []Station `json:"stations"`
}

type TrainType struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Stop struct {
	StationID int64 `json:"station_id"`
	Arrival   Time  `json:"arrival"`
	Departure Time  `json:"departure"`
	Platform  int   `json:"platform"`
}

type Train struct {
	ID             int64          `json:"id"`
	Number         string         `json:"number"`
	TypeID         int64          `json:"type_id"`
	OperationID    int64          `json:"operation_id,omitempty"`
	ServiceDays    uint8          `json:"service_days"`
	Stops          []Stop         `json:"stops"`
	OriginOutside  bool           `json:"origin_outside"`
	DestOutside    bool           `json:"dest_outside"`
	ConnectionType ConnectionType `json:"connection_type"`
	PrevTrainID    int64          `json:"prev_train_id,omitempty"`
	NextTrainID    int64          `json:"next_train_id,omitempty"`
}

type Diagram struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	RouteID int64   `json:"route_id"`
	Trains  []Train `json:"trains"`
}

type Operation struct {
	ID       int64   `json:"id"`
	Number   string  `json:"number"`
	TrainIDs []int64 `json:"train_ids"`
}

type Project struct {
	Version    int         `json:"version"`
	Routes     []Route     `json:"routes"`
	TrainTypes []TrainType `json:"train_types"`
	Diagrams   []Diagram   `json:"diagrams"`
	Operations []Operation `json:"operations"`
}

// CurrentSchemaVersion is the JSON schema version written by this build.
const CurrentSchemaVersion = 1
