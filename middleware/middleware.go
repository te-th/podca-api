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

	"google.golang.org/appengine"
)

// ServeWithContext abstracts a context.Context aware function
type ServeWithContext func(context.Context, http.ResponseWriter, *http.Request)

// ServeHTTP create and returns a http.HandlerFunc, creates a context.Context and invokes the given function
func ServeHTTP(serve ServeWithContext) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := appengine.NewContext(request)
		serve(ctx, writer, request)
	}
}
