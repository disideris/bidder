package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	gock "gopkg.in/h2non/gock.v1"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/bid", BidHandle).Methods("POST")
	return router
}

func TestBidHandleWithResponse(t *testing.T) {
	defer gock.Off()

	gock.New("https://campaigns.apiblueprint.org/campaigns").
		Get("/").
		Reply(200).
		File("response_json")

	requestBody := []byte("{\n  \"id\": \"e7fe51ce4f6376876353ff0961c2cb0d\",\n  \"app\": {\n    \"id\": \"e7fe51ce-4f63-7687-6353-ff0961c2cb0d\",\n    \"name\": \"Morecast Weather\"\n  },\n  \"device\": {\n    \"os\": \"Android\",\n    \"geo\": {\n      \"country\": \"USA\",\n      \"lat\": 0,\n      \"lon\": 0\n    }\n  }\n}")
	request, _ := http.NewRequest("POST", "/bid", bytes.NewBuffer(requestBody))
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")

	expectedBidResponse := BidResponse{"e7fe51ce4f6376876353ff0961c2cb0d", BidResponsePart{"e919799e", 0.39, "<a href=\"http://example.com/click/qbFCjzXR9rkf8qa4\"><img src=\"http://assets.example.com/ad_assets/files/000/000/002/original/banner_300_250.png\" height=\"250\" width=\"300\" alt=\"\"/></a><img src=\"http://example.com/win/qbFCjzXR9rkf8qa4\" height=\"1\" width=\"1\" alt=\"\"/>\r\n"}}
	r := BidResponse{}
	json.Unmarshal(response.Body.Bytes(), &r)

	assert.Equal(t, r, expectedBidResponse, "BidResponse is not as expected")

}

func TestBidHandleWithoutResponse(t *testing.T) {
	defer gock.Off()

	gock.New("https://campaigns.apiblueprint.org/campaigns").
		Get("/").
		Reply(200).
		File("response_json2")

	requestBody := []byte("{\n  \"id\": \"e7fe51ce4f6376876353ff0961c2cb0d\",\n  \"app\": {\n    \"id\": \"e7fe51ce-4f63-7687-6353-ff0961c2cb0d\",\n    \"name\": \"Morecast Weather\"\n  },\n  \"device\": {\n    \"os\": \"Android\",\n    \"geo\": {\n      \"country\": \"USA\",\n      \"lat\": 0,\n      \"lon\": 0\n    }\n  }\n}")
	request, _ := http.NewRequest("POST", "/bid", bytes.NewBuffer(requestBody))
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 204, response.Code, "OK response is expected")
}
