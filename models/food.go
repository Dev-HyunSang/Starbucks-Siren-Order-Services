package models

type Foods struct {
	FoodUUID string `json:"food_uuid"`
	Title    string `json:"title"`
	SubTitle string `json:"subtitle"`
	Info     string `json:"info"`
	Prices   string `json:"prices"`
}
