package api

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
)

func (s *Server) API(method string, header map[string][]string, endURL string, params url.Values, body io.Reader) (response map[string]interface{}, err error) {
	apiURL := "https://api.tdameritrade.com/v1/" + endURL
	if params != nil {
		apiURL = apiURL + "?" + params.Encode()
	}
	if header == nil {
		header = url.Values{}
	}

	_, apiResponse, httpCode, err := api.htclient.ProcessRequest(method, header, apiURL, body)
	if err != nil {
		return nil, err
	}
	if api.debug {
		fmt.Println("API URL: " + apiURL)
		fmt.Println("Status Code: " + strconv.Itoa(httpCode))
		fmt.Println("Response: " + utils.ReturnPrettyPrintJson(apiResponse))
	}
	if httpCode == 204 {
		return map[string]interface{}{"status": "Written", "httpcode": "204"}, nil
	}
	jsonResp, err := utils.GetJson(apiResponse)
	if err != nil {
		return map[string]interface{}{"status": string(apiResponse)}, err
	}

	if httpCode >= 400 {
		aString := make([]string, len(jsonResp.(map[string]interface{})["errors"].([]interface{})))
		for i, v := range jsonResp.(map[string]interface{})["errors"].([]interface{}) {
			aString[i] = v.(string)
		}
		return jsonResp, errors.New(strings.Join(aString, ","))
	}
	// Checking whether the response is of string format rather than JSON format
	stringResp, found := jsonResp.(string)
	if found {
		return map[string]interface{}{"status": stringResp}, nil
	}
	return jsonResp.(map[string]interface{}), nil
}
