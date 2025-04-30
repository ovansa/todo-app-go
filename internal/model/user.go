package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email        string             `json:"email" bson:"email" binding:"required,email" msg:"Email is required and must be valid"`
	FullName     string             `json:"fullName" bson:"fullName" binding:"required,min=3,max=50" msg:"Full name is required and must be between 3 and 50 characters"`
	Password     string             `json:"password,omitempty" bson:"password" binding:"required,min=6" msg:"Password is required and must be at least 6 characters"`
	PasswordHash string             `json:"-" bson:"passwordHash"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type UserRegister struct {
	Email    string `json:"email" bson:"email" binding:"required,email" msg:"Email is required and must be valid"`
	FullName string `json:"fullName" bson:"fullName" binding:"required,min=3,max=50" msg:"Full name is required and must be between 3 and 50 characters"`
	Password string `json:"password,omitempty" bson:"password" binding:"required,min=6" msg:"Password is required and must be at least 6 characters"`
}

type AuthUser struct {
	Email    string `json:"email" bson:"email" binding:"required,email"`
	Password string `json:"password,omitempty" bson:"password" binding:"required,min=6"`
}

func (u *User) HashPassword(pepper string) error {
	combined := u.Password + pepper
	hash, err := bcrypt.GenerateFromPassword([]byte(combined), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	u.Password = ""
	return nil
}

func (u *User) ComparePassword(password, pepper string) error {
	combined := password + pepper
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(combined))
}
