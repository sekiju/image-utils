package resize

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/schollz/progressbar/v3"
	"github.com/sekiju/image_utils/utils"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func Run(s Settings) error {
	imagePaths, err := utils.GetImagesPaths(s.InputPath, s.IncludeSubDirectories)
	if err != nil {
		return err
	}

	if s.OutputPath == "" {
		s.OutputPath = s.InputPath
	}

	err = os.MkdirAll(s.OutputPath, 0o755)
	if err != nil {
		return fmt.Errorf("failed to create folder: %w", err)
	}

	imageChan := make(chan string)
	doneWaitGroup := &sync.WaitGroup{}
	startTime := time.Now()

	for i := 0; i < int(s.Threads); i++ {
		doneWaitGroup.Add(1)
		go worker(imageChan, doneWaitGroup, s)
	}

	bar := progressbar.Default(int64(len(imagePaths)), "Resizing...")
	for _, image := range imagePaths {
		imageChan <- image
		bar.Add(1)
	}

	bar.Close()
	close(imageChan)

	doneWaitGroup.Wait()

	endTime := time.Now()
	totalTime := endTime.Sub(startTime)
	fmt.Printf("Completed! Total time spent: %s\n", totalTime)

	return nil
}

func worker(imagePathChan <-chan string, waitGroup *sync.WaitGroup, s Settings) {
	for imagePath := range imagePathChan {
		executeOnFile(imagePath, s)
	}

	waitGroup.Done()
}

func executeOnFile(inputPath string, s Settings) {
	img, err := utils.ReadImageToStruct(inputPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	size := img.Bounds().Size()
	w := size.X
	h := size.Y

	if s.Percent > 0 {
		p := float32(s.Percent) / 100
		w = int(float32(w) * p)
		h = int(float32(h) * p)
	} else {
		if s.Width != 0 && s.Height != 0 {
			w = int(s.Width)
			h = int(s.Height)
		} else if s.Width != 0 {
			p := float32(s.Width) / float32(w)
			w = int(s.Width)
			h = int(float32(h) * p)
		} else if s.Height != 0 {
			p := float32(s.Height) / float32(h)
			h = int(s.Height)
			w = int(float32(w) * p)
		}
	}

	resizedImg := imaging.Resize(img, w, h, imaging.Box)
	outputPath := filepath.Join(s.OutputPath, filepath.Base(inputPath))

	err = utils.SavePNG(outputPath, resizedImg)
	if err != nil {
		fmt.Print(err)
	}
}
