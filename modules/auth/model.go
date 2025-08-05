package auth

import (
	"gin-boilerplate/modules/primitive"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	primitive.BaseModel
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"uniqueIndex;not null"`
	Password string `json:"-" gorm:"not null"`
}

// TableName returns the table name for User model
func (User) TableName() string {
	return "users"
}

// HashPassword hashes the user's password
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword checks if the provided password matches the user's password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// ToUserInfo converts User to UserInfo (without sensitive data)
func (u *User) ToUserInfo() primitive.UserInfo {
	return primitive.UserInfo{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}
