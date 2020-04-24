package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
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
	ErrQueryOutOfDate       = errors.New("query result is out-of-date")
)

func getQuery(id int) (Query, error) {
	url := fmt.Sprintf("%s/%d", queryAPI, id)

	bodyBytes, err := redashRequest("GET", url)
	logIf(err)

	query := Query{}
	err = json.Unmarshal(bodyBytes, &query)
	logIf(err)

	if query.ID != id {
		log.Errorf("query id %d not equals to %d.\n", query.ID, id)
		err = ErrGetQueryFailed
	}
	return query, err
}

func getQueryResult(id int) (QueryResultDetail, error) {
	url := fmt.Sprintf("%s/%d", queryResualtAPI, id)

	bodyBytes, err := redashRequest("GET", url)
	logIf(err)

	result := QueryResult{}
	err = json.Unmarshal(bodyBytes, &result)
	logIf(err)

	if result.Detail.ID != id {
		log.Errorf("query result id %d not equals to %d.\n", result.Detail.ID, id)
		err = ErrGetQueryResultFailed
	}
	return result.Detail, err
}

func getAlert(id int) (Alert, error) {
	url := fmt.Sprintf("%s/%d", alertAPI, id)

	bodyBytes, err := redashRequest("GET", url)
	logIf(err)

	alert := Alert{}
	err = json.Unmarshal(bodyBytes, &alert)
	logIf(err)

	if alert.ID != id {
		log.Errorf("alert id %d not equals to %d.\n", alert.ID, id)
		err = ErrGetAlertFailed
	}
	return alert, err
}

func isQueryResultFresh(executedAt string) (bool, error) {
	executedAtTime, err := time.Parse(time.RFC3339, executedAt)
	logIf(err)

	interval, err := time.ParseDuration(*RedashProbeInterval)
	logIf(err)

	now := time.Now().UTC()
	lastExecutedTime := now.Add(-interval)

	if executedAtTime.Before(lastExecutedTime) {
		return false, ErrQueryOutOfDate
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
	logIf(err)

	authHeaderValue := fmt.Sprintf("Key %s", *RedashAPIKey)
	req.Header.Add("Authorization", authHeaderValue)

	resp, err := client.Do(req)
	logIf(err)
	log.Infof("[%d] %s", resp.StatusCode, req.URL)

	defer resp.Body.Close()
	log.Debugf("response body: %s\n", resp.StatusCode, resp.Body)

	bytes, err := ioutil.ReadAll(resp.Body)
	logIf(err)

	return bytes, err
}
