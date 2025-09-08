package db

import (
	"context"
	_ "embed"
	"fmt"
	"strings"
)

//go:embed migrations.sql
var migrations string

func (p *Postgres) Migrate(ctx context.Context) error {
	stmts := strings.Split(migrations, ";\n")
	for _, s := range stmts {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if _, err := p.SQL.ExecContext(ctx, s); err != nil {
			return fmt.Errorf("migrate exec failed: %w", err)
		}
	}
	return nil
}
