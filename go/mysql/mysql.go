//
// mysql.go
//
package mysql

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)


const (
    DriverName = "mysql"
)


var (
    myDB *sql.DB
)


//
// Close database.
//
func Close() {
    if myDB != nil {
        myDB.Close()
    }
}


//
// Connect to the database.
//
func Conn() *sql.DB {
    return myDB
}


//
// Initialize with data source name.
//
func InitWithDataSourceName(dataSourceName string) error {
    db, err := sql.Open(DriverName, dataSourceName)
    if err != nil {
        return err
    }

    if err := db.Ping(); err != nil {
        db.Close()
        return err
    }

    myDB = db

    return nil
}
