package api

import (
	"awesomeProject8/mxshop-api/user-web/forms"
	"awesomeProject8/mxshop-api/user-web/global"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"math/rand"
	"strings"
	"time"
	"encoding/json"
)

func Sednnifo()  {
	credential := common.NewCredential(
		global.ServerConfig.TengXuninfo.Apikey,
		global.ServerConfig.TengXuninfo.ApiSecrect,
	)

	cpf := profile.NewClientProfile()

	cpf.HttpProfile.ReqMethod = "POST"

	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"

	client, _ := sms.NewClient(credential, "ap-beijing", cpf)

	/* 实例化一个请求对象，根据调用的接口和实际情况*/
	request := sms.NewSendSmsRequest()

	// 应用 ID 可前往 [短信控制台](https://console.cloud.tencent.com/smsv2/app-manage) 查看
	request.SmsSdkAppId = common.StringPtr("1400441219")

	// 短信签名内容: 使用 UTF-8 编码，必须填写已审核通过的签名
	request.SignName = common.StringPtr("墨韵码站")

	/* 模板 ID: 必须填写已审核通过的模板 ID */
	request.TemplateId = common.StringPtr("754761")

	/* 模板参数: 模板参数的个数需要与 TemplateId 对应模板的变量个数保持一致，若无模板参数，则设置为空*/
	request.TemplateParamSet = common.StringPtrs([]string{GenerateSmsCode(6),"5"})

	/* 下发手机号码，采用 E.164 标准，+[国家或地区码][手机号]
	 * 示例如：+8613711112222， 其中前面有一个+号 ，86为国家码，13711112222为手机号，最多不要超过200个手机号*/
	request.PhoneNumberSet = common.StringPtrs([]string{"+8613256735690"})

	// 通过client对象调用想要访问的接口，需要传入请求对象
	response, err := client.SendSms(request)
	// 处理异常
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	// 非SDK异常，直接失败。实际代码中可以加入其他的处理。
	if err != nil {
		panic(err)
	}
	b, _ := json.Marshal(response.Response)
	// 打印返回的json字符串
	fmt.Printf("%s", b)
}

func SendSmg(ctx *gin.Context)  {
	sendMsgForm := forms.SendSmsForm{}
	if err := ctx.ShouldBind(&sendMsgForm);err != nil{
		HandleValidatorError(ctx,err)
		return
	}
	fmt.Println(sendMsgForm.Mobile)
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d",global.ServerConfig.RedisInfo.Host,global.ServerConfig.RedisInfo.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	code := GenerateSmsCode(6)
	err := rdb.Set(ctx, sendMsgForm.Mobile, code,time.Duration(global.ServerConfig.RedisInfo.Expire) * time.Second).Err()
	if err != nil {
		panic(err)
	}
	ctx.JSON(200,gin.H{
		"msg": "成功了",
		"code": code,
	})
}



// GenerateSmsCode 生成验证码;length代表验证码的长度
func GenerateSmsCode(length int) string {
	numberic := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	rand.Seed(time.Now().Unix())
	var sb strings.Builder
	for i := 0; i < length; i++ {
		fmt.Fprintf(&sb, "%d", numberic[rand.Intn(len(numberic))])
	}
	return sb.String()
}