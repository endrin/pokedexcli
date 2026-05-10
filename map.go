package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pokedexcli/internal/pokecache"
	"time"
)

const (
	MAP_BASE_URL     string        = "https://pokeapi.co/api/v2/location-area/"
	PAGE_SIZE        int           = 20
	CACHE_RETIRE_AGE time.Duration = 5 * time.Minute
)

type Page interface {
	Next() ([]string, error)
	Previous() ([]string, error)
}

type mapPage struct {
	cache    *pokecache.Cache
	next     string
	previous string
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

func newMapPage() mapPage {
	return mapPage{
		cache: pokecache.NewCache(CACHE_RETIRE_AGE),
		next:  fmt.Sprintf("%s?limit=%d", MAP_BASE_URL, PAGE_SIZE),
	}
}

func (mp *mapPage) get(pageUrl string) ([]string, error) {
	var data io.ReadCloser
	if cached, ok := mp.cache.Get(pageUrl); ok {
		data = io.NopCloser(bytes.NewBuffer(cached))
	} else {
		res, err := http.Get(pageUrl)
		if err != nil {
			return nil, err
		}

		if res.StatusCode > 299 {
			res.Body.Close()
			return nil, fmt.Errorf("non-OK HTTP status: %s", res.Status)
		}

		data = res.Body
	}
	defer data.Close()

	decoder := json.NewDecoder(data)
	var mr mapResponse
	if err := decoder.Decode(&mr); err != nil {
		return nil, err
	}

	mp.next = mr.Next
	mp.previous = mr.Previous

	var names []string
	for _, m := range mr.Results {
		names = append(names, m.Name)
	}
	return names, nil
}

func (mp *mapPage) Next() ([]string, error) {
	return mp.get(mp.next)
}

func (mp *mapPage) Previous() ([]string, error) {
	if mp.previous == "" {
		return nil, fmt.Errorf("you're on the first page")
	}
	return mp.get(mp.previous)
}

func addMap(registry commandsRegistry) {
	currentPage := newMapPage()

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
