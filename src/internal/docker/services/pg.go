package services

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type PGService struct {
	Con *pgx.Conn
}

// SanitizeIdentifier ensures the identifier is valid (no injection)
func SanitizeIdentifier(name string) string {
	return pgx.Identifier{name}.Sanitize()
}

func (pg *PGService) CreatePostgresUser(ctx context.Context, username, password string) error {
	sql := fmt.Sprintf(
		`CREATE USER %s WITH PASSWORD $1 NOSUPERUSER NOCREATEDB NOCREATEROLE NOINHERIT`,
		SanitizeIdentifier(username),
	)

	_, err := pg.Con.Exec(ctx, sql, password)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (pg *PGService) CreateDatabase(ctx context.Context, dbName string, owner string) error {
	// 1. Create the database
	createDBSQL := fmt.Sprintf("CREATE DATABASE %s OWNER %s",
		SanitizeIdentifier(dbName),
		SanitizeIdentifier(owner),
	)
	if _, err := pg.Con.Exec(ctx, createDBSQL); err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	// 2. Revoke public CONNECT privilege
	revokeSQL := fmt.Sprintf("REVOKE CONNECT ON DATABASE %s FROM PUBLIC",
		SanitizeIdentifier(dbName),
	)
	if _, err := pg.Con.Exec(ctx, revokeSQL); err != nil {
		return fmt.Errorf("failed to revoke public connect: %w", err)
	}

	return nil
}

func (pg *PGService) DropDatabase(ctx context.Context, dbName string) error {
	sql := fmt.Sprintf("DROP DATABASE IF EXISTS %s", SanitizeIdentifier(dbName))

	_, err := pg.Con.Exec(ctx, sql)
	if err != nil {
		return fmt.Errorf("failed to drop database: %w", err)
	}

	return nil
}

func (pg *PGService) DropPostgresUser(ctx context.Context, username string) error {
	sql := fmt.Sprintf("DROP USER IF EXISTS %s", SanitizeIdentifier(username))

	_, err := pg.Con.Exec(ctx, sql)
	if err != nil {
		return fmt.Errorf("failed to drop user: %w", err)
	}

	return nil
}
