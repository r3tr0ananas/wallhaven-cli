package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func DLSave(url string, folder *string) error {
	if folder == nil {
		folder = &config.SaveFolder
	}

	request, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error while doing a GET request to download the image: %v", err)
	}

	image, ioError := io.ReadAll(request.Body)
	if ioError != nil {
		return fmt.Errorf("Error while reading image: %v", image)
	}

	urlSplit := strings.Split(url, "/")
	name := urlSplit[len(urlSplit)-1]

	file := filepath.Join(*folder, name)

	perm := os.FileMode(0644)

	if err := os.WriteFile(file, image, perm); err != nil {
		return fmt.Errorf("Error while writing file to %v: %v", file, err)
	}

	return nil
}

func SearchAPI(query string, page int) (ImagesResponse, error) {
	var images ImagesResponse
	url, err := url.Parse("https://wallhaven.cc/api/v1/search")
	if err != nil {
		return ImagesResponse{}, fmt.Errorf("huh?")
	}

	m := map[bool]int{false: 0, true: 1}

	categories := fmt.Sprintf(
		"%d%d%d",
		m[config.Searching.Categories.General],
		m[config.Searching.Categories.Anime],
		m[config.Searching.Categories.People],
	)

	params := url.Query()
	params.Add("q", query)
	params.Add("page", fmt.Sprint(page))
	params.Add("categories", categories)
	params.Add("order", config.Searching.Order)
	params.Add("topRange", config.Searching.TopRange)

	if config.Searching.AtLeast != "" {
		params.Add("atleast", config.Searching.AtLeast)
	}

	if config.Searching.Resolutions != nil {
		resolutions := strings.Join(config.Searching.Resolutions, ",")
		params.Add("resolutions", resolutions)
	}

	if config.Searching.Ratios != nil {
		ratios := strings.Join(config.Searching.Ratios, ",")
		params.Add("ratios", ratios)
	}

	url.RawQuery = params.Encode()

	response, err := http.Get(url.String())
	if err != nil {
		return ImagesResponse{}, err
	}
	defer response.Body.Close()

	decodingError := json.NewDecoder(response.Body).Decode(&images)
	if decodingError != nil {
		return ImagesResponse{}, fmt.Errorf("Error while decoding json: %v", decodingError)
	}

	return images, nil
}

func DownloadAPI(id string) error {
	var image ImageResponse
	url := fmt.Sprintf("https://wallhaven.cc/api/v1/w/%v", id)

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	decodingError := json.NewDecoder(response.Body).Decode(&image)
	if decodingError != nil {
		return fmt.Errorf("Error while decoding json: %v", decodingError)
	}

	return DLSave(image.Image.ImageURL, nil)
}

func CollectionsAPI(username string) (CollectionsResponse, error) {
	var collections CollectionsResponse
	url := fmt.Sprintf("https://wallhaven.cc/api/v1/collections/%s", username)

	response, err := http.Get(url)
	if err != nil {
		return CollectionsResponse{}, err
	}
	defer response.Body.Close()

	decodingError := json.NewDecoder(response.Body).Decode(&collections)
	if decodingError != nil {
		return CollectionsResponse{}, fmt.Errorf("Error while decoding json: %v", decodingError)
	}

	return collections, nil
}

func CollectionAPI(username string, id string, page int) (ImagesResponse, error) {
	var images ImagesResponse
	collectionURL, err := url.Parse(
		fmt.Sprintf("https://wallhaven.cc/api/v1/collections/%s/%s", username, id),
	)

	if err != nil {
		return ImagesResponse{}, nil
	}

	query := collectionURL.Query()
	query.Add("page", fmt.Sprint(page))

	collectionURL.RawQuery = query.Encode()

	response, err := http.Get(collectionURL.String())
	if err != nil {
		return ImagesResponse{}, err
	}
	defer response.Body.Close()

	decodingError := json.NewDecoder(response.Body).Decode(&images)
	if decodingError != nil {
		return ImagesResponse{}, fmt.Errorf("Error while decoding json: %v", decodingError)
	}

	return images, nil
}
