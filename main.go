package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nemdull/myapi/controllers"
	"github.com/nemdull/myapi/routers"
	"github.com/nemdull/myapi/services"
	// mysql driver
)

const (
	defaultDBUser     = "docker"
	defaultDBPassword = "docker"
	defaultDBName     = "sampledb"
	defaultDBHost     = "127.0.0.1"
	defaultDBPort     = "3307"
)

func main() {
	dbConn := buildDSN()

	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Printf("fail to create DB handle: %v", err)
		return
	}

	if err := db.Ping(); err != nil {
		log.Printf("fail to connect DB: %v", err)
		return
	}

	if err := ensureSchema(db); err != nil {
		log.Printf("fail to prepare schema: %v", err)
		return
	}

	con := controllers.NewMyAppController(services.NewMyAppService(db))
	r := routers.NewRouter(con)

	log.Println("server start at port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func buildDSN() string {
	user := envOrDefault("DB_USER", defaultDBUser)
	password := envOrDefault("DB_PASSWORD", defaultDBPassword)
	host := envOrDefault("DB_HOST", defaultDBHost)
	port := envOrDefault("DB_PORT", defaultDBPort)
	name := envOrDefault("DB_NAME", defaultDBName)

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, name)
}

func envOrDefault(key, def string) string {
	if val := strings.TrimSpace(os.Getenv(key)); val != "" {
		return val
	}

	return def
}

func ensureSchema(db *sql.DB) error {
	stmts := []string{
		`create table if not exists articles (
			article_id integer unsigned auto_increment primary key,
			title varchar(100) not null,
			contents text not null,
			username varchar(100) not null,
			nice integer not null default 0,
			created_at datetime default current_timestamp
		);`,
		`create table if not exists comments (
			comment_id integer unsigned auto_increment primary key,
			article_id integer unsigned not null,
			message text not null,
			created_at datetime default current_timestamp,
			foreign key (article_id) references articles(article_id)
		);`,
	}

	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}

	return nil
}
