package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Pg struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

func (p *Pg) GetNewDBConnection() *pgxpool.Pool {
	url := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", p.User, p.Password, p.Host, p.Port, p.Name)

	// Set max parallel connections to 20
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatal("Failed to parse pg config")
	}

	config.MaxConns = 20

	dbPool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Fatal("Cannot create database pool\n")
	}

	return dbPool
}
