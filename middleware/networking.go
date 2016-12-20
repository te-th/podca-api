/*
Copyright [2016] TE,TH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package middleware

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

// HTTPClientFacade encapsulates HTTP methods
type HTTPClientFacade interface {
	Get(ctx context.Context, url string) (*http.Response, error)
}

type httpClient struct {
}

// NewHTTPClient creates a new HTTPClientFacade instance
func NewHTTPClient() HTTPClientFacade {
	return &httpClient{}
}

func (httpClient *httpClient) Get(ctx context.Context, url string) (*http.Response, error) {
	client := urlfetch.Client(ctx)
	return client.Get(url)
}
