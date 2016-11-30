package main

import (
    "net"
    "bufio"
    "log"
    "fmt"
    "encoding/json"
    "os"
)

const SIZE int = 8
var cube map[int][SIZE * SIZE]int

//{"face_old": 1, "x_old": 7, "y_old":7, "face_new": 1, "x_new": 7, "y_new":7}
type Move struct {
    FaceO    int     `json:"face_old"`
    Xo       int     `json:"x_old"`
    Yo       int     `json:"y_old"`
    FaceN    int     `json:"face_new"`
    Xn       int     `json:"x_new"`
    Yn       int     `json:"y_new"`
}

func handleConnection(conn net.Conn, id int, connC chan net.Conn, game chan string) {
    defer conn.Close()
    reader := bufio.NewReader(conn)
    c := make(chan string)
    
    if id == 0 {
        go listenP1_writeP2(c, connC, conn, game)
    } else {
        go listenP2_writeP1(c, connC, conn, game)
    }

    for {
        fmt.Fprintf(conn, "OK")
        line,err := reader.ReadBytes('}')
        stringLine := string(line[:])

        if err != nil {
            log.Println(err)
            return
        }
        c <- stringLine
    }
}

func listenP1_writeP2(c chan string, connC chan net.Conn, conn net.Conn, game chan string) {
    var connP2 net.Conn
    connP2 = <- connC
    connC <- conn

    for {
        msg := <- c
        fmt.Fprintf(connP2, msg)
        game <- msg
    }
}

func listenP2_writeP1(c chan string, connC chan net.Conn, conn net.Conn, game chan string) {
    var connP1 net.Conn
    connC <- conn
    connP1 = <- connC

    for {
        msg := <- c
        fmt.Fprintf(connP1, msg)
        game <- msg
    }
}

func save(game chan string) {
    for {
        msg := <- game
        fmt.Println(msg)
        res := Move{}
        json.Unmarshal([]byte(msg), &res)
        fmt.Println(res)
        updateCube(res.FaceO, res.Xo, res.Yo, res.FaceN, res.Xn, res.Yn)
    }
}

func main() {
    initCube()
    listener, err := net.Listen("tcp", "localhost:1234")
    listener, err := net.Listen("tcp", "localhost:"+os.Getenv("PORT"))
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("listening...")

    game := make(chan string)
    connC := make(chan net.Conn)
    acc := 0

    go save(game)

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Println(err)
            continue
        }
        go handleConnection(conn, acc, connC, game)
        acc = acc + 1
    }
}

