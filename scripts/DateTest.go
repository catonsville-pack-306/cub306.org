package main

//go get github.com/arran4/golang-ical

import ("fmt"
    //"bufio"
    //"os"
    //"io"
    "time"
    //"io/ioutil"
    //"bytes"
    //"github.com/arran4/golang-ical"
    //"strings"
    //"flag"
    //"text/template"
    //"syscall"
    //"unsafe"
    )

func main() {
    raw := "20210530T232304Z"
    //time.LoadLocation("GMT-4")
    ans, _ := time.Parse("20060102T150405Z", raw)
    ans = ans.Add(time.Hour * -4)
    fmt.Println (ans)
}