package docker

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"github.com/docker/docker/api/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"os"

	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/mitchellh/go-homedir"
	"io"
)

// Thank you Justin Wilson for your comment on
// https://forums.docker.com/t/how-to-create-registryauth-for-private-registry-login-credentials/29235/2
func PushImage(imageName string, tag string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	registryUrl := os.Getenv("DOCKER_REGISTRY_URL")
	if registryUrl == "" {
		panic("Registry URL cannot be null")
	}

	log.WithField("Registry URL", registryUrl).Info("Found Registry URL")
	fqin := registryUrl + "/" + imageName

	username, password := GetAuthToken(DockerRegistry{
		url: registryUrl,
	})

	authConfig := types.AuthConfig{
		Username: *username,
		Password: *password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	//log.WithField("token", *token).Info("Token:")
	reader, err := cli.ImagePush(ctx, fqin, types.ImagePushOptions{
		RegistryAuth: authStr,
	})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)
	_, _ = ctx, cli
}

// Pull a docker image provided as an argument.
func PullImage(imageName string) {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	registryUrl := os.Getenv("DOCKER_REGISTRY_URL")
	if registryUrl == "" {
		panic("Registry URL cannot be null")
	}

	log.WithField("Registry URL", registryUrl).Info("Found Registry URL")
	fqin := registryUrl + "/" + imageName

	username, password := GetAuthToken(DockerRegistry{
		url: registryUrl,
	})

	authConfig := types.AuthConfig{
		Username: *username,
		Password: *password,
	}

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	reader, err := cli.ImagePull(ctx, fqin, types.ImagePullOptions{
		RegistryAuth: authStr,
	})

	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)
	_, _ = ctx, cli
}

// Build the Docker Context.
// Essentially, Docker uses Tar file as the context.
func GetContext(filePath string) io.Reader {
	// Use homedir.Expand to resolve paths like '~/repos/myrepo'
	filePath, _ = homedir.Expand(filePath)
	log.Println(filePath)
	ctx, _ := archive.TarWithOptions(filePath, &archive.TarOptions{})
	return ctx
}

// Building docker image
// Refer to https://stackoverflow.com/questions/38804313/build-docker-image-from-go-code
func BuildImage(dockerfilePath string, dockerTag string) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	registryUrl := os.Getenv("DOCKER_REGISTRY_URL")
	if registryUrl == "" {
		panic("Registry URL cannot be null")
	}

	log.WithField("Registry URL", registryUrl).Info("Found Registry URL")

	tag := registryUrl + "/" + dockerTag
	buildResp, err := cli.ImageBuild(context.Background(), GetContext(dockerfilePath),
		types.ImageBuildOptions{
			Dockerfile:     "Dockerfile", // optional, is the default
			SuppressOutput: false,
			Tags:           []string{tag},
			PullParent:     true,
		})

	if err != nil {
		log.Fatal(err)
	}

	writeToLog(buildResp.Body)
	log.Println("Build Completed")
}

// https://medium.com/faun/how-to-build-docker-images-on-the-fly-2a1fd696c3fd
// Write LOG to the console.
func writeToLog(reader io.ReadCloser) error {
	defer reader.Close()
	rd := bufio.NewReader(reader)
	for {
		n, _, err := rd.ReadLine()
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		log.Println(string(n))
	}
	return nil
}
