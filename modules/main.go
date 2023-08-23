package modules

import (
	"github.com/sekiju/image_utils/modules/convert"
	"github.com/sekiju/image_utils/modules/filter"
	"github.com/sekiju/image_utils/modules/resize"
	"github.com/sekiju/image_utils/modules/tile"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use: "image_utils",
}

func init() {
	root.DisableSuggestions = true
	root.CompletionOptions.DisableDefaultCmd = true

	root.AddCommand(tile.Command)
	root.AddCommand(resize.Command)
	root.AddCommand(filter.Command)
	root.AddCommand(convert.Command)
}

func Execute() {
	root.Execute()
}
