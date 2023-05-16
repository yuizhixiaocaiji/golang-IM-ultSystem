package main

import (
	"ginchat/router"
	"ginchat/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMySql()
	utils.InitRedis()

	r := router.Router()
	r.Run()
}
