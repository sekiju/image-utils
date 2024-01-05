package convert

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"github.com/sekiju/image_utils/internal/utils"
	"image"
	"image/color"
	"image/jpeg"
	"path/filepath"
	"sync"
	"time"
)

func Run(s Settings) error {
	files, err := utils.GetImagesPaths(s.InputPath, s.IncludeSubDirectories)
	if err != nil {
		return err
	}

	err = utils.CreateDirectoryIfNotExists(s.OutputPath)
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
	img, err := utils.ReadImageToStruct(inputPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	fileName := utils.FileNameWithoutExt(filepath.Base(inputPath))

	if s.Mode != "original" {
		currentMode := getImageMode(img)
		if currentMode == s.Mode {
			fmt.Printf("Skipping conversion for %s. File is already in %s mode.\n", fileName, s.Mode)
			return
		}

		img = convertImageMode(img, s.Mode)
	}

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

func convertImageMode(img image.Image, mode string) image.Image {
	switch mode {
	case "gray":
		return convertToGrayScale(img)
	case "rgb":
		return convertToRGB(img)
	default:
		return img
	}
}

func getImageMode(img image.Image) string {
	switch img.(type) {
	case *image.Gray, *image.Gray16:
		return "gray"
	case *image.RGBA, *image.NRGBA:
		return "rgb"
	default:
		return "unknown"
	}
}

func convertToGrayScale(img image.Image) image.Image {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			grayColor := color.GrayModel.Convert(img.At(x, y))
			grayImg.Set(x, y, grayColor)
		}
	}

	return grayImg
}

func convertToRGB(img image.Image) image.Image {
	bounds := img.Bounds()
	rgbImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgbImg.Set(x, y, img.At(x, y))
		}
	}

	return rgbImg
}
