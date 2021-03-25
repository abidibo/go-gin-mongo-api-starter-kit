package auth

type userSerializer struct{}

type UserData struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Created int64  `json:"created"`
	Role    string `json:"role"`
}

func NewUserSerializer() *userSerializer {
	return &userSerializer{}
}

func (self *userSerializer) Serialize(user *User) UserData {
	userData := UserData{
		ID:      user.ID.Hex(),
		Email:   user.Email,
		Created: user.Created,
		Role:    user.Role,
	}
	return userData
}

func (self *userSerializer) SerializeMany(users *[]User) []UserData {
	var res []UserData
	res = make([]UserData, 0)
	for _, user := range *users {
		res = append(res, self.Serialize(&user))
	}
	return res
}
