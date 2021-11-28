package cmd // import "github.com/docker/docker/client"

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
)

func TestImageHistoryError(t *testing.T) {
	c,createErr := client.NewClientWithOpts(client.WithHTTPClient(newMockClient(errorMock(http.StatusInternalServerError, "Server error"))))
	if createErr != nil {
		t.Fatal(createErr)
	}
	_, err := c.ImageHistory(context.Background(), "nothing")
	if !errdefs.IsSystem(err) {
		t.Fatalf("expected a Server Error, got %[1]T: %[1]v", err)
	}
}

func TestImageHistory(t *testing.T) {
	expectedURL := "/images/image_id/history"
	c,createErr := client.NewClientWithOpts(client.WithHTTPClient(newMockClient(func(r *http.Request) (*http.Response, error) {
		if !strings.HasPrefix(r.URL.Path, expectedURL) {
			return nil, fmt.Errorf("expected URL '%s', got '%s'", expectedURL, r.URL)
		}
		b, err := json.Marshal([]image.HistoryResponseItem{
			{
				ID:   "image_id1",
				Tags: []string{"tag1", "tag2"},
			},
			{
				ID:   "image_id2",
				Tags: []string{"tag1", "tag2"},
			},
		})
		if err != nil {
			return nil, err
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(b)),
		}, nil
	})))
	if createErr != nil {
		t.Fatal(createErr)
	}
	imageHistories, err := c.ImageHistory(context.Background(), "image_id")
	if err != nil {
		t.Fatal(err)
	}
	if len(imageHistories) != 2 {
		t.Fatalf("expected 2 containers, got %v", imageHistories)
	}
}

func newMockClient(doer func(*http.Request) (*http.Response, error)) *http.Client {
	return &http.Client{
		Transport: transportFunc(doer),
	}
}

func errorMock(statusCode int, message string) func(req *http.Request) (*http.Response, error) {
	return func(req *http.Request) (*http.Response, error) {
		header := http.Header{}
		header.Set("Content-Type", "application/json")

		body, err := json.Marshal(&types.ErrorResponse{
			Message: message,
		})
		if err != nil {
			return nil, err
		}

		return &http.Response{
			StatusCode: statusCode,
			Body:       io.NopCloser(bytes.NewReader(body)),
			Header:     header,
		}, nil
	}
}

type transportFunc func(*http.Request) (*http.Response, error)

func (tf transportFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return tf(req)
}

func plainTextErrorMock(statusCode int, message string) func(req *http.Request) (*http.Response, error) {
	return func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: statusCode,
			Body:       io.NopCloser(bytes.NewReader([]byte(message))),
		}, nil
	}
}