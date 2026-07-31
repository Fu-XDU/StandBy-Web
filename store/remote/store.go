package remote

import (
	"sync"
	"time"
)

const pageFloat = "float"

// OnlineWindow 判定设备「在线」的时间窗口：2 小时内有过远程活动即计入。
const OnlineWindow = 2 * time.Hour

type pageRecord struct {
	UpdatedAt int64            `json:"updatedAt"`
	Config    *FloatPageConfig `json:"config,omitempty"`
}

type deviceRecord struct {
	Pages      map[string]*pageRecord `json:"pages"`
	LastSeenAt int64                  `json:"lastSeenAt,omitempty"` // unix ms
}

// Store 配置保存在内存中，并同步落盘到本地文件（默认 Linux: /etc/standby-web/remote-store.json）。
type Store struct {
	mu      sync.RWMutex
	devices map[string]*deviceRecord
	path    string
}

// NewStore 使用默认缓存路径创建并加载已有文件。
func NewStore() *Store {
	return NewStoreWithPath(DefaultPersistPath())
}

// NewStoreWithPath 使用指定路径；path 为空则仅内存、不落盘。
func NewStoreWithPath(path string) *Store {
	s := &Store{
		devices: make(map[string]*deviceRecord),
		path:    path,
	}
	_ = s.load()
	return s
}

func (s *Store) getOrCreateDeviceLocked(deviceID string) *deviceRecord {
	rec, ok := s.devices[deviceID]
	if !ok {
		rec = &deviceRecord{Pages: make(map[string]*pageRecord)}
		s.devices[deviceID] = rec
	}
	return rec
}

func (s *Store) touchLocked(deviceID string) {
	rec := s.getOrCreateDeviceLocked(deviceID)
	rec.LastSeenAt = time.Now().UnixMilli()
}

// CountOnline 返回 within 时间窗口内有过活动的设备数量。
func (s *Store) CountOnline(within time.Duration) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cutoff := time.Now().Add(-within).UnixMilli()
	n := 0
	for _, rec := range s.devices {
		if rec != nil && rec.LastSeenAt >= cutoff {
			n++
		}
	}
	return n
}

// Register 确保设备记录存在（幂等）。
func (s *Store) Register(deviceID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.touchLocked(deviceID)
	_ = s.persistLocked()
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
	s.touchLocked(deviceID)
	page := rec.Pages[pageID]
	if page == nil || page.Config == nil {
		rec.Pages[pageID] = &pageRecord{UpdatedAt: clientAt, Config: cloneFloatConfig(client)}
		if err := s.persistLocked(); err != nil {
			return nil, err
		}
		return &SyncResult{Action: "store_client", UpdatedAt: clientAt}, nil
	}

	if page.UpdatedAt > clientAt {
		if err := s.persistLocked(); err != nil {
			return nil, err
		}
		return &SyncResult{
			Action:    "apply_server",
			UpdatedAt: page.UpdatedAt,
			Config:    cloneFloatConfig(page.Config),
		}, nil
	}

	rec.Pages[pageID] = &pageRecord{UpdatedAt: clientAt, Config: cloneFloatConfig(client)}
	if err := s.persistLocked(); err != nil {
		return nil, err
	}
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

	s.mu.Lock()
	defer s.mu.Unlock()

	rec, ok := s.devices[deviceID]
	if !ok {
		return nil, ErrDeviceNotFound
	}
	page := rec.Pages[pageID]
	if page == nil || page.Config == nil {
		return nil, ErrPageNotFound
	}
	s.touchLocked(deviceID)
	_ = s.persistLocked()
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
	s.touchLocked(deviceID)
	if err := s.persistLocked(); err != nil {
		return nil, err
	}
	return &PageState{
		DeviceID:  deviceID,
		PageID:    pageID,
		UpdatedAt: now,
		Config:    cloneFloatConfig(merged),
	}, nil
}
