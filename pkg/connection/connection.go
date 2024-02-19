package connection

import "net/http"

func Check(URL string) bool {
	_, err := http.Get(URL)
	if err != nil {
		return false
	} else {
		return true
	}
}
