package system

import (
	"net/http"
	"oneinstack/core"
	"oneinstack/internal/services/system"

	"github.com/gin-gonic/gin"
)

func GetSystemInfo(c *gin.Context) {
	info, err := system.GetSystemInfo()
	if err != nil {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	core.HandleSuccess(c, info)
}

func GetSystemMonitor(c *gin.Context) {
	monitor, err := system.GetSystemMonitor()
	if err != nil {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	core.HandleSuccess(c, monitor)
}

func GetLibCount(c *gin.Context) {
	count, err := system.GetLibCount()
	if err != nil {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	core.HandleSuccess(c, count)
}

func GetWebSiteCount(c *gin.Context) {
	count, err := system.GetWebSiteCount()
	if err != nil {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	core.HandleSuccess(c, count)
}

func SystemInfo(c *gin.Context) {
	info, err := system.SystemInfo()
	if err != nil {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	core.HandleSuccess(c, info)
}
