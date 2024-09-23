package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v2"
)

var cfg Config
var page string

func init() {
	var configPath string
	var configFile string
	var defaultWP string

	if runtime.GOOS == "linux" {
		user, err := user.Current()

		if err != nil {
			cli.Exit(err, 0)
		}

		configPath = filepath.Join(user.HomeDir, ".config", "wallhaven-cli")
		configFile = filepath.Join(configPath, "config.toml")
		defaultWP = filepath.Join(user.HomeDir, "Pictures", "wallpapers")

		os.MkdirAll(defaultWP, os.ModePerm)
		os.MkdirAll(configPath, os.ModePerm)

	} else if runtime.GOOS == "windows" {
		// Add Windows Support
	}

	if _, err := toml.DecodeFile(configFile, &cfg); err != nil {
		cfg = Config{
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
	app := &cli.App{
		Name:     "wallhaven",
		Usage:    "Search and download wallpapers from wallhaven",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "r3tr0ananas",
				Email: "ananas@ananas.moe",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "page",
				Value: "1",
			},
		},
		Commands: []*cli.Command{
			{
				Name:      "search",
				Aliases:   []string{"s"},
				Usage:     "Search for wallpapers",
				Args:      true,
				ArgsUsage: "Query",
				Action: func(ctx *cli.Context) error {
					args := ctx.Args()

					if args.Len() == 0 {
						return errors.New("No query given.")
					}

					query := strings.Join(args.Slice(), " ")

					image_urls, err := Search(query)

					if err != nil {
						return err
					}

					selection, err := ShowSelection(image_urls)
					if err != nil {
						return err
					}

					return DirectURL(selection)
				},
			},
			{
				Name:      "download",
				Aliases:   []string{"d"},
				Usage:     "Download wallpaper with given id",
				Args:      true,
				ArgsUsage: "ID",
				Action: func(ctx *cli.Context) error {
					args := ctx.Args()

					if args.Len() == 0 {
						return errors.New("You didn't give an ID.")
					}

					id := args.First()

					if err := Download(id); err != nil {
						return err
					} else {
						log.Println(fmt.Sprint("Image: ", id, " downloaded"))
					}

					return nil
				},
			},
			{
				Name:      "preview",
				Args:      true,
				ArgsUsage: "url",
				Action: func(ctx *cli.Context) error {
					args := ctx.Args()

					url := args.First()

					return Preview(url)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
