package main

import (
	"context"
	"os"

	"github.com/antihax/optional"
	"github.com/banzaicloud/pipeline/client"
)

// First you have to create a Pipeline Bearer token and put it into the TOKEN env variable.
func main() {

	// // Load Root CA cert example
	// CACert := []byte{}

	// caCertPool := x509.NewCertPool()
	// caCertPool.AppendCertsFromPEM([]byte(CACert))

	// clientTLSConfig := &tls.Config{
	// 	RootCAs: caCertPool,
	// }
	// clientTLSConfig.BuildNameToCertificate()
	// transport := &http.Transport{TLSClientConfig: clientTLSConfig}
	// httpClient := &http.Client{Transport: transport}

	// Create a new configuration with a custom root certificate pool
	config := client.NewConfiguration()
	// config.HTTPClient = httpClient

	ctx := context.WithValue(context.Background(), client.ContextAccessToken, os.Getenv("TOKEN"))
	pipeline := client.NewAPIClient(config)

	// Create a new token for a virtual user
	tokenRequest := client.TokenCreateRequest{Name: "drone token", VirtualUser: "banzaicloud/pipeline"}
	tokenResponse, _, err := pipeline.AuthApi.CreateToken(ctx, tokenRequest)

	if err != nil {
		panic(err)
	}

	// Overwrite the existing context token
	ctx = context.WithValue(context.Background(), client.ContextAccessToken, tokenResponse.Token)

	// Create a new Generic secret
	secretRequest := client.CreateSecretRequest{
		Name:   "my-password",
		Type:   "generic",
		Values: map[string]interface{}{"password": "s3cr3t"},
		Tags:   []string{"banzai:hidden"},
	}

	secretResponse, _, err := pipeline.SecretsApi.AddSecrets(ctx, 2, secretRequest, &client.AddSecretsOpts{
		Validate: optional.NewBool(true),
	})
	if err != nil {
		panic(err)
	}

	println("Secret id:", secretResponse.Id)
}
