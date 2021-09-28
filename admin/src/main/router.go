package main

import (
	"github.com/gin-gonic/gin"
	"go-fcm-example/admin/src/define"
	"go-fcm-example/admin/src/service"
	"net"
	"net/http"
	"time"
)

/**
* @Author: caishi13202
* @Date: 2021/9/27 3:00 下午
 */
// 初始化路由
func initRouter(router *gin.Engine) {
	loginUser := make(map[string]string, 16)
	httpClient := createHTTPClient()
	service := &service.Notification{
	}
	service.SetLoginUser(loginUser)
	service.SetHttpClient(httpClient)
	accessToken(service, router)
	listAccount(service, router)
	sendMsg(service, router)
}

// 注册token
func accessToken(service *service.Notification, router *gin.Engine) {
	router.POST("/accessToken", func(c *gin.Context) {
		req := &define.AccessTokenReq{}
		c.ShouldBind(req)
		c.JSON(http.StatusOK, service.Login(req))
	})
}

// 查询当前登陆的用户列表
func listAccount(service *service.Notification, router *gin.Engine) {
	router.GET("/list", func(c *gin.Context) {
		c.JSON(http.StatusOK, service.ListAccount())
	})
}

// 发送消息
func sendMsg(service *service.Notification, router *gin.Engine) {
	router.POST("/send", func(c *gin.Context) {
		req := &define.SendReq{}
		c.ShouldBind(req)
		c.JSON(http.StatusOK, service.Send(req))
	})
}

// createHTTPClient for connection re-use
func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        30,
			MaxIdleConnsPerHost: 30,
			IdleConnTimeout:     time.Duration(5) * time.Second,
		},
		Timeout: time.Millisecond * time.Duration(5000),
	}
	return client
}
