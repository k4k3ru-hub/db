//
// mysql.go
//
package mysql

import (
    "database/sql"
    mysqlDriver "github.com/go-sql-driver/mysql"
)


const (
    DriverName = "mysql"
)


//
// Open with configuration.
//
func OpenWithConfig(cfg *mysqlDriver.Config) (*sql.DB, error) {
    connector, err := mysqlDriver.NewConnector(cfg)
    if err != nil {
        return nil, err
    }

    db := sql.OpenDB(connector)

    if err := ping(db); err != nil {
        db.Close()
    }

    return db, nil
}


//
// Open with data source name.
//
func OpenWithDSN(dsn string) (*sql.DB, error) {
    db, err := sql.Open(DriverName, dsn)
    if err != nil {
        return nil, err
    }

    if err := ping(db); err != nil {
        db.Close()
    }

    return db, nil
}


//
// Ping to the database.
//
func ping(db *sql.DB) error {
    if err := db.Ping(); err != nil {
        db.Close()
        return err
    }
    return nil
}
