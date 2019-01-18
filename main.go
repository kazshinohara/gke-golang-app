/**
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// [START all]
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"contrib.go.opencensus.io/exporter/stackdriver"

	"cloud.google.com/go/profiler"
	"go.opencensus.io/trace"
)

func main() {

	// initialize stackdriver profiler & trace
	initProfiler()
	initTrace()

	// use PORT environment variable, or default to 8080
	port := "8080"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}

	// register functions to handle all requests
	server := http.NewServeMux()
	server.HandleFunc("/hello", hello)
	server.HandleFunc("/slowhello", slowHello)

	// start the web server on port and accept requests
	log.Printf("Server listening on port %s", port)
	nerr := http.ListenAndServe(":"+port, server)
	log.Fatal(nerr)
}

func initProfiler() {
	// Stackdriver Profiler initialization.
	if err := profiler.Start(profiler.Config{
		Service:        "hello-app",
		ServiceVersion: "1.0.0",
		ProjectID:      os.Getenv("GOOGLE_CLOUD_PROJECT"),
		DebugLogging:   true,
	}); err != nil {
		log.Printf("failed to start profiler: %+v", err)
	} else {
		log.Printf("started stackdriver profiler")
	}
}

func initTrace() {
	// Stackdriver Tracer initialization.
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: os.Getenv("GOOGLE_CLOUD_PROJECT"),
	})
	if err != nil {
		log.Printf("failed to initialize stackdriver trace: %+v", err)
	} else {
		trace.RegisterExporter(exporter)
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
		log.Printf("started stackdriver trace")
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	// hello responds to the request with "Hello, world".
	log.Printf("Serving request: %s", r.URL.Path)
	ctx := r.Context()
	_, span := trace.StartSpan(ctx, "hello")
	defer span.End()
	host, _ := os.Hostname()
	fmt.Fprintf(w, "Hello, world!!!!!!\n")
	fmt.Fprintf(w, "Hostname: %s\n", host)
}

func slowHello(w http.ResponseWriter, r *http.Request) {
	// hello responds to the request as well but a bit slow.
	log.Printf("Serving request: %s", r.URL.Path)
	ctx := r.Context()
	ctx, span := trace.StartSpan(ctx, "slowHello")
	defer span.End()
	fibonacci(ctx)
	host, _ := os.Hostname()
	fmt.Fprintf(w, "Hello, world!\n")
	fmt.Fprintf(w, "Hostname: %s\n", host)
}

func fibonacci(ctx context.Context) {
	_, span := trace.StartSpan(ctx, "fibonacci")
	defer span.End()
	prev, next := 0, 1
	for i := 0; i < 1000000; i++ {
		prev, next = next, prev+next
	}
}

// [END all]
