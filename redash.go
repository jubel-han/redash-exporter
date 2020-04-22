package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Alert struct {
	ID       int    `json:"id"`
	State    string `json:"state"`
	UpdateAt string `json:"updated_at"`
}

type Query struct {
	ID                int `json:"id"`
	LastQueryResultID int `json:"latest_query_data_id"`
}

type QueryResult struct {
	Detail QueryResultDetail `json:"query_result"`
}

type QueryResultDetail struct {
	ID          int    `json:"id"`
	RetrievedAt string `json:"retrieved_at"`
}

var (
	client          = &http.Client{}
	alertAPI        = fmt.Sprintf("%s/alerts", *RedashAPIBaseURL)
	queryAPI        = fmt.Sprintf("%s/queries", *RedashAPIBaseURL)
	queryResualtAPI = fmt.Sprintf("%s/query_results", *RedashAPIBaseURL)
)

var (
	err                     error
	ErrGetQueryFailed       = errors.New("get query failed")
	ErrGetQueryResultFailed = errors.New("get query result failed")
	ErrGetAlertFailed       = errors.New("get alert failed")
)

func getQuery(id int) (Query, error) {
	url := fmt.Sprintf("%s/%d", queryAPI, id)

	bodyBytes, err := redashRequest("GET", url)
	logPanic(err)

	query := Query{}
	err = json.Unmarshal(bodyBytes, &query)
	logPanic(err)

	if query.ID != id {
		fmt.Printf("query id %d not equals to %d.\n", query.ID, id)
		err = ErrGetQueryFailed
	}
	return query, err
}

func getQueryResult(id int) (QueryResultDetail, error) {
	url := fmt.Sprintf("%s/%d", queryResualtAPI, id)

	bodyBytes, err := redashRequest("GET", url)
	logPanic(err)

	result := QueryResult{}
	err = json.Unmarshal(bodyBytes, &result)
	logPanic(err)

	if result.Detail.ID != id {
		fmt.Printf("query result id %d not equals to %d.\n", result.Detail.ID, id)
		err = ErrGetQueryResultFailed
	}
	return result.Detail, err
}

func getAlert(id int) (Alert, error) {
	url := fmt.Sprintf("%s/%d", alertAPI, id)

	bodyBytes, err := redashRequest("GET", url)
	logPanic(err)

	alert := Alert{}
	err = json.Unmarshal(bodyBytes, &alert)
	logPanic(err)

	if alert.ID != id {
		fmt.Printf("alert id %d not equals to %d.\n", alert.ID, id)
		err = ErrGetAlertFailed
	}
	return alert, err
}

func isQueryResultFresh(executedAt string) (bool, error) {
	executedAtTime, _ := time.Parse(time.RFC3339, executedAt)
	var err error
	now := time.Now().UTC()
	lastInterval := now.Add(-10 * time.Minute)
	if executedAtTime.Before(lastInterval) {
		return false, errors.New("alert out of date: 5 mins older")
	}
	return true, err
}

func isAlertTriggered(status string) bool {
	if status != "triggered" {
		return false
	}
	return true
}

func redashRequest(method string, url string) ([]byte, error) {
	req, err := http.NewRequest(method, url, nil)
	logPanic(err)

	authHeaderValue := fmt.Sprintf("Key %s", *RedashAPIKey)
	req.Header.Add("Authorization", authHeaderValue)

	resp, err := client.Do(req)
	logPanic(err)
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	logPanic(err)

	return bytes, err
}
