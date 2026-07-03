package remote

import (
	"testing"
	"time"
)

func sampleFloatConfig() *FloatPageConfig {
	return &FloatPageConfig{
		AutoNightMode:      true,
		NightModeRange:     []string{"2025-01-01T00:00:00.000Z", "2025-01-01T06:00:00.000Z"},
		AutoInvisible:      false,
		InvisibleRange:     []string{"2025-01-01T00:00:00.000Z", "2025-01-01T06:00:00.000Z"},
		InvisibleDayEnable: false,
		InvisibleDay:       []bool{false, false, false, false, false, false, false},
		ClockStyle:         "float",
		ColorIndexMap:      map[string]int{"float": 2, "numerical": 0},
		Brightness:         1,
	}
}

func TestSyncPage_timestampMerge(t *testing.T) {
	s := NewStore()
	deviceID := "abcd1234abcd1234abcd1234abcd12"
	client := sampleFloatConfig()

	r1, err := s.SyncPage(deviceID, pageFloat, 100, client)
	if err != nil {
		t.Fatal(err)
	}
	if r1.Action != "store_client" || r1.UpdatedAt != 100 {
		t.Fatalf("unexpected first sync: %+v", r1)
	}

	r2, err := s.SyncPage(deviceID, pageFloat, 50, client)
	if err != nil {
		t.Fatal(err)
	}
	if r2.Action != "apply_server" || r2.UpdatedAt != 100 || r2.Config == nil {
		t.Fatalf("expected server win: %+v", r2)
	}

	updated := *client
	updated.ColorIndexMap = map[string]int{"float": 5, "numerical": 0}
	r3, err := s.SyncPage(deviceID, pageFloat, 200, &updated)
	if err != nil {
		t.Fatal(err)
	}
	if r3.Action != "store_client" || r3.UpdatedAt != 200 {
		t.Fatalf("expected client win: %+v", r3)
	}
}

func TestPatchPageRemote_requiresRegister(t *testing.T) {
	s := NewStore()
	deviceID := "abcd1234abcd1234abcd1234abcd12"
	brightness := 0.5
	_, err := s.PatchPageRemote(deviceID, pageFloat, &FloatPageConfigPatch{Brightness: &brightness})
	if err != ErrDeviceNotFound {
		t.Fatalf("expected not registered: %v", err)
	}
}

func TestPatchPageRemote_partialFields(t *testing.T) {
	s := NewStore()
	deviceID := "abcd1234abcd1234abcd1234abcd12"
	client := sampleFloatConfig()

	s.Register(deviceID)
	_, err := s.SyncPage(deviceID, pageFloat, 500, client)
	if err != nil {
		t.Fatal(err)
	}

	brightness := 0.4
	colorMap := map[string]int{"float": 7, "numerical": 3}
	before := time.Now().UnixMilli()
	state, err := s.PatchPageRemote(deviceID, pageFloat, &FloatPageConfigPatch{
		Brightness:    &brightness,
		ColorIndexMap: &colorMap,
	})
	if err != nil {
		t.Fatal(err)
	}
	if state.UpdatedAt < before {
		t.Fatal("patch timestamp should be current")
	}
	if !floatEqual(state.Config.Brightness, 0.4) || state.Config.ColorIndexMap["float"] != 7 {
		t.Fatalf("unexpected patch result: %+v", state.Config)
	}
	if state.Config.AutoNightMode != client.AutoNightMode {
		t.Fatal("patch should not change autoNightMode")
	}

	r, err := s.SyncPage(deviceID, pageFloat, 500, client)
	if err != nil {
		t.Fatal(err)
	}
	if r.Action != "apply_server" || r.Config.ColorIndexMap["float"] != 7 {
		t.Fatalf("expected patched config on sync: %+v", r)
	}
}

func TestPatchPageRemote_emptyPatch(t *testing.T) {
	s := NewStore()
	deviceID := "abcd1234abcd1234abcd1234abcd12"
	s.Register(deviceID)
	_, err := s.PatchPageRemote(deviceID, pageFloat, &FloatPageConfigPatch{})
	if err != ErrEmptyPatch {
		t.Fatalf("expected empty patch error: %v", err)
	}
}

func TestValidateFloatConfig_ranges(t *testing.T) {
	cfg := sampleFloatConfig()
	if !validateFloatConfig(cfg) {
		t.Fatal("sample config should be valid")
	}

	cfg.Brightness = 1.1
	if validateFloatConfig(cfg) {
		t.Fatal("brightness > 1 should be invalid")
	}
	cfg.Brightness = -0.01
	if validateFloatConfig(cfg) {
		t.Fatal("brightness < 0 should be invalid")
	}
	cfg.Brightness = 1

	cfg.ColorIndexMap = map[string]int{"float": 100}
	if validateFloatConfig(cfg) {
		t.Fatal("color index 100 should be invalid")
	}
	cfg.ColorIndexMap = map[string]int{"float": 0}

	cfg.NightModeRange = []string{"not-a-time", "2025-01-01T06:00:00.000Z"}
	if validateFloatConfig(cfg) {
		t.Fatal("invalid nightModeRange should be rejected")
	}
	cfg.NightModeRange = []string{"2025-01-01T22:00:00.000Z"}
	if validateFloatConfig(cfg) {
		t.Fatal("single-element range should be invalid")
	}

	cfg.InvisibleDay = []bool{true, false, false}
	if validateFloatConfig(cfg) {
		t.Fatal("invisibleDay must have length 7")
	}
}

func TestValidatePatch_partialRanges(t *testing.T) {
	badBrightness := 2.0
	if validatePatch(&FloatPageConfigPatch{Brightness: &badBrightness}) == nil {
		t.Fatal("patch brightness out of range should fail")
	}
	badRange := []string{"invalid", "2025-01-01T08:00:00.000Z"}
	if validatePatch(&FloatPageConfigPatch{NightModeRange: &badRange}) == nil {
		t.Fatal("patch invalid time range should fail")
	}
	goodRange := []string{"2025-01-01T22:00:00.000Z", "2025-01-01T07:00:00.000Z"}
	if validatePatch(&FloatPageConfigPatch{NightModeRange: &goodRange}) != nil {
		t.Fatal("patch valid time range should pass")
	}
}

func TestEqualFloatConfig_fieldByField(t *testing.T) {
	a := sampleFloatConfig()
	b := sampleFloatConfig()
	if !EqualFloatConfig(a, b) {
		t.Fatal("expected equal configs")
	}
	b.ColorIndexMap = map[string]int{"float": 3}
	if EqualFloatConfig(a, b) {
		t.Fatal("expected different configs")
	}
}
