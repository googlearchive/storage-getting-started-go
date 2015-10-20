/*
Copyright 2013 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Binary storage-sample creates a new bucket, performs all of its operations
// within that bucket, and then cleans up after itself if nothing fails along the way.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"
)

const (
	// This can be changed to any valid object name.
	objectName = "test-file"
	// This scope allows the application full control over resources in Google Cloud Storage
	scope = storage.DevstorageFullControlScope
)

var (
	projectID  = flag.String("project", "", "Your cloud project ID.")
	bucketName = flag.String("bucket", "", "The name of an existing bucket within your project.")
	fileName   = flag.String("file", "", "The file to upload.")
)

func fatalf(service *storage.Service, errorMessage string, args ...interface{}) {
	restoreOriginalState(service)
	log.Fatalf("Dying with error:\n"+errorMessage, args...)
}

func restoreOriginalState(service *storage.Service) bool {
	// Delete an object from a bucket.
	if err := service.Objects.Delete(*bucketName, objectName).Do(); err != nil {
		// If the object exists but wasn't deleted, the bucket deletion will also fail.
		fmt.Printf("Could not delete object during cleanup: %v\n\n", err)
	} else {
		fmt.Printf("Successfully deleted %s/%s during cleanup.\n\n", *bucketName, objectName)
	}

	// Delete a bucket in the project
	if err := service.Buckets.Delete(*bucketName).Do(); err != nil {
		fmt.Printf("Could not delete bucket during cleanup: %v\n\n", err)
		fmt.Println("WARNING: Final cleanup attempt failed. Original state could not be restored.\n")
		return false
	}

	fmt.Printf("Successfully deleted bucket %s during cleanup.\n\n", *bucketName)
	return true
}

func main() {
	flag.Parse()
	if *bucketName == "" {
		log.Fatalf("Bucket argument is required. See --help.")
	}
	if *projectID == "" {
		log.Fatalf("Project argument is required. See --help.")
	}
	if *fileName == "" {
		log.Fatalf("File argument is required. See --help.")
	}

	// Authentication is provided by the gcloud tool when running locally, and
	// by the associated service account when running on Compute Engine.
	client, err := google.DefaultClient(context.Background(), scope)
	if err != nil {
		log.Fatalf("Unable to get default client: %v", err)
	}
	service, err := storage.New(client)
	if err != nil {
		log.Fatalf("Unable to create storage service: %v", err)
	}

	// If the bucket already exists and the user has access, warn the user, but don't try to create it.
	if _, err := service.Buckets.Get(*bucketName).Do(); err == nil {
		fmt.Printf("Bucket %s already exists - skipping buckets.insert call.", *bucketName)
	} else {
		// Create a bucket.
		if res, err := service.Buckets.Insert(*projectID, &storage.Bucket{Name: *bucketName}).Do(); err == nil {
			fmt.Printf("Created bucket %v at location %v\n\n", res.Name, res.SelfLink)
		} else {
			fatalf(service, "Failed creating bucket %s: %v", *bucketName, err)
		}
	}

	// List all buckets in a project.
	if res, err := service.Buckets.List(*projectID).Do(); err == nil {
		fmt.Println("Buckets:")
		for _, item := range res.Items {
			fmt.Println(item.Id)
		}
		fmt.Println()
	} else {
		fatalf(service, "Buckets.List failed: %v", err)
	}

	// Insert an object into a bucket.
	object := &storage.Object{Name: objectName}
	file, err := os.Open(*fileName)
	if err != nil {
		fatalf(service, "Error opening %q: %v", *fileName, err)
	}
	if res, err := service.Objects.Insert(*bucketName, object).Media(file).Do(); err == nil {
		fmt.Printf("Created object %v at location %v\n\n", res.Name, res.SelfLink)
	} else {
		fatalf(service, "Objects.Insert failed: %v", err)
	}

	// List all objects in a bucket using pagination
	var objects []string
	pageToken := ""
	for {
		call := service.Objects.List(*bucketName)
		if pageToken != "" {
			call = call.PageToken(pageToken)
		}
		res, err := call.Do()
		if err != nil {
			fatalf(service, "Objects.List failed: %v", err)
		}
		for _, object := range res.Items {
			objects = append(objects, object.Name)
		}
		if pageToken = res.NextPageToken; pageToken == "" {
			break
		}
	}

	fmt.Printf("Objects in bucket %v:\n", *bucketName)
	for _, object := range objects {
		fmt.Println(object)
	}
	fmt.Println()

	// Get an object from a bucket.
	if res, err := service.Objects.Get(*bucketName, objectName).Do(); err == nil {
		fmt.Printf("The media download link for %v/%v is %v.\n\n", *bucketName, res.Name, res.MediaLink)
	} else {
		fatalf(service, "Failed to get %s/%s: %s.", *bucketName, objectName, err)
	}

	if !restoreOriginalState(service) {
		os.Exit(1)
	}
}
