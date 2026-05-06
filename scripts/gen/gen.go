// Package gen generates GORM query interfaces and model structs from the
// live database schema using gorm.io/gen.
//
// Output layout:
//
//	./models/query        — query package (DAO entry point)
//	./models/query/model  — generated model structs
//
// The hand-written entities under ./models/entity remain authoritative for
// service code; gen output is for query-builder convenience and lives in a
// separate package to avoid type collisions.
package gen

import (
	"fmt"

	"gorm.io/gen"

	"github.com/gohive/core/config"
	"github.com/gohive/core/logger"
	dbmysql "github.com/gohive/pkg/mysql"
)

// Options for gen command.
type Options struct {
	ConfigPath   string
	OutPath      string // Query package directory (default: ./models/query)
	ModelPkgPath string // Model package directory (default: <OutPath>/model)
}

// Run generates query and model code based on the live DB schema.
func Run(opts Options) error {
	cfg, err := config.LoadConfig[config.BaseConfig](opts.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	logger.Init(cfg.Log.Level, cfg.Log.Format)

	db, err := dbmysql.OpenGorm(cfg.Database)
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}
	defer dbmysql.CloseGorm(db)

	g := gen.NewGenerator(gen.Config{
		OutPath:           opts.OutPath,
		ModelPkgPath:      opts.ModelPkgPath,
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable:     true,
		FieldCoverable:    false,
		FieldSignable:     false,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
	})

	g.UseDB(db)
	g.ApplyBasic(g.GenerateAllTable()...)
	g.Execute()

	logger.Info("Entity generation completed")
	return nil
}
