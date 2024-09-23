package wallhaven

import (
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
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
	var page string
	var editor string

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
			query := strings.Join(args, " ")

			image_urls, err := Search(query, page)

			if err != nil {
				return err
			}

			selection, err := ShowSelection(image_urls)
			if err != nil {
				return err
			}

			return DirectURL(selection)
		},
	}

	var downloadCmd = &cobra.Command{
		Use:     "download [id]",
		Aliases: []string{"d"},
		Short:   "Download wallpaper with given id",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			if err := Download(id); err != nil {
				return err
			} else {
				log.Printf("Image: %s downloaded", id)
			}

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

	searchCmd.Flags().StringVarP(&page, "page", "p", "1", "Get page.")
	editCmd.Flags().StringVarP(&editor, "editor", "e", "", "Set custom editor.")

	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(previewCmd)
	rootCmd.AddCommand(editCmd)

	rootCmd.Execute()
}
