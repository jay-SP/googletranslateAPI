package cli

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/Jeffail/gabs"
)

// RequestBody holds the data for the translation request
type RequestBody struct {
	SourceLang string // Source language
	TargetLang string // Target language
	SourceText string // Text to be translated
}

// API URL for Google Translate
const translateURL = "https://translate.googleapis.com/translate_a/single"

// RequestTranslate sends a request to the Google Translate API and sends the result to the provided channel
func RequestTranslate(body *RequestBody, resultChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{}
	req, err := createRequest(body)
	if err != nil {
		resultChan <- fmt.Sprintf("Error creating request: %s", err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		resultChan <- fmt.Sprintf("Error making request: %s", err)
		return
	}
	defer res.Body.Close()

	if err := handleRateLimiting(res.StatusCode, resultChan); err != nil {
		resultChan <- err.Error()
		return
	}

	translatedText, err := parseResponse(res.Body)
	if err != nil {
		resultChan <- fmt.Sprintf("Error parsing response: %s", err)
		return
	}

	resultChan <- translatedText
}

// createRequest initializes an HTTP GET request for the Google Translate API
func createRequest(body *RequestBody) (*http.Request, error) {
	req, err := http.NewRequest("GET", translateURL, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("client", "gtx")
	query.Add("sl", body.SourceLang)
	query.Add("tl", body.TargetLang)
	query.Add("dt", "t")
	query.Add("q", body.SourceText)
	req.URL.RawQuery = query.Encode()

	return req, nil
}

// handleRateLimiting checks if the response indicates rate limiting
func handleRateLimiting(statusCode int, resultChan chan<- string) error {
	if statusCode == http.StatusTooManyRequests {
		return fmt.Errorf("Rate limit exceeded. Try again later.")
	}
	if statusCode != http.StatusOK {
		return fmt.Errorf("Received non-OK status code: %d", statusCode)
	}
	return nil
}

// parseResponse extracts the translated text from the API response
func parseResponse(body io.Reader) (string, error) {
	parsedJSON, err := gabs.ParseJSONBuffer(body)
	if err != nil {
		return "", err
	}

	nestOne, err := parsedJSON.ArrayElement(0)
	if err != nil {
		return "", err
	}

	nestTwo, err := nestOne.ArrayElement(0)
	if err != nil {
		return "", err
	}

	translatedStr, err := nestTwo.ArrayElement(0)
	if err != nil {
		return "", err
	}

	translatedText, ok := translatedStr.Data().(string)
	if !ok {
		return "", fmt.Errorf("Unexpected data type for translation")
	}

	return translatedText, nil
}
