package banco

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" //Driver de conaxão com o MySql
)

func Conectar() (*sql.DB, error) {
	stringConaxao := "golang:golang@/devbook?charset=utf8&&parseTime=True&loc=Local"
	db, erro := sql.Open("mysql", stringConaxao)

	if erro != nil {
		return nil, erro
	}

	if erro = db.Ping(); erro != nil {
		return nil, erro
	}

	return db, nil

}
