package main

type Response struct {
	Results []Result `json:"results"`
}

type Result struct {
	Name string `json:"name"`
}
