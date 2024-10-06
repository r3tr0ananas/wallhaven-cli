package main

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

var config Config
var configPath string
var re = regexp.MustCompile(`\(([^)]+)\)`)

func init() {
	if runtime.GOOS == "linux" || runtime.GOOS == "freebsd" {
		usr, err := user.Current()
		if err != nil {
			log.Fatalf("Error while getting current user: %v", err)
			return
		}

		configPath = filepath.Join(usr.HomeDir, ".config", "wallhaven-cli", "config.toml")

		if _, decodeError := toml.DecodeFile(configPath, &config); decodeError != nil {
			config = Config{
				Editor:     "nano",
				SaveFolder: filepath.Join(usr.HomeDir, "Pictures", "wallpapers"),
				TempFolder: filepath.Join("/tmp", "wallhaven-cli"),
				Searching: SearchParams{ // https://wallhaven.cc/help/api#search
					Categories: CategoriesType{
						General: true,
						Anime:   true,
						People:  true,
					},
					Sorting:     "date_added",
					Order:       "desc",
					TopRange:    "1M",
					AtLeast:     "",
					Resolutions: nil,
					Ratios:      nil,
				},
			}

			encodedData, EncodeError := toml.Marshal(config)
			if EncodeError != nil {
				log.Fatalf("Error while encoding default toml data: %v", err)
				return
			}

			os.WriteFile(configPath, encodedData, os.ModePerm)
		}

		os.MkdirAll(config.TempFolder, os.ModePerm)
		os.MkdirAll(config.SaveFolder, os.ModePerm)
	} else {
		log.Fatalf("Your os isn't supported: %v", runtime.GOOS)
		return
	}
}

var page int
var editor string
var all bool

func main() {
	var rootCmd = &cobra.Command{
		Use:   "wallhaven",
		Short: "Search and download wallpapers from wallhaven.",
	}

	var searchCmd = &cobra.Command{
		Use:     "search",
		Aliases: []string{"s"},
		Short:   "Search on wallhaven.",
		Args:    cobra.MinimumNArgs(1),
		RunE:    Search,
	}

	var downloadCmd = &cobra.Command{
		Use:     "download",
		Aliases: []string{"d"},
		Short:   "Download images from wallhaven.",
		Args:    cobra.MinimumNArgs(1),
		RunE:    Download,
	}

	var editCmd = &cobra.Command{
		Use:     "edit",
		Aliases: []string{"e"},
		Short:   "Edit the wallpaper-cli config.",
		RunE:    Edit,
	}

	var collectionCmd = &cobra.Command{
		Use:     "collection",
		Aliases: []string{"c"},
		Short:   "Retrieve a collection from wallhaven.",
		Args:    cobra.MinimumNArgs(1),
		RunE:    Collection,
	}

	var previewCmd = &cobra.Command{
		Use:  "preview",
		Args: cobra.MinimumNArgs(1),
		RunE: Preview,
	}

	searchCmd.Flags().IntVarP(&page, "page", "p", 1, "Specify the page number for paginated results")

	collectionCmd.Flags().IntVarP(&page, "page", "p", 1, "Specify the page number for paginated results")
	collectionCmd.Flags().BoolVarP(&all, "all", "a", false, "Download all images from collection")

	rootCmd.AddCommand(searchCmd, downloadCmd, editCmd, collectionCmd, previewCmd)

	rootCmd.Execute()
}
