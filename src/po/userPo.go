package po

type UserPo struct {
	UserId         int64  `gorm:"primary_key;column:user_id"`
	PassWord       string `gorm:"column:pass_word"`
	NickName       string `gorm:"column:nick_name"`
	ProfilePicture string `gorm:"column:profile_picture"`
}

func (UserPo) TableName() string {
	return "user_info_tab"
}
