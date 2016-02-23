package main

import (
    "github.com/ant0ine/go-json-rest/rest"
    "log"
    "fmt"
    "net/http"
)

func main() {
    api := rest.NewApi()
    api.Use(rest.DefaultDevStack...)
    router, err := rest.MakeRouter(
        rest.Get("/l/:code", GetURL),
        rest.Post("/l", ShortenURL),
    )
    if err != nil {
        log.Fatal(err)
    }
    api.SetApp(router)
    log.Fatal(http.ListenAndServe(":9998", api.MakeHandler()))
}

type Link struct {
    url string
}

func GetURL(w rest.ResponseWriter, r *rest.Request) {
    urlCode := r.PathParam("code")
    fmt.Println(urlCode)
    http.Redirect(w.(http.ResponseWriter), r.Request, "https://google.com/", 301)

}

func ShortenURL(w rest.ResponseWriter, r *rest.Request) {
    link := Link{}
    err := r.DecodeJsonPayload(&link)
    if err != nil {
        rest.Error(w, "link required", 400)
        return
    }
    w.WriteJson(map[string]string{"url": "test"})
}
