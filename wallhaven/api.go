package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func DirectURL(url string) error {
	image, err := http.Get(url)
	if err != nil {
		return errors.New("Failed to GET Request the image.")
	}

	bytes, err := io.ReadAll(image.Body)

	SplitUrl := strings.Split(url, "/")
	FileName := SplitUrl[len(SplitUrl)-1]

	file := filepath.Join(cfg.SaveFolder, FileName)

	if err := os.WriteFile(file, bytes, os.ModePerm); err != nil {
		return errors.New("Failed to save image.")
	}

	return nil
}

func Download(id string) error {
	var data Result
	requestUrl := "https://wallhaven.cc/api/v1/w/" + id

	r, err := http.Get(requestUrl)

	if err != nil {
		return err
	}
	defer r.Body.Close()

	json.NewDecoder(r.Body).Decode(&data)
	if data.Error != "" {
		return errors.New("An image with that ID doesn't exist.")
	}

	return DirectURL(data.Data.FullImage)
}

func Search(query string, page string) ([]string, error) {
	var images SearchResult

	base, err := url.Parse("https://wallhaven.cc/api/v1/search")
	if err != nil {
		return nil, err
	}

	params := base.Query()
	params.Set("q", query)
	params.Set("page", page)
	base.RawQuery = params.Encode()

	request, err := http.NewRequest("GET", base.String(), nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	r, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	json.NewDecoder(r.Body).Decode(&images)

	if len(images.Data) == 0 {
		return nil, errors.New("No images found based on query")
	}

	image_urls := []string{}

	for _, element := range images.Data {
		image_urls = append(image_urls, element.FullImage)
	}

	return image_urls, nil
}
