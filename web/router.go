package web

import (
	"github.com/gin-gonic/gin"
	"oneinstack/web/handler/ftp"
	"oneinstack/web/handler/software"
	"oneinstack/web/handler/storage"
	"oneinstack/web/handler/system"
	"oneinstack/web/handler/user"
	"oneinstack/web/handler/website"
	"oneinstack/web/middleware"
)

func SetupRouter() *gin.Engine {

	r := gin.Default()

	g := r.Group("/v1")

	// 公共路由
	{
		g.POST("/login", user.LoginHandler)
	}

	sys := g.Group("/sys")
	sys.Use(middleware.AuthMiddleware())
	{
		sys.GET("/info", system.GetSystemInfo)
		sys.GET("/monitor", system.GetSystemMonitor)

		sys.POST("/remark/list", system.RemarkList)
		sys.POST("/remark/add", system.AddRemark)
		sys.POST("/remark/update", system.UpdateRemark)
		sys.POST("/remark/del", system.DeleteRemark)

		sys.POST("/dic/list", system.DictionaryList)
		sys.POST("/dic/add", system.AddDictionary)
		sys.POST("/dic/update", system.UpdateDictionary)
		sys.POST("/dic/del", system.DeleteDictionary)
	}

	storageg := g.Group("/storage")
	sys.Use(middleware.AuthMiddleware())
	{
		storageg.POST("/addconn", storage.ADDStorage)
		storageg.POST("/updateconn", storage.UpdateStorage)
		storageg.GET("/connlist", storage.GetStorage)
		storageg.GET("/delconn", storage.DelStorage)
		storageg.POST("/sync", storage.SyncStorage)
		storageg.POST("/liblist", storage.GetLib)
		storageg.POST("/rklist", storage.GetRedisKeys)
	}

	ftpg := g.Group("/ftp")
	sys.Use(middleware.AuthMiddleware())
	{
		ftpg.POST("/list", ftp.ListDirectory)
		ftpg.POST("/create", ftp.CreateFileOrDir)
		ftpg.POST("/upload", ftp.UploadFile)
		ftpg.POST("/download", ftp.DownloadFile)
		ftpg.POST("/delete", ftp.DeleteFileOrDir)
	}

	softg := g.Group("/soft")
	sys.Use(middleware.AuthMiddleware())
	{
		softg.POST("/list", software.GetSoftware)
		softg.GET("/getlog", software.GetLogContent)
		softg.POST("/install", software.RunInstallation)

	}

	websiteg := g.Group("/website")
	sys.Use(middleware.AuthMiddleware())
	{
		websiteg.POST("/list", website.List)
		websiteg.POST("/add", website.Add)
		websiteg.POST("/del", website.Delete)
		websiteg.POST("/update", website.Update)
	}

	return r
}
