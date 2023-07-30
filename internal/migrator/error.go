package migrator

import (
	"errors"
	"log"
)

var (
	// Ошибка миграций БД
	MigratorError = errors.New("migrator")
	// Неизвестный диалект
	UnknownDialect = errors.New("unknown dialect")
)

func logIfErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func logIfErrFunc(f func() error) {
	logIfErr(f())
}
