package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/unrolled/render"
)

var (
	formatter = render.New(render.Options{
		IndentJSON: true,
	})
)

const (
	fakeMatchLocationResult = "/matches/5a003b78-409e-4452-b456-a6f0dcee05bd"
)

func TestCreateMatch(t *testing.T) {
	client := &http.Client{}
	repo := newInMemoryRepository()
	server := httptest.NewServer(http.HandlerFunc(createMatchHandler(formatter, repo)))
	defer server.Close()

	body := []byte("{\n  \"gridsize\": 19,\n  \"players\": [\n    {\n      \"color\": \"white\",\n      \"name\": \"bob\"\n    },\n    {\n      \"color\": \"black\",\n      \"name\": \"alfred\"\n    }\n  ]\n}")

	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error in creating POST request for createMatchHandler: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST to createMatchHandler: %v", err)
	}

	defer res.Body.Close()

	payload, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error parsing response body: %v", err)
	}

	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected response status 201, received %s", res.Status)
	}

	if res.Header["Location"] == nil {
		loc := res.Header["Location"]
		if loc == nil {
			t.Error("Location header is not set")
		} else {
			if !strings.Contains(loc[0], "/matches/") {
				t.Errorf("Location header should contain '/matches/'")
			}
			if len(loc[0]) != len(fakeMatchLocationResult) {
				t.Errorf("Location value does not contain guid of new match")
			}
		}

		var matchResponse newMatchResponse
		err = json.Unmarshal(payload, &matchResponse)
		if err != nil {
			t.Errorf("Could not unmarshal payload into newMatchResponse object")
		}

		if matchResponse.ID == "" || !strings.Contains(loc[0], matchResponse.ID) {
			t.Error("matchResponse.Id does not match Location header")
		}

		// After creating a match, match repository should have 1 item in it.
		matches := repo.getMatches()
		if len(matches) != 1 {
			t.Errorf("Expected a match repo of 1 match, got size %d", len(matches))
		}
	}

	fmt.Printf("Payload: %s", string(payload))
}
