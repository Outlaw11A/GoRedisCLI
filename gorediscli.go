package main

import (
    "github.com/alecthomas/kingpin"
    "log"
    "net"
    "fmt"
    "bufio"
    "os"
    "io"
    "strconv"
)

/**
Define the first byte prefixes for the RESP protocol
See: http://redis.io/topics/protocol
 */
const REDIS_STRING byte = '+'
const REDIS_ERROR byte = '-'
const REDIS_INTEGER byte = ':'
const REDIS_BULK_STRING byte = '$'
const REDIS_ARRAY byte = '*'

var (
    Host = kingpin.Flag("host", "Host").Short('h').String()
    Port = kingpin.Flag("port", "Port").Short('p').Int()
    Command = kingpin.Flag("command", "Command").Short('c').String()
)

type Reply struct {
    Type        byte
    Str         string
    Elements    []Reply
    Nil         bool
}

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
    var hostString string = fmt.Sprintf("%s:%d", *Host, *Port)

    var interactive bool = true
    if len(*Command) != 0 {
        interactive = false
    }

    conn, err := net.Dial("tcp", hostString)
    if err != nil {
        log.Fatalln(err.Error())
    }

    for {
        var input string
        if interactive {
            reader := bufio.NewReader(os.Stdin)
            fmt.Printf("%s> ", hostString)

            inputBytes, err := reader.ReadBytes('\n')
            if err != nil {
                log.Fatalln(err.Error())
            }
            input = string(inputBytes)
        } else {
            input = *Command
        }
        fmt.Fprint(conn, fmt.Sprintf("%s\n", input))
        printReply(readReply(conn))
        if !interactive {
            return
        }
    }
}

func printPrefix(prefix string, output string) {
    fmt.Printf("%s %s\n", prefix, output)
}

// Prints a reply recursively
func printReply(reply Reply) {
    switch reply.Type {
    case REDIS_STRING:
        fmt.Println(reply.Str)
        break
    case REDIS_INTEGER:
        printPrefix("(integer)", reply.Str)
    case REDIS_ERROR:
        printPrefix("(error)", reply.Str)
        break
    case REDIS_BULK_STRING:
        if reply.Nil {
            fmt.Println("(nil)")
        } else {
            fmt.Printf("\"%s\"\n", reply.Str)
        }
        break
    case REDIS_ARRAY:
        for i := 0; i < len(reply.Elements); i++ {
            printReply(reply.Elements[i])
        }
        break
    }
}

// Reads the reply
func readReply(conn io.Reader) Reply {
    scanner := bufio.NewScanner(conn)
    return replyBuilder(scanner)
}

// Builds the reply recursively
func replyBuilder(scanner *bufio.Scanner) Reply {
    var reply Reply
    for scanner.Scan() {
        bytes := scanner.Bytes()
        reply.Type = bytes[0]
        reply.Str = string(bytes[1:])
        reply.Nil = false
        switch reply.Type {
        case REDIS_STRING, REDIS_INTEGER, REDIS_ERROR:
            return reply
        case REDIS_BULK_STRING:
            if reply.Str == "-1" {
                reply.Nil = true
                reply.Str = ""
            } else {
                scanner.Scan()
                reply.Str = string(scanner.Bytes())
            }
            return reply
        case REDIS_ARRAY:
            amount, _ := strconv.Atoi(string(bytes[1:]))
            for i := 0; i < amount; i++ {
                reply.Elements = append(reply.Elements, replyBuilder(scanner))
            }
            return reply
        }
        return reply
    }
    return reply
}