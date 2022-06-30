package forms

type PassWordLoginForm struct {
	Mobile string `json:"mobile" form:"mobile" binding:"required,mobile"`
	PassWord string `json:"password" form:"password" binding:"required,min=3,max=10"`
	Captcha string `json:"captcha" form:"captcha" binding:"required,min=5,max=5"`
	CaptchaId string `json:"captcha_id" form:"captcha_id" binding:"required"`
}

type RegisterForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"`
	Password string `form:"password" json:"password" binding:"required,min=3,max=20"`
	Code string `form:"code" json:"code" binding:"required,min=6,max=6"`

}