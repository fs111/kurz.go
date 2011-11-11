package main

import (
    "web"
    "strings"
    "godis"
    "fmt"
)

const(
    // characters used for short-urls
    SYMBOLS = "0123456789abcdefghijklmnopqrsuvwxyzABCDEFGHIJKLMNOPQRSTUVXYZ"
    // special key in redis, that is our global counter
    COUNTER = "__counter__"
    HTTP = "http"
)

// connecting to redis on localhost, db with id 0 and no password
var (
    redis = godis.New("", 0, "")
)

// function to resolve a shorturl and redirect
func resolve(ctx *web.Context, short string) {
    redirect, _ := redis.Get(short)
    ctx.Redirect(302, redirect.String())
    // TODO needs error handling here
}

// function to shorten and store a url
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

// encodes a number into our *base* representation
// TODO can this be made better with some bitshifting?
func encode(number int64) string{
    const base = int64(len(SYMBOLS))
    rest := number % base
    // strings are a bit weird in go...
    result := string(SYMBOLS[rest])
    if number - rest != 0{
       newnumber := (number - rest ) / base
       result = encode(newnumber) + result
    }
    return result
}

// main function that inits the routes in web.go
func main() {
    web.Get("/store/(.*)", store)
    web.Get("/(.*)", resolve)
    web.Run("0.0.0.0:9999")
}

