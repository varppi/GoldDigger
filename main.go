package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/R00tendo/GoldDigger/pkg/connection"
	"github.com/R00tendo/GoldDigger/pkg/crawler"
	"github.com/R00tendo/GoldDigger/pkg/dirbrute"
	"github.com/R00tendo/GoldDigger/pkg/logs"
)

var Settings struct {
	URL       string
	Keyword   string
	Out       string
	Wordlist  string
	Depth     int
	Threads   int
	Quiet     bool
	NoBrute   bool
	FilesOnly bool
}

var UrlRegex *regexp.Regexp = regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)

var SensitiveFiles []string = []string{".sql", ".bak", ".conf", ".txt", ".json", ".sqlite", ".sqlite3", ".db", ".zip", ".tar", ".rar", ".old", ".gz", ".config", ".pdf"}

func main() {
	fmt.Println()
	logo := "\033[33m" + `	                                                      |
	   ▄████   █████   ██     █████▄   █████▄   ██   ▄████o   ▄████  █████   ██▀███|  
	  ██  ▀█ ██   ██  ██     ██▀ ██▌ |██▀ ██▌  ██  ██  ▀█   ██  ▀█  █   ▀   ██   ██o
	| ██ ▄▄▄  ██   ██  ██    ██   █▌ o██   █▌  ██  ██ ▄▄▄   ██ ▄▄▄  ███  |  ██  ▄█ |
	o ▓█  ██▓ ██   ██  ██    ▓█▄   ▌  ▓█▄  |▌  ██  ▓█  ██▓  ▓█  ██  ▓█  ▄o  █▀ █▄  o
	  ▓███▀|  ████▓  ██████  ████▓    ████▓o   ██  ▓███▀|   ▓███▀   ████|   █▓  ██
	| o |  o   |     o o    |             o       |     o   o      |    o
	o   o      o            o                     o                o    
	` + "\033[0m"
	for _, line := range strings.Split(logo, "\n") {
		fmt.Println(line)
		time.Sleep(70 * time.Millisecond)
	}
	fmt.Println()
	flag.StringVar(&Settings.URL, "u", "", "URL: URL to crawl.")
	flag.StringVar(&Settings.Keyword, "k", "", "Keyword: A keyword the URL must include in it to be considered part of the scope. (default: same domain)")
	flag.StringVar(&Settings.Out, "o", "", "Output: Writes results to a file.")
	flag.StringVar(&Settings.Wordlist, "w", "", "Wordlist: Custom directory bruteforce wordlist")
	flag.IntVar(&Settings.Depth, "d", 2, "Depth: How deep to crawl.")
	flag.IntVar(&Settings.Threads, "t", 10, "Threads: How many threads to use in directory bruteforcing.")
	flag.BoolVar(&Settings.Quiet, "q", false, "Quiet: Does not display the costemic stuff.")
	flag.BoolVar(&Settings.NoBrute, "n", false, "No bruteforce: Doesn't perform directory bruteforcing.")
	flag.BoolVar(&Settings.FilesOnly, "f", false, "Files only: Only shows files in the results.")
	flag.Parse()

	if Settings.URL == "" {
		logs.Error("Argument -u (URL) is required!")
		flag.Usage()
		os.Exit(0)
	}

	logs.Quiet = Settings.Quiet
	dirbrute.Threads = Settings.Threads
	dirbrute.Wordlist = Settings.Wordlist

	if !UrlRegex.MatchString(Settings.URL) {
		logs.Error("Invalid URL!")
		os.Exit(0)
	} else if !connection.Check(Settings.URL) {
		logs.Error("Can't connect to URL!")
		os.Exit(0)
	}

	MainURLInfo, err := url.Parse(Settings.URL)
	if err != nil {
		logs.Error(err.Error())
		os.Exit(0)
	}

	allURLsFound := make(map[string]bool)
	var targetURLs []string = []string{Settings.URL}

	if !Settings.Quiet && !Settings.NoBrute {
		logs.Info("Directory bruteforcing")
	}

	//Dir brute
	if !Settings.NoBrute {
		dirBruteURLs, err := dirbrute.Brute(Settings.URL)
		if err != nil {
			logs.Error(err.Error())
		}
		for _, URL := range dirBruteURLs {
			allURLsFound[URL] = true
		}
	}

	//Crawler
	logs.Info("Starting crawl")
	for depth := 1; depth <= Settings.Depth; depth++ {
		logs.Info("Crawling layer:" + strconv.Itoa(depth))
		var tempTargetURLs []string
		for _, URL := range targetURLs {
			tempTargetURLs = append(tempTargetURLs, crawler.Crawl(URL)...)
		}
		targetURLs = nil
		for _, URL := range tempTargetURLs {
			if Settings.Keyword == "" {
				URLInfo, err := url.Parse(URL)
				if err != nil {
					continue
				}
				if strings.TrimPrefix(URLInfo.Hostname(), "www.") != strings.TrimPrefix(MainURLInfo.Hostname(), "www.") {
					continue
				}
			} else {
				if !strings.Contains(URL, Settings.Keyword) {
					continue
				}
			}
			if !allURLsFound[URL] {
				URLInfo, _ := url.Parse(URL)
				for _, fileExt := range SensitiveFiles {
					if strings.HasSuffix(URLInfo.Path, fileExt) {
						if !Settings.Quiet {
							URL = "\033[31m" + URL + "\033[0m (POTENTIALLY SENSITIVE)"
						}
					}
				}
				targetURLs = append(targetURLs, URL)

				if Settings.FilesOnly {
					if !strings.Contains(URLInfo.Path, ".") {
						continue
					}
				}
				allURLsFound[URL] = true
			}
		}
	}

	//RESULTS
	fmt.Println()
	if len(allURLsFound) == 0 {
		logs.Error("No results :(")
		return
	}
	logs.Success("RESULTS:")
	var outputHandle *os.File
	if Settings.Out != "" {
		outputHandle, err = os.OpenFile(Settings.Out, os.O_WRONLY|os.O_CREATE, 666)
		if err != nil {
			logs.Error(err.Error())
		}
	}

	for URL, _ := range allURLsFound {
		if !Settings.Quiet {
			logs.Success("URL: " + URL)
		} else {
			fmt.Println(URL)
		}
		outputHandle.WriteString(URL + "\n")
	}
}
