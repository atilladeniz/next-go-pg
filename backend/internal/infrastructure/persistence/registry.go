package persistence

// AllEntities returns the GORM-tagged models that AutoMigrate must
// process. Replaces the former internal/domain/registry.go — domain
// types no longer carry persistence concerns.
//
// New entities are added by appending their GORM model here.
func AllEntities() []any {
	return []any{
		&gormUserStats{},
		// AUTO-GENERATED: new GORM models will be added above this line
	}
}
