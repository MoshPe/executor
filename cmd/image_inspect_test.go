package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestImageInspectError(t *testing.T) {
	c,createErr := client.NewClientWithOpts(client.WithHTTPClient(newMockClient(errorMock(http.StatusInternalServerError, "Server error"))))
	if createErr != nil {
		t.Fatal(createErr)
	}
	_, _, err := c.ImageInspectWithRaw(context.Background(), "nothing")
	if !errdefs.IsSystem(err) {
		t.Fatalf("expected a Server Error, got %[1]T: %[1]v", err)
	}
}

func TestImageInspectImageNotFound(t *testing.T) {
	c,createErr := client.NewClientWithOpts(client.WithHTTPClient(newMockClient(errorMock(http.StatusNotFound, "Server error"))))
	if createErr != nil {
		t.Fatal(createErr)
	}
	_, _, err := c.ImageInspectWithRaw(context.Background(), "unknown")
	if err == nil{
		t.Fatalf("expected an imageNotFound error, got %v", err)
	}
}

func TestImageInspectWithEmptyID(t *testing.T) {
	c,createErr := client.NewClientWithOpts(client.WithHTTPClient(newMockClient(func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("should not make request")
		})))
	if createErr != nil {
		t.Fatal(createErr)
	}
	_, _, err := c.ImageInspectWithRaw(context.Background(), "")
	if err == nil {
		t.Fatalf("Expected NotFoundError, got %v", err)
	}
}

func TestImageInspect(t *testing.T) {
	//expectedURL := "/images/image_id/json"
	expectedTags := []string{"tag1", "tag2"}
	c,createErr := client.NewClientWithOpts(client.WithHTTPClient(newMockClient(func(req *http.Request) (*http.Response, error) {
		content, err := json.Marshal(types.ImageInspect{
			ID:       "image_id",
			RepoTags: expectedTags,
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

	imageInspect, _, err := c.ImageInspectWithRaw(context.Background(), "image_id")
	if err != nil {
		t.Fatal(err)
	}
	if imageInspect.ID != "image_id" {
		t.Fatalf("expected `image_id`, got %s", imageInspect.ID)
	}
	if !reflect.DeepEqual(imageInspect.RepoTags, expectedTags) {
		t.Fatalf("expected `%v`, got %v", expectedTags, imageInspect.RepoTags)
	}
}
