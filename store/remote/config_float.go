package remote

import (
	"bytes"
	"encoding/json"
	"math"
	"time"
)

const (
	colorIndexMin        = 0
	colorIndexMax        = 99
	worldSelectionRowMax = 80
	worldSelectionColMax = 160
)

// WorldSelection 世界时钟上被选中的点。
type WorldSelection struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

// FloatPageConfig 与前端钟表配置字段一一对应。
type FloatPageConfig struct {
	AutoNightMode      bool            `json:"autoNightMode"`
	NightModeRange     []string        `json:"nightModeRange"`
	AutoInvisible      bool            `json:"autoInvisible"`
	InvisibleRange     []string        `json:"invisibleRange"`
	InvisibleDayEnable bool            `json:"invisibleDayEnable"`
	InvisibleDay       []bool          `json:"invisibleDay"`
	ClockStyle         string          `json:"clockStyle"`
	ColorIndexMap      map[string]int  `json:"colorIndexMap"`
	Brightness         float64         `json:"brightness"`
	NumericalStocks    []string        `json:"numericalStocks"`
	WorldSelected      *WorldSelection `json:"worldSelected"`
}

// worldSelectionPatch 区分「请求没带这个字段」和「明确取消选中」。
type worldSelectionPatch struct {
	Set   bool
	Value *WorldSelection
}

func (p *worldSelectionPatch) UnmarshalJSON(data []byte) error {
	p.Set = true
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		p.Value = nil
		return nil
	}
	var value WorldSelection
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	p.Value = &value
	return nil
}

// FloatPageConfigPatch 仅包含需要远程修改的字段（指针非 nil 表示要更新）。
type FloatPageConfigPatch struct {
	AutoNightMode      *bool               `json:"autoNightMode"`
	NightModeRange     *[]string           `json:"nightModeRange"`
	AutoInvisible      *bool               `json:"autoInvisible"`
	InvisibleRange     *[]string           `json:"invisibleRange"`
	InvisibleDayEnable *bool               `json:"invisibleDayEnable"`
	InvisibleDay       *[]bool             `json:"invisibleDay"`
	ClockStyle         *string             `json:"clockStyle"`
	ColorIndexMap      *map[string]int     `json:"colorIndexMap"`
	Brightness         *float64            `json:"brightness"`
	NumericalStocks    *[]string           `json:"numericalStocks"`
	WorldSelected      worldSelectionPatch `json:"worldSelected"`
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
		p.ClockStyle == nil &&
		p.ColorIndexMap == nil &&
		p.Brightness == nil &&
		p.NumericalStocks == nil &&
		!p.WorldSelected.Set
}

func defaultFloatConfig() *FloatPageConfig {
	return &FloatPageConfig{
		AutoNightMode:      false,
		NightModeRange:     []string{"2025-01-01T00:00:00.000Z", "2025-01-01T06:00:00.000Z"},
		AutoInvisible:      false,
		InvisibleRange:     []string{"2025-01-01T00:00:00.000Z", "2025-01-01T06:00:00.000Z"},
		InvisibleDayEnable: false,
		InvisibleDay:       []bool{false, false, false, false, false, false, false},
		ClockStyle:         "float",
		ColorIndexMap:      map[string]int{"float": 0, "numerical": 0, "analog": 0, "world": 0},
		Brightness:         1,
		NumericalStocks:    []string{},
		WorldSelected:      nil,
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
	if patch.ClockStyle != nil {
		out.ClockStyle = *patch.ClockStyle
	}
	if patch.ColorIndexMap != nil {
		out.ColorIndexMap = cloneIntMap(*patch.ColorIndexMap)
	}
	if patch.Brightness != nil {
		out.Brightness = *patch.Brightness
	}
	if patch.NumericalStocks != nil {
		out.NumericalStocks = append([]string(nil), (*patch.NumericalStocks)...)
	}
	if patch.WorldSelected.Set {
		out.WorldSelected = cloneWorldSelection(patch.WorldSelected.Value)
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
	if patch.ClockStyle != nil && !validateClockStyle(*patch.ClockStyle) {
		return ErrInvalidConfig
	}
	if patch.ColorIndexMap != nil && !validateColorIndexMap(*patch.ColorIndexMap) {
		return ErrInvalidConfig
	}
	if patch.Brightness != nil && !brightnessInRange(*patch.Brightness) {
		return ErrInvalidConfig
	}
	if patch.NumericalStocks != nil && !validateNumericalStocks(*patch.NumericalStocks) {
		return ErrInvalidConfig
	}
	if patch.WorldSelected.Set && !validateWorldSelection(patch.WorldSelected.Value) {
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
		a.ClockStyle == b.ClockStyle &&
		equalIntMap(a.ColorIndexMap, b.ColorIndexMap) &&
		floatEqual(a.Brightness, b.Brightness) &&
		equalStringSlice(a.NumericalStocks, b.NumericalStocks) &&
		equalWorldSelection(a.WorldSelected, b.WorldSelected)
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
	cp.ColorIndexMap = cloneIntMap(c.ColorIndexMap)
	cp.NumericalStocks = append([]string(nil), c.NumericalStocks...)
	cp.WorldSelected = cloneWorldSelection(c.WorldSelected)
	return &cp
}

func validateFloatConfig(c *FloatPageConfig) bool {
	if c == nil {
		return false
	}
	return validateTimeRange(c.NightModeRange) &&
		validateTimeRange(c.InvisibleRange) &&
		validateInvisibleWeek(c.InvisibleDay) &&
		validateClockStyle(c.ClockStyle) &&
		validateColorIndexMap(c.ColorIndexMap) &&
		brightnessInRange(c.Brightness) &&
		validateNumericalStocks(c.NumericalStocks) &&
		validateWorldSelection(c.WorldSelected)
}

func validateClockStyle(s string) bool {
	return len(s) > 0 && len(s) <= 32
}

func validateWorldSelection(s *WorldSelection) bool {
	if s == nil {
		return true
	}
	return s.Row >= 0 && s.Row < worldSelectionRowMax && s.Col >= 0 && s.Col < worldSelectionColMax
}

func cloneWorldSelection(s *WorldSelection) *WorldSelection {
	if s == nil {
		return nil
	}
	cp := *s
	return &cp
}

func equalWorldSelection(a, b *WorldSelection) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return a.Row == b.Row && a.Col == b.Col
}

func validateNumericalStocks(stocks []string) bool {
	if len(stocks) > 3 {
		return false
	}
	for _, s := range stocks {
		if len(s) == 0 || len(s) > 64 {
			return false
		}
	}
	return true
}

func validateColorIndexMap(m map[string]int) bool {
	for _, v := range m {
		if v < colorIndexMin || v > colorIndexMax {
			return false
		}
	}
	return true
}

func cloneIntMap(m map[string]int) map[string]int {
	if m == nil {
		return nil
	}
	cp := make(map[string]int, len(m))
	for k, v := range m {
		cp[k] = v
	}
	return cp
}

func equalIntMap(a, b map[string]int) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || bv != v {
			return false
		}
	}
	return true
}
