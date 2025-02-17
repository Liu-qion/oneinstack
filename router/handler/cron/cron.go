package cron

import (
	"net/http"
	"oneinstack/core"
	"oneinstack/internal/services/cron"
	"oneinstack/router/input"

	"github.com/gin-gonic/gin"
)

func GetCronList(c *gin.Context) {
	var param input.CronParam
	if err := c.ShouldBindJSON(&param); err != nil {
		core.HandleError(c, http.StatusBadRequest, err, nil)
		return
	}
	rules, err := cron.GetCronList(c, &param)
	if err != nil {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	core.HandleSuccess(c, rules)
}

func AddCron(c *gin.Context) {
	var param input.AddCronParam
	if err := c.ShouldBindJSON(&param); err != nil {
		core.HandleError(c, http.StatusBadRequest, err, nil)
		return
	}
	err := cron.AddCron(c, &param)
	if err != nil {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	core.HandleSuccess(c, nil)
}

func UpdateCron(c *gin.Context) {
	var param input.AddCronParam
	if err := c.ShouldBindJSON(&param); err != nil {
		core.HandleError(c, http.StatusBadRequest, err, nil)
		return
	}
	err := cron.UpdateCron(c, &param)
	if err != nil {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	core.HandleSuccess(c, nil)
}

func DeleteCron(c *gin.Context) {
	var param input.AddCronParam
	if err := c.ShouldBindJSON(&param); err != nil {
		core.HandleError(c, http.StatusBadRequest, err, nil)
		return
	}
	err := cron.DeleteCron(c, param.ID)
	if err != nil {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	core.HandleSuccess(c, nil)
}

func DisableCron(c *gin.Context) {
	var param input.AddCronParam
	if err := c.ShouldBindJSON(&param); err != nil {
		core.HandleError(c, http.StatusBadRequest, err, nil)
		return
	}
	err := cron.DisableCron(c, param.ID)
	if err != nil {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	core.HandleSuccess(c, nil)
}

func EnableCron(c *gin.Context) {
	var param input.AddCronParam
	if err := c.ShouldBindJSON(&param); err != nil {
		core.HandleError(c, http.StatusBadRequest, err, nil)
		return
	}
	err := cron.EnableCron(c, param.ID)
	if err != nil {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	core.HandleSuccess(c, nil)
}
