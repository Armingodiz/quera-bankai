package models

type User struct {
	UserId   uint   `gorm:"primaryKey;autoIncrement"`
	Username string `db:"user_name" json:"user_name"`
	Password string `db:"password" json:"password"`
	Admin    bool   `db:"admin" json:"admin"`
}
