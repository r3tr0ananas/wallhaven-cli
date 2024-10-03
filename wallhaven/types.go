package main

type CategoriesType struct {
	General bool `toml:"general"`
	Anime   bool `toml:"anime"`
	People  bool `toml:"people"`
}

type SearchParams struct {
	Categories  CategoriesType `toml:"categories"`
	Sorting     string         `toml:"sorting"`
	Order       string         `toml:"order"`
	TopRange    string         `toml:"top_range"`
	AtLeast     string         `toml:"at_least"`
	Resolutions []string       `toml:"resolutions"`
	Ratios      []string       `toml:"ratios"`
}

type Config struct {
	Editor     string       `toml:"editor"`
	SaveFolder string       `toml:"save_folder"`
	Searching  SearchParams `toml:"searching"`
}

type Tag struct {
	Name string `json:"name"`
}

type Image struct {
	ImageID    string `json:"id"`
	URL        string `json:"url"`
	Resolution string `json:"resolution"`
	ImageURL   string `json:"path"`
	Tags       []Tag  `json:"tags"`
}

type CollectionType struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type CollectionsResponse struct {
	Collections []CollectionType `json:"id"`
}

type ImageResponse struct {
	Image Image `json:"data"`
}

type ImagesResponse struct {
	Images []Image `json:"data"`
}
