package main

import (
	"net/http"
	"os"
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	UserName string `gorm:"unique; not null" json:"UserName"`
	Password string `gorm:"size:20; not null" json:"Password"`
	Phone    string `gorm:"size:20; not null" json:"Phone"`
}

func main() {

	//数据库
	godotenv.Load("loginPage/config.env") // 路径是 config.env 相对于项目根目录的位置

	mysqlUser := getEnv("MYSQL_USER", "root")          // 默认值 root
	mysqlPwd := getEnv("MYSQL_PASSWORD", "")           // 必须在本地配置，无默认值
	mysqlHost := getEnv("MYSQL_HOST", "127.0.0.1")     // 默认值 127.0.0.1
	mysqlPort := getEnv("MYSQL_PORT", "3306")          // 默认值 3306
	mysqlDB := getEnv("MYSQL_DB", "login_information") // 默认值 login_information

	//拼接DSN
	dsn := mysqlUser + ":" + mysqlPwd + "@tcp(" + mysqlHost + ":" + mysqlPort + ")/" + mysqlDB + "?charset=utf8mb4&parseTime=True&loc=Local"
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("MySQL 连接失败！错误原因：" + err.Error())
	}
	
	db.AutoMigrate(&User{})

	//服务器
	ginServer := gin.Default()
	ginServer.Use(cors.New(cors.Config{

		AllowOrigins:     []string{"http://127.0.0.1:5500"}, //前端地址
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	ginServer.POST("/user/auth", func(context *gin.Context) {

		var req struct {
			Action   string `json:"Action" binding:"required,oneof=login register"` // 必须是login或register
			UserName string `json:"UserName" binding:"required"`
			Password string `json:"Password" binding:"required"`
			Phone    string `json:"Phone"`
		}

		if err := context.ShouldBindJSON(&req); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"msg": "参数错误：登录/注册类型必填，用户名和密码不能为空",
			})
			return
		}

		switch req.Action {

		case "register": //注册

			if req.Phone == "" {
				context.JSON(http.StatusBadRequest, gin.H{"msg": "注册失败，手机号不能为空"})
				return
			}

			resName := db.Where("user_name = ?", req.UserName).First(&User{})
			if resName.RowsAffected > 0 {
				context.JSON(http.StatusOK, gin.H{
					"msg": "注册失败，用户名已存在！",
				})
				return
			}

			resPhoneNo := db.Where("phone = ?", req.Phone).First(&User{})
			if resPhoneNo.RowsAffected > 0 {
				context.JSON(http.StatusOK, gin.H{
					"msg": "注册失败，手机号码已存在！",
				})
				return
			}

			newUser := User{
				UserName: req.UserName,
				Password: req.Password,

				Phone:    req.Phone,
			}
			if err := db.Create(&newUser).Error; err != nil {
				context.JSON(http.StatusInternalServerError, gin.H{
					"msg": "注册失败，保存数据出错",
				})
				return
			}

			context.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg": "注册成功！",
				"userInfo": gin.H{
					"UserName": req.UserName,
					"Password": req.Password,
					"Phone":    req.Phone,
				},
			})

		case "login": // 登录
			var user User
			result := db.Where("user_name = ? AND password = ?", req.UserName, req.Password).First(&user)
			if result.RowsAffected == 0 {
				context.JSON(http.StatusOK, gin.H{
					"msg": "登录失败，用户名/密码不正确!",
				})
				return
			}
			context.JSON(http.StatusOK, gin.H{
				"msg": "登录成功!",
			})
		}
	})

	//修改密码
	ginServer.POST("/user/change-password", func(context *gin.Context) {
		var req struct {
			UserName    string `json:"UserName" binding:"required"`
			Phone       string `json:"Phone" binding:"required"`
			NewPassword string `json:"NewPassword" binding:"required,min=6"`
		}

		if err := context.ShouldBindJSON(&req); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"msg": "请求参数错误, 新密码至少需要6位",
			})
			return
		}

		var u User
		res := db.Where("user_name = ?", req.UserName).First(&u)
		if res.RowsAffected == 0 {
			context.JSON(http.StatusOK, gin.H{"msg": "用户名不存在"})
			return
		}
		if u.Phone != req.Phone {
			context.JSON(http.StatusOK, gin.H{"msg": "手机号不正确"})
			return
		}
		u.Password = req.NewPassword
		if err := db.Save(&u).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"msg": "密码修改失败，系统错误",
			})
			return
		}
		context.JSON(http.StatusOK, gin.H{"msg": "密码修改成功"})
	})

	ginServer.Run(":8080")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}