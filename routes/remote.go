package routes

import (
	"StandBy-Web/controller"

	"github.com/gin-gonic/gin"
)

func addRemoteRoutes(rg *gin.RouterGroup) {
	remoteGroup := rg.Group("/remote")
	remoteGroup.POST("/register", controller.RemoteRegister)
	remoteGroup.POST("/sync", controller.RemoteSync)
	remoteGroup.GET("/config", controller.RemoteGetConfig)
	remoteGroup.PATCH("/config", controller.RemotePatchConfig)
}
