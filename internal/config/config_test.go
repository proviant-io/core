package config

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNewConfig(t *testing.T) {

	content := `
---
db:
  driver: sqlite
  dsn: /path/to/file.sqlite
`

	reader := strings.NewReader(content)

	actual, err := NewConfig(reader)

	assert.NoError(t, err)

	expected := Config{
		Db: DB{
			Driver: "sqlite",
			Dsn:    "/path/to/file.sqlite",
		},
	}

	assert.Equal(t, expected, *actual)
}
