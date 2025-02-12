package safe

import (
	"net/http"
	"oneinstack/core"
	"oneinstack/internal/services/safe"

	"github.com/gin-gonic/gin"
)

func GetFirewallInfo(c *gin.Context) {
	info, err := safe.GetFirewallStatus()
	if err != nil {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	core.HandleSuccess(c, gin.H{"info": info})
}

func GetFirewallStatus(c *gin.Context) {
	status, err := safe.GetFirewallStatus()
	if err != nil {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	core.HandleSuccess(c, gin.H{"status": status})
}

func GetFirewallRules(c *gin.Context) {
	rules, err := safe.GetFirewallPorts()
	if err != nil {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	core.HandleSuccess(c, gin.H{"rules": rules})
}
