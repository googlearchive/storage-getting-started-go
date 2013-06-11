# Google Cloud Storage Go Sample Application

## Description
This is a simple example of calling the Google Cloud Storage APIs in Go.

## Setup Authentication
1) Visit http://cloud.google.com/console to register your application.

2) From the "Project Home" screen, activate access to "Google Cloud Storage API":
- Click on "API Access" in the left column.
- Click the button labeled "Create an OAuth 2.0 client ID".
- Give your application a name and click "Next".
- Select "Installed Application" as the "Application type".
- Select "Other" under "Installed application type".
- Click "Create client ID".

## Prerequisites
1) Run the following commands:
* $ go get code.google.com/p/goauth2/oauth
* $ go get code.google.com/p/google-api-go-client/storage/v1beta2

2) In storage-sample.go, fill in your:
- Bucket name (this bucket will be created and deleted for you - it
      should not yet exist).
- Project ID.
- Client ID (see the documentation here: https://developers.google.com/console/help/index#generatingoauth2).
- Client secret (see the documentation here: https://developers.google.com/console/help/index#generatingoauth2).

## Running the Sample Application
1) Run the application:
  * $ go run storage-sample.go
