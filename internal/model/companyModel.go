package model

var CompanyErrors = map[string]string{
	"PIB":           "PIB is required",
	"MaticniBroj":   "Maticni Broj must be 8 numbers long",
	"Naziv":         "Has to be between 1 and 100 characters long",
	"AdresaSedista": "Has to be between 1 and 100 characters long",
	"Mesto":         "Has to be between 1 and 100 characters long",
	"PostanskiBroj": "Has to be a number shorter than 20 characters",
	"Delatnost":     "Delatnost has to be one of the enum values",
}

// Company
//
// Company represents a registered legal entity. This service
// is built around this model.
//
// It must have a physical place where its headquarters are, denoted by fields Mesto, PostanskiBroj and  Sediste.
// swagger:model company
type Company struct {
	//username of the person who is the owner of the company
	//Example: Južnobačka Oblast
	//Read Only: true
	Vlasnik string `json:"vlasnik" binding:"required"`
	// Unique number which identifies the company for taxes.
	// Required: true
	// Example: 15
	// Unique: true
	PIB int `json:"pib" binding:"required"`
	// Unique number which identifies company in the register.
	// Required: true
	// Pattern: ^\d{8}$
	// Example: 12345678
	// Unique: true
	MaticniBroj string `json:"maticniBroj" binding:"len=8,number"`
	// Full name of the company.
	// Required: true
	// Minimum length: 1
	// Maximum length: 100
	// Example: Labud DOO
	Naziv string `json:"naziv" binding:"min=1,max=100"`
	// Address at which this company's headquarters are.
	// Required: true
	// Minimum length: 1
	// Maximum length: 100
	// Example: Dositejeva 15
	AdresaSedista string `json:"adresaSedista" binding:"min=1,max=100"`
	// Place of this company's headquarters.
	// Required: true
	// Minimum length: 1
	// Maximum length: 100
	// Example: Novi Sad
	Mesto string `json:"mesto" binding:"min=1,max=100"`
	// Area code of this company's address.
	// Required: true
	// Pattern: ^\d{,20}$
	// Example: 21000
	PostanskiBroj string `json:"postanskiBroj" binding:"number,max=20"`
	// Required: true
	// Example: EDUKACIJA
	Delatnost Delatnost `json:"delatnost" binding:"required"`
	// Required: true
	Sediste Nstj `json:"sediste" binding:"required"`
}

// swagger:model nstj
type Nstj struct {
	//Example: RS123
	//Required: true
	Oznaka string `json:"oznaka"`
	//Example: Južnobačka Oblast
	//Read Only: true
	Naziv string `json:"naziv,omitempty"`
}

type CompanyFilter struct {
	OrderBy   string
	Asc       bool
	Page      int
	Mesto     string
	Sediste   string
	Delatnost string
}
