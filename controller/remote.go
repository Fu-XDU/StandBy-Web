package controller

import (
	"net/http"
	"regexp"
	"sync"

	"StandBy-Web/store/remote"

	"github.com/gin-gonic/gin"
)

var (
	remoteStore     *remote.Store
	remoteStoreOnce sync.Once
	deviceIDPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{16,64}$`)
	pageIDPattern   = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,31}$`)
)

func getRemoteStore() *remote.Store {
	remoteStoreOnce.Do(func() {
		remoteStore = remote.NewStore()
	})
	return remoteStore
}

type registerRequest struct {
	DeviceID string `json:"deviceId" binding:"required"`
}

func RemoteRegister(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if !deviceIDPattern.MatchString(req.DeviceID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": remote.ErrInvalidDeviceID.Error()})
		return
	}
	getRemoteStore().Register(req.DeviceID)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

type syncRequest struct {
	DeviceID  string                 `json:"deviceId" binding:"required"`
	PageID    string                 `json:"pageId" binding:"required"`
	UpdatedAt int64                  `json:"updatedAt"`
	Config    *remote.FloatPageConfig `json:"config" binding:"required"`
}

func RemoteSync(c *gin.Context) {
	var req syncRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if !deviceIDPattern.MatchString(req.DeviceID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": remote.ErrInvalidDeviceID.Error()})
		return
	}
	if !pageIDPattern.MatchString(req.PageID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page id"})
		return
	}

	result, err := getRemoteStore().SyncPage(req.DeviceID, req.PageID, req.UpdatedAt, req.Config)
	if err != nil {
		switch err {
		case remote.ErrUnsupportedPage, remote.ErrInvalidConfig:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "sync failed"})
		}
		return
	}
	c.JSON(http.StatusOK, result)
}

func RemoteGetConfig(c *gin.Context) {
	deviceID := c.Query("deviceId")
	pageID := c.Query("pageId")
	if !deviceIDPattern.MatchString(deviceID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": remote.ErrInvalidDeviceID.Error()})
		return
	}
	if !pageIDPattern.MatchString(pageID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page id"})
		return
	}

	state, err := getRemoteStore().GetPage(deviceID, pageID)
	if err != nil {
		switch err {
		case remote.ErrDeviceNotFound, remote.ErrPageNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case remote.ErrUnsupportedPage:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "get config failed"})
		}
		return
	}
	c.JSON(http.StatusOK, state)
}

type patchConfigRequest struct {
	DeviceID string `json:"deviceId" binding:"required"`
	PageID   string `json:"pageId"`
	remote.FloatPageConfigPatch
}

// RemotePatchConfig 公开接口：仅更新请求体中出现的配置项（deviceId 须已注册）。
func RemotePatchConfig(c *gin.Context) {
	var req patchConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if !deviceIDPattern.MatchString(req.DeviceID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": remote.ErrInvalidDeviceID.Error()})
		return
	}
	if req.PageID == "" {
		req.PageID = "float"
	}
	if !pageIDPattern.MatchString(req.PageID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page id"})
		return
	}

	state, err := getRemoteStore().PatchPageRemote(req.DeviceID, req.PageID, &req.FloatPageConfigPatch)
	if err != nil {
		switch err {
		case remote.ErrDeviceNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case remote.ErrEmptyPatch:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case remote.ErrUnsupportedPage, remote.ErrInvalidConfig:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "patch config failed"})
		}
		return
	}
	c.JSON(http.StatusOK, state)
}
