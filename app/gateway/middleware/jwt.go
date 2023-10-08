package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/CocaineCong/grpc-todolist/pkg/ctl"
	"github.com/CocaineCong/grpc-todolist/pkg/e"
	"github.com/CocaineCong/grpc-todolist/pkg/util/jwt"
)

// JWT token验证中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}
		code = 200
		token := c.GetHeader("Authorization")

		if token == "" {
			code = 404
			c.JSON(200, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
				"data":   data,
			})
			c.Abort()
		}
		claims, err := jwt.ParseToken(token)
		if err != nil {
			//fmt.Println("token解析失败")
			code = e.ErrorAuthCheckTokenFail
		} else if time.Now().Unix() > claims.ExpiresAt {
			//fmt.Println("token超时")
			code = e.ErrorAuthCheckTokenTimeout
		}
		if code != e.SUCCESS {
			c.JSON(200, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
				"data":   data,
			})
			c.Abort()
			return
		}
		// 解析之后将token信息，加入到gin的context中
		c.Request = c.Request.WithContext(ctl.NewContext(c.Request.Context(), &ctl.UserInfo{Id: claims.UserID}))
		ctl.InitUserInfo(c.Request.Context())
		c.Next()
	}
}
