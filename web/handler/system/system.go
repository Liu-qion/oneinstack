package system

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"oneinstack/core"
	"oneinstack/internal/services/system"
)

func GetSystemInfo(c *gin.Context) {
	info, err := system.GetSystemInfo()
	if err != nil {
		core.HandleError(c, http.StatusInternalServerError, core.ErrInternalServerError, err)
		return
	}
	core.HandleSuccess(c, info)
}

func GetSystemMonitor(c *gin.Context) {
	monitor := system.GetSystemMonitor()
	core.HandleSuccess(c, monitor)
}
