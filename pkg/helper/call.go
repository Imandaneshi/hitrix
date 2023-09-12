package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"strings"
	"time"
)

// Call helper for api calls
// TODO : make service and use debug service to
func Call(ctx context.Context,
	method,
	url string,
	headers map[string]string,
	timeout time.Duration,
	payload interface{},
	cookies []*http.Cookie) ([]byte, http.Header, int, error) {
	var d []byte
	var e error

	switch v := payload.(type) {
	case string:
		d = []byte(v)
	default:
		d, e = json.Marshal(v)
		if e != nil {
			panic(e)
		}
	}

	var b io.Reader
	b = bytes.NewReader(d)

	method = strings.ToUpper(method)
	if StringInArray(method, "GET", "DELETE") {
		b = nil
	}

	r, err := http.NewRequest(method, url, b)
	if err != nil {
		return nil, nil, 0, err
	}

	for i := range headers {
		r.Header.Set(i, headers[i])
	}

	for i := range cookies {
		r.AddCookie(cookies[i])
	}

	nCtx, cnl := context.WithTimeout(ctx, timeout)
	defer cnl()

	resp, err := http.DefaultClient.Do(r.WithContext(nCtx))
	if err != nil {
		return nil, nil, 0, err
	}

	data, err := io.ReadAll(resp.Body)

	defer func() {
		_ = resp.Body.Close()
	}()

	if err != nil {
		return nil, nil, resp.StatusCode, err
	}

	return data, resp.Header, resp.StatusCode, nil
}

func CallXML(ctx context.Context,
	method,
	url string,
	headers map[string]string,
	timeout time.Duration,
	payload interface{},
	cookies []*http.Cookie) ([]byte, http.Header, int, error) {
	d, err := xml.Marshal(payload)
	if err != nil {
		panic(err)
	}

	var b io.Reader
	b = bytes.NewReader(d)

	method = strings.ToUpper(method)
	if StringInArray(method, "GET", "DELETE") {
		b = nil
	}

	r, err := http.NewRequest(method, url, b)
	if err != nil {
		return nil, nil, 0, err
	}

	for i := range headers {
		r.Header.Set(i, headers[i])
	}

	for i := range cookies {
		r.AddCookie(cookies[i])
	}

	nCtx, cnl := context.WithTimeout(ctx, timeout)
	defer cnl()

	resp, err := http.DefaultClient.Do(r.WithContext(nCtx))
	if err != nil {
		return nil, nil, 0, err
	}

	data, err := io.ReadAll(resp.Body)

	defer func() {
		_ = resp.Body.Close()
	}()

	if err != nil {
		return nil, nil, resp.StatusCode, err
	}

	return data, resp.Header, resp.StatusCode, nil
}
