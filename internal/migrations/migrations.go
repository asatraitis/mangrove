// DO NOT MANUALLY EDIT THIS FILE
// See tools/migration_gen.js for the script used to update this file
package migrations

func (m *Migrator) setMigrations() {
	m.migrations = []Migration{
		Newinitial_20241203101104(),
		// Add new migrations above this line
	}
}
