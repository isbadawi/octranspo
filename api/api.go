// This package interacts with the OC Transpo API. We really only use one
// endpoint, which lists the upcoming trips (for all routes) on a given stop.

// n.b. There is an XML API and a JSON API, but the latter is problematic. For
// example, we're interested in the "routes" field of the response, which lists
// all the routes with upcoming trips for that stop. But instead of "routes"
// being an array, it is omitted if there are no routes, an object if there is
// a single route, and an array if there are two or more. The same dance
// happens with the trips inside each route.

// It doesn't seem possible to cleanly handle this using encoding/json, so we
// stick with the XML API, which is (slightly) more sane.

package api

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	APP_ID  = "0711096f"
	API_KEY = "630f075c3cd02faba16409cc70ac5137"
	API_URL = "https://api.octranspo1.com/v1.2"
)

func request(endpoint string, params url.Values) ([]byte, error) {
	url := fmt.Sprintf("%v/%v", API_URL, endpoint)
	params.Set("appID", APP_ID)
	params.Set("apiKey", API_KEY)
	response, err := http.PostForm(url, params)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

type Route struct {
	Number    string `xml:"RouteNo"`
	Direction string
	Heading   string `xml:"RouteHeading"`
	Trips     []struct {
		Destination string `xml:"TripDestination"`
		StartTime   string `xml:"TripStartTime"`
	} `xml:"Trips>Trip"`
}

func GetNextTripsForStop(stop string) ([]Route, error) {
	params := url.Values{}
	params.Set("stopNo", stop)
	body, err := request("GetNextTripsForStopAllRoutes", params)
	if err != nil {
		return nil, err
	}

	response := struct {
		Error  string  `xml:"Body>GetRouteSummaryForStopResponse>GetRouteSummaryForStopResult>Error"`
		Routes []Route `xml:"Body>GetRouteSummaryForStopResponse>GetRouteSummaryForStopResult>Routes>Route"`
	}{}

	err = xml.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	// TODO(isbadawi): What are the different error codes?
	if response.Error == "10" {
		return nil, fmt.Errorf("invalid stop number")
	}
	return response.Routes, nil
}
