package main

import (
	app "gin_chat/router"
	"gin_chat/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMysql()
	r := app.Router()
	r.Run(":8081")
}
