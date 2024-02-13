package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/jay/go-translate/cli"
)

var wg sync.WaitGroup

var sourceLanguage string
var targetLanguage string
var sourceText string

func init() {
	flag.StringVar(&sourceLanguage, "s", "en", "Source language[en]")
	flag.StringVar(&targetLanguage, "t", "fr", "Target Language")
	flag.StringVar(&sourceText, "st", "", "Text to translate")
}

func main() {
	flag.Parse()

	if flag.NFlag() == 0 {
		fmt.Println("Options: 	")
		flag.PrintDefaults()
		os.Exit(1)
	}

	strChan := make(chan string)

	wg.Add(1)

	reqBody := &cli.RequestBody{
		SourceLang: sourceLanguage,
		TargetLang: targetLanguage,
		SourceText: sourceText,
	}

	go cli.RequestTranslate(reqBody, strChan, &wg)
	processedStr := strings.ReplaceAll(<-strChan, "+", " ")

	fmt.Printf("%s\n", processedStr)
	wg.Wait()

}
