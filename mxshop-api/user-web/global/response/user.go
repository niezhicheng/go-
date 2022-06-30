package response

import (
	"time"
)

type JsonTime time.Time

//func (j JsonTime) MarshalJson() ([]byte,error) {
//	var stmp = fmt.Sprintf("%s",time.Time())
//}

type UserResponse struct {
	Id uint32 `json:"id"`
	NickName string `json:"nick_name"`
	Birthday string `json:"birthday"`
	//Birthday string `json:"birthday"`
	Gender string `json:"gender"`
	Role string `json:"role"`
	Password string `json:"password"`
	Mobile string `json:"mobile"`
}
