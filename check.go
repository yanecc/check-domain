package main

type Check struct {
	Domain struct {
		Availability string `json:"domainAvailability"`
		Name         string `json:"domainName"`
	} `json:"DomainInfo"`
}
