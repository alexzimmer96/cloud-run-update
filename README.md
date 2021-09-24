# Cloud Run Update
[![Go Report Card](https://goreportcard.com/badge/github.com/alexzimmer96/cloud-run-update)](https://goreportcard.com/report/github.com/alexzimmer96/cloud-run-update)

## Project Description

This project contains a Google Cloud Function to automate deployments to Cloud Run.
The function is triggered whenever an image tag is pushed into a Container Registry (or Artifact Registry).
It builds the desired Cloud Run name from the new Digest and updates the configuration.

## Prerequisites

1. [Create the gcr topic](https://cloud.google.com/container-registry/docs/configuring-notifications) for receiving Pub/Sub-Events whenever an image is pushed.
2. Create a service account with access to Google Cloud Run
   1. The Service Accounts needs the Roles `Cloud Run Admin` and `Service Account User`
3. Create your Cloud Run Configurations
4. Update the variables inside the `Makefile` to match your configuration
5. Deploy your function (e.g. using `make deploy`)

## Configuration

|Variable|Description|Example|
|---|---|---|
|`functionName`|Name of the Google Cloud function|`update-cloud-run`|
|`region`|Region the function should run in|`europe-west1`|
|`serviceAccount`|Name of the service account with access to the Cloud Run Admin API|`update-cloud-run@some-cool-project.iam.gserviceaccount.com`|
|`endpoint`|Cloud Run API endpoint that is used. This must match the region of your Cloud Run containers.|`https://europe-west1-run.googleapis.com/`|
|`registry`|Name of the Container Registry holding the images that should automatically delivered to Cloud Run.|`europe-west1-docker.pkg.dev/some-cool-project/my-cloud-run-images`|
|`project`|ID of the project your Cloud Instances are running.|`some-cool-project`|

### Example

You have a Cloud Run Configuration named `my-service-dev` that should always point on the image `europe-west1-docker.pkg.dev/some-cool-project/my-cloud-run-images/my-service:dev`
The configuration inside the `Makefile` should look like this:

```makefile
# Other config vars

endpoint = "https://europe-west1-run.googleapis.com/"
registry = "europe-west1-docker.pkg.dev/some-cool-project/my-cloud-run-images"
project = "some-cool-project"
```

Now, everytime you push a new tag `europe-west1-docker.pkg.dev/some-cool-project/my-cloud-run-images/my-service:dev`, a new revision is created.
This revision points on the new image digest.