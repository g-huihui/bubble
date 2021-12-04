package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
)

var (
	DB *gorm.DB
)

// Todo Model
type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

func initMySQL() (err error) {
	dbInfo := "root:123456@tcp(www.gongyanhui.top:3306)/test_db?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open("mysql", dbInfo)
	if err != nil {
		return
	}
	return DB.DB().Ping()
}

func main() {

	// 创建数据库
	// sql：Create Database bubble;
	// 连接数据库
	err := initMySQL()
	if err != nil {
		panic(err)
	}
	defer DB.Close() //程序退出关闭数据库连接

	//模型绑定
	DB.AutoMigrate(&Todo{})

	r := gin.Default()
	//告诉gin框架模版文件引用的静态文件去哪里找
	r.Static("/static", "static")
	//告诉gin框架去哪里找模版文件
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// v1
	v1Group := r.Group("v1")
	{
		//添加
		v1Group.POST("/todo", func(c *gin.Context) {
			// 前端页面填写代办事项 点击提交 会发请求到这里
			// 1 从请求中把数据拿出来
			var todo Todo
			c.BindJSON(&todo)
			// 2 存入数据库
			// 3 返回响应
			if err = DB.Create(&todo).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, todo)
			}
		})

		//查看 所有代办事项
		v1Group.GET("/todo", func(c *gin.Context) {
			// 查询todo这个表里的所有数据
			var todoList []Todo
			if err = DB.Find(&todoList).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, todoList)
			}
		})

		//查看 某一个代办事项
		v1Group.GET("/todo/:id", func(c *gin.Context) {

		})

		//修改 某一个代办事项
		v1Group.PUT("/todo/:id", func(c *gin.Context) {
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusOK, gin.H{"error": "无效的ID"})
				return
			}
			var todo Todo
			if err = DB.Where("id = ?", id).Find(&todo).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
				return
			}
			c.BindJSON(&todo)
			if err := DB.Save(&todo).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, todo)
			}
		})

		//删除 某一个代办事项
		v1Group.DELETE("/todo/:id", func(c *gin.Context) {
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusOK, gin.H{"error": "无效的ID"})
				return
			}
			if err := DB.Where("id = ?", id).Delete(Todo{}).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{id: "Deleted"})
			}
		})
	}
	r.Run()
}
