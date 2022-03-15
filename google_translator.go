package main

import (
	"fmt"

	"github.com/bregydoc/gtranslate"
)

func tranlate_test() {
	text := "Olá mundo!"
	translated, err := gtranslate.TranslateWithParams(
		text,
		gtranslate.TranslationParams{
			From: "pt",
			To:   "en",
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("pt: %s | en: %s \n", text, translated)
}
