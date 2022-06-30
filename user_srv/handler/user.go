package handler

import (
	"awesomeProject8/user_srv/proto"
	"context"
	"awesomeProject8/user_srv/global"
	"awesomeProject8/user_srv/model"
	"crypto/sha512"
	"github.com/anaskhan96/go-password-encoder"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"fmt"
	"strings"
	"time"
)

type UserServer struct {
}

func ModelToResponse(user model.User) proto.UserInfoResponse {
	// 在grpc 的message 中字段有默认值，你不能随意nill 进去
	userInfoRsp := proto.UserInfoResponse{
		Id:       int32(user.ID),
		Password: user.Password,
		Mobile:   user.Mobile,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     int32(user.Role),
	}
	if user.Birthday != nil{
		userInfoRsp.BirthDay = string(user.Birthday.Unix())
	}
	return userInfoRsp
}


func Paginate(page,pageSize int) func(db *gorm.DB) *gorm.DB {
	return func (db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (s *UserServer) GetUserList(ctx context.Context,info *proto.PageInfo) (*proto.UserListResponse,error) {
	//获取用户列表
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil{
		return nil,result.Error
	}
	rsp := &proto.UserListResponse{
		Total: int32(result.RowsAffected),
	}
	fmt.Println(info.PSize,info.Pn)
	global.DB.Scopes(Paginate(int(info.Pn),int(info.PSize))).Find(&users)
	for _, user := range users {
		userInfoRsp := ModelToResponse(user)
		rsp.Data = append(rsp.Data,&userInfoRsp)
	}

	return rsp,nil
}


//通过手机号拿到用户
func (s *UserServer) GetUserByMobile(ctx context.Context,request *proto.MobileRequest) (*proto.UserInfoResponse,error)  {
	var users model.User
	result := global.DB.Where(&model.User{Mobile: request.Mobile}).First(&users)
	if result.RowsAffected ==0 {
		return nil,status.Errorf(codes.NotFound,"没有此用户")
	}
	if result.Error != nil{
		return nil,result.Error
	}
	userInfoRsp := ModelToResponse(users)
	return &userInfoRsp,nil
}
//通过id 拿到用户数据
func (s *UserServer) GetUserById(ctx context.Context,request *proto.IdRequest) (*proto.UserInfoResponse,error) {
	var users model.User
	results := global.DB.First(&users,request.Id)
	if results.RowsAffected == 0{
		return nil,status.Errorf(codes.NotFound,"没有此用户")
	}
	if results.Error != nil{
		return nil,results.Error
	}
	userInfoRsp := ModelToResponse(users)
	return &userInfoRsp,nil
}


//创建用户

func (s *UserServer) CreateUser(ctx context.Context,info *proto.CreateUserInfo) (*proto.UserInfoResponse,error) {
	var user model.User
	results := global.DB.Where(&model.User{Mobile: info.Mobile}).First(&user)
	if results.RowsAffected == 1 {
		return nil,status.Errorf(codes.AlreadyExists,"此用户已经存在了")
	}
	user.Mobile = info.Mobile
	user.NickName = info.NickName
	//密码加密

	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode(info.Password, options)
	newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s",salt,encodedPwd)
	user.Password = newPassword
	result := global.DB.Create(&user)
	if result.Error != nil{
		return nil,status.Errorf(codes.Internal,result.Error.Error())
	}
	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp,nil
}

func (s *UserServer) UpdateUser(ctx context.Context,info *proto.UpdateUserInfo) (*emptypb.Empty,error) {
	var user model.User
	result := global.DB.First(&user,info.Id)
	if result.RowsAffected == 0 {
		return nil,status.Errorf(codes.NotFound,"没有此用户")
	}
	birthDay := time.Unix(int64(info.BrithDay),0)
	user.NickName = info.NickName
	user.Birthday = &birthDay
	user.Gender = info.Gender
	result = global.DB.Save(user)
	if result.Error != nil{
		return nil,status.Errorf(codes.Internal,"没有保存成功",result.Error.Error())
	}
	return &empty.Empty{},nil
}



//校验密码

func (s *UserServer) CheckPassword(ctx context.Context,info *proto.PasswordCheckInfo) (*proto.CheckResponse,error) {
	fmt.Println("这边来了",info.Password)
	options := &password.Options{16,100,32,sha512.New}
	passwordInfo := strings.Split(info.EncryptedPassword,"$")
	fmt.Println(passwordInfo)
	check := password.Verify(info.Password, passwordInfo[2], passwordInfo[3], options)
	fmt.Println(check) // true
	return &proto.CheckResponse{Success: check},nil
}