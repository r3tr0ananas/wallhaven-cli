package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func search(images ImagesResponse) (string, error) {
	selections := []string{}

	if len(images.Images) == 0 {
		return "", fmt.Errorf("No collections from this user.")
	}

	if images.Meta.LastPage > page {
		selections = append(selections, "Next page -->")
	}

	if page > 1 {
		selections = append(selections, "Previous page <--")
	}

	for _, v := range images.Images {
		preview := fmt.Sprintf("%v (%v)", v.Resolution, v.ImageURL)
		selections = append(selections, preview)
	}

	selection, err := ShowSelection(selections, true)
	if err != nil {
		return "", err
	}

	return selection, nil
}

func Search(cmd *cobra.Command, args []string) error {
	var url string

	query := strings.Join(args, " ")

	for url == "" {
		images, err := SearchAPI(query, page)
		if err != nil {
			return err
		}

		selection, err := search(images)
		if err != nil {
			return err
		}

		if strings.Contains(selection, "-->") {
			page++
		} else if strings.Contains(selection, "<--") {
			page--
		} else {
			url = re.FindStringSubmatch(selection)[1]
		}
	}

	err := DLSave(url, nil)
	if err != nil {
		return err
	}

	fmt.Printf("[download] %v", url)
	return nil
}

func Download(cmd *cobra.Command, args []string) error {
	for _, v := range args {
		err := DownloadAPI(v)
		if err != nil {
			return err
		}

		fmt.Printf("[download] %v", v)
	}

	return nil
}

func Edit(cmd *cobra.Command, args []string) error {
	command := exec.Command(config.Editor, configPath)

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}

func Collection(cmd *cobra.Command, args []string) error {
	username := args[0]
	selections := []string{}

	collections, err := CollectionsAPI(username)
	if err != nil {
		return err
	}

	if len(collections.Collections) == 0 {
		return fmt.Errorf("This user has no collection.")
	}

	for _, v := range collections.Collections {
		selection := fmt.Sprintf("%v (%v)", v.Label, v.ID)
		selections = append(selections, selection)
	}

	SelectedID, err := ShowSelection(selections, false)
	if err != nil {
		return err
	}

	id := re.FindStringSubmatch(SelectedID)[1]

	if all {
		for {
			images, err := CollectionAPI(username, id, page)
			if err != nil {
				return err
			}

			if len(images.Images) == 0 {
				return fmt.Errorf("This collection has no images.")
			}

			for _, v := range images.Images {
				err := DLSave(v.ImageURL, nil)
				if err != nil {
					return err
				}

				fmt.Printf("[download] %v\r\n", v.ImageID)
			}
			if images.Meta.LastPage > page {
				page++
			} else {
				return nil
			}
		}
	}

	var url string

	for url == "" {
		images, err := CollectionAPI(username, id, page)
		if err != nil {
			return err
		}

		for _, v := range images.Images {
			preview := fmt.Sprintf("%v (%v)", v.Resolution, v.ImageURL)
			selections = append(selections, preview)
		}

		selection, err := search(images)
		if err != nil {
			return err
		}

		if strings.Contains(selection, "-->") {
			page++
		} else if strings.Contains(selection, "<--") {
			page--
		} else {
			url = re.FindStringSubmatch(selection)[1]
		}
	}

	ImageURL := re.FindStringSubmatch(url)[0]

	err = DLSave(ImageURL, nil)
	if err != nil {
		return err
	}

	fmt.Printf("[download] %v", ImageURL)
	return nil
}

func ShowSelection(items []string, preview bool) (string, error) {
	input := strings.Join(items, "\n")

	cmd := exec.Command("fzf")
	if preview {
		cmd.Args = append(cmd.Args, "--preview=wallhaven preview {}")
	}

	cmd.Stdin = strings.NewReader(input)

	selected, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(selected)), nil
}

func Preview(cmd *cobra.Command, args []string) error {
	var command *exec.Cmd
	var file string

	text := args[0]

	fzf_preview_lines := os.Getenv("FZF_PREVIEW_LINES")
	fzf_preview_columns := os.Getenv("FZF_PREVIEW_COLUMNS")
	_, exists := os.LookupEnv("KITTY_WINDOW_ID")

	path, _ := exec.LookPath("chafa")

	if strings.Contains(text, "-->") {
		return nil
	} else if strings.Contains(text, "<--") {
		return nil
	}

	url := re.FindStringSubmatch(text)[1]

	if exists {
		command = exec.Command(
			"kitty",
			"icat",
			"--clear",
			"--transfer-mode=memory",
			"--unicode-placeholder",
			"--stdin=no",
			fmt.Sprintf("--place=%sx%s@0x0", fzf_preview_columns, fzf_preview_lines),
			"--scale-up",
			url,
		)
	} else if path != "" {
		err := DLSave(url, &config.TempFolder)
		if err != nil {
			return err
		}

		SplitUrl := strings.Split(url, "/")
		FileName := SplitUrl[len(SplitUrl)-1]

		file = filepath.Join(config.TempFolder, FileName)

		command = exec.Command(
			path,
			file,
			fmt.Sprintf("--size=%sx%s", fzf_preview_columns, fzf_preview_lines),
			"--clear",
		)
	} else {
		return fmt.Errorf("You either need chafa or the kitty terminal.")
	}

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		return fmt.Errorf("failed to execute command: %v", err)
	}

	if file != "" {
		return os.Remove(file)
	}

	return nil
}
