package database

import (
	"api-page/main/src/models"

	"gorm.io/gorm"
)

// Migrate the database schema.
// See: https://gorm.io/docs/migration.html#Auto-Migration
func Migrate(db *gorm.DB) error {
	// Adds the indexing enum type to the database.
	if tx := db.Exec(`DO $$ 
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'indexing') THEN 
			CREATE TYPE indexing AS ENUM ('all', 'follow', 'index', 'indexifembedded', 'max-image-preview', 'max-snippet', 'max-video-preview', 'noai', 'noarchive', 'nocache', 'nofollow', 'noimageai', 'noimageindex', 'noindex', 'noindexifembedded', 'none', 'noodp', 'nosnippet', 'notranslate', 'noydir', 'unavailable_after'); 
		END IF; 
	END $$;`); tx.Error != nil {
		return tx.Error
	}

	// Updated migration set: normalized models.
	err := db.AutoMigrate(
		&models.App{},
		&models.PluginType{},
		&models.ModuleType{},
		&models.Version{},
		&models.Menu{},
		&models.MenuItem{},
		&models.MenuItemIndexing{},
		&models.MenuItemRelation{},
		&models.Module{},
		&models.Page{},
		&models.PageIndexing{},
		&models.PagePartial{},
		&models.PagePartialRow{},
		&models.PagePartialRowColumn{})
	if err != nil {
		return err
	}

	return nil
}
