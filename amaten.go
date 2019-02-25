package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type GiftsResponse struct {
	UserSignedIn bool  `json:"user_signed_in"`
	Gifts        Gifts `json:"gifts"`
	Total        int   `json:"total"`
	AllGiftCount int   `json:"all_gift_count"`
}

type Gifts []struct {
	ID              int    `json:"id"`
	Revision        int    `json:"revision"`
	FaceValue       int    `json:"face_value"`
	Price           int    `json:"price"`
	Type            string `json:"type"`
	Rate            string `json:"rate"`
	IsMine          bool   `json:"is_mine"`
	Cnt             int    `json:"cnt"`
	UsersTotalCount int    `json:"users_total_count"`
	UsersErrorCount int    `json:"users_error_count"`
}

func GetPrise(giftType string, limit int) (*GiftsResponse, error) {
	client := &http.Client{Timeout: time.Duration(10) * time.Second}

	// curl 'https://amaten.com/api/gifts?order=&type=amazon&limit=20&last_id='
	// -H 'User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3713.0 Safari/537.36'
	// -H 'X-Requested-With: XMLHttpRequest'
	values := url.Values{}
	values.Set("order", "")
	values.Set("type", giftType)
	values.Set("limit", strconv.Itoa(limit))
	values.Set("last_id", "")

	req, _ := http.NewRequest("GET", "https://amaten.com/api/gifts?"+values.Encode(), nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3713.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(res.Body)

	var giftRes GiftsResponse
	err = json.Unmarshal(body, &giftRes)

	return &giftRes, err
}
