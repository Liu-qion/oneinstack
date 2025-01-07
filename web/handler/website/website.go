package website

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"oneinstack/core"
	"oneinstack/internal/models"
	"oneinstack/internal/services/website"
)

func List(c *gin.Context) {
	list, err := website.List()
	if err != nil {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	core.HandleSuccess(c, list)
}

func Add(c *gin.Context) {
	input := &models.Website{}
	if err := c.ShouldBindJSON(&input); err != nil {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	err := website.Add(input)
	if err != nil {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	core.HandleSuccess(c, "创建成功")

}

func Update(c *gin.Context) {
	input := &models.Website{}
	if err := c.ShouldBindJSON(&input); err != nil {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	err := website.Update(input)
	if err != nil {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	core.HandleSuccess(c, "更新成功")
}

func Delete(c *gin.Context) {
	input := &models.Website{}
	if err := c.ShouldBindJSON(&input); err != nil {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	err := website.Delete(input.ID)
	if err != nil {
		core.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	core.HandleSuccess(c, "删除成功")
}
