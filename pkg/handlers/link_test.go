package handlers

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestEqResp struct {
	status   int
	respBody string
}

type TestStructMemory struct {
	name         string
	function     func(http.ResponseWriter, *http.Request)
	req          *http.Request
	res          TestEqResp
	resStr       string
	resStatusErr error
}

func TestLinkRepoMemory(t *testing.T) {
	h := NewLinkHandler(nil, "inMemory")

	req1, err := http.NewRequest("GET", "/api/getLink", nil)
	// url.Values{"url": {"link"}}
	if err != nil {
		log.Fatal(err)
	}
	req2, err := http.NewRequest("POST", "/api/newLink", bytes.NewBuffer([]byte(`{"url":"link"}`)))
	if err != nil {
		log.Fatal(err)
	}

	//req3, err := http.NewRequest("POST", "/api/newLink", bytes.NewBuffer([]byte("1")))
	//if err != nil {
	//	log.Fatal(err)
	//}
	req4, err := http.NewRequest("POST", "/api/newLink", bytes.NewBuffer([]byte(`{`)))
	if err != nil {
		log.Fatal(err)
	}
	req5, err := http.NewRequest("GET", "/api/getLink?url=4f0aa52d65", nil)
	if err != nil {
		log.Fatal(err)
	}

	tests := []TestStructMemory{
		{
			name:     "bad link get link",
			function: h.GetLink,
			req:      req1,
			res: TestEqResp{
				status:   500,
				respBody: `{"error":"bad link"}`,
			},
		},
		{
			name:     "ok test new link",
			function: h.NewLink,
			req:      req2,
			res: TestEqResp{
				status:   200,
				respBody: `{"URL":"4f0aa52d65"}`,
			},
		},
		//{
		//	name:     "wrong read body new link",
		//	function: h.NewLink,
		//	req:      req3,
		//	res: TestEqResp{
		//		status:   500,
		//		respBody: `{"error":"some bad in read data"}`,
		//	},
		//},
		{
			name:     "can't unmarshal json-body new link",
			function: h.NewLink,
			req:      req4,
			res: TestEqResp{
				status:   500,
				respBody: `{"error":"some bad in unmarshal data"}`,
			},
		},
		{
			name:     "ok get link",
			function: h.GetLink,
			req:      req5,
			res: TestEqResp{
				status:   200,
				respBody: `{"URL":"link"}`,
			},
		},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		test.function(w, test.req)
		fmt.Println("code = ", w.Code)

		respBody, err := ioutil.ReadAll(w.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("body = ", string(respBody))

		assert.Equal(t, test.res, TestEqResp{status: w.Code, respBody: string(respBody)})
	}
}
