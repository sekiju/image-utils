package filter

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"github.com/sekiju/image_utils/utils"
	"path/filepath"
	"sync"
	"time"
)

type File struct {
	Path  string
	Index int
}

func newFile(p string, i int) *File {
	return &File{Path: p, Index: i}
}

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
	fileChan := make(chan *File)
	startTime := time.Now()

	pn := utils.NewPageName(len(files))
	for i := 0; i < int(s.Threads); i++ {
		wg.Add(1)
		go worker(fileChan, wg, pn, s)
	}

	bar := progressbar.Default(int64(len(files)), "Filtering...")
	for i, filePath := range files {
		fileChan <- newFile(filePath, i)
		bar.Add(1)
	}

	bar.Close()
	close(fileChan)
	wg.Wait()

	endTime := time.Now()
	totalTime := endTime.Sub(startTime)
	fmt.Printf("Completed! Total time spent: %s\n", totalTime)

	return nil
}

func worker(imagePathChan <-chan *File, wg *sync.WaitGroup, pg *utils.PageName, s Settings) {
	for imagePath := range imagePathChan {
		executeOnFile(imagePath, pg, s)
	}

	wg.Done()
}

func executeOnFile(file *File, pg *utils.PageName, s Settings) {
	img, err := utils.ReadImageToStruct(file.Path)
	if err != nil {
		fmt.Println(err)
		return
	}

	width, height := img.Bounds().Dx(), img.Bounds().Dy()

	if len(s.PrimarySide) > 0 {
		if s.PrimarySide == "width" && height > width {
			return
		} else if s.PrimarySide == "height" && width > height {
			return
		}
	}

	if s.MinimalWidth > 0 && width < int(s.MinimalWidth) {
		return
	}

	if s.MinimalHeight > 0 && height < int(s.MinimalHeight) {
		return
	}

	if s.MaximalWidth > 0 && width > int(s.MaximalWidth) {
		return
	}

	if s.MaximalHeight > 0 && height > int(s.MaximalHeight) {
		return
	}

	fileName := filepath.Base(file.Path)
	if s.IndexNaming {
		index := pg.GetName(file.Index)
		fileName = fmt.Sprintf("%s.png", index)
	}

	outputFile := filepath.Join(s.OutputPath, fileName)
	err = utils.CopyFile(file.Path, outputFile)
	if err != nil {
		fmt.Println(err)
	}
}
