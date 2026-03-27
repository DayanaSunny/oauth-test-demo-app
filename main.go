package main

import (
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"
)

var pageTemplate = template.Must(template.New("callback").Parse(`
<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>OAuth Callback</title>
  <style>
    body {
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
      margin: 0;
      padding: 24px;
      background: #f8fafc;
      color: #0f172a;
    }
    .card {
      max-width: 700px;
      margin: 40px auto;
      background: #fff;
      border-radius: 12px;
      box-shadow: 0 8px 24px rgba(15, 23, 42, 0.08);
      padding: 24px;
    }
    h1 {
      margin-top: 0;
      color: #047857;
      font-size: 1.6rem;
    }
    .scope-list {
      margin-top: 16px;
      padding-left: 18px;
    }
    .muted {
      color: #475569;
    }
    code {
      background: #f1f5f9;
      border-radius: 6px;
      padding: 2px 6px;
    }
  </style>
</head>
<body>
  <div class="card">
    <h1>You've been given access to the scopes</h1>
    {{if .Scopes}}
      <ul class="scope-list">
      {{range .Scopes}}
        <li><code>{{.}}</code></li>
      {{end}}
      </ul>
    {{else}}
      <p class="muted">No scopes were provided in the callback.</p>
      <p class="muted">Try: <code>/oauth/callback?scope=openid%20profile%20email</code></p>
    {{end}}
  </div>
</body>
</html>
`))

type callbackData struct {
	Scopes []string
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/oauth/callback?scope=openid%20profile", http.StatusFound)
	})

	mux.HandleFunc("/oauth/callback", oauthCallbackHandler)

	addr := ":8080"
	log.Printf("Server running at http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func oauthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Accept either repeated scope values (?scope=a&scope=b) or a space-delimited scope string.
	rawScopes := r.URL.Query()["scope"]
	if len(rawScopes) == 0 {
		rawScopes = r.URL.Query()["scopes"]
	}

	scopeSet := map[string]struct{}{}
	for _, raw := range rawScopes {
		for _, scope := range strings.Fields(raw) {
			scopeSet[scope] = struct{}{}
		}
	}

	scopes := make([]string, 0, len(scopeSet))
	for scope := range scopeSet {
		scopes = append(scopes, scope)
	}
	sort.Strings(scopes)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := pageTemplate.Execute(w, callbackData{Scopes: scopes}); err != nil {
		http.Error(w, "failed to render response", http.StatusInternalServerError)
	}
}
