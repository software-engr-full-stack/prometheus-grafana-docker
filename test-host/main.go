package main

import (
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "net/http"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func ping(w http.ResponseWriter, req *http.Request) {
    fmt.Fprint(w, "Howdy")
}

func handler(w http.ResponseWriter, req *http.Request) {
    min := 1000
    max := 2000
    fakeDelay := time.Duration(rand.Intn(max-min)+min) * time.Millisecond //nolint:gosec
    bt, err := json.Marshal([]map[string]any{
        map[string]any{"data": map[string]any{
            "fake-delay": fmt.Sprintf("%d seconds", fakeDelay/1000000000), //nolint:gomnd
        }},
    })
    if err != nil {
        eb, perr := json.Marshal(map[string]string{"error (json)": err.Error()})
        if perr != nil {
            panic(perr)
        }
        fmt.Fprint(w, string(eb))
        return
    }

    time.Sleep(fakeDelay)

    fmt.Fprint(w, string(bt))
}

func main() {
    reg := prometheus.NewRegistry()
    reg.MustRegister(totalRequests)
    reg.MustRegister(responseStatus)
    reg.MustRegister(httpDuration)

    mux := http.NewServeMux()

    mux.Handle("/api/v1/resources", prometheusMiddleware(http.HandlerFunc(handler)))
    mux.Handle("/ping", prometheusMiddleware(http.HandlerFunc(ping)))

    promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
    mux.Handle("/metrics", promHandler)

    server := &http.Server{
        Addr:              ":9100",
        ReadHeaderTimeout: 5 * time.Second,
        Handler:           mux,
    }

    err := server.ListenAndServe()
    if err != nil {
        log.Fatal(err)
    }
}
