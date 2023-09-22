package externalapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/KoLLlaka/sobes/internal/logger"
)

type PeopleEnrich interface {
	AddAge(params map[string]string) (int, error)
	AddNation(params map[string]string) (string, error)
	AddGender(params map[string]string) (string, error)
}

type peopleEnrich struct {
	logger logger.MyLogger
}

type requestData struct {
	Count       int       `json:"count,omitempty"`
	Name        string    `json:"name,omitempty"`
	Gender      string    `json:"gender,omitempty"`
	Age         int       `json:"age,omitempty"`
	Probability float64   `json:"probability,omitempty"`
	Country     []country `json:"country,omitempty"`
	Error       string    `json:"error,omitempty"`
}

type country struct {
	CountryId   string  `json:"country_id,omitempty"`
	Probability float64 `json:"probability,omitempty"`
}

func (p peopleEnrich) AddAge(params map[string]string) (int, error) {
	requestData := requestData{}
	data, err := getDataFromApi("https://api.agify.io/", params)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	// fmt.Printf("[data] %s\n", data)

	json.Unmarshal(data, &requestData)

	return requestData.Age, nil
}
func (p peopleEnrich) AddGender(params map[string]string) (string, error) {
	requestData := requestData{}
	data, err := getDataFromApi("https://api.genderize.io/", params)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// fmt.Printf("[data] %s\n", data)

	json.Unmarshal(data, &requestData)

	return requestData.Gender, nil
}
func (p peopleEnrich) AddNation(params map[string]string) (string, error) {
	requestData := requestData{}
	data, err := getDataFromApi("https://api.nationalize.io/", params)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// fmt.Printf("[data] %s\n", data)

	json.Unmarshal(data, &requestData)

	if len(requestData.Country) > 0 {
		probable := requestData.Country[0]
		for _, v := range requestData.Country {
			if probable.Probability < v.Probability {
				probable = v
			}
		}

		return probable.CountryId, nil
	}

	return "", nil
}

func getDataFromApi(url string, params map[string]string) ([]byte, error) {
	var urlParams []string
	for key, value := range params {
		urlParams = append(urlParams, fmt.Sprintf("%s=%s", key, value))
	}

	if urlParams != nil {
		url += "?" + strings.Join(urlParams, "&")
	}

	// fmt.Println("[url]", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)

		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return bodyBytes, nil
}
