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
mode: api
server:
  port: 80
  host: 0.0.0.0
user_content:
  mode: local
  location: /app/user_content/
api:
  gcs:
    json_credential_path: "/app/gcs-creds.json"
    bucket_name: "bucket-name"
    project_id: "project-id"
`

	reader := strings.NewReader(content)

	actual, err := NewConfig(reader)

	assert.NoError(t, err)

	expected := Config{
		Db: DB{
			Driver: "sqlite",
			Dsn:    "/path/to/file.sqlite",
		},
		Mode: "api",
		Server: Server{
			Port: 80,
			Host: "0.0.0.0",
		},
		UserContent: UserContent{
			Mode:     "local",
			Location: "/app/user_content/",
		},
		API: API{
			GCS: GCS{
				BucketName:         "bucket-name",
				ProjectId:          "project-id",
				JsonCredentialPath: "/app/gcs-creds.json",
			},
		},
	}

	assert.Equal(t, expected, *actual)
}
