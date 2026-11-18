package config

import (
	"context"
	"database/sql"
	"menghitung-tekanan/internal/logger"
	"net"
	"net/url"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sirupsen/logrus"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func Connect(cfg PostgresConfig) *sql.DB {
	log := logger.New()

	if cfg.Host == "" || cfg.Port == "" ||
		cfg.User == "" || cfg.Database == "" {

		log.WithFields(logrus.Fields{
			"operation": "config postgre",
			"layer":     "config",
		}).Fatal("konfigurasi database belum lengkap")
		return nil
	}

	if cfg.SSLMode == "" {
		cfg.SSLMode = "disable"
	}

	// Membuat PostgreSQL connection URL
	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   net.JoinHostPort(cfg.Host, cfg.Port),
		Path:   cfg.Database,
	}

	query := dsn.Query()
	query.Set("sslmode", cfg.SSLMode)
	dsn.RawQuery = query.Encode()

	// Membuka koneksi database
	database, err := sql.Open("pgx", dsn.String())
	if err != nil {
		log.WithFields(logrus.Fields{
			"operation": "config postgre",
			"layer":     "config",
		}).WithError(err).Fatal("gagal membuka database")
		return nil
	}

	// Pengaturan connection pool
	database.SetMaxOpenConns(10)
	database.SetMaxIdleConns(5)
	database.SetConnMaxLifetime(time.Hour)
	database.SetConnMaxIdleTime(10 * time.Minute)

	// Memeriksa koneksi
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := database.PingContext(ctx); err != nil {
		database.Close()
		log.WithFields(logrus.Fields{
			"operation": "config postgre",
			"layer":     "config",
		}).WithError(err).Fatal("gagal menghubungkan database")
		return nil
	}

	return database
}
