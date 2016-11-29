package main

import (
    "net"
    "bufio"
    "log"
    "fmt"
    "encoding/json"
    //"os"
)


var cube map[int][64]int

//{"face": 1,"x": 7,"y":7}
type Init struct {
    Face    int     `json:"face"`
    X       int     `json:"x"`
    Y       int     `json:"y"`
}

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
        res := Move{}
        json.Unmarshal([]byte(msg), &res)
        fmt.Println(res)
        updateCube(res.FaceO, res.Xo, res.Yo, res.FaceN, res.Xn, res.Yn)
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
    cube = make(map[int][64]int)
    for i := 0; i < 6; i++ {
        cube[i] = [64]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 
            0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 
            0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 
            0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 
            0, 0, 0, 0, 0}
    }
    fmt.Println(cube)
}


func updateCube(faceO int, xo int, yo int, faceN int, xn int, yn int) {
    coordO := yo * 8 + xo
    fo := cube[faceO]
    fo[coordO] = 0
    cube[faceO] = fo
    fmt.Println(cube)

    coordN := yn * 8 + xn
    fn := cube[faceN]
    fn[coordN] = 1
    cube[faceN] = fn
    fmt.Println(cube)
}

