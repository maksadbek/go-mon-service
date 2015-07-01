package route

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestComputeHMAC(t *testing.T) {
	msgMAC := computeHMAC(TestMsg, TestKey)
	mac := hmac.New(sha256.New, []byte(TestKey))
	mac.Write([]byte(TestMsg))
	expectedMAC := mac.Sum(nil)
	if !hmac.Equal(msgMAC, expectedMAC) {
		t.Errorf("want %v\n, got %v", expectedMAC, msgMAC)
	}
}

func TestCheckMAC(t *testing.T) {
	mac := hmac.New(sha256.New, []byte(TestKey))
	mac.Write([]byte(TestMsg))
	expectedMAC := mac.Sum(nil)
	if !checkMAC(TestMsg, expectedMAC, TestKey) {
		t.Errorf("want %b, got %b", true, false)
	}
}

// TestLogoutHandler tests logout request
func TestLogoutHandler(t *testing.T) {
	/*
		msg = `{
				"user":"newmax",
				"selectedFleetJs":"202",
				"groups":"1,2,3",
				"uid":"testuid"
				}`
	*/
	// signup
	signupMsg := `{
			"hash":"f8cb56593dd08e04cd0f84d796b9cecd",
			"uid":"newmax",
			"user":"newmax"
	}`

	res, err := http.Post("http://localhost:8080/signup", "application/json", strings.NewReader(signupMsg))
	if err != nil {
		t.Error(err)
	}

	if res.Status != "200 OK" {
		t.Errorf("want 200, got %s", res.Status)
	}

	defer res.Body.Close()
	token, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(token))
	// logout
	logoutMsg := `{"token":"` + string(token) + `"}`
	res, err = http.Post("http://localhost:8080/logout", "application/json", strings.NewReader(logoutMsg))
	if err != nil {
		t.Error(err)
	}

	if res.Status != "200 OK" {
		t.Errorf("want 200, got %s", res.Status)
	}
	defer res.Body.Close()
}
