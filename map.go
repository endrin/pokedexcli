package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const PAGE_SIZE int = 20

type Page interface {
	Next() ([]string, error)
	Previous() ([]string, error)
}

type mapPage struct {
	cache  map[int][]string
	offset int
}

type mapResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func (mp *mapPage) get() ([]string, error) {
	if cached, ok := mp.cache[mp.offset]; ok {
		return cached, nil
	}

	pageUrl := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/?limit=%d&offset=%d", PAGE_SIZE, mp.offset)
	res, err := http.Get(pageUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("non-OK HTTP status: %s", res.Status)
	}

	decoder := json.NewDecoder(res.Body)
	var mr mapResponse
	if err := decoder.Decode(&mr); err != nil {
		return nil, err
	}
	var names []string
	for _, m := range mr.Results {
		names = append(names, m.Name)
	}

	mp.cache[mp.offset] = names
	return names, nil
}

func (mp *mapPage) Next() ([]string, error) {
	result, err := mp.get()
	if err == nil {
		mp.offset += PAGE_SIZE
	}
	return result, err
}

func (mp *mapPage) Previous() ([]string, error) {
	if mp.offset < PAGE_SIZE {
		return nil, fmt.Errorf("you're on the first page")
	}
	mp.offset -= PAGE_SIZE
	return mp.get()
}

func addMap(registry commandsRegistry) {
	currentPage := mapPage{
		cache: make(map[int][]string),
	}

	registry.register(cliCommand{
		name:        "map",
		description: "Displays name of 20 Location Areas. Run again to get next 20 areas.",
		callback: func() error {
			names, err := currentPage.Next()
			if err != nil {
				return err
			}
			for _, name := range names {
				fmt.Println(name)
			}

			return nil
		},
	})

	registry.register(cliCommand{
		name:        "mapb",
		description: "Displays name of previous 20 Location Areas. Run again to get previous 20 areas.",
		callback: func() error {
			names, err := currentPage.Previous()
			if err != nil {
				return err
			}
			for _, name := range names {
				fmt.Println(name)
			}

			return nil
		},
	})
}
