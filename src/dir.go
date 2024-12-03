package main

import (
	"io"
	"os"
	"path/filepath"
)

func moveDir(src string, dest string) error {
	// Копируем содержимое папки
	err := copyDir(src, dest)
	if err != nil {
		return err
	}

	// Удаляем исходную папку
	return os.RemoveAll(src)
}

func copyDir(src string, dest string) error {
	// Получаем информацию о исходной папке
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Создаем целевую папку
	err = os.MkdirAll(dest, srcInfo.Mode())
	if err != nil {
		return err
	}

	// Проходим по всем файлам и папкам в исходной папке
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			// Если это папка, рекурсивно копируем ее
			err = copyDir(srcPath, destPath)
			if err != nil {
				return err
			}
		} else {
			// Если это файл, копируем его
			err = copyFile(srcPath, destPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src string, dest string) error {
	// Открываем исходный файл
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// Создаем целевой файл
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	// Копируем содержимое файла
	_, err = io.Copy(out, in)
	return err
}

func isDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
