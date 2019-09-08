/*
Copyright 2019 The Kubernetes Authors.

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

package filters

import (
	"net/http"

	"github.com/opentracing/opentracing-go"

	"k8s.io/apimachinery/pkg/util/trace"
)

func WithTracing(handler http.Handler, tracer opentracing.Tracer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx, span := trace.NewSpan(req.Context(), "WithTracing")
		defer span.Finish()

		req = req.WithContext(ctx)

		handler.ServeHTTP(w, req)
	})
}
