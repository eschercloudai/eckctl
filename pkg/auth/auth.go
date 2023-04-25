package auth

import (
	"context"
	"eckctl/pkg/generated"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func isBase64Encoded(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

func SetAuthorizationHeader(token string) generated.RequestEditorFn {
	if isBase64Encoded(token) {
		return func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))
			return nil
		}
	} else {
		return func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			return nil
		}
	}
}

func InitClient(url string) *generated.Client {
	customHTTPClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	client, err := generated.NewClient(url, generated.WithHTTPClient(customHTTPClient))
	if err != nil {
		log.Fatal()
	}

	return client
}

func getBearer(url string, t string, p string) string {
	client := InitClient(url)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tokenScope := new(generated.TokenScope)
	tokenScope.Project.Id = p

	resp, err := client.PostApiV1AuthTokensToken(ctx, *tokenScope, SetAuthorizationHeader(t))
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected response code: ", resp.StatusCode)
		log.Fatal(err, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	b := generated.Token{}
	err = json.Unmarshal(body, &b)

	if err != nil {
		log.Fatal(err)
	}

	return b.Token
}

func GetToken(url string, username string, password string, project string) string {
	client := InitClient(url)

	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.PostApiV1AuthTokensPassword(ctx, SetAuthorizationHeader(auth))
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected response code: ", resp.StatusCode)
		log.Fatal(err, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	t := generated.Token{}
	err = json.Unmarshal(body, &t)
	if err != nil {
		log.Fatal(err)
	}

	return getBearer(url, t.Token, project)
}
