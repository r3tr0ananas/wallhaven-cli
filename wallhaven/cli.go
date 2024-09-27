package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
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
	fzf_preview_lines := os.Getenv("FZF_PREVIEW_LINES")
	fzf_preview_columns := os.Getenv("FZF_PREVIEW_COLUMNS")
	_, exists := os.LookupEnv("KITTY_WINDOW_ID")

	if !exists {
		return errors.New("You're not using kitty term") // add chafa support
	}

	if strings.Contains(url, "-->") {
		return nil
	} else if strings.Contains(url, "<--") {
		return nil
	}

	cmd := exec.Command(
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

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	return nil
}
