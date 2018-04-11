package main

import (
    "fmt"
    // "strings"
    "strconv"
)

const NOP,NEG,ADD,SUB,MUL,DIV,MOV,JMP = 1,2,3,4,5,6,7,8
const AAT,NR1,NR2,UPP,RGT,DWN,LFT = '@','\'','"','^','>','v','<'

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
    args []arg
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

func (c *core) eval_arg(exp arg, b *board) (val_type string, arg int, addr int, instr cmd, val int) { //evaluates an l_exp
    val_type = "ARG"
    val = exp.val
    for _, p := range exp.prefs {
        if val_type == "ARG" {
            arg = -1
            switch p {
            case AAT:
                addr = c.this
                instr = b.cores[addr].code[val]
                val_type = "CMD"
            case UPP:
                addr = c.up
                instr = b.cores[addr].code[val]
                val_type = "CMD"
            case RGT:
                addr = c.right
                instr = b.cores[addr].code[val]
                val_type = "CMD"
            case DWN:
                addr = c.down
                instr = b.cores[addr].code[val]
                val_type = "CMD"
            case LFT:
                addr = c.left
                instr = b.cores[addr].code[val]
                val_type = "CMD"
            default:
                panic("warning: pref"+strconv.QuoteRune(p)+"called with type ARG")
            }
        } else if val_type == "CMD" {
            switch p {
            case NR1:
                if addr>=0 {
                    val = b.cores[addr].code[val].args[0].val
                    arg = 0
                    val_type = "ARG"
                } else {
                    val = -1
                    arg = -1
                    val_type = "ARG"
                    return
                }
            case NR2:
                if addr>=0 {
                    val = b.cores[addr].code[val].args[1].val
                    arg = 1
                    val_type = "ARG"
                } else {
                    val = -1
                    arg = -1
                    val_type = "ARG"
                    return
                }
            default:
                panic("warning: pref"+strconv.QuoteRune(p)+"called with false type ADDR")
            }
        } else {
            panic("can't call any pref on"+strconv.QuoteRune(p))
        }
    }
    return
}

func (c core) tick(new_b *board, b *board) core { //executes the current instruction
    switch instr := c.code[c.pc]; instr.id {
    case NOP:
        return c
    case NEG:
        if val_type,arg_num,core_index,_,cmd_index := c.eval_arg(instr.args[1], b); val_type=="ARG" {
            if b.cores[core_index].code[cmd_index].args[arg_num].val > 0 {
                new_b.cores[core_index].code[cmd_index].args[arg_num].val = 0
            } else {
                new_b.cores[core_index].code[cmd_index].args[arg_num].val = 1
            }
        } else {fmt.Println("NEG arg2 needs to be of type ARG in line",c.pc); return c}
    case ADD:
        fmt.Println(instr.args)
        if val_type2,arg_num,core_index,_,cmd_index := c.eval_arg(instr.args[1], b); val_type2=="ARG" {
            if val_type1,_,_,_,to_add := c.eval_arg(instr.args[0], b); val_type1=="ARG" {
                new_b.cores[core_index].code[cmd_index].args[arg_num].val += to_add
            } else {fmt.Println("ADD arg1 needs to be of type ARG in line",c.pc); return c}
        } else {fmt.Println("ADD arg2 needs to be of type ARG in line",c.pc); return c}
    case SUB:
        if val_type2,arg_num,core_index,_,cmd_index := c.eval_arg(instr.args[1], b); val_type2=="ARG" {
            if val_type1,_,_,_,to_sub := c.eval_arg(instr.args[0], b); val_type1=="ARG" {
                new_b.cores[core_index].code[cmd_index].args[arg_num].val -= to_sub
            } else {fmt.Println("SUB arg1 needs to be of type ARG in line",c.pc); return c}
        } else {fmt.Println("SUB arg2 needs to be of type ARG in line",c.pc); return c}
    case MUL:
        if val_type2,arg_num,core_index,_,cmd_index := c.eval_arg(instr.args[1], b); val_type2=="ARG" {
            if val_type1,_,_,_,to_mul := c.eval_arg(instr.args[0], b); val_type1=="ARG" {
                new_b.cores[core_index].code[cmd_index].args[arg_num].val *= to_mul
            } else {fmt.Println("MUL arg1 needs to be of type ARG in line",c.pc); return c}
        } else {fmt.Println("MUL arg2 needs to be of type ARG in line",c.pc); return c}
    case DIV:
        if val_type2,arg_num,core_index,_,cmd_index := c.eval_arg(instr.args[1], b); val_type2=="ARG" {
            if val_type1,_,_,_,to_div := c.eval_arg(instr.args[0], b); val_type1=="ARG" {
                new_b.cores[core_index].code[cmd_index].args[arg_num].val /= to_div
            } else {fmt.Println("DIV arg1 needs to be of type ARG in line",c.pc); return c}
        } else {fmt.Println("DIV arg2 needs to be of type ARG in line",c.pc); return c}
    case MOV:
        if val_type2,arg_num,core_index,_,cmd_index := c.eval_arg(instr.args[1], b); val_type2=="ARG" {
            if val_type1,_,_,_,to_mov:= c.eval_arg(instr.args[0], b); val_type1=="ARG" {
                new_b.cores[core_index].code[cmd_index].args[arg_num].val = to_mov
            } else {fmt.Println("Can't MOV command to arg in line",c.pc); return c}
        } else {
            if val_type1,_,_,to_mov,_:= c.eval_arg(instr.args[0], b); val_type1=="CMD" {
                new_b.cores[core_index].code[cmd_index] = to_mov
            } else {fmt.Println("Can't MOV arg to command in line",c.pc); return c}
        }
    case JMP:
        if val_type2,arg_num,core_index,_,cmd_index := c.eval_arg(instr.args[1], b); val_type2=="ARG" {
            if val_type1,_,_,_,to_jmp:= c.eval_arg(instr.args[0], b); val_type1=="ARG" {
                if new_b.cores[core_index].code[cmd_index].args[arg_num].val!=0 {c.pc = to_jmp-1}
            } else {fmt.Println("can't jump (set pc to) a command"); return c}
        } else {fmt.Println("can't read command as condition for jump"); return c}
    default:
        return c
    }
    c.pc = (c.pc+1)%len(c.code)
    return c
}

type board struct { //the whole board with multiple cores on it
    cores []core
}

func (b *board) run() { //lets every core execute it's next instruction
    fmt.Println(b)
    new_b := *b
    for i, c := range b.cores {
        new_b.cores[i] = c.tick(&new_b,b)
    }
    *b = new_b
}

func build_board() (new_b board) { //builds a board from a string
    arg0 := arg{"",1}
    arg1 := arg{"^'",0}
    arg2 := arg{"v'",0}
    args1 := []arg{arg0,arg1}
    args2 := []arg{arg0,arg2}
    cmd1 := cmd{cmd_ids["ADD"],args1}
    cmd2 := cmd{cmd_ids["ADD"],args2}
    code1 := []cmd{cmd1}
    code2 := []cmd{cmd2}
    core1 := core{code1,0,1,0,0,0,0}
    core2 := core{code2,1,1,1,0,1,0}
    cores1 := []core{core1,core2}
    return board{cores1}
}

func main() {
    testboard := build_board()
    for {
        testboard.run()
    }
}
