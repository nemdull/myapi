package services

import "database/sql"

type MyAppService struct {
	DB *sql.DB
}

func NewMyAppService(db *sql.DB) *MyAppService {
	return &MyAppService{DB: db}
}
