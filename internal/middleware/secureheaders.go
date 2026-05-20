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
		// Content Security Policy:
		// - APIs: locked down (no HTML/JS expected).
		// - Swagger UI: allow the assets and inline script/style that swagger-ui uses.
		if isSwaggerRoute(r) {
			w.Header().Set("Content-Security-Policy", swaggerUIContentSecurityPolicy())
		} else {
			w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'; base-uri 'none'")
		}

		next.ServeHTTP(w, r)
	})
}

func isSwaggerRoute(r *http.Request) bool {
	// http-swagger serves UI under /swagger/ and the OpenAPI JSON at /swagger/doc.json.
	return len(r.URL.Path) >= len("/swagger/") && r.URL.Path[:len("/swagger/")] == "/swagger/"
}

func swaggerUIContentSecurityPolicy() string {
	// Swagger UI relies on inline script/style for bootstrapping.
	// We keep it scoped to same-origin resources.
	return "default-src 'self'; base-uri 'none'; frame-ancestors 'none'; " +
		"script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; " +
		"img-src 'self' data:; font-src 'self' data:; connect-src 'self'"
}
