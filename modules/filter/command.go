package filter

import (
	"github.com/spf13/cobra"
	"runtime"
)

var Command = &cobra.Command{
	Use: "filter",
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
	Command.Flags().UintVar(&s.MinimalWidth, "minimal_width", uint(0), "")
	Command.Flags().UintVar(&s.MinimalHeight, "minimal_height", uint(0), "")
	Command.Flags().UintVar(&s.MaximalWidth, "maximal_width", uint(0), "")
	Command.Flags().UintVar(&s.MaximalHeight, "maximal_height", uint(0), "")
	Command.Flags().BoolVar(&s.IncludeSubDirectories, "include_sub_directories", false, "")
	Command.Flags().BoolVar(&s.IndexNaming, "index_naming", false, "Files will be renamed to index: 001.png ... 013.png")
	Command.Flags().StringVar(&s.PrimarySide, "primary_side", "", "Skip if the opposite side is larger than the main side")
}
