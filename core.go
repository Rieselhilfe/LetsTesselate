package main

import (
    "fmt"
    // "strings"
    "strconv"
)

const NOP,NEG,ADD,SUB,MUL,DIV,MOV,JMP = 1,2,3,4,5,6,7,8
const AAT,NR1,NR2,IRN,UPP,RGT,DWN,LFT = '@','#','\'','/','^','>','v','<'

var cmd_ids = map[string]byte{
    "NOP": NOP,
    "NEG": NEG,
    "ADD": ADD,
    "SUB": SUB,
    "MUL": MUL,
    "DIV": DIV,
    "MOV": MOV,
    "JMP": JMP,
}

type arg struct { //an arg, consisting out of prefixes and a numeric value
    prefs string
    val int
}

type cmd struct { //a single command with a name and two arguments
    id byte
    arg1 arg
    arg2 arg
}

type core struct { //one core, containing a program out of commands
    code []cmd
    this int
    up int
    right int
    down int
    left int
    pc int
}

func (c *core) eval_arg(exp arg, b *board) (instr cmd, val int, addr int, val_type string) { //evaluates an l_exp
    val_type = "NUM"
    for _, p := range exp.prefs {
        if val_type == "NUM" {
            switch p {
            case AAT:
                addr = c.this
                val_type = "ADDR"
            case UPP:
                addr = c.up
                val_type = "ADDR"
            case RGT:
                addr = c.right
                val_type = "ADDR"
            case DWN:
                addr = c.down
                val_type = "ADDR"
            case LFT:
                addr = c.left
                val_type = "ADDR"
            default:
                panic("warning: pref"+strconv.QuoteRune(p)+"called with false type")
            }
        } else if val_type == "ADDR" {
            switch p {
            case NR1:
                val = b.cores[addr].code[val].arg1.val
                val_type = "NUM"
            case NR2:
                val = b.cores[addr].code[val].arg2.val
                val_type = "NUM"
            case IRN:
                instr = b.cores[addr].code[val]
                val_type = "CMD"
            default:
                panic("warning: pref"+strconv.QuoteRune(p)+"called with false type")
            }
        } else {
            panic("can't call any pref on"+strconv.QuoteRune(p))
        }
    }
    return
}

func (c core) tick(b *board) core { //executes the current instruction
    switch instr := c.code[c.pc]; instr.id {
    case NOP:
        return c
    case NEG:
        if _,cmd_index,core_index,val_type := c.eval_arg(instr.arg2, b); val_type=="ADDR" {
            if b.cores[core_index].code[cmd_index].arg1.val > 0 {
                b.cores[core_index].code[cmd_index].arg1.val = 0
            } else {
                b.cores[core_index].code[cmd_index].arg1.val = 1
            }
        } else {panic("NEG needs ADDR")}
    default:
        return c
    }
    return c
}

type board struct { //the whole board with multiple cores on it
    cores []core
}

func (b *board) run() { //lets every core execute it's next instruction
    new_b := *b
    for i, c := range b.cores {
        new_b.cores[i] = c.tick(b)
    }
    *b = new_b
}

func build_board() (new_b board) { //builds a board from a string
    return
}

func main() {
    fmt.Println(build_board())
}
