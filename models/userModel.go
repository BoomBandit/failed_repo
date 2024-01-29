package models

import (
	"time"

	valid "github.com/asaskevich/govalidator"
)

type User struct {
	ID        int       `json:"id" gorm:"primarykey"`
	Username  string    `json:"username"  gorm:"unique;not null" valid:"required"`
	Email     string    `json:"email" gorm:"unique;not null" valid:"required,email"`
	Password  string    `json:"password"  gorm:"not null" valid:"required,minstringlength(6)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) Validate() (bool, error) {
	return valid.ValidateStruct(u)
}
