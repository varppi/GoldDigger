package crawler

import (
	"io"
	"net/http"
	"regexp"
)

var UrlRegex *regexp.Regexp = regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)

func Crawl(URL string) []string {
	var URLs []string
	HTTPClient := &http.Client{}
	request, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return URLs
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.3")

	responseIO, err := HTTPClient.Do(request)
	if err != nil {
		return URLs
	}

	responseBytes, _ := io.ReadAll(responseIO.Body)
	responseString := string(responseBytes)
	responseBytes = nil
	URLs = UrlRegex.FindAllString(responseString, -1)
	return URLs
}
