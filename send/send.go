package send

import (
	"bytes"
	"io"
	"net/http"
)

type HTTP interface {
	Send()
}

type Request struct {
	Path    string
	Method  string
	Headers map[string]string
	Body    []byte
}

func (r *Request) Send() ([]byte, error) {

	client := &http.Client{}

	req, err := http.NewRequest(r.Method, r.Path, bytes.NewBuffer(r.Body))
	if err != nil {
		return nil, err
	}

	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}