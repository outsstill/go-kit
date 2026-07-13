# go-kit

Reusable Go toolkit packages for database, Redis, cache, logging, storage, captcha, Elasticsearch, limiter, paginator, hashing, console helpers, and auth.

Publish this repository to your own Git host, then replace the module path in `go.mod`:

```bash
go mod edit -module github.com/<owner>/<repo>
go mod tidy
git tag v0.1.0
```

Consumers can then install it with:

```bash
go get github.com/<owner>/<repo>@v0.1.0
```

## Packages

- `database`: thin `database/sql` manager for multiple named connections.
- `redis`: small RESP client for common Redis commands.
- `cache`: cache interface with memory and file stores.
- `log`: leveled logger with text/json encoders.
- `storage`: storage interface with memory and local filesystem stores.
- `captcha`: simple math captcha generator and in-memory verifier.
- `es`: Elasticsearch HTTP client for index/search/get/delete.
- `limiter`: token bucket and fixed window limiters.
- `paginator`: request/response pagination helpers.
- `hash`: password, HMAC, SHA, MD5, and random token helpers.
- `console`: colored console output and confirmation prompts.
- `auth`: JWT HS256 helpers and password wrapper.


## usage

```go
var err error

app, err := gokit.New(configName)

if err != nil {
    panic(err)
}

// 按需加载
err = gokit.App().Init(
gokit.Kit_Logger,
gokit.Kit_DB,
gokit.Kit_Redis,
gokit.Kit_Cache,
gokit.Kit_JWT,
)

if err != nil {
panic(err)
}

gokit.Set(app)
```