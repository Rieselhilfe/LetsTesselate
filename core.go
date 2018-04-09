package main

import (
    "fmt"
    "strings"
    "strconv"
)

type cmd struct { //a single command with a name and two arguments
    id byte
    arg1_prefs string
    arg1_val int
    arg2_prefs string
    arg2_val int
}

type core struct { //one core, that has a program out of commands and executes them
    code []cmd
    up *core
    right *core
    down *core
    left *core
    pc int
}

func (c core) tick() (c core) { //execute the current instruction
    switch code[pc].id {
    default:
        return c
    }
}

type board struct { //the whole board with multiple cores on it
    cores []core
}

func (b *board) tick() { //let every core execute the next instruction
    new_board := *b
    for i, c := range b.cores {
        new_board[i] = c.tick()
    }
}

func main() {

}
