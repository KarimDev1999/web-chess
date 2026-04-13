package chess

import "math/rand"

type TimeControl struct {
	Base      int `json:"base"`
	Increment int `json:"increment"`
}

func (tc TimeControl) IsTimed() bool {
	return tc.Base > 0
}

var (
	TimeControlBullet1_0      = TimeControl{Base: 60, Increment: 0}
	TimeControlBullet2_1      = TimeControl{Base: 120, Increment: 1}
	TimeControlBullet3_0      = TimeControl{Base: 180, Increment: 0}
	TimeControlBlitz3_0       = TimeControl{Base: 180, Increment: 0}
	TimeControlBlitz3_2       = TimeControl{Base: 180, Increment: 2}
	TimeControlBlitz5_0       = TimeControl{Base: 300, Increment: 0}
	TimeControlBlitz5_3       = TimeControl{Base: 300, Increment: 3}
	TimeControlBlitz10_0      = TimeControl{Base: 600, Increment: 0}
	TimeControlRapid10_0      = TimeControl{Base: 600, Increment: 0}
	TimeControlRapid10_5      = TimeControl{Base: 600, Increment: 5}
	TimeControlRapid15_10     = TimeControl{Base: 900, Increment: 10}
	TimeControlRapid25_0      = TimeControl{Base: 1500, Increment: 0}
	TimeControlClassical30_0  = TimeControl{Base: 1800, Increment: 0}
	TimeControlClassical30_20 = TimeControl{Base: 1800, Increment: 20}
	TimeControlClassical60_30 = TimeControl{Base: 3600, Increment: 30}
)

type TimeControlPreset struct {
	Label string      `json:"label"`
	TC    TimeControl `json:"tc"`
}

func StandardPresets() map[string][]TimeControlPreset {
	return map[string][]TimeControlPreset{
		"bullet": {
			{Label: "1+0", TC: TimeControlBullet1_0},
			{Label: "2+1", TC: TimeControlBullet2_1},
			{Label: "3+0", TC: TimeControlBullet3_0},
		},
		"blitz": {
			{Label: "3+0", TC: TimeControlBlitz3_0},
			{Label: "3+2", TC: TimeControlBlitz3_2},
			{Label: "5+0", TC: TimeControlBlitz5_0},
			{Label: "5+3", TC: TimeControlBlitz5_3},
			{Label: "10+0", TC: TimeControlBlitz10_0},
		},
		"rapid": {
			{Label: "10+0", TC: TimeControlRapid10_0},
			{Label: "10+5", TC: TimeControlRapid10_5},
			{Label: "15+10", TC: TimeControlRapid15_10},
			{Label: "25+0", TC: TimeControlRapid25_0},
		},
		"classical": {
			{Label: "30+0", TC: TimeControlClassical30_0},
			{Label: "30+20", TC: TimeControlClassical30_20},
			{Label: "60+30", TC: TimeControlClassical60_30},
		},
	}
}

type ColorPreference string

const (
	PreferenceWhite  ColorPreference = "white"
	PreferenceBlack  ColorPreference = "black"
	PreferenceRandom ColorPreference = "random"
)

func ResolveColor(pref ColorPreference) Color {
	switch pref {
	case PreferenceWhite:
		return White
	case PreferenceBlack:
		return Black
	default:
		if rand.Intn(2) == 0 {
			return White
		}
		return Black
	}
}
