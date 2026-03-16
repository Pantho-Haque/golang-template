# hubble

## With echo/v4

```go
package main

import (
    "github.com/labstack/echo/v4"
    "github.com/prometheus/client_golang/prometheus/promhttp"

    "magic.pathao.com/beatbox/hubble"
    hubble_echo "magic.pathao.com/beatbox/hubble/echo"
)

func main() {
    c := echo.New()
    c.GET("/health", checkHealth)
    // define /metrics endpoint to expose prometheus metrics
    c.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
    // register metrics
    prom := hubble.NewAPM("test").Histogram(nil)

    // set middleware for echo
    c.Use(hubble_echo.Use(prom))

    c.Start("127.0.0.1:8080")
}

func checkHealth(c echo.Context) error {
    c.String(200, "ok")
    return nil
}
```


## With chi

```go
package main

import (
    "net/http"

    "github.com/go-chi/chi"
    "github.com/prometheus/client_golang/prometheus/promhttp"

    "magic.pathao.com/beatbox/hubble"
    hubble_chi "magic.pathao.com/beatbox/hubble/chi"
)

func main() {
    c := chi.NewMux()

    // register metrics
    prom := hubble.NewAPM("test").Histogram(nil)
    // set middleware for chi
    c.Use(hubble_chi.Use(prom))
    // define /metrics endpoint to expose prometheus metrics
    c.Mount("/metrics", promhttp.Handler())

    c.Get("/health", checkHealth)

    server := &http.Server{
        Addr:    "127.0.0.1:8080",
        Handler: c,
    }

    _ = server.ListenAndServe()
}

func checkHealth(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
}
```

## With chi/v5

```go
package main

import (
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/prometheus/client_golang/prometheus/promhttp"

    "magic.pathao.com/beatbox/hubble"
    hubble_chi "magic.pathao.com/beatbox/hubble/chi/v5"
)

func main() {
    c := chi.NewMux()

    // register metrics
    prom := hubble.NewAPM("test").Histogram(nil)
    // set middleware for chi
    c.Use(hubble_chi.Use(prom))
    // define /metrics endpoint to expose prometheus metrics
    c.Mount("/metrics", promhttp.Handler())

    c.Get("/health", checkHealth)

    server := &http.Server{
        Addr:    "127.0.0.1:8080",
        Handler: c,
    }

    _ = server.ListenAndServe()
}

func checkHealth(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
}
```

## With fiber/v2

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/prometheus/client_golang/prometheus/promhttp"

    "magic.pathao.com/beatbox/hubble"
    hubble_fiber "magic.pathao.com/beatbox/hubble/fiber/v2"
)

func main() {
    app := fiber.New()

    // register metrics
    prom := hubble.NewAPM("test").Histogram(nil)
    // set middleware for fiber
    app.Use(hubble_fiber.Use(prom))
    // define /metrics endpoint to expose prometheus metrics
    app.Mount("/metrics", promhttp.Handler())

    app.Get("/health", func(c *fiber.Ctx) error {
		err := c.SendStatus(200)
		if err != nil {
			return err
		}
		
		return nil
	})

	_ = app.Listen("127.0.0.1:8080")
}

```