package convert

import (
	"github.com/spf13/cobra"
	"runtime"
)

var Command = &cobra.Command{
	Use: "convert",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := s.Validate(); err != nil {
			return err
		}

		return Run(s)
	},
}

var s = Settings{}

func init() {
	Command.Flags().UintVar(&s.Threads, "threads", uint(runtime.NumCPU()), "Number of simultaneously processed images. The default is equal to the number of CPU cores")
	Command.Flags().StringVarP(&s.InputPath, "input", "i", "", "Input folder with folders and images")
	Command.Flags().StringVarP(&s.OutputPath, "output", "o", "", "Output folder")
	Command.Flags().BoolVar(&s.IncludeSubDirectories, "include_sub_directories", false, "")
	Command.Flags().StringVarP(&s.Format, "format", "f", "png", "")
	Command.Flags().UintVarP(&s.Quality, "quality", "q", uint(90), "JPEG image quality")
}
