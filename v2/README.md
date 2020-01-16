## go-loco v2

A simple lib to fetch translation files from [Localise](https://localise.biz/) and save it in a [go-i18n](https://github.com/nicksnyder/go-i18n) compatible format.

## Features
- Multiple project management
- Automatic download new translations 

Checkout some [examples](./examples).

## Basic Usage

```golang
package main

import (
	"context"
	"fmt"

	loco "github.com/lucasmdrs/go-loco/v2"
	"github.com/nicksnyder/go-i18n/i18n"
)

func main() {
	g := loco.Init()

	if err := g.AddProject("MY_LOCALISE_KEY", "./"); err != nil {
		panic(err)
	}

	g.FetchTranslations(context.TODO())

	i18n.MustLoadTranslationFile("en-US.json")
	i18n.MustLoadTranslationFile("pt-BR.json")

	T, _ := i18n.Tfunc("pt-BR")

	fmt.Println(T("hello"))
}
```