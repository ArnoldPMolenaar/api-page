package database

import (
	"api-page/main/src/enums"
	"api-page/main/src/models"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ArnoldPMolenaar/api-utils/database"
	"gorm.io/gorm"
)

var Pg *gorm.DB

// OpenDBConnection Start a new database connection.
// Also tries to migrate the database schema.
func OpenDBConnection() error {
	// Open connection to database.
	db, err := database.PostgresSQLConnection()
	if err != nil {
		return err
	}

	// Migrate the database schema.
	err = Migrate(db)
	if err != nil {
		return err
	}

	// Set the global DB variable.
	Pg = db

	return nil
}

// ReadinessCheck verifies that the database connection is initialized and reachable.
func ReadinessCheck() error {
	if Pg == nil {
		return errors.New("database connection is not initialized")
	}

	sqlDB, err := Pg.DB()
	if err != nil {
		return fmt.Errorf("database sql handle unavailable: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// MigrationReadinessCheck verifies that required tables and enum types exist.
func MigrationReadinessCheck() error {
	if Pg == nil {
		return errors.New("database connection is not initialized")
	}

	requiredTables := []any{
		// Core API entities.
		&models.App{},
		&models.PluginType{},
		&models.ModuleType{},
		&models.Version{},
		&models.FooterRow{},
		&models.FooterRowColumn{},
		&models.Menu{},
		&models.MenuItem{},
		&models.MenuItemIndexing{},
		&models.MenuItemRelation{},
		&models.Module{},
		&models.Page{},
		&models.PageIndexing{},
		&models.PagePartial{},
		&models.PagePartialRow{},
		&models.PagePartialRowColumn{},
	}

	for _, table := range requiredTables {
		if !Pg.Migrator().HasTable(table) {
			return fmt.Errorf("missing required table for %T", table)
		}
	}

	requiredJoinTables := []string{
		"app_module_types",
		"app_plugin_types",
		"footer_row_column_rows",
		"page_partial_row_column_rows",
	}
	for _, table := range requiredJoinTables {
		if !Pg.Migrator().HasTable(table) {
			return fmt.Errorf("missing required join table %s", table)
		}
	}

	indexingLabels, err := getEnumLabels("indexing")
	if err != nil {
		return fmt.Errorf("indexing enum check failed: %w", err)
	}

	expectedIndexingLabels := []string{
		enums.ALL.String(),
		enums.FOLLOW.String(),
		enums.INDEX.String(),
		enums.INDEXIFEMBEDDED.String(),
		enums.MAX_IMAGE_PREVIEW.String(),
		enums.MAX_SNIPPET.String(),
		enums.MAX_VIDEO_PREVIEW.String(),
		enums.NOAI.String(),
		enums.NOARCHIVE.String(),
		enums.NOCACHE.String(),
		enums.NOFOLLOW.String(),
		enums.NOIMAGEAI.String(),
		enums.NOIMAGEINDEX.String(),
		enums.NOINDEX.String(),
		enums.NOINDEXIFEMBEDDED.String(),
		enums.NONE.String(),
		"noodp",
		enums.NOSNIPPET.String(),
		enums.NOTRANSLATE.String(),
		"noydir",
		enums.UNAVAILABLE_AFTER.String(),
	}

	if err := validateEnumLabels("indexing", indexingLabels, expectedIndexingLabels); err != nil {
		return err
	}

	return nil
}

func validateEnumLabels(typeName string, actual, expected []string) error {
	if len(actual) != len(expected) {
		return fmt.Errorf("%s enum labels mismatch: have %v, want %v", typeName, actual, expected)
	}

	for i := range actual {
		if actual[i] != expected[i] {
			return fmt.Errorf("%s enum labels mismatch: have %v, want %v", typeName, actual, expected)
		}
	}

	return nil
}

func getEnumLabels(typeName string) ([]string, error) {
	rows, err := Pg.Raw(`
		SELECT e.enumlabel
		FROM pg_type t
		JOIN pg_enum e ON t.oid = e.enumtypid
		WHERE t.typname = ?
		ORDER BY e.enumsortorder
	`, typeName).Rows()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	labels := make([]string, 0)
	for rows.Next() {
		var label string
		if scanErr := rows.Scan(&label); scanErr != nil {
			return nil, scanErr
		}
		labels = append(labels, label)
	}

	if len(labels) == 0 {
		return nil, fmt.Errorf("enum type %q not found", typeName)
	}

	return labels, nil
}
