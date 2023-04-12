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

type Company struct {
	Vlasnik       string    `json:"vlasnik"`
	PIB           int       `json:"pib" binding:"required"`
	MaticniBroj   string    `json:"maticniBroj" binding:"len=8,number"`
	Naziv         string    `json:"naziv" binding:"min=1,max=100"`
	AdresaSedista string    `json:"adresaSedista" binding:"min=1,max=100"`
	Mesto         string    `json:"mesto" binding:"min=1,max=100"`
	PostanskiBroj string    `json:"postanskiBroj" binding:"number,max=20"`
	Delatnost     Delatnost `json:"delatnost" binding:"required"`
	Sediste       Nstj      `json:"sediste"`
}

type Nstj struct {
	Oznaka string `json:"oznaka"`
	Naziv  string `json:"naziv,omitempty"`
}

type CompanyFilter struct {
	OrderBy   string
	Asc       bool
	Page      int
	Mesto     string
	Sediste   string
	Delatnost string
}
