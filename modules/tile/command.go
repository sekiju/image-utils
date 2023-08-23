package tile

import (
	"github.com/spf13/cobra"
	"runtime"
)

var Command = &cobra.Command{
	Use: "tile",
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
	Command.Flags().UintVarP(&s.Size, "tile", "t", uint(128), "Size of tile")
	Command.Flags().BoolVarP(&s.Minimum, "minimum", "m", false, "Sets a minimum size for tiles. If any tiles are below the specified size, they will not be saved")
	Command.Flags().BoolVarP(&s.SkipBlackWhite, "skip_black_white", "s", false, "Skip black and white tiles")
	Command.Flags().UintVarP(&s.CountPerImage, "count_per_image", "n", uint(0), "The number of tiles to save per image")
}
