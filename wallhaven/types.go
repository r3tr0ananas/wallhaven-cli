package main

type CategoriesType struct {
	General bool `toml:"general"`
	Anime   bool `toml:"anime"`
	People  bool `toml:"people"`
}

type Config struct {
	Editor     string         `toml:"editor"`
	SaveFolder string         `toml:"save_folder"`
	Categories CategoriesType `toml:"categories"`
}

type ResultItem struct {
	ID        string `json:"id"`
	FullImage string `json:"path"`
}

type Results struct {
	Data  []ResultItem `json:"data"`
	Error string       `json:"error"`
}

type Result struct {
	Data  ResultItem `json:"data"`
	Error string     `json:"error"`
}

type CollectionItem struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

type CollectionList struct {
	Data  []CollectionItem `json:"data"`
	Error string           `json:"error"`
}
