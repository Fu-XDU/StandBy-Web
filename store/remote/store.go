package remote

import (
	"sync"
	"time"
)

const pageFloat = "float"

type pageRecord struct {
	UpdatedAt int64            `json:"updatedAt"`
	Config    *FloatPageConfig `json:"config,omitempty"`
}

type deviceRecord struct {
	Pages map[string]*pageRecord `json:"pages"`
}

// Store 全部配置保存在进程内存中，重启后清空。
type Store struct {
	mu      sync.RWMutex
	devices map[string]*deviceRecord
}

func NewStore() *Store {
	return &Store{devices: make(map[string]*deviceRecord)}
}

func (s *Store) getOrCreateDeviceLocked(deviceID string) *deviceRecord {
	rec, ok := s.devices[deviceID]
	if !ok {
		rec = &deviceRecord{Pages: make(map[string]*pageRecord)}
		s.devices[deviceID] = rec
	}
	return rec
}

// Register 确保设备记录存在（幂等）。
func (s *Store) Register(deviceID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.getOrCreateDeviceLocked(deviceID)
}

type SyncResult struct {
	Action    string           `json:"action"`
	UpdatedAt int64            `json:"updatedAt"`
	Config    *FloatPageConfig `json:"config,omitempty"`
}

// SyncPage 按时间戳合并单页配置。
func (s *Store) SyncPage(deviceID, pageID string, clientAt int64, client *FloatPageConfig) (*SyncResult, error) {
	if pageID != pageFloat {
		return nil, ErrUnsupportedPage
	}
	if !validateFloatConfig(client) {
		return nil, ErrInvalidConfig
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	rec := s.getOrCreateDeviceLocked(deviceID)
	page := rec.Pages[pageID]
	if page == nil || page.Config == nil {
		rec.Pages[pageID] = &pageRecord{UpdatedAt: clientAt, Config: cloneFloatConfig(client)}
		return &SyncResult{Action: "store_client", UpdatedAt: clientAt}, nil
	}

	if page.UpdatedAt > clientAt {
		return &SyncResult{
			Action:    "apply_server",
			UpdatedAt: page.UpdatedAt,
			Config:    cloneFloatConfig(page.Config),
		}, nil
	}

	rec.Pages[pageID] = &pageRecord{UpdatedAt: clientAt, Config: cloneFloatConfig(client)}
	return &SyncResult{Action: "store_client", UpdatedAt: clientAt}, nil
}

type PageState struct {
	DeviceID  string           `json:"deviceId"`
	PageID    string           `json:"pageId"`
	UpdatedAt int64            `json:"updatedAt"`
	Config    *FloatPageConfig `json:"config,omitempty"`
}

// GetPage 读取服务端当前配置。
func (s *Store) GetPage(deviceID, pageID string) (*PageState, error) {
	if pageID != pageFloat {
		return nil, ErrUnsupportedPage
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	rec, ok := s.devices[deviceID]
	if !ok {
		return nil, ErrDeviceNotFound
	}
	page := rec.Pages[pageID]
	if page == nil || page.Config == nil {
		return nil, ErrPageNotFound
	}
	return &PageState{
		DeviceID:  deviceID,
		PageID:    pageID,
		UpdatedAt: page.UpdatedAt,
		Config:    cloneFloatConfig(page.Config),
	}, nil
}

// PatchPageRemote 凭已注册 deviceId 局部改配：仅更新 patch 中出现的字段。
func (s *Store) PatchPageRemote(deviceID, pageID string, patch *FloatPageConfigPatch) (*PageState, error) {
	if pageID != pageFloat {
		return nil, ErrUnsupportedPage
	}
	if patch == nil || patch.IsEmpty() {
		return nil, ErrEmptyPatch
	}
	if err := validatePatch(patch); err != nil {
		return nil, err
	}

	now := time.Now().UnixMilli()

	s.mu.Lock()
	defer s.mu.Unlock()

	rec, ok := s.devices[deviceID]
	if !ok {
		return nil, ErrDeviceNotFound
	}

	var base *FloatPageConfig
	page := rec.Pages[pageID]
	if page == nil || page.Config == nil {
		base = defaultFloatConfig()
	} else {
		base = page.Config
	}

	merged := applyPatch(base, patch)
	if !validateFloatConfig(merged) {
		return nil, ErrInvalidConfig
	}

	rec.Pages[pageID] = &pageRecord{UpdatedAt: now, Config: merged}
	return &PageState{
		DeviceID:  deviceID,
		PageID:    pageID,
		UpdatedAt: now,
		Config:    cloneFloatConfig(merged),
	}, nil
}
