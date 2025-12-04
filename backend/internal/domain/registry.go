package domain

// AllEntities returns all domain entities for GORM AutoMigrate.
// When you add a new entity with goca, add it here too.
// This is the ONLY place you need to register new entities.
func AllEntities() []interface{} {
	return []interface{}{
		&UserStats{},
		// Add new entities here:
		// &Product{},
		// &Order{},
		// &Invoice{},
	}
}
