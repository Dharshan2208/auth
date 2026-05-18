package middleware

import "net/http"

func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent MIME sniffing.
		w.Header().Set("X-Content-Type-Options", "nosniff")
		// Prevent clickjacking.
		w.Header().Set("X-Frame-Options", "DENY")
		// Reduce referrer leakage.
		w.Header().Set("Referrer-Policy", "no-referrer")
		// Lock down powerful browser features...just for safety btw deafault only...
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// only for HTTPS....
		if r.TLS != nil {
			// 2 years, include subdomains. Add "preload" only if you plan to submit to the preload list.
			w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		}
		// Content Security Policy (safe for JSON APIs; blocks everything by default).
		w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'; base-uri 'none'")

		next.ServeHTTP(w, r)
	})
}
