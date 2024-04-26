// @Author Bing
// @Date 2024-04-24 0:09:00
// @Desc
package main

import (
	"ginchat/router"
	"ginchat/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()
	r := router.Router()
	r.Run(":8081")
}
