package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/hbagdi/go-unsplash/unsplash"
	"github.com/jessevdk/go-flags"
	"golang.org/x/oauth2"
)

type Options struct {
	Verbose bool `short:"v" long:"verbose" description:"Show verbose debug information"`

	AccessKey string `short:"a" long:"access-key" description:"Unsplash Access Key" required:"true"`

	Tags []string `short:"t" long:"tag" description:"Tags to search for" required:"true"`

	MinResolutionX int `short:"x" long:"min-resolution-x" description:"Minimum resolution width" required:"false"`
	MinResolutionY int `short:"y" long:"min-resolution-y" description:"Minimum resolution height" required:"false"`

	MaxResolutionX int `short:"X" long:"max-resolution-x" description:"Maximum resolution width" required:"false"`
	MaxResolutionY int `short:"Y" long:"max-resolution-y" description:"Maximum resolution height" required:"false"`

	SavePath string `short:"p" long:"save-path" description:"Path to save images to" required:"false"`

	SetCommand string `short:"s" long:"set-command" description:"Command to run to set the wallpaper (%s for relative path, %S for absolute path)" required:"false"`
}

func main() {
	var options Options
	_, err := flags.Parse(&options)
	if err != nil {
		panic(err)
	}

	// Initialize oauth2 config
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "Client-ID " + options.AccessKey},
	)

	// Initialize the http client
	tc := oauth2.NewClient(context.Background(), ts)

	// Now, initialize the Unsplash client using the passed-in keys
	client := unsplash.New(tc)

	// Create search query from tags
	query := ""
	for _, tag := range options.Tags {
		query += tag + " "
	}
	query = strings.TrimSuffix(query, " ")

	// Create random options
	randOpts := &unsplash.RandomPhotoOpt{
		SearchQuery: query,
		Count:       30,
		Orientation: unsplash.Landscape, // Doesn't always work
	}

	// GET!
	photos, resp, err := client.Photos.Random(randOpts)
	if err != nil {
		panic(err)
	}

	// Print some debug info
	if options.Verbose {
		println("Total:", len(*photos))

		println("First page:", resp.FirstPage)
		println("Has next page:", resp.HasNextPage)
		println("Last page:", resp.LastPage)
		println("Next page:", resp.NextPage)
		println("Prev page:", resp.PrevPage)
		println("Rate limit:", resp.RateLimit)
		println("Rate remaining:", resp.RateLimitRemaining)
	}

	attempts := 0
	filePath := ""

	// Iterate through the photos
	for _, img := range *photos {
		attempts++

		// Check if the image meets the resolution requirements
		if !imageMeetsResReq(&img, options.MinResolutionX, options.MinResolutionY, options.MaxResolutionX, options.MaxResolutionY) {
			if options.Verbose {
				println("Image does not meet resolution requirements")
			}
			continue
		}

		// Download the image
		imgUrl := img.Urls.Raw.String()
		if options.Verbose {
			println("Downloading image:", imgUrl)
		}

		// Save the image to the save path
		if options.SavePath == "" {
			// Set to temp folder if no save path is specified (multi-platform support)
			options.SavePath = path.Join(os.TempDir(), "wpchg")

			// Create the temp folder if it doesn't exist
			if _, err := os.Stat(options.SavePath); os.IsNotExist(err) {
				os.Mkdir(options.SavePath, os.ModePerm)
			}

			// Log the save path
			if options.Verbose {
				println("No save path specified, using temp folder:", options.SavePath)
			}
		}

		// Create the file path
		filePath = path.Join(options.SavePath, (*img.ID)+".jpg")

		// Download the image
		resp, err := http.Get(imgUrl)
		if err != nil {
			panic(err)
		}

		// Create the file
		file, err := os.Create(filePath)
		if err != nil {
			panic(err)
		}

		// Copy the image to the file
		io.Copy(file, resp.Body)

		// Close the file
		file.Close()

		// Close the response body
		resp.Body.Close()

		// Log the file path
		if options.Verbose {
			println("Saved image to:", filePath)
		}

		// Stop after the first image
		break
	}

	if filePath == "" {
		if options.Verbose {
			println("No images found that meet the resolution requirements")
			return
		}
	} else {
		if options.Verbose {
			println("Image found!")
		}

		// Run the set command
		if options.SetCommand != "" {
			// Resolve the absolute path
			absPath, err := filepath.Abs(filePath)
			if err != nil {
				panic(err)
			}

			// Replace the %s and %S with the file path
			setCommand := strings.ReplaceAll(options.SetCommand, "%s", filePath)
			setCommand = strings.ReplaceAll(setCommand, "%S", absPath)

			// Split the command into the command and arguments
			splitCommand := strings.Split(setCommand, " ")
			command := splitCommand[0]
			args := splitCommand[1:]

			// Log verbosely
			if options.Verbose {
				println("Running command:", setCommand)
			}

			// Create the command
			cmd := exec.Command(command, args...)

			// If we had a clean way to do cross-platform, we would... but we don't. Here lies the Windows code.
			// // Hide console window on Windows
			// if runtime.GOOS == "windows" {
			// 	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			// }

			// Run the command
			err = cmd.Run()
			if err != nil {
				panic(err)
			}

			if options.Verbose {
				println("Set command ran successfully!")
			}
		}
	}

	if options.Verbose {
		println("Done in attempts:", attempts)
	}
}

func imageMeetsResReq(img *unsplash.Photo, minX int, minY int, maxX int, maxY int) bool {
	// Also check if the image is horizontal
	if *img.Width < *img.Height {
		return false
	}

	if minX > 0 && *img.Width < minX {
		return false
	}

	if minY > 0 && *img.Height < minY {
		return false
	}

	if maxX > 0 && *img.Width > maxX {
		return false
	}

	if maxY > 0 && *img.Height > maxY {
		return false
	}

	return true
}
