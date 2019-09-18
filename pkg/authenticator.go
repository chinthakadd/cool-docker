package docker

import (
	"encoding/base64"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	log "github.com/sirupsen/logrus"
	"strings"
)

// Docker Login
// Only supports ECR Login at this point in time.
// https://docs.aws.amazon.com/sdk-for-go/api/
func GetAuthToken(registry DockerRegistry) (*string, *string) {


	url := registry.url
	registryUrlSplit :=  strings.Split(url, ".")
	registryId := registryUrlSplit[0]

	log.WithField("RegistryID", registryId).Info("Registry ID:")

	// split the docker url.
	authTokenInput := ecr.GetAuthorizationTokenInput{
		RegistryIds: []*string{aws.String(registryId)},
	}

	//895462685128.dkr.ecr.us-west-2.amazonaws.com
	awsConfig := aws.Config{
		Region: aws.String(registryUrlSplit[3]),
	}

	mySession, err := session.NewSession(&awsConfig)
	if err != nil {
		panic(err)
	}
	svc := ecr.New(mySession)

	authTokenOutput, err := svc.GetAuthorizationToken(&authTokenInput)
	if err != nil {
		panic(err)
	}

	authToken := authTokenOutput.AuthorizationData[0].AuthorizationToken
	expiresAt := authTokenOutput.AuthorizationData[0].ExpiresAt
	proxyEndpoint := authTokenOutput.AuthorizationData[0].ProxyEndpoint

	b64DecodedTokenBytes, err := base64.StdEncoding.DecodeString(*authToken)
	// todo: handle error

	base64DecodedToken := string(b64DecodedTokenBytes)
	creds := strings.Split(base64DecodedToken, ":")
	username := creds[0]
	password := creds[1]

	log.WithField("username", username).WithField("password", password).Info("Extracted from ECR")
	log.WithField("Expiry", expiresAt).Info("Expiration")
	log.WithField("Proxy Endpoint", *proxyEndpoint).Info("Proxy Endpoint")

	return &username, &password
}


type DockerRegistry struct {
	url string
}
