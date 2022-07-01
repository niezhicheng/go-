package api

import (
	"awesomeProject8/mxshop-api/user-web/forms"
	"awesomeProject8/mxshop-api/user-web/global"
	"awesomeProject8/mxshop-api/user-web/global/response"
	"awesomeProject8/mxshop-api/user-web/middlewares"
	"awesomeProject8/mxshop-api/user-web/models"
	"awesomeProject8/user_srv/proto"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func HandleGrpcErrorTohttp(err error,c *gin.Context)  {
	//将grpc 的code 转换成http 的状态码
	if err != nil{
		if e,ok := status.FromError(err); ok{
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound,gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError,gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest,gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError,gin.H{
					"msg": "用户服务不可用",
				})
			case codes.AlreadyExists:
				c.JSON(http.StatusAlreadyReported,gin.H{
					"msg": "此用户已经存在了",
				})
			default:
				c.JSON(http.StatusInternalServerError,gin.H{
					"msg": "其他错误",
				})

			}
		}
	}
}

func GetUserList(ctx *gin.Context)  {
	//从注册中心获取到用户服务的信息
	cfg := api.DefaultConfig()
	consulInfo := global.ServerConfig.ConsulInfo
	zap.S().Infof("这是consulinfo",consulInfo.Host,consulInfo.Port)
	cfg.Address = fmt.Sprintf("%s:%d",consulInfo.Host,consulInfo.Port)
	fmt.Println(cfg.Address)
	userSrvHost := ""
	userSrvPort := 0
	client,err := api.NewClient(cfg)
	if err != nil{
		panic(err)
	}
	servicedata,serviceerr := client.Agent().ServicesWithFilter(`Service == "user-srv"`)

	if serviceerr != nil{
		panic(servicedata)
		zap.S().Error("是这边报错了")
	}
	for _, service := range servicedata {
		userSrvHost = service.Address
		userSrvPort = service.Port
		break
	}
	if userSrvHost == "" {
		ctx.JSON(http.StatusBadRequest,gin.H{
			"catcha":"用户服务不可到达",
		})
		return
	}





	userCon,err := grpc.Dial(fmt.Sprintf("%s:%d",userSrvHost,userSrvPort),grpc.WithInsecure())
	if err != nil{
		zap.S().Errorw("[GetUserList] 链接失败了","msg",err.Error())
	}
	s := proto.NewUserClient(userCon)
	pn := ctx.DefaultQuery("pn","0")
	psize := ctx.DefaultQuery("psize","0")
	Psize,_ := strconv.Atoi(psize)
	PnInt,_ := strconv.Atoi(pn)
	cliasa,_ := ctx.Get("claims")
	currentUser := cliasa.(*models.CustomClaims).ID
	zap.S().Infof("访问用户的id是%d",currentUser)

	data,err := s.GetUserList(ctx,&proto.PageInfo{
		Pn:    uint32(PnInt),
		PSize: uint32(Psize),
	})
	if err != nil{
		zap.S().Errorw("[GetUserList] 查询用户列表失败")
		HandleGrpcErrorTohttp(err,ctx)
		return
	}
	results := make([]interface{},0)
	for _, datum := range data.Data {
		//datainfo := make(map[string]interface{})
		var user = response.UserResponse{
			Id:       uint32(datum.Id),
			NickName: datum.NickName,
			Birthday: datum.BirthDay,
			Gender:   datum.Gender,
			Role:     string(datum.Role),
			Password: datum.Password,
			Mobile:   datum.Mobile,
		}
		//datainfo["id"] = datum.Id
		results = append(results,user )
	}
	//for i, datum := range data.Data {
	//	fmt.Println(i,datum)
	//}
	//ceshi := make([]interface{},0)
	//ceshi = append(ceshi, data.Data)
	ctx.JSON(200,results)
}

func removeTopStruct(fileds map[string]string)map[string]string  {
	rsp := map[string]string{}
	for field,err := range fileds{
		rsp[field[strings.Index(field,".")+1:]] = err
	}
	return rsp
}

func HandleValidatorError(ctx *gin.Context,err error)  {
	errs,ok := err.(validator.ValidationErrors)
	if !ok{
		ctx.JSON(http.StatusOK,gin.H{
			"msg": err.Error(),
		})
	}
	ctx.JSON(http.StatusBadRequest,gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
}


//密码登陆
func PasswordLogin(ctx *gin.Context)  {
	//表单验证
	fmt.Println("过来了")
	passwordLoginForm := forms.PassWordLoginForm{}
	if err := ctx.ShouldBind(&passwordLoginForm);err != nil{
		// 如何返回错误信息
		//zap.S().Errorw(err.Error() + "失败的愿意",)
		//zap.S().Errorw("失败了","1")
		HandleValidatorError(ctx,err)
		return
	}

	if !store.Verify(passwordLoginForm.CaptchaId,passwordLoginForm.Captcha,true) {
		ctx.JSON(http.StatusUnauthorized,gin.H{
			"msg": "验证码错误",
		})
		return
	}


	conn, _ := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port),grpc.WithInsecure())
	client := proto.NewUserClient(conn)
	if rsp, err := client.GetUserByMobile(ctx, &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	});err != nil {
		if e,ok := status.FromError(err);ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest,map[string]string{
					"mobile": "用户不存在",
				})
				return

			default:
				ctx.JSON(http.StatusInternalServerError,map[string]string{
					"mobile": "登陆失败",
				})
				return
			}
		}
	}else {
		PassResp,Paserr := client.CheckPassword(ctx,&proto.PasswordCheckInfo{
			Password:          passwordLoginForm.PassWord,
			EncryptedPassword: rsp.Password,
		})
		if Paserr != nil{
			ctx.JSON(http.StatusInternalServerError,map[string]string{
				"msg": "内部错误",
			})
			return
		}else {
			if PassResp.Success {
				//生成token
				//j := middlewares.NewJWT()
				//fmt.Println(j.SigningKey)
				tokenifno := models.CustomClaims{
					ID:             uint(rsp.Id),
					NickName:       rsp.NickName,
					AuthorityId:    uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(),
						//签名的生效时间
						//ExpiresAt: time.Now().Unix() + 60*60*24*30,
						ExpiresAt: time.Now().Unix() + 60 * 60 * 24 * 30,
						//过期时间30天
						Issuer: "xiaonie",
					},
				}
				signedString, err := middlewares.CreateTokenInfo(tokenifno)
				fmt.Println(signedString,"这是token")
				if err != nil{
					zap.S().Error("token 生成失败了")
					ctx.JSON(http.StatusInternalServerError,gin.H{
						"msg": "生成token 失败了",
					})
					return
				}
				ctx.JSON(http.StatusOK,gin.H{
					"msg": "登陆成功",
					"token": signedString,
					"username": rsp.NickName,
					"expired_at": (time.Now().Unix() + 60 * 60 * 24 * 30) * 1000,
				})
			}else {
				ctx.JSON(http.StatusBadRequest,map[string]string{
					"msg": "密码错误",
				})
			}
			return
		}
	}
}

func Register(ctx *gin.Context)  {
	regiserForm := forms.RegisterForm{}
	if err := ctx.ShouldBind(&regiserForm); err != nil{
		HandleValidatorError(ctx,err)
		return
	}
	userConn,err := grpc.Dial(fmt.Sprintf("%s:%d",global.ServerConfig.UserSrvInfo.Host,global.ServerConfig.UserSrvInfo.Port),grpc.WithInsecure())
	if err != nil{
		zap.S().Errorw("[GetUserList] 链接用户服务失败了","msg",err.Error())
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d",global.ServerConfig.RedisInfo.Host,global.ServerConfig.RedisInfo.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	r,err := rdb.Get(ctx,regiserForm.Mobile).Result()
	fmt.Println("这边",r)
	if err == redis.Nil{
		ctx.JSON(http.StatusBadRequest,gin.H{
			"code": "验证码错误",
		})
		return
	}else {
		fmt.Println(r,"这是验证码")
		if r != regiserForm.Code{
			ctx.JSON(http.StatusBadRequest,gin.H{
				"code": "验证码错误",
			})
			return
		}
	}
	fmt.Println(r)
	client := proto.NewUserClient(userConn)
	con, err := client.CreateUser(ctx, &proto.CreateUserInfo{
		NickName: regiserForm.Mobile,
		Password: regiserForm.Password,
		Mobile:   regiserForm.Mobile,
	})
	if err != nil{
		fmt.Println(err.Error(),"看看这边是什么")
		zap.S().Errorw("[注册失败了]",err.Error())
		HandleGrpcErrorTohttp(err,ctx)
		return
	}

	fmt.Println(con)
	tokenifno := models.CustomClaims{
		ID:             uint(con.Id),
		NickName:       con.NickName,
		AuthorityId:    uint(con.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			//签名的生效时间
			//ExpiresAt: time.Now().Unix() + 60*60*24*30,
			ExpiresAt: time.Now().Unix() + 30,
			//过期时间30天
			Issuer: "xiaonie",
		},
	}
	signedString, err := middlewares.CreateTokenInfo(tokenifno)
	fmt.Println(signedString,"这是token")
	if err != nil{
		zap.S().Error("token 生成失败了")
		ctx.JSON(http.StatusInternalServerError,gin.H{
			"msg": "生成token 失败了",
		})
		return
	}
	ctx.JSON(http.StatusOK,gin.H{
		"msg": "注册成功",
		"token": signedString,
		"username": con.NickName,
		"expired_at": (time.Now().Unix() + 60 * 60 * 24 * 30) * 1000,
	})
}

func Ceshiciew(c *gin.Context)  {
	zap.S().Info(c.Params.Get("png"))
	zap.S().Info(c.Query("png"))
	zap.S().Info(c.Request)
	c.JSON(200,gin.H{
		"msg": "成功了",
	})
}