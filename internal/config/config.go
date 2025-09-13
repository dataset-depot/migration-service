package config

import (
	"fmt"
	"log"
	"os"
	"time"
)

type HTTP struct {
	Addr					string
	ReadTimeout		time.Duration
	WriteTimeout	time.Duration
	MaxBodyBytes	int64
} 

type Database struct {
	UseSocket									bool
	InstanceConnectionName		string
	Host											string
	Port											string
	User											string
	Password									string
	Name											string
	MaxOpen										int
	MaxIdle										int
	IdleTime									time.Duration
	LifeTime									time.Duration
}

type Security struct {
	AdminToken	string
}

type CloudSQL struct {
	ProjectID		string
}

type Config struct {
	HTTP			HTTP
	Database	Database
	Security	Security
	CloudSQL	CloudSQL
}

func Load() Config {
	cfg := Config{
		HTTP: HTTP{
			Addr: 				getenv("HTTP_ADDR", ":8080"),
			ReadTimeout:	15 * time.Second,
			WriteTimeout: 10 * time.Minute,
			MaxBodyBytes: 32 << 20,
		},
		Database: Database{
			UseSocket:							os.Getenv("INSTANCE_CONNECTION_NAME") != "",
			InstanceConnectionName: os.Getenv("INSTANCE_CONNECTION_NAME"),
			Host:										getenv("DB_HOST", "127.0.0.1"),
			Port:										getenv("DB_PORT", "5432"),
			User:										must("DB_USER"),
			Password:								must("DB_PASS"),
			Name:										must("DB_NAME"),
			MaxOpen:								getint("DB_MAX_OPEN", 10),
			MaxIdle:								getint("DB_MAX_IDLE", 5),
			IdleTime:								5 * time.Minute,
			LifeTime:								1 * time.Hour,
		},
		Security: Security{
			AdminToken: must("ADMIN_TOKEN"),
		},
		CloudSQL: CloudSQL{
			ProjectID: getenv("GCP_PROJECT_ID", ""),
		},
	}
	return cfg
}

func must(k string) string {
	v := os.Getenv(k); if v == "" {
		log.Fatalf("missing env: %s", k)
	}
	return v
}

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" { return v }
	return d
}

func getint(k string, d int) int {
	if v := os.Getenv(k); v != "" {
		var n int
		_, _ = fmt.Sscanf(v, "%d", &n)
		if n > 0 { return n }
	}
	return d
}

