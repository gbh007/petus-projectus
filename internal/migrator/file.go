package migrator

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"log"
	"sort"
	"strconv"
	"strings"
)

// migrationFile - файл с миграцией
type migrationFile struct {
	// Номер миграции
	Number int
	// Путь до файла
	Path string
	// Название файла
	Name string
}

// getFileList - возвращает список файлов миграций
func getFileList(ctx context.Context, migrationsDir fs.FS) ([]migrationFile, error) {
	migrationFileList, err := fs.ReadDir(migrationsDir, ".")
	if err != nil {
		return nil, err
	}

	list := []migrationFile{}

	for _, migrationFileInfo := range migrationFileList {
		filename := migrationFileInfo.Name()

		// Миграции только в sql файлах
		if !strings.HasSuffix(filename, ".sql") {
			log.Println("Не sql, пропускаю", filename)

			continue
		}

		// Первые 4 символа обозначают номер миграции, если не получилось их обработать то игнорируем файл
		number, err := strconv.Atoi(filename[:4])
		if err != nil {
			log.Println(err)

			continue
		}

		list = append(list, migrationFile{
			Number: number,
			Path:   filename,
			Name:   filename,
		})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Number < list[j].Number
	})

	return list, nil
}

// migrationFromFile - получает данные для применения миграции из файла
func migrationFromFile(ctx context.Context, info migrationFile, migrationsDir fs.FS) (string, string, error) {
	file, err := migrationsDir.Open(info.Path)
	if err != nil {
		return "", "", err
	}

	defer logIfErrFunc(file.Close)

	data, err := io.ReadAll(file)
	if err != nil {
		return "", "", err
	}

	hash := fmt.Sprintf("%x", md5.Sum(data))
	body := string(data)

	return body, hash, nil
}
