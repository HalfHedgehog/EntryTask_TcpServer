package dto

type UserInfoDto struct {
	UserId         int64  `json:"userId"`
	NickName       string `json:"nickName"`
	ProfilePicture string `json:"profilePicture""`
}
