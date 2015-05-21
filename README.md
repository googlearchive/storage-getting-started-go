# Google Cloud Storage Go Sample Application

## Description
This is a simple example of calling the Google Cloud Storage APIs in Go.

## Setup Google Cloud SDK and Authentication

Install the [Google Cloud SDK](https://cloud.google.com/sdk):

```
curl https://sdk.cloud.google.com | bash
```

Once installed, authenicate with your Google account:

```
gcloud auth login
```

## Prerequisites
Install dependencies with ``go get``.

```
$ go get -u golang.org/x/net/context
$ go get -u golang.org/x/oauth2/google
$ go get -u google.golang.org/api/storage/...
```

## Running the Sample Application

You will need a project with billing set up and the Cloud Storage API enabled. You do not need an existing bucket as the sample will create one for you. You will also need a local file to upload.

```
$ go run storage-sample.go --project=<your-project-id> --bucket=<new-bucket-name> --file=<path-to-local-file-to-upload>
```
