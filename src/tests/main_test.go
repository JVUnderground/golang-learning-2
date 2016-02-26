package tests

import "testing"
import "net/http"

func TestIf_style_files_are_accessable(t *testing.T) {
    resp, _ := http.Get("http://localhost:8080/style/style.css")
    if resp.StatusCode != 200 {
        t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
    }
}
