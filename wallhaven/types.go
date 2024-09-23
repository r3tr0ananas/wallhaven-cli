package wallhaven

type Config struct {
	Editor     string `toml:"editor"`
	SaveFolder string `toml:"save_folder"`
}

type ResultItem struct {
	ID        string `json:"id"`
	FullImage string `json:"path"`
	FileType  string `json:"file_type"`
}

type SearchResult struct {
	Data []ResultItem `json:"data"`
}

type Result struct {
	Data  ResultItem `json:"data"`
	Error string     `json:"error"`
}
