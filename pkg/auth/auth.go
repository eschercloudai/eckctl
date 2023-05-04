package auth

import (
	"context"
	"eckctl/pkg/generated"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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

func InitClient(url string) (client *generated.Client, err error) {
	customHTTPClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	client, err = generated.NewClient(url, generated.WithHTTPClient(customHTTPClient))
	if err != nil {
		return
	}

	return
}

func getBearer(url string, t string, p string) (bearer string, err error) {
	client, err := InitClient(url)

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
		err := fmt.Errorf("Unexpected response when obtaining bearer token: %v", resp.StatusCode)
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	b := generated.Token{}
	err = json.Unmarshal(body, &b)

	if err != nil {
		log.Fatal(err)
	}

	return b.Token, err
}

func GetToken(url string, username string, password string, project string) (token string, err error) {
	client, err := InitClient(url)
	if err != nil {
		return
	}

	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.PostApiV1AuthTokensPassword(ctx, SetAuthorizationHeader(auth))
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected response code when authenticating: %s %v", err, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	t := generated.Token{}
	err = json.Unmarshal(body, &t)
	if err != nil {
		return
	}

	token, err = getBearer(url, t.Token, project)
	if err != nil {
		return
	}
	return
}
