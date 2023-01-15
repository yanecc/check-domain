package main

type Balance struct {
	Data []struct {
		ProductID int `json:"product_id"`
		Product   struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"product"`
		Credits int `json:"credits"`
	} `json:"data"`
}
