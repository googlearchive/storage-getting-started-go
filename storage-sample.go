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
	"net/http"
	"os"

	"code.google.com/p/goauth2/oauth"
	"code.google.com/p/google-api-go-client/storage/v1beta2"
)

const (
	// Change these variable to match your personal information.
	bucketName   = "YOUR_BUCKET_NAME"
	projectID    = "YOUR_PROJECT_ID"
	clientId     = "YOUR_CLIENT_ID"
	clientSecret = "YOUR_CLIENT_SECRET"

	fileName   = "/usr/share/dict/words" // The name of the local file to upload.
	objectName = "english-dictionary"    // This can be changed to any valid object name.

	// For the basic sample, these variables need not be changed.
	scope      = storage.DevstorageFull_controlScope
	authURL    = "https://accounts.google.com/o/oauth2/auth"
	tokenURL   = "https://accounts.google.com/o/oauth2/token"
	entityName = "allUsers"
	redirectURL = "urn:ietf:wg:oauth:2.0:oob"
)

var (
	cacheFile = flag.String("cache", "cache.json", "Token cache file")
	code      = flag.String("code", "", "Authorization Code")

	// For additional help with OAuth2 setup,
	// see http://goo.gl/cJ2OC and http://goo.gl/Y0os2

	// Set up a configuration boilerplate.
	config = &oauth.Config{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		Scope:        scope,
		AuthURL:      authURL,
		TokenURL:     tokenURL,
		TokenCache:   oauth.CacheFile(*cacheFile),
		RedirectURL:  redirectURL,
	}
)

func fatalf(service *storage.Service, errorMessage string, args ...interface{}) {
	restoreOriginalState(service)
	log.Fatalf("Dying with error:\n"+errorMessage, args...)
}

func restoreOriginalState(service *storage.Service) bool {
	succeeded := true

	// Delete an object from a bucket.
	if err := service.Objects.Delete(bucketName, objectName).Do(); err == nil {
		fmt.Printf("Successfully deleted %s/%s during cleanup.\n\n", bucketName, objectName)
	} else {
		// If the object exists but wasn't deleted, the bucket deletion will also fail.
		fmt.Printf("Could not delete object during cleanup: %v\n\n", err)
	}

	// Delete a bucket in the project
	if err := service.Buckets.Delete(bucketName).Do(); err == nil {
		fmt.Printf("Successfully deleted bucket %s during cleanup.\n\n", bucketName)
	} else {
		succeeded = false
		fmt.Printf("Could not delete bucket during cleanup: %v\n\n", err)
	}

	if !succeeded {
		fmt.Println("WARNING: Final cleanup attempt failed. Original state could not be restored.\n")
	}
	return succeeded
}

func main() {
	flag.Parse()

	// Set up a transport using the config
	transport := &oauth.Transport{
		Config:    config,
		Transport: http.DefaultTransport,
	}

	token, err := config.TokenCache.Token()
	if err != nil {
		if *code == "" {
			url := config.AuthCodeURL("")
			fmt.Println("Visit URL to get a code then run again with -code=YOUR_CODE")
			fmt.Println(url)
			os.Exit(1)
		}

		// Exchange auth code for access token
		token, err = transport.Exchange(*code)
		if err != nil {
			log.Fatal("Exchange: ", err)
		}
		fmt.Printf("Token is cached in %v\n", config.TokenCache)
	}
	transport.Token = token

	httpClient := transport.Client()
	service, err := storage.New(httpClient)

	// If the bucket already exists and the user has access, warn the user, but don't try to create it.
	if _, err := service.Buckets.Get(bucketName).Do(); err == nil {
		fmt.Printf("Bucket %s already exists - skipping buckets.insert call.", bucketName)
	} else {
		// Create a bucket.
		if res, err := service.Buckets.Insert(projectID, &storage.Bucket{Name: bucketName}).Do(); err == nil {
			fmt.Printf("Created bucket %v at location %v\n\n", res.Name, res.SelfLink)
		} else {
			fatalf(service, "Failed creating bucket %s: %v", bucketName, err)
		}
	}

	// List all buckets in a project.
	if res, err := service.Buckets.List(projectID).Do(); err == nil {
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
	file, err := os.Open(fileName)
	if err != nil {
		fatalf(service, "Error opening %q: %v", fileName, err)
	}
	if res, err := service.Objects.Insert(bucketName, object).Media(file).Do(); err == nil {
		fmt.Printf("Created object %v at location %v\n\n", res.Name, res.SelfLink)
	} else {
		fatalf(service, "Objects.Insert failed: %v", err)
	}

	// List all objects in a bucket
	if res, err := service.Objects.List(bucketName).Do(); err == nil {
		fmt.Printf("Objects in bucket %v:\n", bucketName)
		for _, object := range res.Items {
			fmt.Println(object.Name)
		}
		fmt.Println()
	} else {
		fatalf(service, "Objects.List failed: %v", err)
	}

	// Insert ACL for an object.
	// This illustrates the minimum requirements.
	objectAcl := &storage.ObjectAccessControl{
		Bucket: bucketName, Entity: entityName, Object: objectName, Role: "READER",
	}
	if res, err := service.ObjectAccessControls.Insert(bucketName, objectName, objectAcl).Do(); err == nil {
		fmt.Printf("Result of inserting ACL for %v/%v:\n%v\n\n", bucketName, objectName, res)
	} else {
		fatalf(service, "Failed to insert ACL for %s/%s: %v.", bucketName, objectName, err)
	}

	// Get ACL for an object.
	if res, err := service.ObjectAccessControls.Get(bucketName, objectName, entityName).Do(); err == nil {
		fmt.Printf("Users in group %v can access %v/%v as %v.\n\n",
			res.Entity, bucketName, objectName, res.Role)
	} else {
		fatalf(service, "Failed to get ACL for %s/%s: %v.", bucketName, objectName, err)
	}

	// Get an object from a bucket.
	if res, err := service.Objects.Get(bucketName, objectName).Do(); err == nil {
		fmt.Printf("The media download link for %v/%v is %v.\n\n", bucketName, res.Name, res.MediaLink)
	} else {
		fatalf(service, "Failed to get %s/%s: %s.", bucketName, objectName, err)
	}

	if !restoreOriginalState(service) {
		os.Exit(1)
	}

}
