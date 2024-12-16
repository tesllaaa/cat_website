package util

import (
	"log"
	"os"
)

// CreateDirectory Создание директорий для статей и временного сохранения файлов
func CreateDirectory() {
	dirName := "tmp"
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err := os.Mkdir(dirName, 0755)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	dirName = "articles"
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err := os.Mkdir(dirName, 0755)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
