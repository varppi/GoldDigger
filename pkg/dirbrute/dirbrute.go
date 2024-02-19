package dirbrute

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/R00tendo/GoldDigger/pkg/logs"
)

var Threads int
var Wordlist string

func Brute(URL string) ([]string, error) {
	if Wordlist == "" {
		Wordlist = "wordlists/common.txt"
	}
	URL = strings.TrimSuffix(URL, "/")
	resultChan := make(chan string, 5000)
	wordlistHandle, err := os.Open(Wordlist)
	if err != nil {
		return nil, err
	}

	var threadsRunning int
	var counter int
	wordlistScanner := bufio.NewScanner(wordlistHandle)
	for wordlistScanner.Scan() {
		for threadsRunning > Threads {
			time.Sleep(50 * time.Millisecond)
		}
		threadsRunning++
		go request(URL+"/"+wordlistScanner.Text(), resultChan, &threadsRunning, &counter)
	}
	for threadsRunning != 0 {
		time.Sleep(500 * time.Millisecond)
	}

	var results []string
	for {
		select {
		case url := <-resultChan:
			results = append(results, url)
		default:
			goto exit
		}
	}
exit:
	logs.Success("Directory bruteforcing complete!")

	return results, nil
}

func request(URL string, resultChan chan string, threadsRunning *int, counter *int) {
	response, err := http.Get(URL)
	if err != nil {
		return
	}
	if response.StatusCode != 404 {
		resultChan <- URL
	}

	*counter++
	fmt.Print("\033[36m[i] Requests sent:", *counter, "\033[0m\r")

	*threadsRunning--
}
