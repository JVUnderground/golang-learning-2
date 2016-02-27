package tests

import (
    "testing"
    "net/http"
    "io/ioutil"
    "log"
    "regexp"
)

func TestIf_style_files_are_accessable(t *testing.T) {
    resp, _ := http.Get("http://localhost:8080/style")
    if resp.StatusCode != 200 {
        t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
    }
}

func TestIf_server_returns_a_404_page_with_garbage_url(t *testing.T) {
    resp, _ := http.Get("http://localhost:8080/this/is/garbage")
    htmlData, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
    match, _ := regexp.MatchString("404 page not found", string(htmlData))
    if match != true {
        t.Fatalf("Did not find a suitable 404 page.")
    }
}
