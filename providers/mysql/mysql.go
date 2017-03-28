package mysql

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	"os"

	"github.com/VividCortex/mysqlerr"
	"github.com/alioygur/fb-tinder-app/service"
	"github.com/alioygur/goutil"
)

const (
	usersTbl       = `users`
	reactionsTbl   = `reactions`
	creditsTbl     = `credits`
	friendshipsTbl = `friendships`
	matchesTbl     = `matches`
	abusesTbl      = `abuses`
	imagesTbl      = `images`
)

type (
	repository struct {
		db *sql.DB
		tx *sql.Tx
	}

	queryer interface {
		Exec(query string, args ...interface{}) (sql.Result, error)
		Query(query string, args ...interface{}) (*sql.Rows, error)
		QueryRow(query string, args ...interface{}) *sql.Row
	}
)

// New instances new sql repository
func New(db *sql.DB) service.SQLRepository {
	return &repository{db: db}
}

// Begin starts a transaction
func (r *repository) Begin() error {
	tx, err := r.db.Begin()
	if err != nil {
		return errors.WithStack(err)
	}
	r.tx = tx
	return nil
}

func (r *repository) Rollback() error {
	defer func() {
		r.tx = nil
	}()
	return errors.WithStack(r.tx.Rollback())
}

func (r *repository) Commit() error {
	defer func() {
		r.tx = nil
	}()
	return errors.WithStack(r.tx.Commit())
}

func (r *repository) sess() queryer {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *repository) insert(tbl string, v interface{}) (uint64, error) {
	ss, err := goutil.NewSQLStruct(v)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	sql := fmt.Sprintf(`insert into %s SET %s=?`, tbl, strings.Join(ss.Columns("id"), "=?, "))
	res, err := r.sess().Exec(sql, ss.Values("id")...)
	if err != nil {
		return 0, handleErr(err)
	}
	id, err := res.LastInsertId()
	return uint64(id), errors.WithStack(err)
}

func (r *repository) oneBy(v interface{}, tbl string, w string, args ...interface{}) error {
	ss, err := goutil.NewSQLStruct(v)
	if err != nil {
		return err
	}
	q := fmt.Sprintf("SELECT %s FROM %s WHERE %s", strings.Join(ss.Columns(), ","), tbl, w)
	err = r.sess().QueryRow(q, args...).Scan(ss.Ptrs()...)
	return handleErr(err)
}

func (r *repository) existsBy(tbl string, w string, args ...interface{}) (bool, error) {
	var exists bool
	q := fmt.Sprintf("SELECT COUNT(id) > 0 FROM %s WHERE %s", tbl, w)
	err := r.sess().QueryRow(q, args...).Scan(&exists)
	return exists, errors.WithStack(err)
}

func (r *repository) update(v interface{}, tbl string, w string, args ...interface{}) error {
	ss, err := goutil.NewSQLStruct(v)
	if err != nil {
		return err
	}
	q := fmt.Sprintf("UPDATE %s SET %s=? WHERE %s", tbl, strings.Join(ss.Columns("id"), "=?, "), w)

	args = append(ss.Values("id"), args...)
	_, err = r.db.Exec(q, args...)
	return handleErr(err)
}

func handleErr(err error) error {
	if err == sql.ErrNoRows {
		err = service.NewErr(service.NotFoundErrCode, err)
		return errors.WithStack(err)
	}

	if e, ok := err.(*mysqldriver.MySQLError); ok {
		switch int(e.Number) {
		case mysqlerr.ER_NO_REFERENCED_ROW_2: // foreign key errors
			return errors.WithStack(service.NewErr(service.NotFoundErrCode, err))
		case mysqlerr.ER_DUP_ENTRY:
			return errors.WithStack(service.NewErr(service.AlreadyExistsErrCode, err))
		}
	}

	return errors.WithStack(err)
}

func boolPtr(b bool) *bool {
	return &b
}

// ConnectToDB try connect to mysql server.
// It uses MYSQL_URL env, if it doesn't exists then uses fallback.
func ConnectToDB() (*sql.DB, error) {
	url, err := parseDSN(os.Getenv("MYSQL_URL"))
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := db.Ping(); err != nil {
		return nil, errors.WithStack(err)
	}
	return db, nil
}

// ApplyScheme applies db scheme
func ApplyScheme(db *sql.DB, scheme string) error {
	if scheme == "" {
		scheme = "./scheme.sql"
	}
	file, err := ioutil.ReadFile(scheme)
	if err != nil {
		return errors.WithStack(err)
	}
	queries := strings.Split(string(file), ";")

	tx, err := db.Begin()
	if err != nil {
		return errors.WithStack(err)
	}
	for _, q := range queries {
		// skip empty queries
		if q == "" {
			continue
		}
		_, err := tx.Exec(q)
		if err != nil {
			tx.Rollback()
			return errors.WithStack(err)
		}
	}
	return errors.WithStack(tx.Commit())
}

func parseDSN(dsn string) (string, error) {
	durl, err := url.Parse(dsn)
	if err != nil {
		return "", errors.WithStack(err)
	}
	user := durl.User.Username()
	password, _ := durl.User.Password()
	host := durl.Host
	if host == "" {
		host = "localhost:3306"
	}
	dbname := durl.Path // like: /path

	return fmt.Sprintf("%s:%s@tcp(%s)%s?charset=utf8&parseTime=True&loc=Local", user, password, host, dbname), nil
}

func tables(sess *sql.DB) ([]string, error) {
	rows, err := sess.Query(`SHOW TABLES`)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, errors.WithStack(err)
		}
		tables = append(tables, table)
	}
	return tables, nil
}

// TruncateTables truncates all the tables in db
func TruncateTables(sess *sql.DB) error {
	tables, err := tables(sess)
	if err != nil {
		return err
	}
	sess.Exec(`SET foreign_key_checks = ?`, 0)
	defer sess.Exec(`SET foreign_key_checks = ?`, 1)
	for _, t := range tables {
		if _, err := sess.Exec(`truncate ` + t); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
