package migrator

import (
	"database/sql"
	"fmt"
	"gameapp/repository/mysql"

	migrate "github.com/rubenv/sql-migrate"
)

type Migrator struct {
	dialect    string
	dbConfig   mysql.Config
	migrations *migrate.FileMigrationSource
	tableName  string
}

func New(dbConfig mysql.Config) Migrator {
	migrations := &migrate.FileMigrationSource{
		Dir: "./repository/mysql/migrations",
	}

	return Migrator{
		dbConfig:   dbConfig,
		dialect:    "mysql",
		migrations: migrations,
		tableName:  "migrations", // Default migration table name
	}
}

// SetMigrationTable sets a custom migration table name
func (m *Migrator) SetMigrationTable(tableName string) {
	m.tableName = tableName
}

func (m Migrator) getDB() (*sql.DB, error) {
	return sql.Open(m.dialect, fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true",
		m.dbConfig.Username, m.dbConfig.Password, m.dbConfig.Host, m.dbConfig.Port, m.dbConfig.DBName))
}

// Up applies migrations. If limit is 0, applies all pending migrations.
// If limit > 0, applies up to that many migrations.
func (m Migrator) Up(limit int) {
	db, err := m.getDB()
	if err != nil {
		panic(fmt.Errorf("can't open mysql db: %v", err))
	}
	defer db.Close()

	migrate.SetTable(m.tableName)

	var n int
	if limit == 0 {
		n, err = migrate.Exec(db, m.dialect, m.migrations, migrate.Up)
	} else {
		n, err = migrate.ExecMax(db, m.dialect, m.migrations, migrate.Up, limit)
	}

	if err != nil {
		panic(fmt.Errorf("can't apply migrations: %v", err))
	}
	fmt.Printf("Applied %d migrations!\n", n)
}

// Down rolls back migrations. If limit is 0, rolls back all migrations.
// If limit > 0, rolls back up to that many migrations.
func (m Migrator) Down(limit int) {
	db, err := m.getDB()
	if err != nil {
		panic(fmt.Errorf("can't open mysql db: %v", err))
	}
	defer db.Close()

	migrate.SetTable(m.tableName)

	var n int
	if limit == 0 {
		n, err = migrate.Exec(db, m.dialect, m.migrations, migrate.Down)
	} else {
		n, err = migrate.ExecMax(db, m.dialect, m.migrations, migrate.Down, limit)
	}

	if err != nil {
		panic(fmt.Errorf("can't rollback migrations: %v", err))
	}
	fmt.Printf("Rollbacked %d migrations!\n", n)
}

func (m Migrator) Status() {
	db, err := m.getDB()
	if err != nil {
		panic(fmt.Errorf("can't open mysql db: %v", err))
	}
	defer db.Close()

	migrate.SetTable(m.tableName)

	records, err := migrate.GetMigrationRecords(db, m.dialect)
	if err != nil {
		panic(fmt.Errorf("can't get migration records: %v", err))
	}

	migrations, err := m.migrations.FindMigrations()
	if err != nil {
		panic(fmt.Errorf("can't find migrations: %v", err))
	}

	fmt.Println("Migration Status:")
	fmt.Println("================")
	fmt.Printf("Table: %s\n\n", m.tableName)

	// Create a map of applied migrations for quick lookup
	appliedMap := make(map[string]bool)
	for _, record := range records {
		appliedMap[record.Id] = true
	}

	// Check each migration file
	for _, migration := range migrations {
		status := "[ ]"
		if appliedMap[migration.Id] {
			status = "[✓]"
		}
		fmt.Printf("%s %s\n", status, migration.Id)
	}

	fmt.Printf("\nTotal migrations: %d\n", len(migrations))
	fmt.Printf("Applied: %d\n", len(records))
	fmt.Printf("Pending: %d\n", len(migrations)-len(records))
}
