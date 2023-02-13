package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"log"
	// "strconv"
	"reflect"
)

// define security const var
const (
	authKey = "secret_6mlUHLWvSR03MxuS2fcUWzVYd9lbg9tJFZeHLWPVwHF"
	dbID    = "042070db669c497db454f41cb4e079e6"
)

func main() {

	client := &http.Client{}

	// define notion api with dbID
	url := fmt.Sprintf("https://api.notion.com/v1/databases/%s/query", dbID)

	// uncertain the type of value
	data := map[string]interface{}{
		"filter": map[string]interface{}{
			"and": []interface{}{
				map[string]interface{}{
					"property": "Tags",
					"multi_select": map[string]interface{}{
						"contains": "review",
					},
				},
				map[string]interface{}{
					"property": "Tags",
					"multi_select": map[string]interface{}{
						"does_not_contain": "d1",
					},
				},
				map[string]interface{}{
					"property": "Tags",
					"multi_select": map[string]interface{}{
						"does_not_contain": "d7",
					},
				},
				map[string]interface{}{
					"property": "Tags",
					"multi_select": map[string]interface{}{
						"does_not_contain": "d21",
					},
				},
				map[string]interface{}{
					"property": "Tags",
					"multi_select": map[string]interface{}{
						"does_not_contain": "d30",
					},
				},
			},
		},
	}


	// convert map to json
	payload, _ := json.Marshal(data)

	req, _ := http.NewRequest("POST", url, bytes.NewReader(payload))

	req.Header.Add("accept", "application/json")
	req.Header.Add("Notion-Version", "2022-06-28")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", "Bearer secret_6mlUHLWvSR03MxuS2fcUWzVYd9lbg9tJFZeHLWPVwHF")


	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	// 後面利用 defer 來確定是否有收到
	defer resp.Body.Close()


	// 將 response 裡的 Body 也就是網頁內容讀取出來
	sitemap, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}


	// Unmarshal response
	var result map[string]interface{}
	err = json.Unmarshal(sitemap, &result)
	if err != nil {
		fmt.Println("Error unmarshaling API response:", err)
		return
	}

	// type Person struct {
	// 	Id   int    `json:"id"`
	// }



	// var id map[string]int
	// err = json.Unmarshal(result, &id)
	fmt.Println("Type of x:", reflect.TypeOf(result))

	//TODO: Add "D1" tag
	pages := result["results"].([]interface{})
	for _, page := range pages {

		pageMap := page.(map[string]interface{})

		// list := 
		// fmt.Println("page:", pageMap)
		id, ok := pageMap["id"].(string)

		properties := pageMap["properties"].(map[string]interface{})
		tags := properties["Tags"].(map[string]interface{})
		LIS, ok := tags["multi_select"].([]interface{})

		fmt.Println("Type of LIS:", reflect.TypeOf(LIS))

		fmt.Println("%s", LIS)

		if ok {
			for _, v := range LIS {
				lisMap := v.(map[string]interface{})
				if name, ok := lisMap["name"].(string); ok {
					fmt.Println("Name:", name)
				}
			}
		}

		// Define the payload for the API request
		payload := map[string]interface{}{
			"properties": map[string]interface{}{
				"Tags": map[string]interface{}{
					"multi_select": []interface{}{
						map[string]interface{}{
							"name": "D1",
						},
					},
				},
			},
		}

		// 將原有的 tags 放進去
		for _, elem := range LIS {
			payload["properties"].(map[string]interface{})["Tags"].(map[string]interface{})["multi_select"] = append(payload["properties"].(map[string]interface{})["Tags"].(map[string]interface{})["multi_select"].([]interface{}), elem)
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Error marshalling payload:", err)
			return
		}

		url := fmt.Sprintf("https://api.notion.com/v1/pages/%s/", id)

		// Define the API request
		req, err := http.NewRequest("PATCH",url, bytes.NewBuffer(payloadBytes))
		if err != nil {
			fmt.Println("Error creating API request:", err)
			return
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("accept", "application/json")
		req.Header.Add("Notion-Version", "2022-06-28")
		req.Header.Add("authorization", "Bearer secret_6mlUHLWvSR03MxuS2fcUWzVYd9lbg9tJFZeHLWPVwHF")

		// Execute the API request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error executing API request:", err)
			return
		}
		defer resp.Body.Close()

		// Read the API response
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading API response:", err)
			return
		}
		fmt.Println("API response:", string(body))
	}
}