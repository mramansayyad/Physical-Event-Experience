package middleware

import (
	"net/http"
)

// SecurityHeaders acts as a rigorous Zero-Trust middleware, explicitly mapping HSTS,
// Content-Security-Policy (CSP), and X-Content-Type-Options to lock down client interactions.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Strict Transport Security (HSTS) - demands TLS connections natively for 2 years
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

		// Preempts MIME type sniffing to prevent malicious file injection masking
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Locks embedding the platform within iframe architectures (Clickjacking defense)
		w.Header().Set("X-Frame-Options", "DENY")

		// Highly restrictive CSP completely denying external script/resource evaluation unless natively declared
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self'; script-src 'self'; object-src 'none'")

		next.ServeHTTP(w, r)
	})
}
