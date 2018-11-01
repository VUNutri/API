package recipe

type Recipe struct {
	ID int
	Title string `json:"title"`
	Category string `json:"category"`
	Time int `json:"time"`
	Products []
}