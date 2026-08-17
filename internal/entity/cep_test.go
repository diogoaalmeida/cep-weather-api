package entity_test

import (
	"testing"

	"github.com/diogoaalmeida/cep-weather-api/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestIsValidCEP(t *testing.T) {
	cases := map[string]bool{
		"01310100":  true,
		"1310100":   false,
		"013101000": false,
		"0131010a":  false,
		"":          false,
		"01310-100": false,
	}

	for cep, want := range cases {
		assert.Equal(t, want, entity.IsValidCEP(cep), "cep=%q", cep)
	}
}
