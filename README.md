# Google Cloud Storage Go Sample Application

## Description
This is a simple example of calling the Google Cloud Storage APIs in Go.

## Setup Authentication
1) Visit http://cloud.google.com/console to register your application.

2) In order to create (or find) the credentials for your application:
- Visit https://cloud.google.com/console
- Select a project that has Google Cloud Storage enabled (create such a project).
- From that project's page, click on the "APIs" section.
- Click the "REGISTER APP" button (or, if you have an existing Native app, you can instead follow the "All registered apps" link to select the app, and skip the next step).
- Name your application, and select "Native" as the platform, and register your app.
- Expand the "OAuth 2.0 Client ID" section to see your client ID and secret.


## Prerequisites
1) Run the following commands:
* $ go get code.google.com/p/goauth2/oauth
* $ go get code.google.com/p/google-api-go-client/storage/v1beta2

2) In storage-sample.go, fill in your:
- Bucket name (this bucket will be created and deleted for you - it
      should not yet exist).
- Project ID.
- Client ID (see the steps outlined above to find this).
- Client secret (see the steps outlined above to find this).


## Running the Sample Application
1) Run the application (on the first run, you will be prompted to go through the OAuth2 flow):
  * $ go run storage-sample.go
