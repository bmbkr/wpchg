# WPChg - Unsplash Wallpaper Changer

[![Version](https://img.shields.io/github/v/tag/bmbkr/wpchg?label=version&color=green)](https://niea.me/wpchg/releases/latest)
![Language](https://img.shields.io/badge/language-Go-blue)
![License](https://img.shields.io/badge/license-GPLv3-orange)

## Table of Contents
- [Overview](#overview)
- [Installation](#installation)
- [Usage](#usage)
- [Options](#options)
- [File Structure](#file-structure)
- [Build](#build)
- [Contributing](#contributing)
- [License](#license)

## Overview
WPChg is a simple and efficient tool for grabbing wallpapers from Unsplash. Written in Go, it's designed for customization and ease-of-use.

## Installation
To install WPChg, you can either download the binaries or compile from source.

### From Binaries
Download the appropriate binary from the latest [Release](https://niea.me/wpchg/releases/latest).

### From Source
```
git clone https://github.com/yourusername/wpchg.git
cd wpchg
make build
```

## Usage
Run the application with the following command:

```
./wpchg [OPTIONS]

# For example:
./wpchg -a [Unsplash Access Key] -x 1920 -y 1080 -p ~/.wallpapers \
 -s "sh -c ~/.setWallpaper.sh %S" -t nature -t trees -t rain
```

## Options
- `-v` or `--verbose`: Shows verbose debug information.
- `-a` or `--access-key`: Unsplash Access Key. (Required)
- `-t` or `--tag`: Tags to search for on Unsplash. (Required)
- `-x` or `--min-resolution-x`: Minimum resolution width. (Optional)
- `-y` or `--min-resolution-y`: Minimum resolution height. (Optional)
- `-X` or `--max-resolution-x`: Maximum resolution width. (Optional)
- `-Y` or `--max-resolution-y`: Maximum resolution height. (Optional)
- `-p` or `--save-path`: Path to save downloaded images to. (Optional)
- `-s` or `--set-command`: Command to run to set the wallpaper. Use `%s` for relative path and `%S` for absolute path. (Optional)

## File Structure
```
- Makefile
- cmd/
  - wpchg/
    - main.go
- go.mod
- go.sum
```

## Build
The project uses a Makefile for building and installing. For more information, check out the Makefile itself.

## Contributing
Feel free to submit a pull request or create an issue.

## License
GNU General Public License v3.0. See [LICENSE](./LICENSE)
