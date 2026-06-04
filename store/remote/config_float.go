package remote

import (
	"math"
	"time"
)

const (
	floatColorIndexMin = 0
	floatColorIndexMax = 13
)

// FloatPageConfig 与前端 Float 钟表配置字段一一对应。
type FloatPageConfig struct {
	AutoNightMode      bool     `json:"autoNightMode"`
	NightModeRange     []string `json:"nightModeRange"`
	AutoInvisible      bool     `json:"autoInvisible"`
	InvisibleRange     []string `json:"invisibleRange"`
	InvisibleDayEnable bool     `json:"invisibleDayEnable"`
	InvisibleDay       []bool   `json:"invisibleDay"`
	SelectedColorIndex int      `json:"selectedColorIndex"`
	Brightness         float64  `json:"brightness"`
}

// FloatPageConfigPatch 仅包含需要远程修改的字段（指针非 nil 表示要更新）。
type FloatPageConfigPatch struct {
	AutoNightMode      *bool     `json:"autoNightMode"`
	NightModeRange     *[]string `json:"nightModeRange"`
	AutoInvisible      *bool     `json:"autoInvisible"`
	InvisibleRange     *[]string `json:"invisibleRange"`
	InvisibleDayEnable *bool     `json:"invisibleDayEnable"`
	InvisibleDay       *[]bool   `json:"invisibleDay"`
	SelectedColorIndex *int      `json:"selectedColorIndex"`
	Brightness         *float64  `json:"brightness"`
}

func (p *FloatPageConfigPatch) IsEmpty() bool {
	if p == nil {
		return true
	}
	return p.AutoNightMode == nil &&
		p.NightModeRange == nil &&
		p.AutoInvisible == nil &&
		p.InvisibleRange == nil &&
		p.InvisibleDayEnable == nil &&
		p.InvisibleDay == nil &&
		p.SelectedColorIndex == nil &&
		p.Brightness == nil
}

func defaultFloatConfig() *FloatPageConfig {
	return &FloatPageConfig{
		AutoNightMode:      false,
		NightModeRange:     []string{"2025-01-01T00:00:00.000Z", "2025-01-01T06:00:00.000Z"},
		AutoInvisible:      false,
		InvisibleRange:     []string{"2025-01-01T00:00:00.000Z", "2025-01-01T06:00:00.000Z"},
		InvisibleDayEnable: false,
		InvisibleDay:       []bool{false, false, false, false, false, false, false},
		SelectedColorIndex: 0,
		Brightness:         1,
	}
}

func applyPatch(base *FloatPageConfig, patch *FloatPageConfigPatch) *FloatPageConfig {
	out := cloneFloatConfig(base)
	if patch.AutoNightMode != nil {
		out.AutoNightMode = *patch.AutoNightMode
	}
	if patch.NightModeRange != nil {
		out.NightModeRange = append([]string(nil), (*patch.NightModeRange)...)
	}
	if patch.AutoInvisible != nil {
		out.AutoInvisible = *patch.AutoInvisible
	}
	if patch.InvisibleRange != nil {
		out.InvisibleRange = append([]string(nil), (*patch.InvisibleRange)...)
	}
	if patch.InvisibleDayEnable != nil {
		out.InvisibleDayEnable = *patch.InvisibleDayEnable
	}
	if patch.InvisibleDay != nil {
		out.InvisibleDay = append([]bool(nil), (*patch.InvisibleDay)...)
	}
	if patch.SelectedColorIndex != nil {
		out.SelectedColorIndex = *patch.SelectedColorIndex
	}
	if patch.Brightness != nil {
		out.Brightness = *patch.Brightness
	}
	return out
}

func validatePatch(patch *FloatPageConfigPatch) error {
	if patch.NightModeRange != nil && !validateTimeRange(*patch.NightModeRange) {
		return ErrInvalidConfig
	}
	if patch.InvisibleRange != nil && !validateTimeRange(*patch.InvisibleRange) {
		return ErrInvalidConfig
	}
	if patch.InvisibleDay != nil && !validateInvisibleWeek(*patch.InvisibleDay) {
		return ErrInvalidConfig
	}
	if patch.SelectedColorIndex != nil && !colorIndexInRange(*patch.SelectedColorIndex) {
		return ErrInvalidConfig
	}
	if patch.Brightness != nil && !brightnessInRange(*patch.Brightness) {
		return ErrInvalidConfig
	}
	return nil
}

func floatEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-6
}

func equalStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func equalBoolSlice(a, b []bool) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// EqualFloatConfig 逐字段比较配置，不用 JSON 字符串比较。
func EqualFloatConfig(a, b *FloatPageConfig) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.AutoNightMode == b.AutoNightMode &&
		equalStringSlice(a.NightModeRange, b.NightModeRange) &&
		a.AutoInvisible == b.AutoInvisible &&
		equalStringSlice(a.InvisibleRange, b.InvisibleRange) &&
		a.InvisibleDayEnable == b.InvisibleDayEnable &&
		equalBoolSlice(a.InvisibleDay, b.InvisibleDay) &&
		a.SelectedColorIndex == b.SelectedColorIndex &&
		floatEqual(a.Brightness, b.Brightness)
}

func brightnessInRange(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0) && v >= 0 && v <= 1
}

// validateTimeRange 起止时间段：2 个可解析的时钟时刻（与前端 ISO 字符串一致）。
func validateTimeRange(pair []string) bool {
	if len(pair) != 2 {
		return false
	}
	for _, s := range pair {
		if !isValidClockTime(s) {
			return false
		}
	}
	return true
}

var clockTimeLayouts = []string{
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02T15:04:05Z07:00",
	"15:04:05",
	"15:04",
}

func isValidClockTime(s string) bool {
	if s == "" {
		return false
	}
	for _, layout := range clockTimeLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			h, m, sec := t.Clock()
			return h >= 0 && h <= 23 && m >= 0 && m <= 59 && sec >= 0 && sec <= 59
		}
	}
	return false
}

// validateInvisibleWeek 全天关闭显示：固定 7 天（周日–周六），均为布尔值。
func validateInvisibleWeek(days []bool) bool {
	return len(days) == 7
}

func cloneFloatConfig(c *FloatPageConfig) *FloatPageConfig {
	if c == nil {
		return nil
	}
	cp := *c
	cp.NightModeRange = append([]string(nil), c.NightModeRange...)
	cp.InvisibleRange = append([]string(nil), c.InvisibleRange...)
	cp.InvisibleDay = append([]bool(nil), c.InvisibleDay...)
	return &cp
}

func validateFloatConfig(c *FloatPageConfig) bool {
	if c == nil {
		return false
	}
	return validateTimeRange(c.NightModeRange) &&
		validateTimeRange(c.InvisibleRange) &&
		validateInvisibleWeek(c.InvisibleDay) &&
		colorIndexInRange(c.SelectedColorIndex) &&
		brightnessInRange(c.Brightness)
}

func colorIndexInRange(i int) bool {
	return i >= floatColorIndexMin && i <= floatColorIndexMax
}
