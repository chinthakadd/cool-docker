# cool-docker

My Attempt To Build a Docker CLI less Docker Client based on Go. Currently only AWS ECR Registry is 
supported.

This CLI is built using Docker SDK and Amazon ECR SDK and currently it works connecting to the docker daemon running on local machine (/var/run/docker.sock).

## Install

```
go get github.com/chinthakadd/cool-docker
```

## Usage

**PREREQUISITE**

Please set the `DOCKER_REGISTRY_URL` environment variable with the docker registry url.

For help on all the commands available

```
cool-docker --help
```

To Build an image:

```
cool-docker build {DirWhereDockerFileExist} {ImageName}:{Tag(optional)}
```

To Pull an Image:

```
cool-docker pull {ImageName}:{Tag(optional)}
```

To Push an Image to registry:

```
cool-docker push {ImageName} {Tag}
```

## References


- https://godoc.org/github.com/docker/docker/client
- https://docs.docker.com/develop/sdk/examples/
- https://medium.com/faun/how-to-build-docker-images-on-the-fly-2a1fd696c3fd
- https://forums.docker.com/t/how-to-create-registryauth-for-private-registry-login-credentials/29235/2
- https://docs.aws.amazon.com/sdk-for-go/api/aws/session/

# LICENSE

[![BSD Zero Clause License](http://img.shields.io/badge/license-0BSD-blue.svg)](http://landley.net/toybox/license.html)

This is distributed under the [BSD Zero Clause License](http://landley.net/toybox/license.html).
