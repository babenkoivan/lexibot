package testkit

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func MockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *gorm.DB) {
	conn, mock, err := sqlmock.New()
	require.NoError(t, err)

	db, err := gorm.Open(mysql.New(mysql.Config{Conn: conn, SkipInitializeWithVersion: true}))
	require.NoError(t, err)

	return conn, mock, db
}
