// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var (
	testContactGroup = ContactGroup{
		CID:               "/contact_group/1234",
		LastModfiedBy:     "/user/1234",
		LastModified:      1483041636,
		AggregationWindow: 300,
		Contacts: ContactGroupContacts{
			External: []ContactGroupContactsExternal{
				ContactGroupContactsExternal{
					Info:   "12125550100",
					Method: "sms",
				},
				ContactGroupContactsExternal{
					Info:   "bert@example.com",
					Method: "xmpp",
				},
				ContactGroupContactsExternal{
					Info:   "ernie@example.com",
					Method: "email",
				},
			},
			Users: []ContactGroupContactsUser{
				ContactGroupContactsUser{
					Info:    "snuffy@example.com",
					Method:  "email",
					UserCID: "/user/1234",
				},
				ContactGroupContactsUser{
					Info:    "12125550199",
					Method:  "sms",
					UserCID: "/user/4567",
				},
			},
		},
		Escalations: []ContactGroupEscalation{
			ContactGroupEscalation{
				After:           900,
				ContactGroupCID: "/contact_group/4567",
			},
			ContactGroupEscalation{},
			ContactGroupEscalation{},
			ContactGroupEscalation{},
			ContactGroupEscalation{},
		},
		Name:      "FooBar",
		Reminders: []int{10, 0, 0, 15, 30},
	}
)

func testContactGroupServer() *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/contact_group/1234" {
			switch r.Method {
			case "GET":
				ret, err := json.Marshal(testContactGroup)
				if err != nil {
					panic(err)
				}
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, string(ret))
			case "PUT":
				defer r.Body.Close()
				b, err := ioutil.ReadAll(r.Body)
				if err != nil {
					panic(err)
				}
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, string(b))
			case "DELETE":
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json")
			default:
				w.WriteHeader(404)
				fmt.Fprintln(w, fmt.Sprintf("not found: %s %s", r.Method, path))
			}
		} else if path == "/contact_group" {
			switch r.Method {
			case "GET":
				c := []ContactGroup{testContactGroup}
				ret, err := json.Marshal(c)
				if err != nil {
					panic(err)
				}
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, string(ret))
			case "POST":
				defer r.Body.Close()
				_, err := ioutil.ReadAll(r.Body)
				if err != nil {
					panic(err)
				}
				ret, err := json.Marshal(testContactGroup)
				if err != nil {
					panic(err)
				}
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, string(ret))
			default:
				w.WriteHeader(404)
				fmt.Fprintln(w, fmt.Sprintf("not found: %s %s", r.Method, path))
			}
		} else {
			w.WriteHeader(404)
			fmt.Fprintln(w, fmt.Sprintf("not found: %s %s", r.Method, path))
		}
	}

	return httptest.NewServer(http.HandlerFunc(f))
}

func TestFetchContactGroup(t *testing.T) {
	server := testContactGroupServer()
	defer server.Close()

	ac := &Config{
		TokenKey: "abc123",
		TokenApp: "test",
		URL:      server.URL,
	}
	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	t.Log("without CID")
	{
		expectedError := errors.New("Invalid contact group CID ")
		_, err := apih.FetchContactGroup(CIDType(""))
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}

	t.Log("with valid CID")
	{
		cid := CIDType("/contact_group/1234")
		contactGroup, err := apih.FetchContactGroup(cid)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(contactGroup)
		expectedType := "*api.ContactGroup"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}

		if contactGroup.CID != testContactGroup.CID {
			t.Fatalf("CIDs do not match: %+v != %+v\n", contactGroup, testContactGroup)
		}
	}

	t.Log("with invalid CID")
	{
		expectedError := errors.New("Invalid contact group CID /invalid")
		_, err := apih.FetchContactGroup(CIDType("/invalid"))
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}
}

func TestFetchContactGroups(t *testing.T) {
	server := testContactGroupServer()
	defer server.Close()

	ac := &Config{
		TokenKey: "abc123",
		TokenApp: "test",
		URL:      server.URL,
	}
	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	contactGroups, err := apih.FetchContactGroups()
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(contactGroups)
	expectedType := "[]api.ContactGroup"
	if actualType.String() != expectedType {
		t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
	}

}

func TestCreateContactGroup(t *testing.T) {
	server := testContactGroupServer()
	defer server.Close()

	var apih *API

	ac := &Config{
		TokenKey: "abc123",
		TokenApp: "test",
		URL:      server.URL,
	}
	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	contactGroup, err := apih.CreateContactGroup(&testContactGroup)
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(contactGroup)
	expectedType := "*api.ContactGroup"
	if actualType.String() != expectedType {
		t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
	}
}

func TestUpdateContactGroup(t *testing.T) {
	server := testContactGroupServer()
	defer server.Close()

	var apih *API

	ac := &Config{
		TokenKey: "abc123",
		TokenApp: "test",
		URL:      server.URL,
	}
	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	t.Log("valid ContactGroup")
	{
		contactGroup, err := apih.UpdateContactGroup(&testContactGroup)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(contactGroup)
		expectedType := "*api.ContactGroup"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("Test with invalid CID")
	{
		expectedError := errors.New("Invalid contact group CID /invalid")
		x := &ContactGroup{CID: "/invalid"}
		_, err := apih.UpdateContactGroup(x)
		if err == nil {
			t.Fatal("Expected an error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}
}

func TestDeleteContactGroup(t *testing.T) {
	server := testContactGroupServer()
	defer server.Close()

	var apih *API

	ac := &Config{
		TokenKey: "abc123",
		TokenApp: "test",
		URL:      server.URL,
	}
	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	t.Log("valid ContactGroup")
	{
		_, err := apih.DeleteContactGroup(&testContactGroup)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}
	}

	t.Log("Test with invalid CID")
	{
		expectedError := errors.New("Invalid contact group CID /invalid")
		x := &ContactGroup{CID: "/invalid"}
		_, err := apih.UpdateContactGroup(x)
		if err == nil {
			t.Fatal("Expected an error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}
}