package feed

import (
	"bufio"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const (
	endpoint = "http://localhost:8080"
	apiKey   = "DEMO_KEY"
)

var localHost = endpoint

func TestMain(m *testing.M) {
	handle, err := os.Open("testdata/news.json")
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	text := ""
	scanner := bufio.NewScanner(handle)
	for scanner.Scan() {
		text += scanner.Text()
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if strings.Contains(req.URL.String(), "/news") {
			_, _ = rw.Write([]byte(text))
		} else {
			rw.WriteHeader(404)
		}
	}))
	localHost = server.URL

	code := m.Run()
	defer server.Close()
	os.Exit(code)
}
