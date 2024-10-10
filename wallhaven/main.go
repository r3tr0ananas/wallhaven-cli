package main

import (
	"fmt"
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
var configFile string
var re = regexp.MustCompile(`\(([^)]+)\)`)

func init() {
	if runtime.GOOS == "linux" || runtime.GOOS == "freebsd" {
		usr, err := user.Current()
		if err != nil {
			log.Fatalf("Error while getting current user: %v", err)
			return
		}

		configPath := filepath.Join(usr.HomeDir, ".config", "wallhaven-cli")
		configFile = filepath.Join(configPath, "config.toml")

		if err := os.MkdirAll(configPath, os.ModePerm); err != nil {
			log.Fatalf("Failed to create config path %s: %v", configPath, err)
			os.Exit(1)
		}

		if _, decodeError := toml.DecodeFile(configFile, &config); decodeError != nil {
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

			WriteError := os.WriteFile(configFile, encodedData, os.ModePerm)
			if WriteError != nil {
				fmt.Print(err)
				os.Exit(1)
			}
		}

		if err := os.MkdirAll(config.TempFolder, os.ModePerm); err != nil {
			log.Fatalf("Failed to create temp folder %s: %v", config.TempFolder, err)
			os.Exit(1)
		}

		if err := os.MkdirAll(config.SaveFolder, os.ModePerm); err != nil {
			log.Fatalf("Failed to create save folder %s: %v", config.SaveFolder, err)
			os.Exit(1)
		}
	} else {
		log.Fatalf("Your os isn't supported: %s", runtime.GOOS)
		return
	}
}

var page int
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
		Short:   "Edit the wallhaven-cli config.",
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

	if err := rootCmd.Execute(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
