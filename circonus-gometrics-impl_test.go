// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package circonusgometrics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testServer() *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		// fmt.Printf("%s %s\n", r.Method, r.URL.String())
		switch r.URL.Path {
		case "/metrics_endpoint": // submit metrics
			switch r.Method {
			case "POST":
				fallthrough
			case "PUT":
				defer r.Body.Close()
				b, err := ioutil.ReadAll(r.Body)
				if err != nil {
					panic(err)
				}
				var ret []byte
				var r interface{}
				err = json.Unmarshal(b, &r)
				if err != nil {
					ret, err = json.Marshal(err)
					if err != nil {
						panic(err)
					}
				} else {
					ret, err = json.Marshal(r)
					if err != nil {
						panic(err)
					}
				}
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, string(ret))
			default:
				w.WriteHeader(500)
				fmt.Fprintln(w, "unsupported method")
			}
		default:
			msg := fmt.Sprintf("not found %s", r.URL.Path)
			w.WriteHeader(404)
			fmt.Fprintln(w, msg)
		}
	}

	return httptest.NewServer(http.HandlerFunc(f))
}

func TestNewImpl(t *testing.T) {
	t.Log("no API token, submission URL only")
	cfg := &Config{}
	cfg.CheckManager.Check.SubmissionURL = "http://127.0.0.1:56104/blah/blah"

	cm, err := New(cfg)
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	trap, err := cm.check.GetTrap()
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	if trap.URL.String() != cfg.CheckManager.Check.SubmissionURL {
		t.Fatalf("Expected '%s' == '%s'", trap.URL.String(), cfg.CheckManager.Check.SubmissionURL)
	}

	trap, err = cm.check.GetTrap()
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	if trap.URL.String() != cfg.CheckManager.Check.SubmissionURL {
		t.Fatalf("Expected '%s' == '%s'", trap.URL.String(), cfg.CheckManager.Check.SubmissionURL)
	}
}

func TestFlush(t *testing.T) {
	server := testServer()
	defer server.Close()

	submissionURL := server.URL + "/metrics_endpoint"

	t.Log("Already flushing")
	{
		cfg := &Config{}
		cfg.CheckManager.Check.SubmissionURL = submissionURL
		cm, err := New(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}

		cm.flushing = true
		cm.Flush()
	}

	t.Log("No metrics")
	{
		cfg := &Config{}
		cfg.CheckManager.Check.SubmissionURL = submissionURL
		cm, err := New(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}

		cm.Flush()
	}

	t.Log("counter")
	{
		cfg := &Config{}
		cfg.CheckManager.Check.SubmissionURL = submissionURL
		cm, err := New(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}

		cm.Set("foo", 30)

		cm.Flush()
	}

	t.Log("gauge")
	{
		cfg := &Config{}
		cfg.CheckManager.Check.SubmissionURL = submissionURL
		cm, err := New(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}

		cm.SetGauge("foo", 30)

		cm.Flush()
	}

	t.Log("histogram")
	{
		cfg := &Config{}
		cfg.CheckManager.Check.SubmissionURL = submissionURL
		cm, err := New(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}

		cm.Timing("foo", 30.28)

		cm.Flush()
	}

	t.Log("text")
	{
		cfg := &Config{}
		cfg.CheckManager.Check.SubmissionURL = submissionURL
		cm, err := New(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}

		cm.SetText("foo", "bar")

		cm.Flush()
	}
}
