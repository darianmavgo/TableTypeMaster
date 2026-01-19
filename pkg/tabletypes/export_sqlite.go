package tabletypes

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// ExportToSQLite exports the table types data to a SQLite database file.
func ExportToSQLite(filename string) error {
	os.Remove(filename) // Remove if exists

	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return err
	}
	defer db.Close()

	createTableSQL := `CREATE TABLE table_types (
		name TEXT,
		file_ext TEXT,
		mimetype TEXT,
		magic_number TEXT,
		list_tables TEXT,
		list_columns TEXT,
		list_column_types TEXT
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare("INSERT INTO table_types(name, file_ext, mimetype, magic_number, list_tables, list_columns, list_column_types) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, tt := range AllTableTypes {
		_, err = stmt.Exec(tt.Name, tt.FileExt, tt.MimeType, tt.MagicNumber, tt.ListTables, tt.ListColumns, tt.ListColumnTypes)
		if err != nil {
			return err
		}
	}

	return nil
}
