package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(dsn string) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Gagal parse config: %v", err)
	}

	// PAKSA Postgres pakai WIB setiap kali koneksi baru dibuka
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		_, err := conn.Exec(ctx, "SET TIME ZONE 'Asia/Jakarta';")
		return err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("gagal konek ke postgres: %v", err)
	}

	return pool
}