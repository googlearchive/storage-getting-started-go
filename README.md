# Google Cloud Storage Go Sample Application

## Description
This is a simple example of calling the Google Cloud Storage APIs in Go.

## Setup Authentication
1. Visit https://code.google.com/apis/console/ to register your application.
2. From the "Project Home" screen, activate access to "Google Cloud Storage API":
   A. Click on "API Access" in the left column.
   B. Click the button labeled "Create an OAuth 2.0 client ID".
   C. Give your application a name and click "Next".
   D. Select "Installed Application" as the "Application type".
   E. Select "Other" under "Installed application type".
   F. Click "Create client ID".

## Prerequisites
1. Run the following command:
   * $ go get code.google.com/p/google-api-go-client/storage/v1beta2
2. In storage-sample.go, fill in your:
   A. Bucket name (this bucket will be created and deleted for you - it
      should not yet exist).
   B. Project ID.
   C. Client ID (in the "API Access" tab of https://code.google.com/apis/console/)
   D. Client secret (in the "API Access" tab of https://code.google.com/apis/console/)

## Running the Sample Application
1. Run the application:
  * $ go run storage-sample.go
