package auth

import (
	"crypto/md5"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive" // for BSON ObjectID
)

// User the user model
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Email    string             `json:"email"`
	Password string             `json:"password"`
	Created  int64              `json:"created"`
	Role     string             `json:"role"`
}

func (self *User) isAnonymous() bool {
	return self.ID.IsZero()
}

func (self *User) SetPassword(password string) {
	md5Password := []byte(password)
	self.Password = fmt.Sprintf("%x", md5.Sum(md5Password))
}

func (self *User) Save() (bool, error) {
	userService := new(UserService) // @TODO factory method
	result, err := userService.Save(self)
	return result, err
}

func (self *User) Delete() (bool, error) {
	userService := new(UserService) // @TODO factory method
	result, err := userService.Delete(self)
	return result, err
}
