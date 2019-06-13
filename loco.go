package loco

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const endpointTemplate = "https://localise.biz:443/api/export/%s.json?key=%s"

// Translation represents a translation from Loco
// as structure that is compatible with go-i18n
type Translation struct {
	ID          string `json:"id"`
	Translation string `json:"translation"`
}

// FetchTranslations creates translations files from a Loco Project
// in the destination path it will try to
// fetch all available languages if none is provided
func FetchTranslations(locoKey, destinationPath string, langs ...string) error {
	if len(langs) > 0 {
		return fetchLangSpecifics(locoKey, destinationPath, langs)
	}
	return fetchAllLang(locoKey, destinationPath)
}

func fetchSpecificLang(locoKey, destinationPath, lang string) error {
	uri := fmt.Sprintf(endpointTemplate, "locale/"+lang, locoKey)
	locoData, err := sendLocoRequest(uri)
	if err != nil {
		return err
	}

	translationList := keyValueToTranslation(locoData)
	err = writeToFile(translationList, destinationPath, lang)
	if err != nil {
		return err
	}

	return nil
}

func sendLocoRequest(url string) (map[string]interface{}, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var locoData map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&locoData)
	if err != nil {
		return nil, err
	}

	return locoData, nil

}

func fetchLangSpecifics(locoKey, destinationPath string, langs []string) error {
	for _, lang := range langs {
		err := fetchSpecificLang(locoKey, destinationPath, lang)
		if err != nil {
			return err
		}
	}
	return nil
}

func fetchAllLang(locoKey, destinationPath string) error {
	uri := fmt.Sprintf(endpointTemplate, "all", locoKey)
	locoData, err := sendLocoRequest(uri)
	if err != nil {
		return err
	}
	for lang, keys := range locoData {
		translationList := keyValueToTranslation(keys.(map[string]interface{}))
		err := writeToFile(translationList, destinationPath, lang)
		if err != nil {
			return err
		}
	}
	return nil
}

func keyValueToTranslation(assets map[string]interface{}) []Translation {
	list := make([]Translation, 0)
	for id, translation := range assets {
		list = append(list, Translation{ID: id, Translation: translation.(string)})
	}
	return list
}

func writeToFile(translationList []Translation, path, lang string) error {
	preparePath(&path)
	output, err := json.Marshal(translationList)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fmt.Sprintf(path, lang), output, 0644)
}

func preparePath(path *string) {
	if !strings.HasSuffix(*path, "/") {
		*path = *path + "/"
	}
	*path = *path + "%s.json"
}
