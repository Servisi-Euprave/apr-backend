package model

var CompanyErrors = map[string]string{
	"PIB":           "PIB is required",
	"Naziv":         "Has to be between 1 and 100 characters long",
	"AdresaSedista": "Has to be between 1 and 100 characters long",
	"Mesto":         "Has to be between 1 and 100 characters long",
	"PostanskiBroj": "Has to be a number shorter than 20 characters",
	"Delatnost":     "Delatnost has to be one of the enum values",
	"Vlasnik":       "JMBG is required",
	"Sediste":       "Sediste is required",
}

// Company
//
// Company represents a registered legal entity. This service
// is built around this model.
//
// It must have a physical place where its headquarters are, denoted by fields Mesto, PostanskiBroj and  Sediste.
// swagger:model company
type Company struct {
	//JMBG of the person who is the owner of the company
	//Example: 1234567891234
	//Read Only: true
	Vlasnik string `json:"vlasnik" binding:"required,len=13"`
	// Unique number which identifies the company for taxes.
	// Required: true
	// Example: 15
	// Unique: true
	PIB int `json:"pib"`
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
	PostanskiBroj string `json:"postanskiBroj" binding:"required,number,max=20"`
	// Required: true
	// Example: EDUKACIJA
	Delatnost Delatnost `json:"delatnost" binding:"required"`
	// Required: true
	Sediste Nstj `json:"sediste"`
	// Password used for authentication
	// Required: true
	// Minimum length: 12
	// Maximum length: 72
	Password string `json:"password,omitempty" binding:"min=12,max=72,required"`
	// True if company is likvidirana
	Likvidirana bool `json:"likvidirana"`
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
