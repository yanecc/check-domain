package main

type Config struct {
	ApiKey       string
	AccurateMode bool
	UseWhois     bool
}

type Check struct {
	Domain struct {
		Availability string `json:"domainAvailability"`
		Name         string `json:"domainName"`
	} `json:"DomainInfo"`
}

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
