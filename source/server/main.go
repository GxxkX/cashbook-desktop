package main

import (
	"cashbook-server/controller"
	"cashbook-server/dao"
	"cashbook-server/util"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

func main() {

	dao.InitDb()

	router := gin.Default()
	store := memstore.NewStore([]byte("secret_for_cashbook"))
	router.Use(Cors())
	router.Use(sessions.Sessions("bookAuthenticated", store))
	port := ":131"

	root := router.Group("/")
	root.Use(Cors())
	root.GET("/captcha/:img", controller.CaptchaHandle)

	api := router.Group("/api")
	api.Use(Cors())
	api.GET("/book/:key", controller.GetBook)
	api.POST("/book", controller.CreateBook)
	api.GET("/book/list", controller.GetBookList)
	api.GET("/server", controller.GetServerInfo)
	api.GET("/captcha", controller.Captcha)

	adminApi := api.Group("/admin")
	adminApi.Use(openBook())
	{
		adminApi.POST("/book/changeKey", controller.ChangeKey)
		// 字典相关
		adminApi.GET("/dist/:type", controller.GetDistList)
		adminApi.GET("/dist", controller.GetDistPage)
		adminApi.POST("/dist", controller.AddDist)
		adminApi.PUT("/dist/:id", controller.UpdateDist)
		adminApi.DELETE("/dist/:id", controller.DeleteDist)
		// 分析图表相关
		adminApi.POST("/analysis/dailyLine", controller.GetDailyLine)
		adminApi.POST("/analysis/typePie", controller.GetTypePie)
		adminApi.POST("/analysis/payTypeBar", controller.GetPayTypeBar)
		adminApi.POST("/analysis/monthBar", controller.MonthBar)
		// 流水相关
		adminApi.GET("/flow/getAll", controller.GetAll)
		adminApi.GET("/flow/getAllByMon", controller.GetAllByMon)
		adminApi.POST("/flow/importFlows", controller.ImportFlows)
		adminApi.GET("/flow", controller.GetFlowsPage)
		adminApi.POST("/flow", controller.AddFlow)
		adminApi.PUT("/flow/:id", controller.UpdateFlow)
		adminApi.DELETE("/flow/:id", controller.DeleteFlow)
		// 计划相关
		adminApi.GET("/plans/:month", controller.GetPlan)
		adminApi.POST("/plans/:overwrite", controller.SetPlan)
		// 在线同步相关
		adminApi.POST("/online/upload", controller.Upload)
		adminApi.POST("/online/download", controller.Download)
	}
	fmt.Println("-------- 服务启动成功：" + port + " --------")
	err := router.Run(port)
	util.CheckErr(err)
}

func openBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		if sessions.Default(c).Get("bookKey") == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success":      false,
				"errorMessage": "请输入账本密钥！",
			})
			c.Abort()
			return
		}
		bookKey := sessions.Default(c).Get("bookKey").(string)
		if len(bookKey) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success":      false,
				"errorMessage": "请输入账本密钥！",
			})
			c.Abort()
			return
		}

		if dao.GetBook(bookKey).Id == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success":      false,
				"errorMessage": "账本不存在！",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
