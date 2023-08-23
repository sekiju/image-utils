package utils

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func CopyFile(source string, destination string) error {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, input)
	return err
}

var acceptedImageExtensions = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
}

func GetImagesPaths(folder string, includeSubDirectories bool) ([]string, error) {
	fileList := make([]string, 0)

	if includeSubDirectories {
		err := filepath.Walk(folder, func(root string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			fileExt := filepath.Ext(root)
			if _, ok := acceptedImageExtensions[fileExt]; !ok {
				return nil
			}

			fileList = append(fileList, root)
			return nil
		})

		if err != nil {
			return fileList, err
		}
	} else {
		files, err := ioutil.ReadDir(folder)
		if err != nil {
			return fileList, err
		}

		for _, file := range files {
			if file.Mode().IsRegular() {
				fileExt := filepath.Ext(file.Name())
				if _, ok := acceptedImageExtensions[fileExt]; !ok {
					continue
				}

				fileList = append(fileList, filepath.Join(folder, file.Name()))
			}
		}
	}

	return fileList, nil
}

func FileNameWithoutExt(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}

func CreateDirectoryIfNotExists(p string) error {
	if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(p, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create folder: %w", err)
		}
	}

	return nil
}
