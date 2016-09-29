# Step by Step: continuous deployment of Go microservices with Docker, Semaphore and Google Cloud
Leverage the speed of Go, the isolation of Docker and the container engine in Google Cloud with automated testing and deployment after just a few clicks in Semaphore.

![intro](https://cdn.pbrd.co/images/6CIeybgZg.png)

## Introduction

In this tutorial, we'll see how we can start from a microservice written in Go, create a Docker container for it and use Semaphore for continuous integration and deployment to Google Container Registry (GCR). 

Because Semaphore natively support all of these technologies, enabling this scenario is just a few clicks away and can be a big time-saver for your Go deployments. 

The workflow covered is a common one:

* You push new code to GitHub or Bitbucket
* Semaphore runs your test suite
* If tests pass, Semaphore builds your container image
* Semaphore pushes your code to GCR
* The registry ships your container to production

Full source-code for this tutorial is [available at Bitbucket](https://bitbucket.org/theplastictoy/go-semaphore-docker).

## Prerequisites

We assume you're familiar, at a basic level, with the technologies used. That means you:

* Can write a simple Go app
* Know what Docker is and the purpose behind it
* Can write simple automated tests for a Go app
* Have used Google Container Registry before

Some programs are also required in your local machine before we move on. Please make sure that:

* Go is installed. You can [get it here](https://golang.org/dl/).
* Docker is up and running. You can [get it here](https://docs.docker.com/engine/installation/).

At any time, free free to [contact us](https://semaphoreci.com/contact) if you need assistance.

## 1. Create a Go microservice

We'll be creating a small microservice written in Go that returns the current server date in [RFC822Z](https://golang.org/pkg/time/#pkg-constants) format. Let's keep the codebase simple so we can focus on the Docker/Semaphore/Google integration. 

For the rest of this tutorial, we'll assume the current folder is `$GOPATH`. Create a `src/servertime/main.go` file with the following code:

```
package main

import (
  "time"
  "encoding/json"
  "net/http"
)

type ServiceResult struct {
  FormattedTime string
}

func currentTime() string {
  return time.Now().Format(time.RFC822Z)
}

func handler(w http.ResponseWriter, r *http.Request) {
  t := currentTime()
  sr := ServiceResult{t}
  js, err := json.Marshal(sr)

  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

func main() {
  http.HandleFunc("/", handler)
  http.ListenAndServe(":8080", nil)
}
```

Now install and run the `servertime` microservice:

```
go install servertime
./bin/servertime
```

The microservice should be running on port 8080, so let's `curl -X GET http://localhost:8080/` and confirm that we have an RFC822Z time as the result: `{"FormattedTime":"23 Sep 16 11:39 +0100"}`

## 2. Testing the Go microservice

We want to make sure that our microservice always returns the date in the correct format, so let's write a small unit test to enforce it. 

Let's create a `$GPATH/src/servertime/main_test.go` with the following contents:

```
package main

import (
  "testing"
  "time"
)

func TestCurrentTime(t *testing.T) {
  result := currentTime()
  // Check that err is null for RFC822Z time format
  _, err := time.Parse(time.RFC822Z, result)
  if err != nil {
    t.Fail()
  }
}
``` 

We can make sure our test works by running 

```
go install servertime
go test servertime
```

At this point, we have a Go microservice ready to be deployed. Let's dive into that.


## 3. Dockerize the Go microservice

The increasing adoption of Docker pushed major cloud service providers ([Google](https://cloud.google.com/container-registry/), [Amazon](https://aws.amazon.com/ecs/) and [Microsoft](https://azure.microsoft.com/en-us/services/container-service/), among many others) to fully support dockerized applications in their offerings. By using Docker, our applications become an environment that can predictably run anywhere without code changes. 

To create the docker container for our microservice, we'll be using the official Go docker image. The `Dockerfile` should look like this:

```
FROM golang
ADD . /go/src/servertime
RUN go install servertime
ENTRYPOINT /go/bin/servertime
EXPOSE 8080
```

The code above fetches the official Go image, installs and runs our microservice and exposes the 8080 port, which is everything we need. Let's test it locally:

```
docker build -t servertime .
docker run --publish 8080:8080 --name servertime --rm servertime
```

The service should be available at `http://localhost:8080` and we're now ready to add continuous integration and deployment.

## 4. Set-up continuous integration in Semaphore

Semaphore has built-in support for Go, Docker and Google Container Registry. Let's leverage Semaphore integration to set-up our continuous integration environment. You'll need a Semaphore account, so if you haven't done it, [create one for free](https://semaphoreci.com/users/sign_up).

Start by [creating a new project in Semaphore](https://semaphoreci.com/docs/adding-github-bitbucket-project-to-semaphore.html) and adding the servertime repo, which can be either on GitHub or Bitbucket. 

Choose `Docker` as the project's platform and [follow our instructions on configuring Docker and GCR](https://semaphoreci.com/docs/docker/setting-up-continuous-integration-for-docker-project.html). When configuring GCR, remember that your email is the service account one from Google, which should be something like `<name>@<project_name>.iam.gserviceaccount.com`.

Once Docker and GCR are configured, Semaphore analyses the project's code and detects that we're using Go, so it installs default commands for building and testing our microservice.

![Project Setup](https://cdn.pbrd.co/images/7JMQMwVIK.png)

Let's build with these settings and we should have a passing build. 

![Build](https://cdn.pbrd.co/images/7JPLCeDRH.png)

While the default commands added by Semaphore build and test our Go code, we're still not using our Docker container or pushing to GCR. Let's go for it in the next section.

## 5. Continuously push to Google Container Registry

Because Semaphore natively support Docker, we have access to the `docker` command line. With that, we can build the image in Semaphore and push it to Google.

In the project's main page, click on Set Up Deployment. 

![Setup Deployment](https://cdn.pbrd.co/images/7JTuY9pCt.png)

Select "Generic Deployment".

![Generic Deployment](https://cdn.pbrd.co/images/LZwc6x6J.png)

Select 'Automatic' to enable continuous deployment.

![Automatic Deployment](https://cdn.pbrd.co/images/7JVWCnQ7n.png)

Let's use the `master`branch. Semaphore will trigger a deployment after successful build of this branch.

Every deployment needs to build the Docker image, tag it and upload it to GCR. This translates to the following build commands (replace `go-semaphore-demo-144214` with your Google project id)

```
docker build -t servertime
docker tag servertime gcr.io/go-semaphore-demo-144214/servertime
docker push gcr.io/go-semaphore-demo-144214/servertime
```

Skip ssh configuration and name your server as you wish. Proceed to deploy and you should now have the docker container in GCR. 

![GCR](https://cdn.pbrd.co/images/M0XIH3KP.png)

**TIP:** You can add `.semaphore-cache` to the `.dockerignore` file in your repo to avoid including the project's dependency cache in Docker's build context. You can also [save and reuse Docker images](https://semaphoreci.com/docs/docker.html).

## Next steps

[Let us know](https://semaphoreci.com/contact) how Semaphore is making it easier for you to build great stuff and what you're building with it. We are always on the hunt for great ways to build software.

Feel free to explore more scenarios for Go, Docker or GCR in our [Documentation](https://semaphoreci.com/docs/) and [Community Tutorials](https://semaphoreci.com/community). You're more than welcome to submit a tutorial yourself. Our community pages are [open-source on GitHub](https://github.com/renderedtext/semaphore-docs-new?__hstc=233161921.52904f9128a55f51771332247e541fb8.1474644404421.1474888169701.1474894848837.4&__hssc=233161921.10.1474894848837&__hsfp=803116496).

Google Container Registry and Container Engine are evolving rapidly, so stay tuned in their [docs](https://cloud.google.com/container-engine/docs/).


