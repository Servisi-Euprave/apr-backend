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

// Person
//
// Person represents a physical person.
//
//swagger:model user
type Person struct {
	// Required: true
	// Maximum length: 100
	// Example: Petar
	Name string `json:"name"     binding:"required,max=100"`
	// Required: true
	// Maximum length: 100
	// Example: Petrovic
	Lastname string `json:"lastname" binding:"required,max=100"`
	Jmbg     string `json:"jmbg"     binding:"number,len=13,required"`
}
