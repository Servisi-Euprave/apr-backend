package model

var UserErrors = map[string]string{
	"Phone":    "Phone must follow the E.164 standard",
	"Email":    "Must be a valid email",
	"Sex":      "Must be either MALE or FEMALE",
	"Address":  "Cannot be longer than 100 characters",
	"Name":     "Required, cannot be longer than 100 characters",
	"Lastname": "Required, cannot be longer than 100 characters",
	"Username": "Must use only letters and numbers, and be between 4 and 20 characters",
	"Jmbg":     "Must be 13 numbers long",
	"Password": "Must be between 12 and 72 characters",
}

type User struct {
	Phone    string `json:"phone,omitempty" binding:"omitempty,e164"`
	Email    string `json:"email"    binding:"omitempty,email"`
	Sex      string `json:"sex"      binding:"sex"`
	Address  string `json:"address,omitempty" binding:"omitempty,max=100"`
	Name     string `json:"name"     binding:"required,max=100"`
	Lastname string `json:"lastname" binding:"required,max=100"`
	Password string `json:"password" binding:"min=12,max=72,required"`
	Username string `json:"username" binding:"alphanum,min=4,max=20,required"`
	Jmbg     string `json:"jmbg"     binding:"number,len=13,required"`
}
