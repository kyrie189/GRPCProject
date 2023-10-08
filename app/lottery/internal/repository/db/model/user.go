package model

import "gorm.io/gorm"

type Lottery struct {
	gorm.Model
	UserId int32 `gorm:"column:user_id;comment:'用户id'" json:"userid"`
	LuckyTime int64 `gorm:"column:lucky_time;comment:'中间时间'" json:"created"`
	Award string `gorm:"column:award;comment:'奖品'" json:"award"`
}

func (*Lottery) TableName() string {
	return "lottery"
}

// SetPassword 加密密码
func (l *Lottery) SetPassword(password string) error {
	//bytes, err := bcrypt.GenerateFromPassword([]byte(password), consts.PassWordCost)
	//if err != nil {
	//	return err
	//}
	//l.PasswordDigest = string(bytes)
	return nil
}

// CheckPassword 检验密码
func (l *Lottery) CheckPassword(password string) bool {
	//err := bcrypt.CompareHashAndPassword([]byte(l.PasswordDigest), []byte(password))
	return false
}
