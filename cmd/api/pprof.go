package main

import (
	"log"
	"net/http"
	"net/http/pprof"
)

// StartPprofServer boots an isolated HTTP server specifically dedicated to Go profiling.
// It is explicitly isolated on port :6060 to prevent public boundary routing ingress.
func StartPprofServer() {
	mux := http.NewServeMux()
	
	// Map the standard pprof profile handlers independently
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	log.Println("Diagnostic Profiler active identically and isolated on Internal Port :6060")
	// Start explicitly separated from the public multiplexer
	if err := http.ListenAndServe("127.0.0.1:6060", mux); err != nil {
		log.Fatalf("Pprof server shutdown unexpectedly: %v", err)
	}
}
