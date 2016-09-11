package jsonapi

import (
	"testing"

	"github.com/labstack/echo/engine"
	"github.com/stretchr/testify/assert"
)

func TestParseRequestError(t *testing.T) {
	invalidAccept := constructRequest("GET", "")
	invalidAccept.Header().Set("Accept", "foo")

	invalidContentType := constructRequest("GET", "")
	invalidContentType.Header().Set("Content-Type", "foo")

	missingContentType := constructRequest("POST", "foo")

	list := []engine.Request{
		invalidAccept,
		invalidContentType,
		missingContentType,
		constructRequest("PUT", ""),
		constructRequest("GET", ""),
		constructRequest("POST", ""),
		constructRequest("GET", "/"),
		constructRequest("GET", "foo/bar/baz/qux"),
		constructRequest("GET", "foo/bar/baz/qux/quux"),
		constructRequest("GET", "foo?page[number]=bar"),
		constructRequest("GET", "foo?page[size]=bar"),
		constructRequest("GET", "foo?page[number]=1"),
		constructRequest("GET", "foo?page[size]=1"),
		constructRequest("GET", "foo?page[number]=bar&page[number]=baz"),
		constructRequest("GET", "foo?page[size]=bar&page[size]=baz"),
		constructRequest("PATCH", "foo"),
	}

	for _, r := range list {
		req, err := ParseRequest(r, "")
		assert.Error(t, err)
		assert.Nil(t, req)
	}
}

func TestParseRequestPrefix(t *testing.T) {
	list := map[string]string{
		"bar":           "",
		"/bar":          "",
		"foo/bar":       "foo",
		"/foo/bar":      "foo",
		"foo/bar/":      "/foo",
		"/foo/bar/":     "foo/",
		"baz/foo/bar/":  "/baz/foo",
		"/baz/foo/bar/": "baz/foo/",
	}

	for url, prefix := range list {
		r := constructRequest("GET", url)

		req, err := ParseRequest(r, prefix)
		assert.NoError(t, err)
		assert.Equal(t, "bar", req.ResourceType)
	}
}

func TestParseRequestResource(t *testing.T) {
	r := constructRequest("GET", "foo")

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
	}, req)
}

func TestParseRequestResourceID(t *testing.T) {
	r := constructRequest("GET", "foo/1")

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       FindResource,
		ResourceType: "foo",
		ResourceID:   "1",
	}, req)
}

func TestParseRequestRelatedResource(t *testing.T) {
	r := constructRequest("GET", "foo/1/bar")

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:          GetRelatedResources,
		ResourceType:    "foo",
		ResourceID:      "1",
		RelatedResource: "bar",
	}, req)
}

func TestParseRequestRelationship(t *testing.T) {
	r := constructRequest("GET", "foo/1/relationships/bar")

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       GetRelationship,
		ResourceType: "foo",
		ResourceID:   "1",
		Relationship: "bar",
	}, req)
}

func TestParseRequestIntent(t *testing.T) {
	list := []struct {
		method string
		url    string
		intent Intent
		doc    bool
	}{
		{"GET", "/posts", ListResources, false},
		{"GET", "/posts/1", FindResource, false},
		{"POST", "/posts", CreateResource, true},
		{"PATCH", "/posts/1", UpdateResource, true},
		{"DELETE", "/posts/1", DeleteResource, false},
		{"GET", "/posts/1/author", GetRelatedResources, false},
		{"GET", "/posts/1/relationships/author", GetRelationship, false},
		{"PATCH", "/posts/1/relationships/author", SetRelationship, true},
		{"POST", "/posts/1/relationships/comments", AppendToRelationship, true},
		{"DELETE", "/posts/1/relationships/comments", RemoveFromRelationship, true},
	}

	for _, entry := range list {
		r := constructRequest(entry.method, entry.url)
		r.Header().Set("Content-Type", MediaType)

		req, err := ParseRequest(r, "")
		assert.NoError(t, err)
		assert.Equal(t, entry.intent, req.Intent)
		assert.Equal(t, entry.doc, req.Intent.DocumentExpected())
		assert.Equal(t, entry.url, req.Self())
		assert.Equal(t, entry.method, req.Intent.RequestMethod())
	}
}

func TestParseRequestInclude(t *testing.T) {
	r1 := constructRequest("GET", "foo?include=bar,baz")

	req, err := ParseRequest(r1, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
		Include:      []string{"bar", "baz"},
	}, req)

	r2 := constructRequest("GET", "foo?include=bar&include=baz,qux")

	req, err = ParseRequest(r2, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
		Include:      []string{"bar", "baz", "qux"},
	}, req)
}

func TestParseRequestSorting(t *testing.T) {
	r1 := constructRequest("GET", "foo?sort=bar,baz")

	req, err := ParseRequest(r1, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
		Sorting:      []string{"bar", "baz"},
	}, req)

	r2 := constructRequest("GET", "foo?sort=bar&sort=baz,qux")

	req, err = ParseRequest(r2, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
		Sorting:      []string{"bar", "baz", "qux"},
	}, req)
}

func TestParseRequestPage(t *testing.T) {
	r := constructRequest("GET", "foo?page[number]=1&page[size]=2")

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
		PageNumber:   1,
		PageSize:     2,
	}, req)
}

func TestParseRequestFields(t *testing.T) {
	r := constructRequest("GET", "foo?fields[foo]=bar,baz")

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
		Fields: map[string][]string{
			"foo": {"bar", "baz"},
		},
	}, req)
}

func TestParseRequestFilters(t *testing.T) {
	r := constructRequest("GET", "foo?filter[foo]=bar,baz")

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
		Filters: map[string][]string{
			"foo": {"bar", "baz"},
		},
	}, req)
}

func TestZeroIntentRequestMethod(t *testing.T) {
	assert.Empty(t, Intent(0).RequestMethod())
}

func BenchmarkParseRequest(b *testing.B) {
	r := constructRequest("GET", "foo/1")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ParseRequest(r, "")
	}
}

func BenchmarkParseRequestFilterAndSort(b *testing.B) {
	r := constructRequest("GET", "foo/1?sort=bar&filter[baz]=qux")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ParseRequest(r, "")
	}
}
