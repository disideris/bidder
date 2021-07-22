package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Appinfo struct for holding the Appinfo attributes
type AppInfo struct {
	ID        string
	Name      string
	Category  string
	Publisher string
}

func unmarshallAppInfo(body []byte) (AppInfo, error) {
	var appInfo AppInfo
	err := json.Unmarshal(body, &appInfo)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	return appInfo, err
}

m := make(map[string]AppInfo)

func getAppInfo(id string) AppInfo {

	if val, ok := m[id]; ok {
		return val
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "/appinfo/{id}", nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the service X")
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	appInfo, _ := unmarshallAppInfo(respBody)
	m[id] = appInfo
	return appInfo

}
