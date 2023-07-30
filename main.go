package main

import (
	"log"
	"restful-api/auth"
	"restful-api/handler"
	"restful-api/user"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main()  {
	dsn := "laravel:password@tcp(127.0.0.1:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}
	userRepository := user.NewRepository(db)
	
	userService := user.NewService(userRepository)
	authService := auth.NewService()

	
	userHandler := handler.NewUserHandler(userService,authService)
	
	router := gin.Default()
	api := router.Group("/api/v1")
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.LoginUser)
	api.POST("/email-checkers",userHandler.CheckEmailAvailability)
	api.POST("/avatars",userHandler.UploadAvatar)
	router.Run()
}

