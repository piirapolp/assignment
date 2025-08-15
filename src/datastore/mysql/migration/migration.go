package migration

import (
	"fmt"
	"net/url"
	"sort"

	"assignment/global"
	"assignment/logger"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Migration struct {
	Number uint `gorm:"primary_key"`
	Name   string

	Forwards func(db *gorm.DB) error `gorm:"-"`
}

var Migrations []*Migration

func Migrate(dryRun bool, number int, forceMigrate bool, isTest bool) error {
	dbHost := viper.GetString("Database.Host")
	dbPort := viper.GetString("Database.Port")
	dbUser := viper.GetString("Database.Username")
	dbPassword := viper.GetString("Database.Password")
	dbName := viper.GetString("Database.DatabaseName")

	if dryRun {
		logger.Logger.Infof("=== DRY RUN ===")
	}

	// check for duplicate migration Number
	migrationIDs := make(map[uint]struct{})
	for _, migration := range Migrations {
		if _, ok := migrationIDs[migration.Number]; ok {
			err := fmt.Errorf("duplicate migration Number found: %d", migration.Number)
			logger.Logger.Errorf("Unable to apply migrations, err: %+v", err)
			return err
		}

		migrationIDs[migration.Number] = struct{}{}
	}

	sort.Slice(Migrations, func(i, j int) bool {
		return Migrations[i].Number < Migrations[j].Number
	})

	connStr := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=%s&multiStatements=true",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
		url.QueryEscape(global.TimeZone),
	)

	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		logger.Logger.Errorf("unable to connect db: %+v", err)
		return errors.Wrap(err, "unable to connect db")
	}

	// Force Migrate Zone
	if forceMigrate {
		logger.Logger.Infof("=== FORCE MIGRATE ===")
		if err := db.Migrator().DropTable(&Migration{}); err != nil {
			return errors.Wrap(err, "unable to drop migrations table")
		}
	}

	// Make sure Migration table is there
	logger.Logger.Debugf("ensuring migrations table is present")
	if err := db.AutoMigrate(&Migration{}); err != nil {
		return errors.Wrap(err, "unable to automatically migrate migrations table")
	}

	var latest Migration
	if err := db.Order("number desc").First(&latest).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrap(err, "unable to find latest migration")
	}

	noMigrationsApplied := latest.Number == 0

	if noMigrationsApplied && len(Migrations) == 0 {
		logger.Logger.Infof("no migrations to apply")
		return nil
	}

	if latest.Number >= Migrations[len(Migrations)-1].Number {
		logger.Logger.Infof("no migrations to apply")
		return nil
	}

	if number == -1 {
		number = int(Migrations[len(Migrations)-1].Number)
	}

	if uint(number) <= latest.Number && latest.Number > 0 {
		logger.Logger.Infof("no migrations to apply, specified number is less than or equal to latest migration; backwards migrations are not supported")
		return nil
	}

	for _, migration := range Migrations {
		if migration.Number > uint(number) {
			break
		}

		if migration.Number <= latest.Number {
			continue
		}

		if latest.Number > 0 {
			logger.Logger.Infof("continuing migration starting from %d", migration.Number)
		}

		migrationLogger := logger.Logger.With(
			"migration_number", migration.Number,
		)

		migrationLogger.Infof("applying migration %q", migration.Name)

		if dryRun {
			continue
		}

		tx := db.Begin()

		if err := migration.Forwards(tx); err != nil {
			logger.Logger.Errorf("unable to apply migration, rolling back. err: %+v", err)
			if err := tx.Rollback().Error; err != nil {
				logger.Logger.Errorf("unable to rollback... err: %+v", err)
			}
			break
		}

		if err := tx.Commit().Error; err != nil {
			logger.Logger.Errorf("unable to commit transaction... err: %+v", err)
			break
		}

		// Create migration record
		if err := db.Create(migration).Error; err != nil {
			logger.Logger.Errorf("unable to create migration record. err: %+v", err)
			break
		}
	}

	return nil
}
