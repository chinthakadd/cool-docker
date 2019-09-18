package main

import (
	"github.com/alecthomas/kingpin"
	"github.com/chinthakadd/cool-docker/pkg"
	log "github.com/sirupsen/logrus"
	"os"
)

// nolint
var (
	// application information
	name    = "cool-docker"
	version = "0.0.1"
	desc    = "My Attempt To Build a Docker CLI less Docker Client"

	app = kingpin.New(name, desc).Version(version).Author("chinthakadd")

	pullImage = app.Command("pull", "Pull a Docker Image")
	pullImageName = pullImage.Arg("imageName", "Name of the docker image").Required().String()

	pushImage = app.Command("push", "Pull a Docker Image")
	pushImageName = pushImage.Arg("imageName", "Name of the docker image").Required().String()
	pushImageDockerTag      = pushImage.Arg("tag", "Image tag").Required().String()

	buildImage     = app.Command("build", "Build a Docker Image")
	dockerFilePath = buildImage.Arg("path", "Location of Dockerfile").Required().String()
	dockerTag      = buildImage.Arg("tag", "Image tag").Required().String()
)

// gets invoked automatically in go.
func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

// example to build my own command
// https://github.com/alecthomas/kingpin/blob/master/_examples/completion/main.go
// https://github.com/alecthomas/kingpin#complex-example
// TODO: Connect to a remote daemon
func main() {

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case pullImage.FullCommand():
		println(*pullImageName)
		docker.PullImage(*pullImageName)

	case pushImage.FullCommand():
		println(*pushImageName)
		docker.PushImage(*pushImageName, *pushImageDockerTag)

	case buildImage.FullCommand():
		log.Println(*dockerFilePath)

		docker.BuildImage(*dockerFilePath, *dockerTag)
		log.Println("COMPLETED")
	}
}

