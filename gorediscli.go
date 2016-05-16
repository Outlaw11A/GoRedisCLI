package main

import (
    "github.com/alecthomas/kingpin"
    "log"
    "net"
    "fmt"
    "bufio"
    "os"
)

var (
    Host = kingpin.Flag("host", "Host").Short('h').String()
    Port = kingpin.Flag("port", "Port").Short('p').Int()
)

func main() {
    log.SetFlags(0)

    // Get the host and port
    kingpin.Parse()
    if len(*Host) == 0 {
        *Host = "127.0.0.1"
    }
    if *Port == 0 {
        *Port = 6379
    }

    conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", *Host, *Port))
    if err != nil {
        log.Fatalln(err.Error())
    }

    for {
        reader := bufio.NewReader(os.Stdin)
        fmt.Print("> ")

        input, err := reader.ReadBytes('\n')
        if err != nil {
            log.Fatalln(err.Error())
        }

        fmt.Fprint(conn, fmt.Sprintf("%s\n", input))
        reply, err := bufio.NewReader(conn).ReadString('\n')
        fmt.Println(reply)
    }
}