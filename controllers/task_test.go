package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func Test_Task_Import(t *testing.T) {
	/*
	 */
	cases := []struct {
		path          string
		requestMethod string
		requestBody   []byte
	}{
		{
			"http://127.0.0.1:8080/api/task/import",
			"POST",
			[]byte(`{"configPath" : "examples/v2/example.yaml"}`),
		},
	}
	for _, tc := range cases {
		var Response Response
		req, err := http.NewRequest(tc.requestMethod, tc.path, bytes.NewBuffer(tc.requestBody))
		req.Header.Set("Content-Type", "application/json")

		if err != nil {
			log.Fatal(err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal([]byte(body), &Response)
		if Response.Code != -1 && Response.Code != 0 {
			t.Fail()
		}
	}
}

func Test_Task_Action(t *testing.T) {
	/*
	 */
	cases := []struct {
		path          string
		requestMethod string
		requestBody   []byte
	}{
		{
			"http://127.0.0.1:8080/api/task/action",
			"POST",
			[]byte(`{"taskID" : "0", "taskAction": "stop"}`),
		},
		{
			"http://127.0.0.1:8080/api/task/action",
			"POST",
			[]byte(`{"taskAction": "stopAll"}`),
		},
		{
			"http://127.0.0.1:8080/api/task/action",
			"POST",
			[]byte(`{"taskID" : "0", "taskAction": "query"}`),
		},
		{
			"http://127.0.0.1:8080/api/task/action",
			"POST",
			[]byte(`{"taskAction": "queryAll"}`),
		},
	}
	for _, tc := range cases {
		var Response Response
		req, err := http.NewRequest(tc.requestMethod, tc.path, bytes.NewBuffer(tc.requestBody))
		req.Header.Set("Content-Type", "application/json")

		if err != nil {
			log.Fatal(err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal([]byte(body), &Response)
		if Response.Code != -1 && Response.Code != 0 {
			t.Fail()
		}
	}
}
