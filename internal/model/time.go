package model

import (
	"fmt"
	"strconv"
	"strings"
)

// Time は 15秒を 1 ティックとする内部時刻単位。
// 1日 = 24*60*4 = 5760 ティック。翌日超えは 5760 以上で表現する。
type Time int

const (
	SecondsPerTick = 15
	TicksPerMinute = 60 / SecondsPerTick
	TicksPerHour   = 60 * 60 / SecondsPerTick
	TicksPerDay    = 24 * 60 * 60 / SecondsPerTick
)

// Seconds は Time を秒に変換する。
func (t Time) Seconds() int { return int(t) * SecondsPerTick }

// Minutes は Time を分（切り捨て）に変換する。
func (t Time) Minutes() int { return int(t) / TicksPerMinute }

// String は "HH:MM" 形式の文字列を返す（15秒未満は切り捨て）。
func (t Time) String() string {
	sec := t.Seconds()
	h := (sec / 3600) % 24
	m := (sec / 60) % 60
	return fmt.Sprintf("%02d:%02d", h, m)
}

// ParseTime は "HH:MM" または "HH:MM:SS" を Time に変換する。
// 秒は 15秒単位にスナップする。
func ParseTime(s string) (Time, bool) {
	parts := strings.Split(s, ":")
	var h, m, sec int
	var err error
	switch len(parts) {
	case 2:
		h, err = strconv.Atoi(parts[0])
		if err != nil {
			return 0, false
		}
		m, err = strconv.Atoi(parts[1])
		if err != nil {
			return 0, false
		}
	case 3:
		h, err = strconv.Atoi(parts[0])
		if err != nil {
			return 0, false
		}
		m, err = strconv.Atoi(parts[1])
		if err != nil {
			return 0, false
		}
		sec, err = strconv.Atoi(parts[2])
		if err != nil {
			return 0, false
		}
	default:
		return 0, false
	}
	if h < 0 || h > 23 || m < 0 || m > 59 || sec < 0 || sec > 59 {
		return 0, false
	}
	totalSec := h*3600 + m*60 + sec
	ticks := (totalSec + SecondsPerTick/2) / SecondsPerTick
	return Time(ticks), true
}
