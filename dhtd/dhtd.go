package main

import (
    "github.com/JGrotex/GrovePi/Software/Go/grovepi"
    "fmt"
    "net/http"
    "log"
    "math"
    "errors"
)

const (
  dht_port = grovepi.D2
  http_port = ":8080"
)



type dht_data struct {
    temperature, humidity float64
}



func get_dht_data() (dht_data, error) {
    var g grovepi.GrovePi
    g = *grovepi.InitGrovePi(0x04)
    defer g.CloseDevice()

    t, h, err := g.ReadDHT(dht_port)

    if err != nil {
        log.Println(err)
        return dht_data{}, errors.New(fmt.Sprintf("%s", err))
    }

    return dht_data{temperature: float64(t), humidity: float64(h)}, nil
}



func http_handler_f(w http.ResponseWriter, r *http.Request) {
    data, err := get_dht_data()
    if err != nil {
        return
    }
    fmt.Fprintf(w, `{ "temperature": %.0f, "humidity": %.0f }`, math.Round(data.temperature*1.8+32), data.humidity)
}



func http_handler_c(w http.ResponseWriter, r *http.Request) {
    data, err := get_dht_data()
    if err != nil {
        return
    }
    fmt.Fprintf(w, `{ "temperature": %.0f, "humidity": %.0f }`, data.temperature, data.humidity)
}



func start_http() {
    log.SetFlags(log.Flags()|log.Lshortfile)
    log.Println("Starting webserver")
    http.HandleFunc("/f", http_handler_f)
    http.HandleFunc("/c", http_handler_c)
    http.ListenAndServe(http_port, nil)
}



func main() {
    start_http()
}

