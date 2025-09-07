package repositories_test

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var testDB *sql.DB

var (
	dbUser     = "docker"
	dbPassword = "docker"
	dbDatabase = "sampledb"
	dbConn     = fmt.Sprintf("%s:%s@tcp(127.0.0.1:3307)/%s?parseTime=true", dbUser, dbPassword, dbDatabase)
)

func connectDB() error {
	var err error
	testDB, err = sql.Open("mysql", dbConn)
	if err != nil {
		return err
	}
	return nil
}

func setupTestData() error {
	content, err := os.ReadFile("./testdata/setupDB.sql")
	if err != nil {
		return err
	}
	statements := strings.Split(string(content), ";")
	for _, stmt := range statements {
		if strings.TrimSpace(stmt) == "" {
			continue
		}
		_, err := testDB.Exec(stmt)
		if err != nil {
			return fmt.Errorf("failed to execute statement: %s\nError: %w", stmt, err)
		}
	}
	return nil
}

func cleanupDB() error {
	content, err := os.ReadFile("./testdata/cleanupDB.sql")
	if err != nil {
		return err
	}
	statements := strings.Split(string(content), ";")
	for _, stmt := range statements {
		if strings.TrimSpace(stmt) == "" {
			continue
		}
		_, err := testDB.Exec(stmt)
		if err != nil {
			return fmt.Errorf("failed to execute statement: %s\nError: %w", stmt, err)
		}
	}
	return nil
}

// 全テスト共通の前処理を書く
func setup() error {
	if err := connectDB(); err != nil {
		return err
	}
	if err := cleanupDB(); err != nil {
		fmt.Println("cleanup", err)
		return err
	}
	if err := setupTestData(); err != nil {
		fmt.Println("setup")
		return err
	}
	return nil
}

// 前テスト共通の後処理を書く
func teardown() {
	cleanupDB()
	testDB.Close()
}

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	m.Run()

	teardown()
}
