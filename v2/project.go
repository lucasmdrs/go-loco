package loco

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"reflect"
)

type authVerifyResponse struct {
	project `json:"project"`
}
type project struct {
	key        string
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	assetsPath string
	lastResult translationResponse
}

func getProjectInformation(key string) (project, error) {
	res, err := http.Get(fmt.Sprintf(baseURL+authEndpoint+authParameter, key))
	if err != nil {
		return project{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return project{}, errors.New("Invalid Key")
	}
	var data authVerifyResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return project{}, err
	}
	data.project.key = key
	data.project.lastResult = make(translationResponse)
	return data.project, nil
}

func (p *project) fetchProjectTranslations(notifier chan interface{}) error {
	uri := fmt.Sprintf(baseURL+endpointTemplate+authParameter, "all", p.key)
	res, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("Failed to request project translations")
	}

	var locoData translationResponse
	if err := json.NewDecoder(res.Body).Decode(&locoData); err != nil {
		return err
	}

	if reflect.DeepEqual(locoData, p.lastResult) {
		return nil
	}

	for lang, keys := range locoData {
		if reflect.DeepEqual(keys, p.lastResult[lang]) {
			continue
		}

		json, err := keys.toJSONTranslationList()
		if err != nil {
			return err
		}

		if err := writeToFile(json, p.assetsPath, lang); err != nil {
			return err
		}

		p.lastResult[lang] = keys
		log.Println(p.Name, ":", lang)
	}

	notifier <- nil
	return nil
}

func writeToFile(content []byte, destination, lang string) error {
	destination = path.Join(destination, "%s.json")
	return ioutil.WriteFile(fmt.Sprintf(destination, lang), content, 0644)
}
