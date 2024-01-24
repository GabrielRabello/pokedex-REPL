package dto

type Location struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Region      region        `json:"region"`
	Names       []names       `json:"names"`
	GameIndices []gameIndices `json:"game_indices"`
	Areas       []areas       `json:"areas"`
}
type region struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type language struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type names struct {
	Name     string   `json:"name"`
	Language language `json:"language"`
}
type generation struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type gameIndices struct {
	GameIndex  int        `json:"game_index"`
	Generation generation `json:"generation"`
}
type areas struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
