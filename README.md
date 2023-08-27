# image_utils

## Modules

Each module has 4 basic commands:
- `threads`: Number of simultaneously processed images. The default is equal to the number of CPU cores
- `input`: Input folder with folders and images
- `output`: Output folder (in some commands not required)
- `include_sub_directories`: Images will also be searched in subfolders. Off by default

### Convert
Converts images from input folder to another format. Supports `.psd`
- `format`: Export format. Values: `png` or `jpeg`
- `quality`: Parameter for `format: jpeg` responsible for final quality. Default: `90`

### Filter
Filters images by parameters
- Filters to check the resolution: `minimal_width`, `minimal_height`, `maximal_width`, `maximal_height`. Value in pixels, example: `--minimal_width=700`
- `index_naming`: Files will be renamed to index: 001.png ... 013.png
- `primary_side`: Skip if the opposite side is larger than the main side. Values: `width` or `height`

### Resize
- `width`.  Value in pixels, example: `--width=700`
- `height`.  Value in pixels, example: `--height=700`
- `percent`. Factor resize by percent. If the setting is enabled, it has the highest priority over `width` or `height`. Write a number from 1 to 100 inclusive. Example: `-p=50`
- `filter`: Image resize filter. Values: `nearest_neighbor, box, linear, hermite, mitchell_netravali, catmull_rom, bspline, gaussian, bartlett, lanczos, hann, hamming, blackman, welch, cosine`. Default value: `box`

### Tile
Image tiling script
- `tile`: Size of tile. Default: `128`
- `minimum`. Sets a minimum size for tiles. If any tiles are below the specified size, they will not be saved. Default: `false`
- `skip_black_white`. Skip black and white tiles. Default: `false`
- `count_per_image`: The number of tiles to save per image. Default: `0 aka. unlimited`
