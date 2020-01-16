package loco

import "encoding/json"

// translation represents a translation from Loco
// as structure that is compatible with go-i18n

type keyValue map[string]string

type translationResponse map[string]keyValue

func (kv *keyValue) toTranslationList() []translation {
	list := make([]translation, 0)
	for id, t := range *kv {
		list = append(list, translation{ID: id, Translation: t})
	}
	return list
}

func (kv *keyValue) toJSONTranslationList() ([]byte, error) {
	return json.Marshal(kv.toTranslationList())
}

type translation struct {
	ID          string `json:"id"`
	Translation string `json:"translation"`
}
