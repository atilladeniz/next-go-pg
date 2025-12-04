package domain

// AllEntities returns all domain entities for GORM AutoMigrate.
// New entities are automatically added by 'make goca-feature'.
func AllEntities() []any {
	return []any{
		&UserStats{},
		// AUTO-GENERATED: New entities will be added above this line
	}
}
