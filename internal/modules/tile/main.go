package tile

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/schollz/progressbar/v3"
	"github.com/sekiju/image_utils/internal/utils"
	"image"
	_ "image/jpeg"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func Run(s Settings) error {
	imagePaths, err := utils.GetImagesPaths(s.InputPath, s.IncludeSubDirectories)
	if err != nil {
		return err
	}

	isOutPathNotSetted := false
	if s.OutputPath == "" {
		s.OutputPath = s.InputPath
		isOutPathNotSetted = true
	}

	err = utils.CreateDirectoryIfNotExists(s.OutputPath)
	if err != nil {
		return err
	}

	startTime := time.Now()
	imageChan := make(chan string)
	wg := &sync.WaitGroup{}

	for i := 0; i < int(s.Threads); i++ {
		wg.Add(1)
		go worker(imageChan, wg, isOutPathNotSetted, s)
	}

	bar := progressbar.Default(int64(len(imagePaths)), "Tiling...")
	for _, imagePath := range imagePaths {
		imageChan <- imagePath
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

func worker(imagePathChan <-chan string, waitGroup *sync.WaitGroup, isOutPathNotSetted bool, s Settings) {
	for imagePath := range imagePathChan {
		executeOnFile(imagePath, s)
		if isOutPathNotSetted {
			err := os.Remove(imagePath)
			if err != nil {
				fmt.Println("Error deleting the original image:", err)
			}
		}
	}

	waitGroup.Done()
}

func executeOnFile(inputPath string, s Settings) {
	img, err := utils.ReadImageToStruct(inputPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	tileSize := int(s.Size)
	numTiles := int(s.CountPerImage)

	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	tilesPerRow := width / tileSize
	tilesPerCol := height / tileSize
	if numTiles == 0 {
		numTiles = tilesPerRow * tilesPerCol
	}

	baseName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	tilesSaved := 0

	for row := 0; row < tilesPerRow; row++ {
		for col := 0; col < tilesPerCol; col++ {
			if tilesSaved >= numTiles {
				return
			}

			x := col * tileSize
			y := row * tileSize

			if s.Minimum && (x+tileSize > width || y+tileSize > height) {
				continue
			}

			tile := imaging.Crop(img, image.Rect(x, y, x+tileSize, y+tileSize))

			if s.SkipBlackWhite && utils.IsFullyBlackOrWhiteImage(tile) {
				continue
			}

			tileFileName := fmt.Sprintf("%s_tile_%d_%d.png", baseName, row, col)
			tileFilePath := filepath.Join(s.OutputPath, tileFileName)

			err = utils.SavePNG(tileFilePath, tile)
			if err != nil {
				fmt.Print(err)
				return
			}

			tilesSaved++
		}
	}
}
