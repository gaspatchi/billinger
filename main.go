package main

import (
	"fmt"

	"./lib/prometheus"
	"./middlewares/verifyauth"
	"./routes/create"
	"./routes/get"
	"./routes/verify"
	"./utils/startup"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var address string

func init() {
	startup.InitService()
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Print(err)
	}
	address = fmt.Sprintf("%s:%d", viper.GetString("server.address"), viper.GetInt("server.port"))
}

func main() {
	router := gin.Default()
	router.POST("/verify", verify.VerifyRefill)
	router.GET("/get", auth.VerifyAuth, get.GetStatement)
	router.POST("/create", auth.VerifyAuth, create.CreateToken)
	router.GET("/metrics", metrics.GetMetrics)

	ginpprof.Wrap(router)
	router.Run(address)
}
