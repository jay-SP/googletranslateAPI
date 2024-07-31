package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/jay/go-translate/cli"
)

var (
	sourceLanguage string
	targetLanguage string
	sourceText     string
)

func init() {
	flag.StringVar(&sourceLanguage, "s", "en", "Source language (default: en)")
	flag.StringVar(&targetLanguage, "t", "fr", "Target language (default: fr)")
	flag.StringVar(&sourceText, "st", "", "Text to translate")
}

func main() {
	flag.Parse()

	if sourceText == "" {
		fmt.Println("Error: Text to translate is required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	translatedText, err := translateText(sourceLanguage, targetLanguage, sourceText)
	if err != nil {
		fmt.Printf("Error translating text: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(processTranslation(translatedText))
}

func translateText(sourceLang, targetLang, text string) (string, error) {
	var wg sync.WaitGroup
	strChan := make(chan string)

	wg.Add(1)
	defer wg.Wait()
	defer close(strChan)

	reqBody := &cli.RequestBody{
		SourceLang: sourceLang,
		TargetLang: targetLang,
		SourceText: text,
	}

	go cli.RequestTranslate(reqBody, strChan, &wg)

	// Wait for the translation to complete
	return <-strChan, nil
}

func processTranslation(text string) string {
	return strings.ReplaceAll(text, "+", " ")
}
