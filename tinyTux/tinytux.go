package main

import (
    "github.com/ant0ine/go-json-rest/rest"
    "bitbucket.org/tebeka/base62"
    "gopkg.in/gorp.v1"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "log"
    "fmt"
    "net/http"
    "time"
)

var dbmap = initDb()

func main() {
    // delete any existing rows
    err := dbmap.TruncateTables()
    checkErr(err, "TruncateTables failed")

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
    fmt.Println("App started!")
    log.Fatal(http.ListenAndServe(":9998", api.MakeHandler()))

    dbmap.Db.Close()
}

type Url struct {
    URL string
}

func GetURL(w rest.ResponseWriter, r *rest.Request) {
    urlCode := r.PathParam("code")
    fmt.Println(urlCode)
    http.Redirect(w.(http.ResponseWriter), r.Request, "https://google.com/", 301)

}

func ShortenURL(w rest.ResponseWriter, r *rest.Request) {
    newURL := Url{}
    err := r.DecodeJsonPayload(&newURL)
    if err != nil {
        rest.Error(w, "link required", 400)
        return
    }
    
    // Create the row
    link := newLink(newURL.URL)
    err = dbmap.Insert(&link)
    checkErr(err, "Insert failed")

    // Encode the id with base62
    urlCode := base62.Encode(uint64(link.Id))

    w.WriteJson(map[string]string{"url": "tux.sh/l/"+urlCode})
}

func newLink(url string) Link {
    return Link{
        Created: time.Now().UnixNano(),
        URL: url,
    }
}

func initDb() *gorp.DbMap {
    // Connect to db using stdlib sql driver
    db, err := sql.Open("mysql", "test:testPASS@/tux_sh")
    checkErr(err, "sql.Open failed")

    dbmap := gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

    // Add the table
    dbmap.AddTableWithName(Link{}, "links").SetKeys(true, "Id")

    // Create the tables
    err = dbmap.CreateTablesIfNotExists()
    checkErr(err, "Create tables failed")

    return &dbmap
}

func checkErr(err error, msg string) {
    if err != nil {
        log.Fatalln(msg, err)
    }
}
