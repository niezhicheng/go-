package middlewares

import (
	"awesomeProject8/mxshop-api/user-web/global"
	"awesomeProject8/mxshop-api/user-web/models"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func JWTAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.Request.Header.Get("x-token")
		if token == "" {
			context.JSON(http.StatusUnauthorized,map[string]string{
				"msg": "请登陆",
			})
			context.Abort()
			return
		}
		j := NewJWT()
		claims,err := j.ParseToken(token)
		if err != nil{
			if err == TokenExpired{
				if err == TokenExpired{
					context.JSON(http.StatusUnauthorized,map[string]string{
						"msg": "授权已经过期",
					})
					context.Abort()
					return
				}
			}
			context.JSON(http.StatusUnauthorized,"未登陆")
			context.Abort()
			return
		}
		context.Set("claims",claims)
		context.Set("userId",claims.Id)
		context.Next()
	}
}

var (
	TokenExpired = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed = errors.New("That's not even a token")
	TokenInvalid = errors.New("Couldn't handle this token:")

)

type JWT struct {
	SigningKey []byte
}
func NewJWT() *JWT {
	return &JWT{SigningKey: []byte("9zZ#2gg+7pGO9C12NahNF#5KIEc#HQ5E")}
}

func (j *JWT) CreateTOken(claims *models.CustomClaims) (string,error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256,claims)
	//s := []byte("9zZ#2gg+7pGO9C12NahNF#5KIEc#HQ5E")
	return token.SignedString([]byte("9zZ#2gg+7pGO9C12NahNF#5KIEc#HQ5E"))
}

func (j *JWT) ParseToken(tokenString string) (*models.CustomClaims,error) {
	token,err := jwt.ParseWithClaims(tokenString,&models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey,nil
	})
	if err != nil{
		if ve,ok := err.(*jwt.ValidationError); ok{
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil,TokenMalformed
			}else if ve.Errors&jwt.ValidationErrorExpired != 0{
				return nil,TokenExpired
			}
		}
	}
	if token != nil{
		if clasims,ok := token.Claims.(*models.CustomClaims);ok && token.Valid{
			return clasims, nil
		}
		return nil, TokenInvalid
	}else {
		return nil, TokenInvalid
	}
}

func (j *JWT) RefreshToken(tokenString string) (string,error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0,0)
	}
	token,err := jwt.ParseWithClaims(tokenString,&models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey,nil
	})
	if err != nil{
		return "", err
	}
	if claims,ok := token.Claims.(*models.CustomClaims);ok &&token.Valid{
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateTOken(claims)
	}
	return "",TokenInvalid
}

func CreateTokenInfo(claims models.CustomClaims) (string,error) {
	fmt.Println("这便是",global.ServerConfig.JWTInfo.SigningKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s := []byte("9zZ#2gg+7pGO9C12NahNF#5KIEc#HQ5E")
	signedString, err := token.SignedString(s)
	if err != nil{
		return "", err
	}
	return signedString, err
}



