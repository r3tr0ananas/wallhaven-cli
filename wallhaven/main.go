package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

var cfg Config
var configFile string

func init() {
	var configPath string
	var defaultWP string

	if runtime.GOOS == "linux" {
		user, err := user.Current()

		if err != nil {
			os.Exit(0)
		}

		configPath = filepath.Join(user.HomeDir, ".config", "wallhaven-cli")
		configFile = filepath.Join(configPath, "config.toml")
		defaultWP = filepath.Join(user.HomeDir, "Pictures", "wallpapers")

		os.MkdirAll(defaultWP, os.ModePerm)
		os.MkdirAll(configPath, os.ModePerm)

	} else if runtime.GOOS == "windows" {
		// Windows support ain't happening
	}

	if _, err := toml.DecodeFile(configFile, &cfg); err != nil {
		cfg = Config{
			Editor:     "nano",
			SaveFolder: defaultWP,
			Categories: CategoriesType{
				General: true,
				Anime:   true,
				People:  true,
			},
		}

		log.Println("Default config folder:", cfg.SaveFolder)

		if b, err := toml.Marshal(cfg); err == nil {
			if writeErr := os.WriteFile(configFile, b, os.ModePerm); writeErr != nil {
				log.Fatalf("Failed to write config file: %v", writeErr)
			} else {
				log.Println("Created config file at", configFile)
			}
		} else {
			log.Fatalf("Failed to marshal config: %v", err)
		}

	}
}

func main() {
	var page int
	var editor string
	var downloadAll bool

	var rootCmd = &cobra.Command{
		Use:   "wallhaven",
		Short: "Search and download wallpapers from wallhaven",
	}

	var searchCmd = &cobra.Command{
		Use:     "search [query]",
		Aliases: []string{"s"},
		Short:   "Search on wallhaven",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var url string

			query := strings.Join(args, " ")

			for url == "" {
				selections := []string{"Next page -->", "Previous Page <--"}
				image_urls, err := Search(query, page)

				if err != nil {
					return err
				}

				selections = append(selections, image_urls...)

				selection, err := ShowSelection(selections, true)
				if err != nil {
					return err
				}

				if strings.Contains(selection, "-->") {
					page++
				} else if strings.Contains(selection, "<--") {
					page--
				} else {
					url = selection
				}
			}

			return DirectURL([]string{url})
		},
	}

	var downloadCmd = &cobra.Command{
		Use:     "download [id/ids]",
		Aliases: []string{"d"},
		Short:   "Download wallpaper with given id",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Download(args); err != nil {
				return err
			} else {
				log.Printf("Done downloading")
			}

			return nil
		},
	}

	var collectionCmd = &cobra.Command{
		Use:   "collection [username]",
		Short: "Download images from a selection",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			username := args[0]

			collections, err := GetCollections(username)
			if err != nil {
				return err
			}

			selection, err := ShowSelection(collections, false)
			if err != nil {
				return err
			}

			re := regexp.MustCompile(`\(([^)]+)\)`)

			matches := re.FindStringSubmatch(selection)

			if len(matches) > 1 {
				id := matches[1]

				images, err := GetCollectionImages(username, id)
				if err != nil {
					return err
				}

				if !downloadAll {
					selectedImage, err := ShowSelection(images, true)

					images = []string{selectedImage}

					if err != nil {
						return err
					}
				}

				err = DirectURL(images)
				if err != nil {
					return err
				}
			} else {
				return errors.New("uh?")
			}

			log.Println("Done downloading the images/image")

			return nil
		},
	}

	var previewCmd = &cobra.Command{
		Use:   "preview [url]",
		Short: "Preview image.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			url := args[0]

			return Preview(url)
		},
	}

	var editCmd = &cobra.Command{
		Use:   "edit",
		Short: "Edit wallhaven's config.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if editor != "" {
				cfg.Editor = editor
			}

			configCmd := exec.Command(cfg.Editor, configFile)

			configCmd.Stdin = os.Stdin
			configCmd.Stdout = os.Stdout
			configCmd.Stderr = os.Stderr

			return configCmd.Run()
		},
	}

	searchCmd.Flags().IntVarP(&page, "page", "p", 1, "Get page.")
	editCmd.Flags().StringVarP(&editor, "editor", "e", "", "Set custom editor.")
	collectionCmd.Flags().BoolVarP(&downloadAll, "all", "a", false, "Download all images from the collection.")

	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(collectionCmd)
	rootCmd.AddCommand(previewCmd)
	rootCmd.AddCommand(editCmd)

	rootCmd.Execute()
}
