package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB interface {
	Connection() *gorm.DB
}

type SQLite struct {
	c *gorm.DB
}

func (d *SQLite) Connection() *gorm.DB {
	return d.c
}

func NewSQLite(sqliteLocation string) (DB, error){

	d := &SQLite{}

	var err error

	d.c, err = gorm.Open(sqlite.Open(sqliteLocation), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return d, nil
}

type MySQL struct {
	c *gorm.DB
}

func (d *MySQL) Connection() *gorm.DB {
	return d.c
}

func NewMySQL(dsn string) (DB, error){

	d := &MySQL{}

	var err error

	d.c, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return d, nil
}