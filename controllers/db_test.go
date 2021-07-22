package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestDBConnect(t *testing.T) {
	var Response Response
	cases := []struct {
		path          string
		requestMethod string
		requestBody   []byte
	}{
		{
			"http://127.0.0.1:8080/api/db/connect",
			"POST",
			[]byte(`{"username": "root",
					"password": "nebula",
					"host": "127.0.0.1:9669"}`),
		},
		{
			"http://127.0.0.1:8080/api/db/connect",
			"POST",
			[]byte(`{"username": "user1",
					"password": "password",
					"host": "127.0.0.1:9669"}`),
		},
	}

	for _, tc := range cases {
		req, err := http.NewRequest(tc.requestMethod, tc.path, bytes.NewBuffer(tc.requestBody))
		if err != nil {
			t.Fail()
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fail()
		}

		defer req.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fail()
		}

		if err := json.Unmarshal([]byte(body), &Response); err != nil {
			t.Fail()
		}

		if Response.Code != -1 && Response.Code != 0 {
			t.Fail()
		}
	}
}

func TestDBExecute(t *testing.T) {
	cases := []struct {
		path          string
		requestMethod string
		requestBody   []byte
	}{
		{
			"http://127.0.0.1:8080/api/db/exec",
			"POST",
			[]byte(`{"username" : "user",
					"password" : "password",
					"host" : "127.0.0.1:9669",
					"gql" : "SHOW SPACES11;"}`),
		},
	}
	for _, tc := range cases {
		var Response Response
		req, err := http.NewRequest(tc.requestMethod, tc.path, bytes.NewBuffer(tc.requestBody))
		if err != nil {
			t.Fail()
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fail()
		}

		if err := json.Unmarshal([]byte(body), &Response); err != nil {
			t.Fail()
		}

		if Response.Code != -1 && Response.Code != 0 {
			t.Fail()
		}
	}
}
