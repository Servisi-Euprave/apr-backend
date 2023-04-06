package model

type User struct {
	Phone    string `json:"phone,omitempty" binding:"omitempty,e164"`
	Email    string `json:"email"    binding:"omitempty,email"`
	Sex      string `json:"sex"      binding:"sex"`
	Address  string `json:"address"  binding:"omitempty,max=100"`
	Name     string `json:"name"     binding:"required,max=100"`
	Lastname string `json:"lastname" binding:"required,max=100"`
	Password []byte `json:"password" binding:"min=12,max=72,required"`
	Username string `json:"username" binding:"alphanum,required"`
	Jmbg     string `json:"jmbg"     binding:"number,required"`
}
