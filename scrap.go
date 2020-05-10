package main

import (	
	"net/http"
	"github.com/PuerkitoBio/goquery"
)

// Word struct for import
type Word struct {
	Origin string `json:"origin"`
	Translate string `json:"translate"`
}
// GetData data from website
func GetData() ([]Word, error) {
	url := "https://studynow.ru/dicta/allwords" 

	// Get the HTML
	resp, err := http.Get(url)
	if err != nil {
		return []Word{}, err
	}
	defer resp.Body.Close()

	// Convert HTML into goquery document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return []Word{}, err
	}

	// Save each 
	var words = []Word{}
	doc.Find("#wordlist tbody tr").Each(func(i int, s *goquery.Selection) {
		word := Word{
			Origin: s.Find("td:nth-child(2)").Text(),
			Translate: s.Find("td:nth-child(3)").Text(),
		}
		words = append(words, word)
	})
	return words, nil
}

