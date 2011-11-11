package main

import (
    "web"
    "strings"
    "godis"
    "fmt"
)

const(
    SYMBOLS = "0123456789abcdefghijklmnopqrsuvwxyzABCDEFGHIJKLMNOPQRSTUVXYZ"
    COUNTER = "__counter__"
    HTTP = "http"
)

var (
    redis = godis.New("", 0, "")
)

func resolve(ctx *web.Context, short string) {
    redirect, _ := redis.Get(short)
    ctx.Redirect(302, redirect.String())
}

func store(ctx *web.Context, url string){
    if ! strings.HasPrefix(url, HTTP){
        url = fmt.Sprintf("%s://%s", HTTP, url)
    }
    ctr, _ := redis.Incr(COUNTER)
    encoded := encode(ctr)
    redis.Set(encoded, url)
    request := ctx.Request
    ctx.WriteString( fmt.Sprintf("%s://%s/%s\n", HTTP, request.Host, encoded))
}


func encode(number int64) string{
    const base = int64(len(SYMBOLS))
    rest := number % base
    result := string(SYMBOLS[rest])
    if number - rest != 0{
       newnumber := (number - rest ) / base
       result = encode(newnumber) + result
    }
    return result
}


func main() {
    web.Get("/store/(.*)", store)
    web.Get("/(.*)", resolve)
    web.Run("0.0.0.0:9999")
}

