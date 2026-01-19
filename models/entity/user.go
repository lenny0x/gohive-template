package entity

type User struct {
	BaseModel
	Username string `gorm:"size:64;uniqueIndex" json:"username"`
	Email    string `gorm:"size:128;uniqueIndex" json:"email"`
	Password string `gorm:"size:256" json:"-"`
	Nickname string `gorm:"size:64" json:"nickname"`
	Avatar   string `gorm:"size:256" json:"avatar"`
	Status   int    `gorm:"default:1" json:"status"` // 1: active, 0: inactive
}

func (User) TableName() string {
	return "users"
}
