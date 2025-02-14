package main_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

const (
	url = "http://localhost:8080/key/test_key"

	statusNotFound   = "404 Not Found"
	statusBadRequest = "400 Bad Request"
	statusCreated    = "201 Created"
	statusOK         = "200 OK"
	statusNoContent  = "204 No Content"

	get    = "GET"
	post   = "POST"
	put    = "PUT"
	delete = "DELETE"
)

func TestGetNotExistingKey(t *testing.T) {
	client := &http.Client{}
	req, _ := http.NewRequest(get, url+"_temp", bytes.NewBufferString(""))
	resp, err := client.Do(req)
	if err != nil || resp == nil || resp.Status != statusNotFound {
		t.Error("TestGetNotExistingKey() failed", err, resp)
	}
	defer resp.Body.Close()
}

func TestStoreEmptyValue(t *testing.T) {
	client := &http.Client{}
	req, _ := http.NewRequest(post, url, bytes.NewBufferString(""))
	resp, err := client.Do(req)
	if err != nil || resp == nil || resp.Status != statusBadRequest {
		t.Error("TestStoreEmptyValue() failed", err, resp)
	}
	defer resp.Body.Close()
}

func TestAppendEmptyValue(t *testing.T) {
	client := &http.Client{}
	req, _ := http.NewRequest(post, url+"?ttl=1", bytes.NewBufferString("hello"))
	resp, err := client.Do(req)
	if err != nil || resp == nil || resp.Status != statusCreated {
		t.Error("TestAppendEmptyValue() failed", err, resp)
	}
	resp.Body.Close()

	req, _ = http.NewRequest(put, url, bytes.NewBufferString(""))
	resp, err = client.Do(req)
	if err != nil || resp == nil || resp.Status != statusBadRequest {
		t.Error("TestAppendEmptyValue() failed", err, resp)
	}
	defer resp.Body.Close()
}

func TestAppendNotExistingKey(t *testing.T) {
	client := &http.Client{}
	req, _ := http.NewRequest(put, url+"_temp", bytes.NewBufferString(" world"))
	resp, err := client.Do(req)
	if err != nil || resp == nil || resp.Status != statusNotFound {
		t.Error("TestAppendNotExistingKey() failed", err, resp)
	}
	defer resp.Body.Close()
}

func TestDeleteNotExistingKey(t *testing.T) {
	client := &http.Client{}
	req, _ := http.NewRequest(delete, url+"_temp", bytes.NewBufferString(""))
	resp, err := client.Do(req)
	if err != nil || resp == nil || resp.Status != statusNotFound {
		t.Error("TestDeleteNotExistingKey() failed", err, resp)
	}
	defer resp.Body.Close()
}

func TestStoreAndGet(t *testing.T) {
	client := &http.Client{}
	req, _ := http.NewRequest(post, url+"?ttl=1", bytes.NewBufferString("hello"))
	resp, err := client.Do(req)
	if err != nil || resp == nil || resp.Status != statusCreated {
		t.Error("TestStoreAndGet() failed", err, resp)
	}
	resp.Body.Close()

	req, _ = http.NewRequest(get, url, bytes.NewBufferString(""))
	resp, err = client.Do(req)
	if err != nil || resp == nil || resp.Status != statusOK {
		t.Error("TestStoreAndGet() failed", err, resp)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil || data == nil || string(data) != "hello" {
		t.Error("TestStoreAndGet() failed", err, resp)
	}
}

func TestStoreAndDelete(t *testing.T) {
	client := &http.Client{}
	req, _ := http.NewRequest(post, url+"?ttl=1", bytes.NewBufferString("hello"))
	resp, err := client.Do(req)
	if err != nil || resp == nil || resp.Status != statusCreated {
		t.Error("TestStoreAndDelete() failed", err, resp)
	}
	resp.Body.Close()

	req, _ = http.NewRequest(delete, url, bytes.NewBufferString(""))
	resp, err = client.Do(req)
	if err != nil || resp == nil || resp.Status != statusNoContent {
		t.Error("TestStoreAndDelete() failed", err, resp)
	}
	resp.Body.Close()

	req, _ = http.NewRequest(get, url, bytes.NewBufferString(""))
	resp, _ = client.Do(req)
	if err != nil || resp == nil || resp.Status != statusNotFound {
		t.Error("TestStoreAndDelete() failed", err, resp)
	}
	defer resp.Body.Close()
}

func TestStoreAndAppend(t *testing.T) {
	client := &http.Client{}
	req, _ := http.NewRequest(post, url+"?ttl=1", bytes.NewBufferString("hello"))
	resp, err := client.Do(req)
	if err != nil || resp == nil || resp.Status != statusCreated {
		t.Error("TestStoreAndAppend() failed", err, resp)
	}
	resp.Body.Close()

	req, _ = http.NewRequest(put, url, bytes.NewBufferString(" world"))
	resp, err = client.Do(req)
	if err != nil || resp == nil || resp.Status != statusOK {
		t.Error("TestStoreAndAppend() failed", err, resp)
	}
	resp.Body.Close()

	req, _ = http.NewRequest(get, url, bytes.NewBufferString(""))
	resp, err = client.Do(req)
	if err != nil || resp == nil || resp.Status != statusOK {
		t.Error("TestStoreAndAppend() failed", err, resp)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil || data == nil || string(data) != "hello world" {
		t.Error("TestStoreAndAppend() failed", err, resp)
	}
}

func TestTTL(t *testing.T) {
	client := &http.Client{}
	req, _ := http.NewRequest(post, url+"?ttl=1", bytes.NewBufferString("hello"))
	resp, err := client.Do(req)
	if err != nil || resp == nil || resp.Status != statusCreated {
		t.Error("TestTTL() failed", err, resp)
	}
	resp.Body.Close()

	time.Sleep(2 * time.Second)

	req, _ = http.NewRequest(get, url, bytes.NewBufferString(""))
	resp, err = client.Do(req)
	if err != nil || resp == nil || resp.Status != statusNotFound {
		t.Error("TestTTL() failed", err, resp)
	}
	defer resp.Body.Close()
}
