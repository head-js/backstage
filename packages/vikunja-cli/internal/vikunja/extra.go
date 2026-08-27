package vikunja

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// ListProjectTasksByBucket returns the tasks of a project that are in the given kanban bucket.
func (a *Adapter) ListProjectTasksByBucket(projectID, bucketID string) (any, error) {
	projectID, _, err := parseProjectID(projectID)
	if err != nil {
		return nil, err
	}
	_, bucketIDValue, err := parseID("bucket", bucketID)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("filter", "bucket_id = "+strconv.FormatInt(bucketIDValue, 10))
	query.Set("per_page", "1000")
	query.Set("expand", "buckets")
	path := "/projects/" + projectID + "/tasks?" + query.Encode()
	return a.Do(context.Background(), http.MethodGet, path, nil)
}
