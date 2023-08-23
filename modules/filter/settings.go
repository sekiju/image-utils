package filter

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"runtime"
)

type Settings struct {
	InputPath             string `json:"input_path"`
	OutputPath            string `json:"output_path"`
	Threads               uint   `json:"threads"`
	MinimalWidth          uint   `json:"minimal_width"`
	MinimalHeight         uint   `json:"minimal_height"`
	IncludeSubDirectories bool   `json:"include_sub_directories"`
	IndexNaming           bool   `json:"index_naming"`
}

func (s Settings) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.InputPath, validation.Required),
		validation.Field(&s.OutputPath, validation.Required),
		validation.Field(&s.Threads, validation.Min(uint(1)), validation.Max(uint(runtime.NumCPU()))),
		validation.Field(&s.MinimalWidth, validation.Min(uint(1))),
		validation.Field(&s.MinimalHeight, validation.Min(uint(1))),
	)
}
