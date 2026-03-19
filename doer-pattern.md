# Doer Pattern — Implementation Guide

A step-by-step reference for implementing and extending the doer pattern in a Go + Gin + fx app.

---

## What is the Doer Pattern?

The doer pattern breaks a single handler's logic into small, ordered, single-responsibility tasks called **doers**. Each doer reads from and writes into a shared `ctx` struct. The handler just kicks off the chain and reads the final result.

```
Request → Handler → build ctx → doer chain → ctx filled → Response
                                   ↓
                          Doer1 → Doer2 → Doer3
                          (each reads & writes ctx)
```

---

## Step 1 — Define `core.Ctx` (the shared data carrier)

`core.Ctx` holds everything every feature needs: DB, config, logger, and request-scoped values like `UserId`.

```go
// internal/core/context.go
package core

import (
    "time"
    "go.uber.org/zap"
    "yourapp/internal/config"
    "yourapp/internal/stores"
)

type Ctx struct {
    Store  *stores.StoreHolder
    Config *config.Config
    Log    *zap.Logger
    Now    time.Time

    // request-scoped — set per request in the handler
    UserId int32
}

// NewCtx is registered as an fx provider (singleton)
func NewCtx(
    store *stores.StoreHolder,
    config *config.Config,
    logger *zap.Logger,
) *Ctx {
    return &Ctx{
        Store:  store,
        Config: config,
        Log:    logger,
    }
}
```

> **Key rule:** `NewCtx` is a singleton — created once at startup by fx.
> Fields like `UserId` or `Now` are request-scoped. Set them per request inside the handler (see Step 4).

---

## Step 2 — Define the `Doer` interface and runner

```go
// internal/core/doer.go
package core

type DoCtx struct {
    IsExit  bool  // set true inside a doer to stop the chain early
    NxtDoer Doer  // set to jump forward to a specific doer
}

type Doer interface {
    Do(ctx *Ctx, doCtx *DoCtx) error
}

type Doers []Doer

func (ds Doers) Do(ctx *Ctx, doCtx *DoCtx) error {
    for _, d := range ds {
        if doCtx != nil && doCtx.NxtDoer != nil && doCtx.NxtDoer != d {
            continue
        }
        doCtx.NxtDoer = nil

        if err := d.Do(ctx, doCtx); err != nil {
            return err // stops the whole chain
        }

        if doCtx != nil && doCtx.IsExit {
            break
        }
    }
    return nil
}
```

> `Doers.Do()` is just a for-loop. Every doer gets called in order. If any returns an error, the chain stops immediately.

---

## Step 3 — Wire dependencies with fx

### 3a. Connect your DB and build the StoreHolder

```go
// internal/stores/userStore.go
type UserStore struct {
    db *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
    return &UserStore{db: db}
}
```

```go
// internal/stores/modulo.go
type StoreHolder struct {
    UserStore *UserStore
}

// fx injects *gorm.DB automatically from conn.ConnectPostgres
func NewStoreHolder(db *gorm.DB) *StoreHolder {
    return &StoreHolder{
        UserStore: NewUserStore(db),
    }
}
```

```go
// internal/conn/postgres.go
func ConnectPostgres(cfg *config.Config) (*gorm.DB, error) {
    return gorm.Open(postgres.Open(cfg.DB.DSN), &gorm.Config{})
}
```

### 3b. Register everything in main.go

```go
fx.New(
    fx.Provide(
        pkg.CustomLogger,
        config.LoadConfig,
        conn.ConnectPostgres,   // provides *gorm.DB
        stores.NewStoreHolder,  // needs *gorm.DB → injected automatically
        core.NewCtx,            // needs StoreHolder, Config, Logger → injected
        handlers.NewUserHandler,
        services.NewUserService,
        GinHttpServer,
        api.SetupRoutes,
    ),
    fx.Invoke(...),
).Run()
```

> fx resolves the dependency graph automatically. You just declare what each constructor needs as parameters.

---

## Step 4 — Write the handler

The handler's job:

1. Parse request params
2. Make a per-request copy of `ctx` and set request-scoped fields
3. Build a `ResponseCtx` (the output carrier)
4. Kick off the doer chain
5. Serve the result

```go
// internal/api/handlers/userHandler.go
package handlers

import (
    "time"
    "yourapp/internal/core"
    userService "yourapp/internal/services/user"
    "yourapp/pkg/utils"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type UserHandler interface {
    GetUsers(c *gin.Context)
}

type userHandler struct {
    ctx *core.Ctx // singleton injected at startup
}

func NewUserHandler(ctx *core.Ctx) UserHandler {
    return &userHandler{ctx: ctx}
}

func (h *userHandler) GetUsers(c *gin.Context) {
    // shallow copy so we can set request-scoped fields safely
    ctx := *h.ctx
    ctx.Now = time.Now()
    // ctx.UserId = parseUserID(c) // set any request-scoped fields here

    resCtx := userService.ResponseCtx{}
    if err := userService.New(&resCtx).Do(&ctx, &core.DoCtx{}); err != nil {
        ctx.Log.Error("error getting users", zap.Error(err))
        utils.ServeErr(c, err)
        return
    }

    utils.ServeData(c, resCtx)
}
```

> **Why shallow copy?** `h.ctx` is a singleton. If you set `h.ctx.UserId = X` directly, it leaks across concurrent requests. Copy first, then mutate.

---

## Step 5 — Create a service package

Each feature (users, parcels, etc.) gets its own package with three files:

```
internal/services/user/
  ├── constructor.go   // New() + Do() — builds and runs the doer chain
  ├── data.go          // ResponseCtx — the output struct for this feature
  └── tsk.go           // individual doers (one per task)
```

### constructor.go

```go
package userService

import "yourapp/internal/core"

type user struct {
    ctx *ResponseCtx
}

// New wraps ResponseCtx and returns a core.Doer
func New(ctx *ResponseCtx) core.Doer {
    return &user{ctx: ctx}
}

// Do declares the ordered chain for this feature
func (s *user) Do(ctx *core.Ctx, doCtx *core.DoCtx) error {
    doers := core.Doers{
        &FetchUser{resCtx: s.ctx},
        // add more doers here as the feature grows
    }
    return doers.Do(ctx, doCtx)
}
```

### data.go

```go
package userService

import "yourapp/internal/models"

// ResponseCtx is the output carrier — doers write results here
// Export fields so encoding/json can serialize them
type ResponseCtx struct {
    Name string        `json:"name"`
    Data []models.User `json:"data"`
}
```

> **Common mistake:** lowercase fields in `ResponseCtx` means `json.Marshal` ignores them → empty `{}` response. Always export and tag.

### tsk.go

```go
package userService

import (
    "fmt"
    "yourapp/internal/core"
)

type FetchUser struct {
    resCtx *ResponseCtx
}

func (fu *FetchUser) Do(ctx *core.Ctx, doCtx *core.DoCtx) error {
    users, err := ctx.Store.UserStore.GetFirstTenUsers()
    if err != nil {
        return fmt.Errorf("users not found: %w", err)
    }
    fu.resCtx.Name = "users"
    fu.resCtx.Data = users
    return nil
}
```

---

## How data flows through the chain

```
handler creates:
  ctx       → *core.Ctx  (shared deps + request values)
  resCtx    → *ResponseCtx  (empty, will be filled)

handler calls:
  userService.New(&resCtx).Do(&ctx, &core.DoCtx{})
                  ↑                 ↑
          resCtx pointer         ctx pointer
          stored in user{}       passed to every doer

doer chain runs:
  FetchUser.Do(&ctx, doCtx)
    reads  → ctx.Store.UserStore   (from core.Ctx)
    writes → fu.resCtx.Data        (into ResponseCtx)

handler reads:
  resCtx.Data  ← already filled because pointer was shared
```

---

## How to add a new feature

1. Create `internal/services/yourfeature/` with `constructor.go`, `data.go`, `tsk.go`
2. Define `ResponseCtx` with exported, json-tagged fields
3. Write each task as a struct implementing `Do(*core.Ctx, *core.DoCtx) error`
4. List tasks in order in `constructor.go`'s `Do()` method
5. Create a handler, inject `*core.Ctx` via constructor, shallow-copy per request

If the feature needs fields beyond `core.Ctx`, embed it:

```go
// for feature-specific shared state between doers
type FeatureCtx struct {
    core.Ctx           // embedded — access all core fields directly
    SomeIntermediateValue string
}
```

Then pass `*FeatureCtx` through the chain instead of `*core.Ctx`.

---

## How to add a new store

```go
// 1. Define the store
type OrderStore struct { db *gorm.DB }
func NewOrderStore(db *gorm.DB) *OrderStore { return &OrderStore{db: db} }
func (s *OrderStore) GetOrder(id int) (*models.Order, error) { ... }

// 2. Add to StoreHolder
type StoreHolder struct {
    UserStore  *UserStore
    OrderStore *OrderStore   // ← add here
}

func NewStoreHolder(db *gorm.DB) *StoreHolder {
    return &StoreHolder{
        UserStore:  NewUserStore(db),
        OrderStore: NewOrderStore(db),  // ← init here
    }
}

// 3. Use in any doer
func (d *FetchOrder) Do(ctx *core.Ctx, doCtx *core.DoCtx) error {
    order, err := ctx.Store.OrderStore.GetOrder(int(ctx.OrderId))
    ...
}
```

---

## Quick reference — common mistakes

| Mistake                                      | Symptom                          | Fix                                                 |
| -------------------------------------------- | -------------------------------- | --------------------------------------------------- |
| Handler interface signature mismatch         | compile error                    | Method signature must exactly match the interface   |
| `Ctx` value not `*Ctx` in `Doer` interface   | inconsistent behaviour           | Use `*Ctx` everywhere                               |
| Mutating singleton `h.ctx` directly          | data leak across requests        | Shallow copy: `ctx := *h.ctx` then mutate           |
| Unexported `ResponseCtx` fields              | empty `{}` JSON response         | Capitalize fields + add `json:` tags                |
| `Now` set in `NewCtx`                        | stale timestamp                  | Set `ctx.Now = time.Now()` in handler after copying |
| Store field not initialized in `StoreHolder` | nil pointer panic                | Always init every store in `NewStoreHolder`         |
| DB not passed into store                     | nil pointer panic on first query | Pass `*gorm.DB` into every store constructor        |
