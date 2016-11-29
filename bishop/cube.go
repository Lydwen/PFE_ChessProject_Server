package main

import (
	"fmt"
)

func initCube() {
    cube = make(map[int][SIZE * SIZE]int)
    var face [SIZE * SIZE]int
    for j := 0; j < SIZE * SIZE; j++ {
        face[j] = 0
    }
    for i := 0; i < 6; i++ {
        cube[i] = face
    }
    fmt.Println(cube)
}


func updateCube(faceO int, xo int, yo int, faceN int, xn int, yn int) {
    coordO := yo * SIZE + xo
    fo := cube[faceO]
    fo[coordO] = 0
    cube[faceO] = fo
    fmt.Println(cube)

    coordN := yn * SIZE + xn
    fn := cube[faceN]
    fn[coordN] = 1
    cube[faceN] = fn
    fmt.Println(cube)
}