package resize

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"runtime"
)

type Settings struct {
	InputPath             string `json:"input_path"`
	OutputPath            string `json:"output_path"`
	Threads               uint   `json:"threads"`
	IncludeSubDirectories bool   `json:"include_sub_directories"`
	Width                 uint   `json:"width"`
	Height                uint   `json:"height"`
	Percent               uint   `json:"percent"`
	ResampleFilter        string `json:"resample_filter"`
}

func (s Settings) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.InputPath, validation.Required),
		validation.Field(&s.Threads, validation.Min(uint(1)), validation.Max(uint(runtime.NumCPU()))),
		validation.Field(&s.Width, validation.Min(uint(1))),
		validation.Field(&s.Height, validation.Min(uint(1))),
		validation.Field(&s.Percent, validation.Min(uint(1)), validation.Max(uint(100))),
		validation.Field(&s.ResampleFilter, validation.In("nearest_neighbor", "box", "linear", "hermite", "mitchell_netravali", "catmull_rom", "bspline", "gaussian", "bartlett", "lanczos", "hann", "hamming", "blackman", "welch", "cosine")),
	)
}
