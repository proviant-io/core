package i18n

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestLanguageChange(t *testing.T) {

	l := NewFileLocalizer()

	m := Message{
		Template: "product with id %d not found",
		Params:   []interface{}{1},
	}

	assert.Equal(t, "product with id 1 not found", l.T(m, En))
	assert.Equal(t, "продукт под номером 1 не найден",  l.T(m, Ru))
	assert.Equal(t, "product with id 1 not found", l.T(m, Es))

}
