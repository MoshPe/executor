package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
	"io"
	"net/http"
	"testing"
)

func TestImageListError(t *testing.T) {
	c,createErr := client.NewClientWithOpts(client.WithHTTPClient(newMockClient(errorMock(http.StatusInternalServerError, "Server error"))))
	if createErr != nil {
		t.Fatal(createErr)
	}
	_, err := c.ImageList(context.Background(), types.ImageListOptions{})
	if !errdefs.IsSystem(err) {
		t.Fatalf("expected a Server Error, got %[1]T: %[1]v", err)
	}
}

func TestImageList(t *testing.T) {
	//expectedURL := "/images/json"

	noDanglingfilters := filters.NewArgs()
	noDanglingfilters.Add("dangling", "false")

	filters := filters.NewArgs()
	filters.Add("label", "label1")
	filters.Add("label", "label2")
	filters.Add("dangling", "true")

	listCases := []struct {
		options             types.ImageListOptions
		expectedQueryParams map[string]string
	}{
		{
			options: types.ImageListOptions{},
			expectedQueryParams: map[string]string{
				"all":     "",
				"filter":  "",
				"filters": "",
			},
		},
		{
			options: types.ImageListOptions{
				Filters: filters,
			},
			expectedQueryParams: map[string]string{
				"all":     "",
				"filter":  "",
				"filters": `{"dangling":{"true":true},"label":{"label1":true,"label2":true}}`,
			},
		},
		{
			options: types.ImageListOptions{
				Filters: noDanglingfilters,
			},
			expectedQueryParams: map[string]string{
				"all":     "",
				"filter":  "",
				"filters": `{"dangling":{"false":true}}`,
			},
		},
	}
	for _, listCase := range listCases {
			c,createErr := client.NewClientWithOpts(client.WithHTTPClient(newMockClient(func(req *http.Request) (*http.Response, error) {
				query := req.URL.Query()
				for key, expected := range listCase.expectedQueryParams {
					actual := query.Get(key)
					if actual != expected {
						return nil, fmt.Errorf("%s not set in URL query properly. Expected '%s', got %s", key, expected, actual)
					}
				}
				content, err := json.Marshal([]types.ImageSummary{
					{
						ID: "image_id2",
					},
					{
						ID: "image_id2",
					},
				})
				if err != nil {
					return nil, err
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(content)),
				}, nil
			})))
			if createErr != nil {
				t.Fatal(createErr)
			}
			images, err := c.ImageList(context.Background(), listCase.options)
			if err != nil {
				t.Fatal(err)
			}
			if len(images) != 2 {
				t.Fatalf("expected 2 images, got %v", images)
			}
		}
	}

func TestImageListApiBefore125(t *testing.T) {
	expectedFilter := "image:tag"
	c,createErr := client.NewClientWithOpts(client.WithHTTPClient(newMockClient(func(req *http.Request) (*http.Response, error)  {
			query := req.URL.Query()
			actualFilter := query.Get("filters")
			if actualFilter != expectedFilter {
				return nil, fmt.Errorf("filter not set in URL query properly. Expected '%s', got %s", expectedFilter, actualFilter)
			}
			actualFilters := query.Get("filters")
			if actualFilters != "" {
				return nil, fmt.Errorf("filters should have not been present, were with value: %s", actualFilters)
			}
			content, err := json.Marshal([]types.ImageSummary{
				{
					ID: "image_id2",
				},
				{
					ID: "image_id2",
				},
			})
			if err != nil {
				return nil, err
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(content)),
			}, nil
		})))
	if createErr != nil{
		t.Fatal(createErr)
	}

	filters := filters.NewArgs()
	filters.Add("reference", "image:tag")

	options := types.ImageListOptions{
		Filters: filters,
	}

	images, err := c.ImageList(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}
	if len(images) != 2 {
		t.Fatalf("expected 2 images, got %v", images)
	}
}