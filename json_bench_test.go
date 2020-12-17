package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type JSON struct {
	Comment  string `json:"comment"`
	Name     string `json:"name"`
	Optional bool   `json:"optional"`
	Type     string `json:"type"`
}

const (
	j = `{
   "name":"id",
   "type":"UUID",
   "optional":true,
   "comment":"unique identifier of the domain. generated on create, never reused"
}`
)

func unmarshal() error {
	var target JSON
	body, err := ioutil.ReadAll(strings.NewReader(j))
	if err != nil {
		return err
	}
	return json.Unmarshal(body, &target)
}

func decode() error {
	var target JSON
	return json.NewDecoder(strings.NewReader(j)).Decode(&target)
}

func BenchmarkHTTPUnmarshal(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		data, readErr := ioutil.ReadAll(r.Body)
		if readErr != nil {
			return
		}
		_ = r.Body.Close()

		var target JSON
		if jsonErr := json.Unmarshal(data, &target); jsonErr != nil {
			return
		}
	}))

	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		req, err := http.NewRequest(http.MethodPost, server.URL, strings.NewReader(j))
		if err != nil {
			b.Error(err)
		}

		_, respErr := http.DefaultClient.Do(req)
		if respErr != nil {
			b.Error(respErr)
		}
	}
	b.StopTimer()
}

func BenchmarkHTTPJSONDecode(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var target JSON
		if err := json.NewDecoder(r.Body).Decode(&target); err != nil {
			b.Error(err)
		}
		_ = r.Body.Close()
	}))

	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		req, err := http.NewRequest(http.MethodPost, server.URL, strings.NewReader(j))
		if err != nil {
			b.Error(err)
		}

		_, respErr := http.DefaultClient.Do(req)
		if respErr != nil {
			b.Error(respErr)
		}
	}
	b.StopTimer()
}

func BenchmarkJSONUnmarshal(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		err := unmarshal()
		if err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
}

func BenchmarkJSONDecode(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		err := decode()
		if err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
}

func BenchmarkJSONUnmarshalParallel(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := unmarshal()
			if err != nil {
				b.Error(err)
			}

		}
	})
	b.StopTimer()
}

func BenchmarkJSONDecodeParallel(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := decode()
			if err != nil {
				b.Error(err)
			}
		}
	})
	b.StopTimer()
}
