package main

import (
    "database/sql"
    "log"
    "net/http"
    "encoding/json"
    "time"

    "github.com/gorilla/mux"
    "bitbucket.org/tebeka/base62"
    "gopkg.in/gorp.v1"
    _ "github.com/go-sql-driver/mysql"
)

var dbmap = initDb()

func main() {
    defer dbmap.Db.Close()

    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/{code}/", GetURL).Methods("GET")
    router.HandleFunc("/", ShortenURL).Methods("POST")

    log.Fatal(http.ListenAndServe(":9998", router))
}

type Url struct {
    URL string
}

func GetURL(w http.ResponseWriter, r *http.Request) {
    link := Link{}
    urlCode := base62.Decode(mux.Vars(r)["code"])
    err := dbmap.SelectOne(&link, "SELECT * FROM links WHERE id = :id",
        map[string]interface{} {"id": urlCode})
    checkErr(err, "Failed to find link with urlCode")
    http.Redirect(w, r, link.URL, 301)

}

func ShortenURL(w http.ResponseWriter, r *http.Request) {
    link := Url{}
    err := json.NewDecoder(r.Body).Decode(&link) 
    if err != nil {
        http.Error(w, "link required", 400)
        return
    }
    
    // Create the row
    shortLink := newLink(link.URL)
    err = dbmap.Insert(&shortLink)
    checkErr(err, "Insert failed")

    // Encode the id with base62
    url := Url{
        URL: "https://tux.sh/l/" + base62.Encode(uint64(shortLink.Id)),
    }
    
    json.NewEncoder(w).Encode(url)
}

func newLink(url string) Link {
    return Link{
        Created: time.Now().Unix(),
        URL: url,
    }
}

func initDb() *gorp.DbMap {
    // Connect to db using stdlib sql driver
    db, err := sql.Open("mysql", "test:testPASS@/tux_sh")
    checkErr(err, "sql.Open failed")
    dbmap := gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

    dbmap.AddTableWithName(Link{}, "links").SetKeys(true, "Id")

    err = dbmap.CreateTablesIfNotExists()
    checkErr(err, "Create tables failed")

    return &dbmap
}

func checkErr(err error, msg string) {
    if err != nil {
        log.Fatalln(msg, err)
    }
}
