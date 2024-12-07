package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

func main() {
	//创建数据库
	//连接数据库

	r := gin.Default()
	//使用静态文件
	r.Static("/static", "static")
	//匹配模板
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	v1Group := r.Group("v1")
	{
		//待办事项
		//添加
		v1Group.POST("/todo", func(c *gin.Context) {

		})
		//查看所有待办事项
		v1Group.GET("/todo", func(c *gin.Context) {

		})
		//查看某一个待办事项
		v1Group.GET("/todo/:id", func(c *gin.Context) {

		})
		//修改
		v1Group.PUT("/todo/:id", func(c *gin.Context) {

		})
		//删除
		v1Group.DELETE("/todo/:id", func(c *gin.Context) {

		})
	}

	r.Run()
}
