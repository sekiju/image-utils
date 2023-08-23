package resize

import (
	"github.com/spf13/cobra"
	"runtime"
)

var Command = &cobra.Command{
	Use: "resize",
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
	Command.Flags().StringVarP(&s.OutputPath, "output", "o", "", "Output folder. If not specified, the images will overwrite the originals")
	Command.Flags().BoolVar(&s.IncludeSubDirectories, "include_sub_directories", false, "")
	Command.Flags().UintVar(&s.Width, "width", 0, "")
	Command.Flags().UintVar(&s.Height, "height", 0, "")
	Command.Flags().UintVarP(&s.Percent, "percent", "p", 0, "")
	Command.Flags().StringVarP(&s.ResampleFilter, "filter", "f", "box", "")
}
