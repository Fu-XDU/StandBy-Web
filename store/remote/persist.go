package remote

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	configDirName  = "standby-web"
	storeFileName  = "remote-store.json"
	envConfigDir   = "STANDBY_CONFIG_DIR"
	linuxSystemDir = "/etc/standby-web"
)

type persistedStore struct {
	Devices map[string]*deviceRecord `json:"devices"`
}

// DefaultPersistPath 返回远程配置缓存文件路径。
// 优先 STANDBY_CONFIG_DIR；Linux 默认 /etc/standby-web；否则 ~/.config/standby-web。
func DefaultPersistPath() string {
	if dir := os.Getenv(envConfigDir); dir != "" {
		return filepath.Join(dir, storeFileName)
	}
	if runtime.GOOS == "linux" {
		return filepath.Join(linuxSystemDir, storeFileName)
	}
	return filepath.Join(userConfigDir(), storeFileName)
}

func userConfigDir() string {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return filepath.Join(dir, configDirName)
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return filepath.Join(".", "data", configDirName)
	}
	return filepath.Join(home, ".config", configDirName)
}

func (s *Store) load() error {
	if s.path == "" {
		return nil
	}
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	if len(data) == 0 {
		return nil
	}

	var snap persistedStore
	if err := json.Unmarshal(data, &snap); err != nil {
		return fmt.Errorf("parse remote store cache: %w", err)
	}
	if snap.Devices == nil {
		s.devices = make(map[string]*deviceRecord)
		return nil
	}

	devices := make(map[string]*deviceRecord, len(snap.Devices))
	for deviceID, rec := range snap.Devices {
		if rec == nil {
			devices[deviceID] = &deviceRecord{Pages: make(map[string]*pageRecord)}
			continue
		}
		pages := make(map[string]*pageRecord, len(rec.Pages))
		for pageID, page := range rec.Pages {
			if page == nil {
				continue
			}
			pages[pageID] = &pageRecord{
				UpdatedAt: page.UpdatedAt,
				Config:    cloneFloatConfig(page.Config),
			}
		}
		devices[deviceID] = &deviceRecord{
			Pages:      pages,
			LastSeenAt: rec.LastSeenAt,
		}
	}
	s.devices = devices
	return nil
}

// persistLocked 将当前内存快照原子写入文件；调用方须已持有写锁。
func (s *Store) persistLocked() error {
	if s.path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	snap := persistedStore{Devices: make(map[string]*deviceRecord, len(s.devices))}
	for deviceID, rec := range s.devices {
		pages := make(map[string]*pageRecord, len(rec.Pages))
		for pageID, page := range rec.Pages {
			if page == nil {
				continue
			}
			pages[pageID] = &pageRecord{
				UpdatedAt: page.UpdatedAt,
				Config:    cloneFloatConfig(page.Config),
			}
		}
		snap.Devices[deviceID] = &deviceRecord{
			Pages:      pages,
			LastSeenAt: rec.LastSeenAt,
		}
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	dir := filepath.Dir(s.path)
	tmp, err := os.CreateTemp(dir, storeFileName+".tmp.*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() { _ = os.Remove(tmpName) }()

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Chmod(0o644); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpName, s.path); err != nil {
		return err
	}
	return nil
}
