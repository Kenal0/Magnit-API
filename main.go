package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type SearchRequest struct {
	IncludeAdultGoods bool       `json:"includeAdultGoods"`
	Pagination        Pagination `json:"pagination"`
	Sort              Sort       `json:"sort"`
	Term              string     `json:"term"`
	StoreCode         string     `json:"storeCode"`
	StoreType         string     `json:"storeType"`
	CatalogType       string     `json:"catalogType"`
}

type Pagination struct {
	Limit  int
	Offset int
}

type Sort struct {
	Order string
	Type  string
}

func main() {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("What goods you are looking for?")

	search, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input: ", err)
		return
	}

	payload := SearchRequest{
		IncludeAdultGoods: true,
		Pagination: Pagination{
			Offset: 0,
			Limit:  36,
		},
		Sort: Sort{
			Order: "desc",
			Type:  "popularity",
		},
		Term:        search,
		StoreCode:   "992301",
		StoreType:   "dostavka",
		CatalogType: "3",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://magnit.ru/webgate/v2/goods/search",
		bytes.NewReader(jsonData),
	)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("Status:", resp.Status)

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", " ")
	if err != nil {
		panic(err)
	}

	fmt.Println(prettyJSON.String())
}
