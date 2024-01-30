package dto

type BatchResult struct {
	Count    int       `json:"count"`
	Next     *string   `json:"next"`
	Previous *string   `json:"previous"`
	Results  []results `json:"results"`
}

type results struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
