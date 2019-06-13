## go-loco

A simple lib to fetch translation files from [Localise](https://localise.biz/) and save it in a [go-i18n](https://github.com/nicksnyder/go-i18n) compatible format.

## Usage

```golang
package main

import (
	"fmt"

	loco "github.com/lucasmdrs/go-loco"
	"github.com/nicksnyder/go-i18n/i18n"
)

func main() {
	err := loco.FetchTranslations("MY_LOCALISE_KEY", "./", "pt-BR", "en-US")
	if err != nil {
		panic(err)
	}

	i18n.MustLoadTranslationFile("en-US.json")
	i18n.MustLoadTranslationFile("pt-BR.json")

	T, _ := i18n.Tfunc("pt-BR")

	fmt.Println(T("hello"))
}
```