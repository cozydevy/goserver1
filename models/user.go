package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Email          string `gorm:"unique_index; not null"`
	Password       string `gorm:"not null"`
	Name           string `gorm:"not null"`
	Avatar         string
	Role           string `gorm:"default:'Member'; not null"`
	Progress       []Progress
	ProgressDetail []ProgressDetail
}

func (u *User) Promote() {
	u.Role = "Admin"
}

func (u *User) Demote() {
	u.Role = "Member"
}

func (u *User) GenerateEncryptedPassword() string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	return string(hash)
}

func (u *User) GenerateEncryptedPassword2(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hash)
}
