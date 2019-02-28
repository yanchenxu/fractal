package sqlite

import (
	"strings"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"database/sql"
)

type SqliteDB struct {
	dbfile string
	db     *sql.DB
}

func New(dbfile string) (*SqliteDB, error) {
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(CREATE_SQL)
	if err != nil {
		db.Close()
		return nil, err
	}
	return &SqliteDB{
		dbfile: dbfile,
		db:     db,
	}, nil
}

func (db *SqliteDB) Insert(table string, field []string, values [][]interface{}) error {
	vs := strings.Repeat("?,", len(field))
	vs = vs[0:len(vs)-1]
	vfmt := " values(" + vs + ")"
	sql := fmt.Sprintf("INSERT INTO %s(%s) %s", table, strings.Join(field, ","), vfmt)
	value := make([]interface{}, len(values))
	for _, v := range values {
		value = append(value, v)
	}
	stmt, err := db.db.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(value...)
	if err != nil {
		return err
	}
	return nil
}

func (db *SqliteDB) Delete(table string, limit map[string]interface{}) error {
	var limits []string
	var value []interface{}
	for k, v := range limit {
		limits = append(limits, fmt.Sprintf("%s=?", k))
		value = append(value, v)
	}
	sql := fmt.Sprintf("delete from %s where %s", table, strings.Join(limits, ","))
	_, err := db.db.Exec(sql, value)
	if err != nil {
		return err
	}
	return nil
}

func (db *SqliteDB) Query(table string, fields []string, limit map[string]interface{}) ([][]interface{}, error) {
	var limits []string
	var value []interface{}
	for k, v := range limit {
		limits = append(limits, fmt.Sprintf("%s=?", k))
		value = append(value, v)
	}
	sql := fmt.Sprintf("SELECT %s FROM %s where %s", table, strings.Join(fields, ","),
		strings.Join(limits, " and "))
	rows, err := db.db.Query(sql, value)
	if err != nil {
		return nil, err
	}
	var results [][]interface{}
	for rows.Next() {
		keys, _ := rows.Columns()
		length := len(keys)
		var result []interface{}
		for i := 0; i < length; i++ {
			var tmp interface{}
			result = append(result, tmp)
		}
		err = rows.Scan(result...)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (db *SqliteDB) Close() error {
	return db.Close()
}
