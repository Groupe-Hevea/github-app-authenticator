package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/cristalhq/jwt/v3"
)

// Version is the CLI version number - injected at build time
var Version string

// InstallationAccessToken is a response from the Github APIs
type InstallationAccessToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

func main() {
	argsWithoutProg := os.Args[1:]

	if argsWithoutProg[0] == "-v" {
		fmt.Printf("github-app-authenticator - version %s\n", Version)
		os.Exit(0)
	}

	if len(argsWithoutProg) != 3 {
		log.Fatalln("Usage : ./cmd github_app_id private_key_path app_installation_id")
	}

	githubAppID := argsWithoutProg[0]
	privateKeyPemPath := argsWithoutProg[1]
	githubAppInstallationID := argsWithoutProg[2]

	priv, err := ReadPEMPrivateKey(privateKeyPemPath)
	if err != nil {
		log.Panic(err)
	}

	signer, err := jwt.NewSignerRS(jwt.RS256, priv)
	if err != nil {
		log.Panic("failed to create jwt signer")
	}

	claims := jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(10 * time.Minute))),
		Issuer:    githubAppID,
	}

	builder := jwt.NewBuilder(signer)
	token, err := builder.Build(claims)
	if err != nil {
		log.Panic("failed to generate token")
	}

	client := &http.Client{}
	url := fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", githubAppInstallationID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Panic("failed to generate new HTTP request")
	}

	req.Header.Add("Authorization", fmt.Sprintf(`Bearer %s`, token.String()))
	req.Header.Add("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		log.Panic("HTTP request failed")
	}
	defer resp.Body.Close()

	var response InstallationAccessToken

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 201 {
		log.Fatalf("failed to create installation access token: got status code %d \n", resp.StatusCode)
	}
	json.Unmarshal([]byte(body), &response)

	fmt.Println(response.Token)
	os.Exit(0)
}

// ReadPEMPrivateKey reads the given file path and creates a new RSA PrivateKey
func ReadPEMPrivateKey(pemPath string) (*rsa.PrivateKey, error) {
	privateKey, err := ioutil.ReadFile(pemPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("Failed to parse PEM block")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}
