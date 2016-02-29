package main

type Link struct {
    Id      uint64   `db:"id",primarykey,autoincrement`
    Created int64   `db:"created"`
    URL     string  `db:"url",size:1024`
}
