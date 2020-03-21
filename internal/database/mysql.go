package database

import (
	"context"
	"database/sql"

	mysql "github.com/go-sql-driver/mysql"
)

type MySQLDB struct {
	dataSourceName string
	dbName         string
	Option         *DBOption
	Conn           *sql.DB
}

func init() {
	Register("mysql", func(dataSourceName, dbName string) Database {
		return &MySQLDB{
			dataSourceName: dataSourceName,
			dbName:         dbName,
			Option:         &DBOption{},
		}
	})
}

func (db *MySQLDB) Open() error {
	var conn *sql.DB

	var connString string
	if db.dbName != "" {
		cfg, err := mysql.ParseDSN(db.dataSourceName)
		if err != nil {
			return err
		}
		cfg.DBName = db.dbName
		connString = cfg.FormatDSN()
	} else {
		connString = db.dataSourceName
	}
	conn, err := sql.Open("mysql", connString)
	if err != nil {
		return err
	}

	conn.SetMaxIdleConns(DefaultMaxIdleConns)
	if db.Option.MaxIdleConns != 0 {
		conn.SetMaxIdleConns(db.Option.MaxIdleConns)
	}
	conn.SetMaxOpenConns(DefaultMaxOpenConns)
	if db.Option.MaxOpenConns != 0 {
		conn.SetMaxOpenConns(db.Option.MaxOpenConns)
	}
	db.Conn = conn
	return nil
}

func (db *MySQLDB) Close() error {
	return db.Conn.Close()
}

func (db *MySQLDB) Databases() ([]string, error) {
	rows, err := db.Conn.Query("SHOW DATABASES")
	if err != nil {
		return nil, err
	}
	databases := []string{}
	for rows.Next() {
		var database string
		if err := rows.Scan(&database); err != nil {
			return nil, err
		}
		databases = append(databases, database)
	}
	return databases, nil
}

func (db *MySQLDB) Tables() ([]string, error) {
	rows, err := db.Conn.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	tables := []string{}
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	return tables, nil
}

func (db *MySQLDB) DescribeTable(tableName string) ([]*ColumnDesc, error) {
	rows, err := db.Conn.Query("DESC " + tableName)
	if err != nil {
		return nil, err
	}
	tableInfos := []*ColumnDesc{}
	for rows.Next() {
		var tableInfo ColumnDesc
		err := rows.Scan(
			&tableInfo.Name,
			&tableInfo.Type,
			&tableInfo.Null,
			&tableInfo.Key,
			&tableInfo.Default,
			&tableInfo.Extra,
		)
		if err != nil {
			return nil, err
		}
		tableInfos = append(tableInfos, &tableInfo)
	}
	return tableInfos, nil
}

func (db *MySQLDB) Exec(ctx context.Context, query string) (sql.Result, error) {
	return db.Conn.ExecContext(ctx, query)
}

func (db *MySQLDB) Query(ctx context.Context, query string) (*sql.Rows, error) {
	return db.Conn.QueryContext(ctx, query)
}

func (db *MySQLDB) SwitchDB(dbName string) error {
	db.dbName = dbName
	return nil
}
