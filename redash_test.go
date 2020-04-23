package main

import (
	"fmt"
	"testing"
)

func TestGetResponse(t *testing.T) {
	t.Log("starting process: \n")
	query, _ := getQuery(281)

	t.Log(query.ID, query.LastQueryResultID)

	result, _ := getQueryResult(query.LastQueryResultID)
	t.Log(result.ID, result.RetrievedAt)

	alert, _ := getAlert(42)

	isTriggered := isAlertTriggered(alert.State)
	t.Log(isTriggered)

	if isTriggered {
		fmt.Print("redash alert is triggered")
	} else {
		fmt.Print("redash alert is normal")
	}

	isFresh, err := isQueryResultFresh(result.RetrievedAt)
	logIf(err)

	if isFresh {
		fmt.Print("redash query is fresh")
	} else {
		fmt.Print("redash query is out-of-date")
	}
}
