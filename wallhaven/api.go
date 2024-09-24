package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func DirectURL(url []string) error {
	for _, v := range url {
		image, err := http.Get(v)
		if err != nil {
			return errors.New("Failed to GET Request the image.")
		}

		bytes, err := io.ReadAll(image.Body)

		SplitUrl := strings.Split(v, "/")
		FileName := SplitUrl[len(SplitUrl)-1]

		file := filepath.Join(cfg.SaveFolder, FileName)

		if err := os.WriteFile(file, bytes, os.ModePerm); err != nil {
			return errors.New("Failed to save image.")
		}
	}

	return nil
}

func Download(ids []string) error {
	var urls []string
	for _, id := range ids {
		var data Result
		requestUrl := fmt.Sprintf("https://wallhaven.cc/api/v1/w/%s", id)

		r, err := http.Get(requestUrl)

		if err != nil {
			return err
		}
		defer r.Body.Close()

		json.NewDecoder(r.Body).Decode(&data)
		if data.Error != "" {
			return fmt.Errorf("An image with that ID \"%s\" doesn't exist. %v", id, data.Error)
		}

		urls = append(urls, data.Data.FullImage)
	}

	return DirectURL(urls)
}

func Search(query string, page string) ([]string, error) {
	var images Results

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

func GetCollections(username string) ([]string, error) {
	var collectionList CollectionList

	collectionURL := fmt.Sprintf("https://wallhaven.cc/api/v1/collections/%s", username)

	r, err := http.Get(collectionURL)

	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	json.NewDecoder(r.Body).Decode(&collectionList)

	if collectionList.Error != "" {
		return nil, errors.New("This user doesn't exist.")
	}

	if len(collectionList.Data) == 0 {
		return nil, errors.New("This user doesn't have any collection.")
	}

	collections := []string{}

	for _, element := range collectionList.Data {
		label := fmt.Sprintf("%s (%d)", element.Label, element.ID)

		collections = append(collections, label)
	}

	return collections, nil
}

func GetCollectionImages(username string, id string) ([]string, error) {
	var results Results

	collectionURL := fmt.Sprintf("https://wallhaven.cc/api/v1/collections/%s/%s", username, id)

	r, err := http.Get(collectionURL)

	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	json.NewDecoder(r.Body).Decode(&results)

	if results.Error != "" {
		return nil, errors.New("There is no such ID.")
	}

	images := []string{}

	for _, element := range results.Data {
		images = append(images, element.FullImage)
	}

	return images, nil
}
