package main

/*
# fsrouter

A Go code generator that mirrors your `api/` filesystem into HTTP routes using [gorilla/mux](https://github.com/gorilla/mux).

## Installation

```bash
go install github.com/aquaticalf/fsrouter
```

## Usage

1. Structure your handlers under api/:

```
api/
  users/            # First-level directory becomes route group
    get.go          # exports func Get(w http.ResponseWriter, r *http.Request)
    [userId]/       # Dynamic parameter with [param] syntax
      get.go        # exports func Get(...)
      post.go       # exports func Post(...)
  auth/             # Another route group
    login/
      post.go       # exports func Post(...)
```

2. In your main.go, add:

```go
//go:generate fsrouter \
//    -api=./api \
//    -out=routes_gen.go \
//    -pkg=main \
//    -importPREFIX=yourmodule/api \
//    -middleware=yourmodule/middleware \
//    -middlewares="loggingMiddleware,authMiddleware,corsMiddleware" \
//    -groupMiddlewares='{"users":"authMiddleware","admin":"adminAuthMiddleware,loggingMiddleware"}' \
//    -notFound=customHandlers.NotFound

package main

import (
    "log"
    "net/http"
)

func main() {
    r := RegisterRoutes() // generated
    log.Fatal(http.ListenAndServe(":3000", r))
}
```

3. Generate and build:

```bash
go generate ./...
go build ./...
```

Routes will be wired up from `api/` files like `get.go`, `post.go`, etc.

## Features

- Automatic route registration from file system
- Dynamic parameters with `[param]` folder syntax
- Extensive Middleware Support
  - Multiple global middlewares for all routes
  - Group-specific middlewares for route groups
  - Flexible configuration via command line flags
- Route grouping via first-level directories
  - Each top-level directory becomes a subrouter
  - Example: `/api/users/...` becomes a group
- Custom 404 handler support
  - Specify with `-notFound=package.Handler`
  - Default JSON 404 handler included

## Command Line Options

| Flag | Description | Default |
|------|-------------|---------|
| `-api` | Directory of API handlers | `api` |
| `-out` | Output file path | `routes_gen.go` |
| `-pkg` | Package name for generated file | `main` |
| `-importPREFIX` | Import path prefix for API handlers | (required) |
| `-middleware` | Package containing middleware functions | (optional) |
| `-middlewares` | Comma-separated list of middleware functions to apply globally | `loggingMiddleware` |
| `-groupMiddlewares` | JSON mapping of group to middleware list | (optional) |
| `-notFound` | Custom 404 handler (format: `package.Handler`) | (default handler used) |

## Setting Up Middleware

### Global Middleware

You can specify multiple global middlewares using the `-middlewares` flag:

```bash
fsrouter -middlewares="loggingMiddleware,authMiddleware,corsMiddleware"
```

### Group-Specific Middleware

There are two ways to set up group-specific middleware:

1. Using the `-groupMiddlewares` flag with a JSON string:

```bash
fsrouter -groupMiddlewares='{"users":"authMiddleware,rateLimit","admin":"adminAuthMiddleware"}'
```

2. Editing the generated code (will be overwritten on regeneration):

```go
// After generation, you can manually edit:
usersRouter.Use(authMiddleware)
adminRouter.Use(adminAuthMiddleware)
```

### Creating Custom Middleware

Define your middleware functions in your application code:

```go
// In your middleware package
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Authentication logic here
        token := r.Header.Get("Authorization")
        if !validateToken(token) {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

## Extending

The generated code can be further customized as needed, but keep in mind that regeneration will overwrite your changes.

Use the command-line flags whenever possible to avoid manual edits.

## Example Generated Router

```go
// Example for a simple API structure
func RegisterRoutes() *mux.Router {
    r := mux.NewRouter()

    // Global middleware
    r.Use(loggingMiddleware)
    r.Use(authMiddleware)
    r.Use(corsMiddleware)

    // Route groups
    usersRouter := r.PathPrefix("/users").Subrouter()
    usersRouter.Use(authMiddleware)  // Group-specific middleware

    adminRouter := r.PathPrefix("/admin").Subrouter()
    adminRouter.Use(adminAuthMiddleware)
    adminRouter.Use(loggingMiddleware)

    authRouter := r.PathPrefix("/auth").Subrouter()

    // Routes
    usersRouter.HandleFunc("", users.Get).Methods("GET")
    usersRouter.HandleFunc("/{userId}", users_userId.Get).Methods("GET")
    adminRouter.HandleFunc("/dashboard", admin_dashboard.Get).Methods("GET")
    authRouter.HandleFunc("/login", auth_login.Post).Methods("POST")

    // 404 Handler
    r.NotFoundHandler = http.HandlerFunc(defaultNotFoundHandler)

    return r
}
```
*/
