package myjwt

import (
	"douyin/pkg/common"
	"douyin/pkg/constants"
	"douyin/pkg/errno"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	jwt4 "github.com/golang-jwt/jwt/v4"
	"net/http"
	"strconv"
	"time"
)

const (
	IdentityKey       = "userId"
	SecretKey         = "lym_hrh"
	TokenValidTime    = 7 * 24 * time.Hour
	Token             = "token"
	HeadAuthorization = "Authorization"

	Issuer = "douyin_lym_lr_hrh"
)

func ReleaseToken(userId int64) (string, error) {
	claims := jwt4.MapClaims{
		IdentityKey: userId,
		"exp":       time.Now().Add(TokenValidTime).Unix(),
		"iat":       time.Now().Unix(),
		"iss":       Issuer,
	}

	token := jwt4.NewWithClaims(jwt4.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

// Authenticate 使用这个，上面的没有，但这个依赖于GinJWTMiddleware，因为TimeFunc()，IdentityKey
//但其实也这在constants里写死它，一样的
func Authenticate(c *gin.Context, tokenString string) error {
	claims := jwt4.MapClaims{}
	if _, err := jwt4.ParseWithClaims(tokenString, &claims, func(token *jwt4.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	}); err != nil {
		return fmt.Errorf("parse token failed,question:%v", err)
	}
	if claims["exp"] == nil {
		return fmt.Errorf("claims[exp] == nil")
	}

	if v, ok := claims["exp"].(float64); !ok {
		return fmt.Errorf("claims[exp] is not float64,but %T %[1]v", v)
	}

	if int64(claims["exp"].(float64)) < time.Now().Unix() {
		return errors.New("token is expired")
	}

	userId, ok := claims[IdentityKey].(float64)
	if !ok {
		return errno.IdTypeErr.WithMessage(fmt.Sprintf("%T %#[1]v,id is not exist in claims", userId))
	}

	c.Set(IdentityKey, int64(userId))
	return nil
}
func Unauthorized(c *gin.Context, err error) {
	c.Abort()
	c.JSON(http.StatusOK, common.Response{
		StatusCode: http.StatusUnauthorized,
		StatusMsg:  "you have not the authority," + err.Error(),
	})
}
func Authorizator(c *gin.Context) bool {
	tokenIdentity := ""
	if tokenIdentity = c.Query(constants.UserIdQuery); tokenIdentity == "" {
		if tokenIdentity = c.PostForm(constants.UserIdQuery); tokenIdentity == "" {
			return true
		}
	}
	tokenId, err := strconv.ParseInt(tokenIdentity, 10, 64)
	if err != nil {
		return false
	}

	return tokenId == c.Keys[IdentityKey]
}
func JWTMiddleWareImpl(c *gin.Context) {
	var err error
	token := ""
	if token = c.Query(Token); token == "" {
		if token = c.PostForm(Token); token == "" {
			if token = c.GetHeader(HeadAuthorization); token == "" {
				Unauthorized(c, errors.New("can not find the token"))
				return
			}
		}
	}

	if err = Authenticate(c, token); err != nil {
		Unauthorized(c, err)
		return
	}
	if !Authorizator(c) {
		Unauthorized(c, errors.New("the tokenIdentity is not equal with the identity from request"))
		return
	}

	c.Next()
}

func IdentityHandler(c *gin.Context) (int64, error) {
	id, exists := c.Get(IdentityKey)
	if !exists {
		return 0, errors.New("userId not in context")
	}

	userId, ok := id.(int64)
	if !ok {
		return 0, errno.IdTypeErr
	}

	return userId, nil
}
