package handler

import (
	"filestore-server/common"
	"filestore-server/util"
	"fmt"
	"log"
	"net/http"
	"time"

	dbcli "filestore-server/service/dbproxy/client"

	"github.com/gin-gonic/gin"
)

// Authorize : http请求拦截器
func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Request.FormValue("username")
		token := c.Request.FormValue("token")

		//验证登录token是否有效
		if len(username) < 3 || !IsTokenValid(token, username) {
			// w.WriteHeader(http.StatusForbidden)
			// token校验失败则跳转到登录页面
			c.Abort()
			resp := util.NewRespMsg(
				int(common.StatusTokenInvalid),
				"token无效",
				nil,
			)
			c.JSON(http.StatusOK, resp)
			return
		}
		c.Next()
	}
}

// IsTokenValid : token是否有效
func IsTokenValid(token string, username string) bool {
	if len(token) != 40 {
		log.Println("token invalid: " + token)
		return false
	}
	// example，假设token的有效期为1天   (根据同学们反馈完善, 相对于视频有所更新)
	tokenTS := token[32:40]
	if util.Hex2Dec(tokenTS) < time.Now().Unix()-86400 {
		log.Println("token expired: " + token)
		return false
	}
	// example, IsTokenValid方法增加传入参数username
	dbToken, err := dbcli.GetUserToken(username)
	if err != nil || dbToken != token {
		return false
	}

	return true
}

// GenToken : 生成token
func GenToken(username string) string {
	// 40位字符:md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}
