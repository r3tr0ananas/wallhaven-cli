package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

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

func Preview(url string) error {
	var cmd *exec.Cmd
	var file string

	fzf_preview_lines := os.Getenv("FZF_PREVIEW_LINES")
	fzf_preview_columns := os.Getenv("FZF_PREVIEW_COLUMNS")
	_, exists := os.LookupEnv("KITTY_WINDOW_ID")

	path, _ := exec.LookPath("chafa")

	if strings.Contains(url, "-->") {
		return nil
	} else if strings.Contains(url, "<--") {
		return nil
	}

	if exists {
		cmd = exec.Command(
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
		DirectURL([]string{url}, &cfg.TmpFolder)

		SplitUrl := strings.Split(url, "/")
		FileName := SplitUrl[len(SplitUrl)-1]

		file = filepath.Join(cfg.TmpFolder, FileName)

		cmd = exec.Command(
			path,
			file,
			fmt.Sprintf("--size=%sx%s", fzf_preview_columns, fzf_preview_lines),
			"--clear",
		)
	} else {
		return errors.New("You either need chafa or the kitty terminal.")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	if file != "" {
		return os.Remove(file)
	}

	return nil
}
