package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateSummaryFields(t *testing.T) {
	summary := summaryInfo{Requested: 10, Responded: 9}
	timeTaken := 3 * time.Second
	totalTransferred := int64(30720) // 30 KB
	totalTimeAllRequests := 27 * time.Second

	result := createSummary("https://example.com:8080/path", timeTaken, totalTransferred, totalTimeAllRequests, 5, 1024, summary)

	var s summaryInfo
	if err := json.Unmarshal([]byte(result), &s); err != nil {
		t.Fatalf("failed to unmarshal summary: %v", err)
	}

	if s.Hostname != "example.com" {
		t.Errorf("Hostname: got %q, want %q", s.Hostname, "example.com")
	}
	if s.Port != "8080" {
		t.Errorf("Port: got %q, want %q", s.Port, "8080")
	}
	if s.ConcurrencyLevel != 5 {
		t.Errorf("ConcurrencyLevel: got %d, want %d", s.ConcurrencyLevel, 5)
	}
	if s.DocumentLength != 1024 {
		t.Errorf("DocumentLength: got %d, want %d", s.DocumentLength, 1024)
	}
	if s.CompletedRequests != 9 {
		t.Errorf("CompletedRequests: got %d, want %d", s.CompletedRequests, 9)
	}
	if s.FailedRequests != 1 {
		t.Errorf("FailedRequests: got %d, want %d", s.FailedRequests, 1)
	}
	if s.Rps != 3 { // 9 / 3s = 3
		t.Errorf("Rps: got %d, want %d", s.Rps, 3)
	}
	if s.TransferRate != 10 { // 30720 / 3s / 1024 = 10 KB/s
		t.Errorf("TransferRate: got %d, want %d", s.TransferRate, 10)
	}
}

func TestCreateSummaryNoPort(t *testing.T) {
	summary := summaryInfo{Requested: 1, Responded: 1}
	result := createSummary("https://example.com/", time.Second, 1024, time.Second, 1, 512, summary)

	var s summaryInfo
	if err := json.Unmarshal([]byte(result), &s); err != nil {
		t.Fatalf("failed to unmarshal summary: %v", err)
	}
	if s.Port != "" {
		t.Errorf("Port: got %q, want empty string for implicit port", s.Port)
	}
}

func TestCreateSummaryZeroResponded(t *testing.T) {
	// TimePerRequest should not panic when Responded == 0
	summary := summaryInfo{Requested: 5, Responded: 0}
	result := createSummary("https://example.com/", time.Second, 0, 0, 1, 0, summary)
	if result == "" {
		t.Error("expected non-empty result even with zero responses")
	}
}

func TestCheckLinkSuccess(t *testing.T) {
	body := "hello world"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, body)
	}))
	defer server.Close()

	c := make(chan responseInfo, 1)
	checkLink(context.Background(), &http.Client{}, server.URL, c)

	resp := <-c
	if resp.status != http.StatusOK {
		t.Errorf("status: got %d, want %d", resp.status, http.StatusOK)
	}
	if resp.bytes != int64(len(body)) {
		t.Errorf("bytes: got %d, want %d", resp.bytes, len(body))
	}
	if resp.duration <= 0 {
		t.Error("expected positive duration")
	}
}

func TestCheckLinkNonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	c := make(chan responseInfo, 1)
	checkLink(context.Background(), &http.Client{}, server.URL, c)

	resp := <-c
	if resp.status != http.StatusNotFound {
		t.Errorf("status: got %d, want %d", resp.status, http.StatusNotFound)
	}
}

func TestCheckLinkBadURL(t *testing.T) {
	c := make(chan responseInfo, 1)
	checkLink(context.Background(), &http.Client{}, "http://localhost:0", c)

	resp := <-c
	if resp.status != -1 {
		t.Errorf("status: got %d, want -1 for unreachable host", resp.status)
	}
}

func TestCheckLinkCancelledContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel before sending

	c := make(chan responseInfo, 1)
	checkLink(ctx, &http.Client{}, server.URL, c)

	resp := <-c
	if resp.status != -1 {
		t.Errorf("status: got %d, want -1 for cancelled context", resp.status)
	}
}

func TestCheckLinkTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &http.Client{Timeout: 50 * time.Millisecond}
	c := make(chan responseInfo, 1)
	checkLink(context.Background(), client, server.URL, c)

	resp := <-c
	if resp.status != -1 {
		t.Errorf("status: got %d, want -1 for timed-out request", resp.status)
	}
}
