package config

import (
	"gopkg.in/yaml.v3"
	"io"
)

type DB struct {
	Driver string `yaml:"driver"`
	Dsn    string `yaml:"dsn"`
}

type UserContent struct {
	Mode     string `yaml:"mode"`
	Location string `yaml:"location"`
}

type API struct {
	GCS GCS
}

type GCS struct {
	JsonCredentialPath string `yaml:"json_credential_path"`
	BucketName         string `yaml:"bucket_name"`
	ProjectId          string `yaml:"project_id"`
}

type Server struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type Config struct {
	Db          DB          `yaml:"db"`
	Mode        string      `yaml:"mode"`
	Server      Server      `yaml:"server"`
	UserContent UserContent `yaml:"user_content"`
	API         API         `yaml:"api"`
	APM         APM         `yaml:"apm"`
}

type APM struct {
	Vendor          string `yaml:"vendor"`
	LicenseKey      string `yaml:"license_key"`
	ApplicationName string `yaml:"application_name"`
}

const DbDriverSqlite = "sqlite"
const DbDriverMysql = "mysql"

const ModeWeb = "web"
const ModeApi = "api"

const UserContentModeLocal = "local"
const UserContentModeS3 = "s3"
const UserContentModeGCS = "gcs"

func NewConfig(r io.Reader) (*Config, error) {

	cfg := &Config{}

	decoder := yaml.NewDecoder(r)
	err := decoder.Decode(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
