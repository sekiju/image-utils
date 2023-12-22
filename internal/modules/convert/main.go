package convert

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	utils2 "github.com/sekiju/image_utils/internal/utils"
	"image/jpeg"
	"path/filepath"
	"sync"
	"time"
)

func Run(s Settings) error {
	files, err := utils2.GetImagesPaths(s.InputPath, s.IncludeSubDirectories)
	if err != nil {
		return err
	}

	err = utils2.CreateDirectoryIfNotExists(s.OutputPath)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	imageChan := make(chan string)
	startTime := time.Now()

	for i := 0; i < int(s.Threads); i++ {
		wg.Add(1)
		go worker(imageChan, wg, s)
	}

	bar := progressbar.Default(int64(len(files)), "Converting...")
	for _, filePath := range files {
		imageChan <- filePath
		bar.Add(1)
	}

	bar.Close()
	close(imageChan)
	wg.Wait()

	endTime := time.Now()
	totalTime := endTime.Sub(startTime)
	fmt.Printf("Completed! Total time spent: %s\n", totalTime)

	return nil
}

func worker(imagePathChan <-chan string, wg *sync.WaitGroup, s Settings) {
	for imagePath := range imagePathChan {
		executeOnFile(imagePath, s)
	}

	wg.Done()
}

func executeOnFile(inputPath string, s Settings) {
	img, err := utils2.ReadImageToStruct(inputPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	fileName := utils2.FileNameWithoutExt(filepath.Base(inputPath))

	if s.Format == "jpeg" {
		outputPath := filepath.Join(s.OutputPath, fmt.Sprintf("%s.jpeg", fileName))

		err = utils2.SaveJPEG(outputPath, img, &jpeg.Options{
			Quality: int(s.Quality),
		})
		if err != nil {
			fmt.Print(err)
		}
	} else {
		outputPath := filepath.Join(s.OutputPath, fmt.Sprintf("%s.png", fileName))
		err = utils2.SavePNG(outputPath, img)
		if err != nil {
			fmt.Print(err)
		}
	}
}
