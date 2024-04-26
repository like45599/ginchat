// @Author Bing
// @Date 2024-04-24 0:19:00
// @Desc
package service

import (
	"github.com/gin-gonic/gin"
	"html/template"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	ind, err := template.ParseFiles("D:/GAOLIKE/goproject/msbProject/ginchat/view/user/test.html")
	if err != nil {
		panic(err)
	}
	ind.ExecuteTemplate(c.Writer, "/user/test.shtml", nil)
	//c.JSON(200, gin.H{
	//	"message": "welcome!!",
	//})
}

func ToRegister(c *gin.Context) {
	ind, err := template.ParseFiles("view/user/register.html")
	if err != nil {
		panic(err)
	}
	ind.ExecuteTemplate(c.Writer, "/user/register.shtml", nil)
	//c.JSON(200, gin.H{
	//	"message": "welcome!!",
	//})
}
