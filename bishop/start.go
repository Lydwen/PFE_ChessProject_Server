package main

import (
    "net"
    "bufio"
    "log"
    "fmt"
    "encoding/json"
    //"os"
)


var cube map[int][49]int

//{"face": 1,"x": 1,"y":1}
type Msg struct {
    Face    int     `json:"face"`
    X       int     `json:"x"`
    Y       int     `json:"y"`
}

func handleConnection(conn net.Conn, id int, connC chan net.Conn, game chan string) {
    defer conn.Close()
    reader := bufio.NewReader(conn)
    // nickname,err := reader.ReadString('\n')
    // if err != nil {
    //     log.Println(err)
    //     return
    // }
    // fmt.Printf("j'ai dit à %s 'hello'\n", nickname)
    // fmt.Fprintf(conn, "hello %s", nickname)

    c := make(chan string)
    if id == 0 {
        go listenP1_writeP2(c, connC, conn, game)
    } else {
        go listenP2_writeP1(c, connC, conn, game)
    }

    for {
        line,err := reader.ReadString('\n')
        if err != nil {
            log.Println(err)
            return
        }
        line = line
        c <- line
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
        res := Msg{}
        json.Unmarshal([]byte(msg), &res)
        fmt.Println(res)
        updateCube(res.Face, res.X, res.Y)
    }
}

func main() {
    initCube()
    listener, err := net.Listen("tcp", "localhost:1234")
    //listener, err := net.Listen("tcp", "localhost:"+os.Getenv("PORT"))
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


//      >>>>>>>>>>>> gestion CUBE <<<<<<<<<<<<<

func initCube() {
    cube = make(map[int][49]int)
    for i := 0; i < 6; i++ {
        cube[i] = [49]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 
            0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 
            0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 
            0, 0, 0, 0, 0, 0}
    }
    fmt.Println(cube)
}


func updateCube(face int, x int, y int) {
    coord := y * 7 + x
    f := cube[face]
    f[coord] = 1
    cube[face] = f
    fmt.Println(cube)
}

