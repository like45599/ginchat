// @Author Bing
// @Date 2024-04-24 0:19:00
// @Desc
package service

import (
	"context"
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// GetUserList
// @Summary 所有用户
// @Tags 用户模块
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()

	c.JSON(200, gin.H{
		"code":    0, // 0成功， -1失败
		"message": "查询成功",
		"data":    data,
	})
}

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	//user.Name = c.Query("name")
	//password := c.Query("password")
	//repassword := c.Query("repassword")
	user.Name = c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	repassword := c.Request.FormValue("repassword")

	salt := fmt.Sprintf("%06d", rand.Int31())

	data := models.FindUserByName(user.Name)
	if user.Name == "" || password == "" || repassword == "" {
		c.JSON(-1, gin.H{
			"code":    -1, // 0成功， -1失败
			"message": "用户名或密码不能为空",
			"data":    user,
		})
		return
	}
	if data.Name != "" {
		c.JSON(-1, gin.H{
			"code":    -1, // 0成功， -1失败
			"message": "用户名已注册",
			"data":    user,
		})
		return
	}

	if password != repassword {
		c.JSON(-1, gin.H{
			"code":    -1, // 0成功， -1失败
			"message": "两次密码不一致",
			"data":    user,
		})
		return
	}
	//user.PassWord = password
	user.PassWord = utils.MakePassword(password, salt)
	user.Salt = salt
	models.CreateUser(user)
	c.JSON(200, gin.H{
		"code":    0, // 0成功， -1失败
		"message": "新增用户成功！",
		"data":    user,
	})
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id query string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"code":    0, // 0成功， -1失败
		"message": "删除用户成功！",
		"data":    user,
	})
}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")
	fmt.Println("update :", user)

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code":    0, // 0成功， -1失败
			"message": "修改参数不匹配！",
			"data":    user,
		})
	} else {
		models.UpdateUser(user)
		c.JSON(200, gin.H{
			"code":    0, // 0成功， -1失败
			"message": "修改用户成功！",
			"data":    user,
		})
	}

}

// FindUserByNameAndPwd
// @Summary 根据用户名和密码查询
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	data := models.UserBasic{}
	name := c.Query("name")
	password := c.Query("password")
	user := models.FindUserByName(name)

	fmt.Println(name)
	fmt.Println(user)
	if user.Name == "" {
		c.JSON(200, gin.H{
			"code":    -1, // 0成功， -1失败
			"message": "该用户不存在",
			"data":    data,
		})
		return
	}

	flag := utils.ValidPassword(password, user.Salt, user.PassWord)
	if !flag {
		c.JSON(200, gin.H{
			"code":    -1, // 0成功， -1失败
			"message": "密码不正确",
			"data":    data,
		})
		return
	}

	pwd := utils.MakePassword(password, user.Salt)
	data = models.FindUserByNameAndPwd(name, pwd)

	c.JSON(200, gin.H{
		"code":    0, // 0成功， -1失败
		"message": "登录成功",
		"data":    data,
	})
}

// UserLogin
// @Summary 根据用户名和密码查询
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/login [post]
func UserLogin(c *gin.Context) {
	name := c.Request.FormValue("name")
	password := c.Request.FormValue("password")

	user := models.FindUserByName(name) // 修改这里，只通过用户名查找用户
	if user.ID == 0 {                   // 使用 ID 作为用户存在的判断，因为Name字段可能是空字符串
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    -1,
			"message": "该用户不存在",
		})
		return
	}

	// 检查密码是否正确
	if !utils.ValidPassword(password, user.Salt, user.PassWord) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    -1,
			"message": "密码不正确",
		})
		return
	}

	// 更新身份验证令牌
	str := fmt.Sprintf("%d", time.Now().Unix())
	temp := utils.MD5Encode(str)
	utils.DB.Model(&user).Update("Identity", temp)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
		"data":    user,
	})
}

// 防止跨域站点伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 使用 upgrader 来升级 HTTP 到 WebSocket
func SendMsg(c *gin.Context) {
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(ws, c)
}

func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	// 创建一个单独的 goroutine 来处理来自 Redis 的订阅消息
	go func() {
		for {
			fmt.Println("Waiting for Redis messages...")
			msg, err := utils.Subscribe(c, utils.PublishKey)
			if err != nil {
				fmt.Println("Error subscribing to Redis channel:", err)
				return // 或者考虑重试逻辑
			}
			fmt.Println("Received message from Redis:", msg)
			tm := time.Now().Format("2006-01-02 15:04:05")
			m := fmt.Sprintf("[ws][%s]: %s", tm, msg)
			err = ws.WriteMessage(websocket.TextMessage, []byte(m))
			if err != nil {
				fmt.Println("Error sending message via WebSocket:", err)
				return
			}
		}
	}()

	// 循环读取 WebSocket 客户端发送的消息
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("Error reading from WebSocket:", err)
			break
		}
		fmt.Println("Received message from WebSocket client:", string(msg))

		// 可选择发布这个消息到 Redis
		ctx := context.Background()
		err = utils.Publish(ctx, utils.PublishKey, string(msg))
		if err != nil {
			fmt.Println("Failed to publish message to Redis:", err)
			continue
		}
	}
}

func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
