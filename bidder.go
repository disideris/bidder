package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Geo stuct for holding Country's latitude and longitude
type Geo struct {
	Country string `json:"country"`
	Lat     int    `json:"lat"`
	Lon     int    `json:"lon"`
}

// App struct for holding Id and name of application
type App struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Device struct for holding Operating system and Geographic Location sub struct
type Device struct {
	Os  string `json:"os"`
	Geo Geo    `json:"geo"`
}

// BidRequest struct for holding bid's id, Application and Device sub structs
type BidRequest struct {
	ID     string `json:"id"`
	App    App    `json:"app"`
	Device Device `json:"device"`
}

// BidResponse struct for holding the bidder's response
type BidResponse struct {
	ID  string          `json:"id"`
	Bid BidResponsePart `json:"bid"`
}

// BidResponsePart struct for holding campaignId, price and adm that is needed for BidResponse
type BidResponsePart struct {
	CampaignID string  `json:"campaignId"`
	Price      float64 `json:"price"`
	Adm        string  `json:"adm"`
}

// CampaignResponse struct for holding the campaign list when hitting the campaign api
type CampaignResponse struct {
	CampaignList []Campaign
}

// Campaign struct for holding the fileds for each campaing returned from the camapign api
type Campaign struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Price             float64  `json:"price"`
	Adm               string   `json:"adm"`
	TargetedCountries []string `json:"targetedCountries"`
}

// Helper function that Unmarshals json response from campaign api get calls
func getCampaigns(body []byte) (*CampaignResponse, error) {
	var s []Campaign
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	var cr = new(CampaignResponse)
	cr.CampaignList = s
	return cr, err
}

// Helper function that finds the campaign with the minimum price
func findMinCampaignPrice(campaigns []Campaign, bidrequest BidRequest) (Campaign, bool) {
	var bestPriceCapmaign Campaign
	validCampaign := false
	price := 10000.0
	for _, campaign := range campaigns {
		campaignCoutries := campaign.TargetedCountries
		for _, campaignCountry := range campaignCoutries {
			if campaignCountry == bidrequest.Device.Geo.Country && campaign.Price < price {
				price = campaign.Price
				bestPriceCapmaign = campaign
				validCampaign = true
			}
		}
	}
	return bestPriceCapmaign, validCampaign
}

// BidHandle Http handler function that handles the /bid path
func BidHandle(w http.ResponseWriter, r *http.Request) {
	var bidrequest BidRequest
	_ = json.NewDecoder(r.Body).Decode(&bidrequest)

	client := &http.Client{}
	//req, _ := http.NewRequest("GET", "http://demo2797226.mockable.io/", nil)
	req, _ := http.NewRequest("GET", "https://campaigns.apiblueprint.org/campaigns", nil)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Errored when sending request to the campaign server")
		return
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	campaigns, _ := getCampaigns(respBody)
	bestPriceCampaign, isCampaignValid := findMinCampaignPrice(campaigns.CampaignList, bidrequest)

	if isCampaignValid {
		bidresponse := BidResponse{bidrequest.ID, BidResponsePart{bestPriceCampaign.ID, bestPriceCampaign.Price, bestPriceCampaign.Adm}}
		jsonResponse, err := json.Marshal(bidresponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
		fmt.Println("200")

	} else {
		w.WriteHeader(204)
		fmt.Println("204")
	}
}

// Greeting Http handler that handles with a simple greeting the / path
func Greeting(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Blue Banana Api!")
}

// main function to boot up everything
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", Greeting).Methods("GET")
	router.HandleFunc("/bid", BidHandle).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}
