package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
)

var (
	DB *gorm.DB
)

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

func (Todo) TableName() string {
	return "todo"
}

func initMySQL() (err error) {
	dsn := "root:@tcp(localhost:3306)/xiangmu?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}
	//是否成功连接
	sqlDB, err := DB.DB()
	if err != nil {
		defer sqlDB.Close()
		log.Fatalf("failed to get *sql.DB: %v", err)
		return err
	}
	log.Println("连接成功")
	return nil
}
func main() {
	//创建数据库
	//连接数据库
	err := initMySQL()
	if err != nil {
		panic(err)
	}

	//模型绑定
	DB.AutoMigrate(&Todo{})

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
			//拿出请求中的数据
			var todo Todo
			c.BindJSON(&todo)
			//存入数据库
			if err = DB.Create(&todo).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{
					"error": err.Error(),
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"code": 200,
					"msg":  "success",
					"data": todo,
				})
			}
			//响应
		})
		//查看所有待办事项
		v1Group.GET("/todo", func(c *gin.Context) {
			var todoList []Todo
			if err = DB.Find(&todoList).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{
					"error": err.Error(),
				})
			} else {
				c.JSON(http.StatusOK, todoList)
			}
		})
		//查看某一个待办事项
		v1Group.GET("/todo/:id", func(c *gin.Context) {

		})
		//修改
		v1Group.PUT("/todo/:id", func(c *gin.Context) {
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusOK, gin.H{
					"error": "id不存在",
				})
				return
			}
			var todo Todo
			if err = DB.Where("id=?", id).First(&todo).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.BindJSON(&todo)
			if err = DB.Save(&todo).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{
					"error": err.Error(),
				})
			} else {
				c.JSON(http.StatusOK, todo)
			}
		})
		//删除
		v1Group.DELETE("/todo/:id", func(c *gin.Context) {
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusOK, gin.H{
					"error": "无效id",
				})
				return
			}
			if err = DB.Where("id=?", id).Delete(Todo{}).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{
					"error": err.Error(),
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					id: "delete",
				})
			}
		})
	}

	r.Run()
}
