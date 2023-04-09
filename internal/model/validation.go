package model

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

const (
	male   = "MALE"
	female = "FEMALE"
)

// type Delatnost int64

// // UnmarshalJSON implements json.Unmarshaler
// func (d Delatnost) UnmarshalJSON(b []byte) error {
// 	dStr := string(b)
// 	if len(b) == 1 {
// 		d = Delatnost(b[0])
// 		return nil
// 	}
// 	d, err := parseDelatnost(dStr[1 : len(dStr)-1])
// 	return err
// }

// func parseDelatnost(dStr string) (Delatnost, error) {
// 	switch dStr {
// 	case "ODBRANA":
// 		return Odbrana, nil
// 	case "EKONOMSKI_I_FINANSIJSKI_ODNOSI":
// 		return Ekonomski, nil
// 	case "EDUKACIJA":
// 		return Edukacija, nil
// 	case "OPSTE_JAVNE_USLUGE":
// 		return JavneUsluge, nil
// 	case "ZDRAVSTVO":
// 		return Zdravstvo, nil
// 	}
// 	return 0, fmt.Errorf("%w: %s", ErrInvalidDelatnost, dStr)
// }

var ErrInvalidDelatnost = errors.New("Invalid value for delatnost")

type Delatnost string

type delatnostRegistry struct {
	Odbrana     Delatnost
	Ekonomski   Delatnost
	Edukacija   Delatnost
	JavneUsluge Delatnost
	Zdravstvo   Delatnost
	delatnosti  []Delatnost
}

func newDelatnostRegistry() *delatnostRegistry {
	odbrana := Delatnost("ODBRANA")
	ekonomski := Delatnost("EKONOMSKI_I_FINANSIJSKI_ODNOSI")
	edukacija := Delatnost("EDUKACIJA")
	javneUsluge := Delatnost("OPSTE_JAVNE_USLUGE")
	zdravstvo := Delatnost("ZDRAVSTVO")

	return &delatnostRegistry{
		Odbrana:     odbrana,
		Ekonomski:   ekonomski,
		Edukacija:   edukacija,
		JavneUsluge: javneUsluge,
		Zdravstvo:   zdravstvo,
		delatnosti:  []Delatnost{odbrana, ekonomski, edukacija, javneUsluge, zdravstvo},
	}
}

func (delatnost Delatnost) String() string {
	return string(delatnost)
}

func (d *Delatnost) UnmarshalJSON(b []byte) error {
	dStr := string(b)
	tmpd, err := Delatnosti.Parse(dStr[1 : len(dStr)-1])
	*d = tmpd
	return err
}

func (delatnostRegistry delatnostRegistry) List() []Delatnost {
	return delatnostRegistry.delatnosti
}

func (delatnostRegistry delatnostRegistry) Parse(s string) (Delatnost, error) {
	for _, delatnost := range delatnostRegistry.List() {
		if delatnost.String() == s {
			return delatnost, nil
		}
	}
	return "", ErrInvalidDelatnost
}

var Delatnosti = newDelatnostRegistry()

func ValidateSex(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	if str != male && str != female {
		return false
	}
	return true
}
