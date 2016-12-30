// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package circonusgometrics_test

import (
	"errors"
	"testing"

	cgm "github.com/circonus-labs/circonus-gometrics"
)

func TestNew(t *testing.T) {

	t.Log("invalid config (none)")
	{
		expectedError := errors.New("invalid configuration (nil)")
		_, err := cgm.New(nil)
		if err == nil || err.Error() != expectedError.Error() {
			t.Fatalf("Expected an '%#v' error, got '%#v'", expectedError, err)
		}
	}

	t.Log("no API token, no submission URL")
	{
		cfg := &cgm.Config{}
		expectedError := errors.New("invalid check manager configuration (no API token AND no submission url)")
		_, err := cgm.New(cfg)
		if err == nil || err.Error() != expectedError.Error() {
			t.Fatalf("Expected an '%#v' error, got '%#v'", expectedError, err)
		}
	}

	t.Log("no Log, Debug = true")
	{
		cfg := &cgm.Config{
			Debug: true,
		}
		cfg.CheckManager.Check.SubmissionURL = "http://127.0.0.1:56104/blah/blah"
		_, err := cgm.New(cfg)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}
	}

	t.Log("flush interval [good]")
	{
		cfg := &cgm.Config{
			Interval: "30s",
		}
		cfg.CheckManager.Check.SubmissionURL = "http://127.0.0.1:56104/blah/blah"
		_, err := cgm.New(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}
	}
	t.Log("flush interval [bad]")
	{
		cfg := &cgm.Config{
			Interval: "thirty seconds",
		}
		expectedError := errors.New("time: invalid duration thirty seconds")
		_, err := cgm.New(cfg)
		if err == nil {
			t.Fatal("expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %v got '%v'", expectedError, err)
		}
	}

	t.Log("reset counters [good(true)]")
	{
		cfg := &cgm.Config{
			ResetCounters: "true",
		}
		cfg.CheckManager.Check.SubmissionURL = "http://127.0.0.1:56104/blah/blah"
		_, err := cgm.New(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}
	}
	t.Log("reset counters [good(1)]")
	{
		cfg := &cgm.Config{
			ResetCounters: "1",
		}
		cfg.CheckManager.Check.SubmissionURL = "http://127.0.0.1:56104/blah/blah"
		_, err := cgm.New(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}
	}
	t.Log("reset counters [bad(yes)]")
	{
		cfg := &cgm.Config{
			ResetCounters: "yes",
		}
		expectedError := errors.New("strconv.ParseBool: parsing \"yes\": invalid syntax")
		_, err := cgm.New(cfg)
		if err == nil {
			t.Fatal("expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %v got '%v'", expectedError, err)
		}
	}

	t.Log("reset gauges [good(true)]")
	{
		cfg := &cgm.Config{
			ResetGauges: "true",
		}
		cfg.CheckManager.Check.SubmissionURL = "http://127.0.0.1:56104/blah/blah"
		_, err := cgm.New(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}
	}
	t.Log("reset gauges [good(1)]")
	{
		cfg := &cgm.Config{
			ResetGauges: "1",
		}
		cfg.CheckManager.Check.SubmissionURL = "http://127.0.0.1:56104/blah/blah"
		_, err := cgm.New(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}
	}
	t.Log("reset gauges [bad(yes)]")
	{
		cfg := &cgm.Config{
			ResetGauges: "yes",
		}
		expectedError := errors.New("strconv.ParseBool: parsing \"yes\": invalid syntax")
		_, err := cgm.New(cfg)
		if err == nil {
			t.Fatal("expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %v got '%v'", expectedError, err)
		}
	}

	t.Log("reset histograms [good(true)]")
	{
		cfg := &cgm.Config{
			ResetHistograms: "true",
		}
		cfg.CheckManager.Check.SubmissionURL = "http://127.0.0.1:56104/blah/blah"
		_, err := cgm.New(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}
	}
	t.Log("reset histograms [good(1)]")
	{
		cfg := &cgm.Config{
			ResetHistograms: "1",
		}
		cfg.CheckManager.Check.SubmissionURL = "http://127.0.0.1:56104/blah/blah"
		_, err := cgm.New(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}
	}
	t.Log("reset histograms [bad(yes)]")
	{
		cfg := &cgm.Config{
			ResetHistograms: "yes",
		}
		expectedError := errors.New("strconv.ParseBool: parsing \"yes\": invalid syntax")
		_, err := cgm.New(cfg)
		if err == nil {
			t.Fatal("expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %v got '%v'", expectedError, err)
		}
	}

	t.Log("reset text metrics [good(true)]")
	{
		cfg := &cgm.Config{
			ResetText: "true",
		}
		cfg.CheckManager.Check.SubmissionURL = "http://127.0.0.1:56104/blah/blah"
		_, err := cgm.New(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}
	}
	t.Log("reset text metrics [good(1)]")
	{
		cfg := &cgm.Config{
			ResetText: "1",
		}
		cfg.CheckManager.Check.SubmissionURL = "http://127.0.0.1:56104/blah/blah"
		_, err := cgm.New(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}
	}
	t.Log("reset text metrics [bad(yes)]")
	{
		cfg := &cgm.Config{
			ResetText: "yes",
		}
		expectedError := errors.New("strconv.ParseBool: parsing \"yes\": invalid syntax")
		_, err := cgm.New(cfg)
		if err == nil {
			t.Fatal("expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %v got '%v'", expectedError, err)
		}
	}
}
