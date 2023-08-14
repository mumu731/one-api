package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"one-api/common"
	"strings"
)

type JwtBody struct {
	Username string `json:"username"`
}
type JwtAuthBody struct {
	Token string `json:"token"`
}

type MyClaims struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Role     int    `json:"role"`
	jwt.StandardClaims
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	//解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return common.JwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// JWTAuthMiddleware 基于JWT认证中间件
func JWTAuthMiddleware(minRole int) func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code":    2003,
				"message": "请求头中的auth为空",
			})
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)

		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code":    2004,
				"message": "请求头中的auth格式错误",
			})
			//阻止调用后续的函数
			c.Abort()
			return
		}
		mc, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    2005,
				"message": "无效的token",
			})
			c.Abort()
			return
		}
		println(mc.Role)
		println(minRole)

		if mc.Role < minRole {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "无权进行此操作，权限不足",
			})
			c.Abort()
			return
		}
		//将当前请求的username信息保存到请求的上下文c上
		c.Set("username", mc.Username)
		c.Set("role", mc.Role)
		c.Set("id", mc.Id)

		c.Next()
	}

}
