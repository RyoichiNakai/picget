package downloader

import (
	"fmt"
	"net/http"
)

type RangeAccess interface {
	HasRangeAceessHeader (url string) bool
	HTTPRangeRequest (req http.Request) http.Response
}

func Range() {
	fmt.Println("Range")
}

