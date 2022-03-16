package service

import (
	"TcpServer/src/Util"
	"TcpServer/src/dto"
	"TcpServer/src/enum"
	"TcpServer/src/global"
	"TcpServer/src/po"
	"TcpServer/src/rpc/userRpc"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"log"
	"strconv"
	"time"
)

type UserServe struct {
	userRpc.UnimplementedSearchServiceServer
}

func (*UserServe) SignIn(ctx context.Context, req *userRpc.UserDto) (*userRpc.SignInRes, error) {
	userId := Util.GetGlobalId()
	user := po.UserPo{
		UserId:         userId,
		PassWord:       req.Password,
		NickName:       req.NickName,
		ProfilePicture: req.ProfilePicture,
	}
	global.DBHelper.Create(user)
	return &userRpc.SignInRes{Flag: true, Message: "注册成功", UserId: userId}, nil
}

// LogIn 登陆
func (UserServe) LogIn(ctx context.Context, req *userRpc.UserDto) (*userRpc.LogInRes, error) {
	//查询是否有当前用户
	var userList []po.UserPo
	global.DBHelper.Where("user_id = ? and pass_word = ?", req.UserID, req.Password).Find(&userList)
	var res userRpc.LogInRes
	if len(userList) > 0 {
		token := Util.CreateToken(req.UserID)
		global.RedisHelper.Set(ctx, string(enum.LogInPre)+strconv.FormatInt(req.UserID, 10), req.UserID, 30*time.Minute)
		res = userRpc.LogInRes{
			Flag:    true,
			Message: "登陆成功",
			Token:   token,
		}
	} else {
		res = userRpc.LogInRes{
			Flag:    false,
			Message: "用户名或密码错误",
		}
	}
	return &res, nil
}

func (UserServe) GetUserInfo(ctx context.Context, req *userRpc.UserDto) (*userRpc.UserInfoRes, error) {
	userId := req.UserID
	//查询缓存如果查询到了就直接返回
	if cash, flag := findUserInfoFromRedis(userId); flag {
		//构建返回结果
		userInfo := userRpc.UserDto{
			UserID:         cash.UserId,
			NickName:       cash.NickName,
			ProfilePicture: cash.ProfilePicture,
		}
		return &userRpc.UserInfoRes{
			Flag:    true,
			Message: "查询成功",
			UserDto: &userInfo,
		}, nil
	}
	//查询数据库
	return findUserInfoFromDb(userId)
}

//从redis获取数据
func findUserInfoFromRedis(userId int64) (*dto.UserInfoDto, bool) {
	cash := dto.UserInfoDto{}
	val, err := global.RedisHelper.Get(context.Background(), string(enum.GetInfoPre)+strconv.FormatInt(userId, 10)).Result()
	if err != nil || val == "" {
		return &cash, false
	}
	err = json.Unmarshal([]byte(val), &cash)
	if err != nil {
		return &cash, false
	}
	return &cash, true
}

//从数据库获取用户信息
func findUserInfoFromDb(userId int64) (*userRpc.UserInfoRes, error) {
	//从数据库里查询
	var userList []po.UserPo
	global.DBHelper.Where("user_id = ?", userId).First(&userList)
	if len(userList) > 0 {
		user := userList[0]
		cash := dto.UserInfoDto{
			UserId:         user.UserId,
			NickName:       user.NickName,
			ProfilePicture: user.ProfilePicture,
		}

		//缓存结果
		err := global.RedisHelper.Set(context.Background(), string(enum.GetInfoPre)+strconv.FormatInt(userId, 10), cash, 10*time.Minute).Err()
		if err != nil {
			log.Printf("缓存错误：%s", err)
		}
		//构建返回结果
		userInfo := userRpc.UserDto{
			UserID:         user.UserId,
			NickName:       user.NickName,
			ProfilePicture: user.ProfilePicture,
		}

		return &userRpc.UserInfoRes{
			Flag:    true,
			Message: "查询成功",
			UserDto: &userInfo,
		}, nil
	} else {
		return &userRpc.UserInfoRes{
			Flag:    false,
			Message: "用户信息不存在",
		}, nil
	}
}

func (UserServe) UpdateUserInfo(ctx context.Context, req *userRpc.UserDto) (*userRpc.UpdateInfoRes, error) {
	user := po.UserPo{
		UserId:         req.UserID,
		PassWord:       req.Password,
		NickName:       req.NickName,
		ProfilePicture: req.ProfilePicture,
	}

	//删除缓存
	err := global.RedisHelper.Del(context.Background(), string(enum.GetInfoPre)+strconv.FormatInt(user.UserId, 10)).Err()
	if err != nil {
		fmt.Println("redis异常")
		err = nil
	}

	//更新数据库
	err = global.DBHelper.Debug().Model(&user).Omit("user_id").Where("user_id = ?", user.UserId).Updates(user).Error
	if err != nil {
		fmt.Println("数据库异常")
		return &userRpc.UpdateInfoRes{
			Flag:    false,
			Message: "更新失败",
		}, nil
	}

	//再次删除缓存保证缓存一致
	err = global.RedisHelper.Del(context.Background(), string(enum.GetInfoPre)+strconv.FormatInt(user.UserId, 10)).Err()
	if err != nil {
		fmt.Println("redis异常")
	}
	return &userRpc.UpdateInfoRes{
		Flag:    true,
		Message: "更新成功",
	}, nil
}
