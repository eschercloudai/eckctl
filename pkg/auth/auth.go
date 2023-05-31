// https://github.com/eschercloudai/unikorn/blob/main/examples/serverclient.go
package auth

import (
	"bytes"
	"context"
	"crypto/tls"
	"eckctl/pkg/generated"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

// NewClient is a helper to abstract away client authentication.
func NewClient(server string, accessToken string) (*generated.ClientWithResponses, error) {
	return generated.NewClientWithResponses(server, generated.WithHTTPClient(httpClient(false)), generated.WithRequestEditorFn(bearerTokenInjector(accessToken)))
}

func GetToken(server string, u string, p string, project string) (accessToken string, err error) {
	token, err := oauth2Authenticate(server, u, p)
	if err != nil {
		return
	}

	scopedToken, err := getScopedToken(token, server, project)
	if err != nil {
		return
	}

	accessToken = scopedToken.AccessToken

	return
}

func httpClient(insecure bool) *http.Client {
	if insecure {
		return &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	} else {
		return &http.Client{}
	}
}

// JSONReader implments io.Reader that does lazy JSON marshaling.
type JSONReader struct {
	data interface{}
	buf  *bytes.Buffer
}

func NewJSONReader(data interface{}) *JSONReader {
	return &JSONReader{
		data: data,
	}
}

func (r *JSONReader) Read(p []byte) (int, error) {
	if r.buf == nil {
		data, err := json.Marshal(r.data)
		if err != nil {
			return 0, err
		}

		r.buf = bytes.NewBuffer(data)
	}

	return r.buf.Read(p)
}

// bearerTokenInjector is a handy function that augments the clients to add
func bearerTokenInjector(token string) generated.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+token)

		return nil
	}
}


// Login via oauth2's password grant flow.  But you should never do this.
// See https://tools.ietf.org/html/rfc6749#section-4.3.
func oauth2Authenticate(server string, u string, p string) (*oauth2.Token, error) {

	client := httpClient(false)

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, client)

	config := &oauth2.Config{
		Endpoint: oauth2.Endpoint{
			TokenURL: server + "/api/v1/auth/oauth2/tokens",
		},
	}

	return config.PasswordCredentialsToken(ctx, u, p)
}

// getScopedToken exchanges a token for one with a new project scope.
func getScopedToken(token *oauth2.Token, server string, projectID string) (*generated.Token, error) {
	client, err := NewClient(server, token.AccessToken)
	if err != nil {
		return nil, err
	}

	scope := &generated.TokenScope{
		Project: generated.TokenScopeProject{
			Id: projectID,
		},
	}

	response, err := client.PostApiV1AuthTokensTokenWithBodyWithResponse(context.TODO(), "application/json", NewJSONReader(scope))
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != 200 {
		return nil, fmt.Errorf("%w: unable to scope token", err)
	}

	return response.JSON200, nil
}

