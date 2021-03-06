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

package app

import (
	"net/http"

	"errors"

	"github.com/gorilla/mux"
	"github.com/te-th/podca-api/api"
)

func init() {
	router := mux.NewRouter()

	handler := api.Handler()
	if len(api.Handler()) == 0 {
		panic(errors.New("No handlers registered. Will exit."))
	}

	for path, handler := range handler {
		router.HandleFunc(path, handler)
	}
	http.Handle("/", router)
}
