package resize

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/schollz/progressbar/v3"
	utils2 "github.com/sekiju/image_utils/internal/utils"
	"path/filepath"
	"sync"
	"time"
)

var AcceptedFilters = map[string]imaging.ResampleFilter{
	"nearest_neighbor":   imaging.NearestNeighbor,
	"box":                imaging.Box,
	"linear":             imaging.Linear,
	"hermite":            imaging.Hermite,
	"mitchell_netravali": imaging.MitchellNetravali,
	"catmull_rom":        imaging.CatmullRom,
	"bspline":            imaging.BSpline,
	"gaussian":           imaging.Gaussian,
	"bartlett":           imaging.Bartlett,
	"lanczos":            imaging.Lanczos,
	"hann":               imaging.Hann,
	"hamming":            imaging.Hamming,
	"blackman":           imaging.Blackman,
	"welch":              imaging.Welch,
	"cosine":             imaging.Cosine,
}

func Run(s Settings) error {
	imagePaths, err := utils2.GetImagesPaths(s.InputPath, s.IncludeSubDirectories)
	if err != nil {
		return err
	}

	if s.OutputPath == "" {
		s.OutputPath = s.InputPath
	}

	err = utils2.CreateDirectoryIfNotExists(s.OutputPath)
	if err != nil {
		return err
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
	img, err := utils2.ReadImageToStruct(inputPath)
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

	resizedImg := imaging.Resize(img, w, h, AcceptedFilters[s.ResampleFilter])
	outputPath := filepath.Join(s.OutputPath, filepath.Base(inputPath))

	err = utils2.SavePNG(outputPath, resizedImg)
	if err != nil {
		fmt.Print(err)
	}
}
