package main

import (
	"fmt"

	"github.com/bregydoc/gtranslate"
)

func main() {
	text := "Na tarefa anterior de alterações da tela de resultados do teste interativo por celular, prevemos um campo de informações chamado “Detalhes do componente” [ver na foto abaixo]. Este campo ainda não foi amplamente validado em termos de informações e organização dessas informações. "
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
