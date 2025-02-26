import fs from 'fs';

generateMigrationFile();


function generateMigrationFile() {
    if (process.argv[2] === undefined || process.argv[2] === '' || process.argv[2] === null) {
        console.error('Please provide a name for the migration file');
        console.error('Example: task migration-new -- UsersTable');
        return;
    }
    const version = getDate();
    const name = `${process.argv[2]}_${version}`;

    // console.log(createMigrationString(name, version));
    fs.writeFile(`../internal/migrations/${name}.go`, createMigrationString(name, version), (err) => {
        if (err) {
        console.error(err);
        return;
        }
        console.log('Migration file created successfully!');
    });

    // add migration to migrations.go
    addToMigrationList(createMigrationListItemString(name));
}

// returns date in format YYYYMMDDHHMMSS
function getDate() {
    const now = new Date();
    const year = now.getFullYear();

    // getMonth() returns 0-11, so add 1 and pad with '0' if needed
    const month = String(now.getMonth() + 1).padStart(2, '0');

    const day = String(now.getDate()).padStart(2, '0');
    const hour = String(now.getHours()).padStart(2, '0');
    const minute = String(now.getMinutes()).padStart(2, '0');
    const second = String(now.getSeconds()).padStart(2, '0');

    return `${year}${month}${day}${hour}${minute}${second}`;
}

function addToMigrationList(migrationName) {
    // read the migrations.go file
    fs.readFile('../internal/migrations/migrations.go', 'utf8', (err, data) => {
        if (err) {
        console.error(err);
        return;
        }

        // append the new migration to the list
        const migrationList = data.split('\n');
        // Find the comment in go file and add a new line there
        const commentLineNumber = migrationList.findIndex((line) => line.includes('// Add new migrations above this line'));
        migrationList.splice(commentLineNumber, 0, migrationName);
        

        // write the new migration list back to the file
        fs.writeFile('../internal/migrations/migrations.go', migrationList.join('\n'), (err) => {
        if (err) {
            console.error(err);
            return;
        }
        console.log('Migration added to migrations.go successfully!');
        });
    });
}

function createMigrationString(name, version) {
    return ` // Migration generated by tools/migration_gen.js
package migrations

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type ${name} struct {
	version int
}

func New${name}() Migration {
	return &${name}{
		version: ${version},
	}
}

func (m *${name}) Version() int {
	return m.version
}

func (m *${name}) Up(tx pgx.Tx) error {
	_, err := tx.Exec(context.Background(), \`
		-- SQL migration
	\`)
	return err
}
func (m *${name}) Down(tx pgx.Tx) error {
	_, err := tx.Exec(context.Background(), \`
        -- SQL migration revert
    \`)
	return err
}`
}
function createMigrationListItemString(name) {
    return `New${name}(),`
}