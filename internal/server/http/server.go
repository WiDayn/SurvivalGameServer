package http

import (
	"SurvivalGame/internal/controller"
	"SurvivalGame/internal/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 处理跨域请求,支持options访问
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

func StartHTTP() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(Cors())
	middleware.LoadMiddleware(router)
	controller.InitController(router)

	srv := &http.Server{
		Addr:    ":7890",
		Handler: router,
	}

	err := srv.ListenAndServe()
	if err != nil {
		return
	}
}
