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

// User
//
// User represents a physical person who is registered in this service.
// APR service stores user credentials and personal data.
//
//swagger:model user
type User struct {
	// Example: +381123123
	Phone string `json:"phone,omitempty" binding:"omitempty,e164"`
	// Example: user@example.com
	Email string `json:"email"    binding:"omitempty,email"`
	// Required: true
	// Example: MALE
	Sex string `json:"sex"      binding:"sex"`
	// Maximum length: 100
	// Example: Dositejeva 15
	Address string `json:"address,omitempty" binding:"omitempty,max=100"`
	// Required: true
	// Maximum length: 100
	// Example: Petar
	Name string `json:"name"     binding:"required,max=100"`
	// Required: true
	// Maximum length: 100
	// Example: Petrovic
	Lastname string `json:"lastname" binding:"required,max=100"`
	// Required: true
	// Minimum length: 12
	// Maximum length: 72
	Password string `json:"password" binding:"min=12,max=72,required"`
	// Required: true
	// Pattern: ^[a-zA-Z0-9]{4,20}$
	// Unique: true
	Username string `json:"username" binding:"alphanum,min=4,max=20,required"`
	// Required: true
	// Pattern: ^\13{8}$
	// Unique: true
	// Example: 1234567891234
	Jmbg string `json:"jmbg"     binding:"number,len=13,required"`
}
