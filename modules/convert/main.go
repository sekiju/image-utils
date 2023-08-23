package convert

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"github.com/sekiju/image_utils/utils"
	"image/jpeg"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func Run(s Settings) error {
	files, err := utils.GetImagesPaths(s.InputPath, false)
	if err != nil {
		return err
	}

	if _, err := os.Stat(s.OutputPath); os.IsNotExist(err) {
		err := os.Mkdir(s.OutputPath, 0755)
		if err != nil {
			fmt.Println(err)
		}
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
	img, err := utils.ReadImageToStruct(inputPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	fileName := utils.FileNameWithoutExt(filepath.Base(inputPath))

	if s.Format == "jpeg" {
		outputPath := filepath.Join(s.OutputPath, fmt.Sprintf("%s.jpeg", fileName))

		err = utils.SaveJPEG(outputPath, img, &jpeg.Options{
			Quality: int(s.Quality),
		})
		if err != nil {
			fmt.Print(err)
		}
	} else {
		outputPath := filepath.Join(s.OutputPath, fmt.Sprintf("%s.png", fileName))
		err = utils.SavePNG(outputPath, img)
		if err != nil {
			fmt.Print(err)
		}
	}
}
