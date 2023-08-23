package convert

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"runtime"
)

type Settings struct {
	InputPath             string `json:"input_path"`
	OutputPath            string `json:"output_path"`
	Threads               uint   `json:"threads"`
	IncludeSubDirectories bool   `json:"include_sub_directories"`
	Format                string `json:"format"`
	Quality               uint   `json:"quality"`
}

func (s Settings) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.InputPath, validation.Required),
		validation.Field(&s.OutputPath, validation.Required),
		validation.Field(&s.Threads, validation.Min(uint(1)), validation.Max(uint(runtime.NumCPU()))),
		validation.Field(&s.Format, validation.In("png", "jpeg")),
		validation.Field(&s.Quality, validation.Min(uint(1)), validation.Max(uint(100))),
	)
}
