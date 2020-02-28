package auth

import (
	"encoding/json"

	"elementser/src/common"
	"fmt"

	help "elementser/src/help"
	"github.com/gin-gonic/gin"
	"time"
)

type User struct{}

type UserData struct {
	Menus        []string `json:"roles"`        // 菜单
	Introduction string   `json:"introduction"` // 介绍
	Avatar       string   `json:"avatar"`       // 图标
	Name         string   `json:"name"`         // 姓名
}

// 用户登录
func (User) Login(c *gin.Context) {
	requestData, err := c.GetRawData() //从request payload中获取参数

	if err != nil {
		common.ResErrSrv(c, err)
		return
	}

	var requestMap map[string]string
	err = json.Unmarshal(requestData, &requestMap)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	username := requestMap["username"]
	password := requestMap["password"]
	fmt.Println("in msg:", username, password)
	if username == "" || password == "" {
		common.ResFail(c, "用户名或密码不能为空")
		return
	}

	exptime := time.Now().Add(time.Duration(3600) * time.Second).Unix()

	token, err := help.GenerateJwtToken(&help.JwtPayload{
		AppID:     "id",
		ExpiresAt: exptime,
	}, true)

	println(token.AccessToken)

	resData := make(map[string]interface{})
	resData["token"] = token.AccessToken
	common.ResSuccess(c, &resData)
}

func (User) Info(c *gin.Context) {
	var menus = []string{"admin"}
	resData := UserData{Menus: menus, Name: "小王"}
	resData.Avatar = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	common.ResSuccess(c, &resData)
}

func (User) Logout(c *gin.Context) {
	common.ResSuccess(c, "success")
}
