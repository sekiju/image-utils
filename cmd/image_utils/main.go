package main

import (
	"github.com/sekiju/image_utils/internal/modules/convert"
	"github.com/sekiju/image_utils/internal/modules/filter"
	"github.com/sekiju/image_utils/internal/modules/resize"
	"github.com/sekiju/image_utils/internal/modules/tile"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use: "image_utils",
}

func main() {
	root.Execute()
}

func init() {
	root.DisableSuggestions = true
	root.CompletionOptions.DisableDefaultCmd = true

	root.AddCommand(tile.Command)
	root.AddCommand(resize.Command)
	root.AddCommand(filter.Command)
	root.AddCommand(convert.Command)
}
