package tile

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"runtime"
)

type Settings struct {
	InputPath             string `json:"input_path"`
	OutputPath            string `json:"output_path"`
	Threads               uint   `json:"threads"`
	IncludeSubDirectories bool   `json:"include_sub_directories"`
	Size                  uint   `json:"size"`
	Minimum               bool   `json:"minimum"`
	SkipBlackWhite        bool   `json:"skip_black_white"`
	CountPerImage         uint   `json:"count_per_image"`
}

// todo: Grayscale mode

func (s Settings) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.InputPath, validation.Required),
		validation.Field(&s.Threads, validation.Min(uint(1)), validation.Max(uint(runtime.NumCPU()))),
		validation.Field(&s.Size, validation.Min(uint(8))),
	)
}
