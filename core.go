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
    for _, p := range exp.prefs {
        if val_type == "ARG" {
            arg = 0
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
                val = b.cores[addr].code[val].args[1].val
                arg = 1
                val_type = "ARG"
            case NR2:
                val = b.cores[addr].code[val].args[2].val
                arg = 2
                val_type = "ARG"
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
        if val_type,arg_num,core_index,_,cmd_index := c.eval_arg(instr.args[2], b); val_type=="ARG" {
            if b.cores[core_index].code[cmd_index].args[arg_num].val > 0 {
                new_b.cores[core_index].code[cmd_index].args[arg_num].val = 0
            } else {
                new_b.cores[core_index].code[cmd_index].args[arg_num].val = 1
            }
        } else {panic("NEG arg2 needs to be of type ARG")}
    case ADD:
        if val_type2,arg_num,core_index,_,cmd_index := c.eval_arg(instr.args[2], b); val_type2=="ARG" {
            if val_type1,_,_,_,to_add := c.eval_arg(instr.args[1], b); val_type1=="ARG" {
                new_b.cores[core_index].code[cmd_index].args[arg_num].val += to_add
            } else {panic("ADD arg1 needs to be of type ARG")}
        } else {panic("ADD arg2 needs to be of type ARG")}
    case SUB:
        if val_type2,arg_num,core_index,_,cmd_index := c.eval_arg(instr.args[2], b); val_type2=="ARG" {
            if val_type1,_,_,_,to_sub := c.eval_arg(instr.args[1], b); val_type1=="ARG" {
                new_b.cores[core_index].code[cmd_index].args[arg_num].val -= to_sub
            } else {panic("SUB arg1 needs to be of type ARG")}
        } else {panic("SUB arg2 needs to be of type ARG")}
    case MUL:
        if val_type2,arg_num,core_index,_,cmd_index := c.eval_arg(instr.args[2], b); val_type2=="ARG" {
            if val_type1,_,_,_,to_mul := c.eval_arg(instr.args[1], b); val_type1=="ARG" {
                new_b.cores[core_index].code[cmd_index].args[arg_num].val *= to_mul
            } else {panic("MUL arg1 needs to be of type ARG")}
        } else {panic("MUL arg2 needs to be of type ARG")}
    case DIV:
        if val_type2,arg_num,core_index,_,cmd_index := c.eval_arg(instr.args[2], b); val_type2=="ARG" {
            if val_type1,_,_,_,to_div := c.eval_arg(instr.args[1], b); val_type1=="ARG" {
                new_b.cores[core_index].code[cmd_index].args[arg_num].val /= to_div
            } else {panic("DIV arg1 needs to be of type ARG")}
        } else {panic("DIV arg2 needs to be of type ARG")}
    case MOV:
        if val_type2,arg_num,core_index,_,cmd_index := c.eval_arg(instr.args[2], b); val_type2=="ARG" {
            if val_type1,_,_,_,to_mov:= c.eval_arg(instr.args[1], b); val_type1=="ARG" {
                new_b.cores[core_index].code[cmd_index].args[arg_num].val = to_mov
            } else {panic("Can't MOV command to arg")}
        } else {
            if val_type1,_,_,to_mov,_:= c.eval_arg(instr.args[1], b); val_type1=="CMD" {
                new_b.cores[core_index].code[cmd_index] = to_mov
            } else {panic("Can't MOV arg to command")}
        }
    case JMP:
        if val_type2,arg_num,core_index,_,cmd_index := c.eval_arg(instr.args[2], b); val_type2=="ARG" {
            if val_type1,_,_,_,to_jmp:= c.eval_arg(instr.args[1], b); val_type1=="ARG" {
                if new_b.cores[core_index].code[cmd_index].args[arg_num].val!=0 {c.pc = to_jmp}
            } else {panic("can't jump (set pc to) a command")}
        } else {panic("can't read command as condition for jump")}
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
        new_b.cores[i] = c.tick(&new_b,b)
    }
    *b = new_b
}

func build_board() (new_b board) { //builds a board from a string
    return
}

func main() {
    fmt.Println(build_board())
}
