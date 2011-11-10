package main

import (
    "github.com/hoisie/web.go"
    "strings"
    "godis"
)

var counter = 0

var redis = godis.New("", 0, "")

func resolve(ctx *web.Context, short string) {
    redirect, _ := redis.Get(short)
    ctx.Redirect(302, redirect.String())
}

func store(ctx *web.Context, url string) string {
    if ! strings.HasPrefix(url, "http"){
        url = "http://" + url
    }
    request := ctx.Request
    counter += 1
    encoded := encode(counter)
    redis.Set(encoded, url)
    return request.Proto + request.Host + "/" + encoded + "\n"

}


func encode(number int) string{
    const symbols string = "0123456789abcdefghijklmnopqrsuvwxyzABCDEFGHIJKLMNOPQRSTUVXYZ"
    const base int = len(symbols)
    rest := number % base
    result := string(symbols[rest])
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

