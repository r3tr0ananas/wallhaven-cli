package main

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

var config Config

func init() {
	if runtime.GOOS == "linux" || runtime.GOOS == "freebsd" {
		usr, err := user.Current()
		if err != nil {
			log.Fatalf("Error while getting current user: %v", err)
			return
		}

		configPath := filepath.Join(usr.HomeDir, ".config", "wallhaven-cli", "config.toml")

		if _, decodeError := toml.DecodeFile(configPath, &config); decodeError != nil {
			config = Config{
				Editor:     "nano",
				SaveFolder: filepath.Join(usr.HomeDir, "Pictures", "wallpapers"),
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
}
