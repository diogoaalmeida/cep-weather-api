package entity

import "regexp"

var cepPattern = regexp.MustCompile(`^\d{8}$`)

func IsValidCEP(cep string) bool {
	return cepPattern.MatchString(cep)
}
