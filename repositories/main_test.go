package repositories_test

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"testing"
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
	file, err := os.Open("./testdata/setupDB.sql")
	if err != nil {
		return err
	}
	defer file.Close()

	cmd := exec.Command("mysql", "-h", "127.0.0.1", "-P", "3307", "-u", "docker", "sampledb", "--password=docker")
	cmd.Stdin = file

	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func cleanupDB() error {
	file, err := os.Open("./testdata/cleanupDB.sql")
	if err != nil {
		return err
	}
	defer file.Close()

	cmd := exec.Command("mysql", "-h", "127.0.0.1", "-P", "3307", "-u", "docker", "sampledb", "--password=docker")
	cmd.Stdin = file

	err = cmd.Run()
	if err != nil {
		return err
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
		os.Exit(1)
	}

	m.Run()

	teardown()
}
