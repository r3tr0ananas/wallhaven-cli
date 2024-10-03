package main

import (
	"strings"

	"github.com/spf13/cobra"
)

func Search(cmd *cobra.Command, args []string) error {
	var url string

	query := strings.Join(args, " ")

	for url == "" {
		selections := []string{"Next page -->", "Previous Page <--"}
		images, err := SearchAPI(query, page)
	}
}

func Download(cmd *cobra.Command, args []string) error {
	return nil
}

func Edit(cmd *cobra.Command, args []string) error {
	return nil
}

func Collection(cmd *cobra.Command, args []string) error {
	return nil
}
