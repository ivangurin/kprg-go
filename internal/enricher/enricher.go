package enricher

import (
	"encoding/json"
	"fmt"
	"io"
	"kprg/internal/cacher"
	"net/http"
	"sort"
)

type Enricher struct {
	cacher cacher.Cacher
}

type age struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

type gender struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
}

type nationality struct {
	Count   int    `json:"count"`
	Name    string `json:"name"`
	Country []struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"country"`
}

func NewEnricher(cacher cacher.Cacher) *Enricher {

	enricher :=
		Enricher{
			cacher: cacher,
		}

	return &enricher

}

func (e *Enricher) GetAge(name string) (int, error) {

	key := fmt.Sprintf("age/%s", name)

	exists, err := e.cacher.Has(key)
	if err != nil {
		return 0, err
	}

	if exists {

		value, err := e.cacher.Get(key)
		if err != nil {
			return 0, err
		}

		return value.(int), nil

	}

	resp, err := http.Get(fmt.Sprintf("https://api.agify.io/?name=%s", name))
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	age := age{}

	err = json.Unmarshal(body, &age)
	if err != nil {
		return 0, err
	}

	return age.Age, nil

}

func (e *Enricher) GetGender(name string) (string, error) {

	key := fmt.Sprintf("gender/%s", name)

	exists, err := e.cacher.Has(key)
	if err != nil {
		return "", err
	}

	if exists {

		value, err := e.cacher.Get(key)
		if err != nil {
			return "", err
		}

		return value.(string), nil

	}

	resp, err := http.Get(fmt.Sprintf("https://api.genderize.io/?name=%s", name))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	gender := gender{}

	err = json.Unmarshal(body, &gender)
	if err != nil {
		return "", err
	}

	return gender.Gender, nil

}

func (e *Enricher) GetNationality(name string) (string, error) {

	key := fmt.Sprintf("gender/%s", name)

	exists, err := e.cacher.Has(key)
	if err != nil {
		return "", err
	}

	if exists {

		value, err := e.cacher.Get(key)
		if err != nil {
			return "", err
		}

		return value.(string), nil

	}

	resp, err := http.Get(fmt.Sprintf("https://api.nationalize.io/?name=%s", name))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	nationality := nationality{}

	err = json.Unmarshal(body, &nationality)
	if err != nil {
		return "", err
	}

	if len(nationality.Country) > 0 {

		sort.Slice(nationality.Country, func(i, j int) bool {
			return nationality.Country[i].Probability > nationality.Country[j].Probability
		})

		return nationality.Country[0].CountryID, nil

	}

	return "", nil

}
