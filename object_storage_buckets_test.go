package linodego_test

import (
	"context"

	. "github.com/linode/linodego"

	"testing"
)

var (
	testObjectStorageBucketCreateOpts = ObjectStorageBucketCreateOptions{
		Cluster: "us-east-1",
		Label:   "linodego-test-bucket",
	}
)

func TestCreateObjectStorageBucket(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}
	_, bucket, teardown, err := setupObjectStorageBucket(t, "fixtures/TestCreateObjectStorageBucket")
	defer teardown()

	if err != nil {
		t.Errorf("Error creating Object Storage Bucket, got error %v", err)
	}

	expected := testObjectStorageBucketCreateOpts

	if bucket.Label != expected.Label ||
		bucket.Cluster != expected.Cluster {
		t.Errorf("Object Storage Bucket did not match CreateOptions")
	}
}

func TestGetObjectStorageBucket_missing(t *testing.T) {
	client, bucket, teardown, err := setupObjectStorageBucket(t, "fixtures/TestGetObjectStorageBucket_missing")
	defer teardown()

	sameLabel := bucket.Label
	differentCluster := "us-west-1"

	i, err := client.GetObjectStorageBucket(context.Background(), differentCluster, sameLabel)
	if err == nil {
		t.Errorf("should have received an error requesting a missing objectStorageBucket, got %v", i)
	}
	e, ok := err.(*Error)
	if !ok {
		t.Errorf("should have received an Error requesting a missing objectStorageBucket, got %v", e)
	}

	if e.Code != 404 {
		t.Errorf("should have received a 404 Code requesting a missing objectStorageBucket, got %v", e.Code)
	}
}

func TestGetObjectStorageBucket_found(t *testing.T) {
	client, bucket, teardown, err := setupObjectStorageBucket(t, "fixtures/TestGetObjectStorageBucket_found")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	i, err := client.GetObjectStorageBucket(context.Background(), bucket.Cluster, bucket.Label)
	if err != nil {
		t.Errorf("Error getting objectStorageBucket, expected struct, got %v and error %v", i, err)
	}
	if i.Label != bucket.Label {
		t.Errorf("Expected a specific objectStorageBucket, but got a different one %v", i)
	}
	expected := testObjectStorageBucketCreateOpts

	if bucket.Label != expected.Label ||
		bucket.Cluster != expected.Cluster {
		t.Errorf("Object Storage Bucket did not match CreateOptions")
	}
}
func TestListObjectStorageBuckets(t *testing.T) {
	client, _, teardown, err := setupObjectStorageBucket(t, "fixtures/TestListObjectStorageBucket")
	defer teardown()

	i, err := client.ListObjectStorageBuckets(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing objectStorageBuckets, expected struct, got error %v", err)
	}
	if len(i) == 0 {
		t.Errorf("Expected a list of objectStorageBuckets, but got none %v", i)
	}
}

func setupObjectStorageBucket(t *testing.T, fixturesYaml string) (*Client, *ObjectStorageBucket, func(), error) {
	t.Helper()

	client, fixtureTeardown := createTestClient(t, fixturesYaml)
	bucket, err := client.CreateObjectStorageBucket(context.Background(), testObjectStorageBucketCreateOpts)

	if err != nil {
		t.Errorf("Error creating test Bucket: %s", err)
	}

	teardown := func() {
		if err := client.DeleteObjectStorageBucket(context.Background(), bucket.Cluster, bucket.Label); err != nil {
			if t != nil {
				t.Errorf("Error deleting test Bucket: %s", err)
			}
		}
		fixtureTeardown()
	}

	return client, bucket, teardown, err
}
