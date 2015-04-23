# Google Cloud Storage Go Sample Application

## Description
This is a simple example of calling the Google Cloud Storage APIs in Go.

## Setup Authentication
1) Visit http://cloud.google.com/console to register your application.

2) In order to create (or find) the credentials for your application:
- Visit https://cloud.google.com/console
- Select a project that has Google Cloud Storage enabled (create such a project).
- From the project's page, under "APIs & auth", click on "Credentials".
- Under the "OAuth" secion, click on "Create new Client ID".
- Select "Service Account" and make sure that "JSON Key" is selected.
- Click on "Create Client ID" and download the JSON file.


## Prerequisites
1) Run the following commands:
* $ go get -u golang.org/x/net/context
* $ go get -u golang.org/x/oauth2/google
* $ go get -u google.golang.org/api/storage


2) In storage-sample.go, fill in your:
- Bucket name (this bucket will be created and deleted for you - it
      should not yet exist).
- Project ID.


## Running the Sample Application
1) Run the application (on the first run, you will be prompted to go through the OAuth2 flow):
  * $ go run storage-sample.go -creds <your-service-account-info>.json
