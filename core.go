package main

import (
    "fmt"
    "strings"
    "strconv"
    "time"
    "io/ioutil"
    "os"
    "regexp"
    // "github.com/gen2brain/raylib-go/raylib"
)

const (
    screen_width = 800
	screen_height = 600

    NOP = 0
    NEG = 1
    ADD = 2
    SUB = 3
    MUL = 4
    DIV = 5
    MOV = 6
    JMP = 7
    OUT = 8
    AND = 9
    IOR = 10
    XOR = 11

    AAT = '@'
    NR1 = '\''
    NR2 = '"'
    UPP = '^'
    RGT = '>'
    DWN = 'v'
    LFT = '<'

    INACTIVE_COLOR = "#1e1e37"
)

var (
    cycles = 0
    text_mode = true
)

//BACKEND
//BACKEND
//BACKEND
//BACKEND
//BACKEND
//BACKEND

func mod(a int, b int) int {
    if a > 0 {
        return a%b
    } else if a < 0 {
        return b-(-1*a)%b
    }
    return 0
}

func deepcopy_board(a *board, b *board) {
    b.Cores = make([]core, len(a.Cores))
    for i, ithcore := range a.Cores {
        b.Cores[i] = ithcore
        b.Cores[i].Code = make([]cmd, len(ithcore.Code))
        for j, jthcmd := range ithcore.Code {
            b.Cores[i].Code[j].Id = jthcmd.Id
            b.Cores[i].Code[j].Args = make([]arg, len(jthcmd.Args))
            for k,ktharg := range jthcmd.Args {
                b.Cores[i].Code[j].Args[k] = ktharg
            }
        }
    }
}


var cmd_ids = map[string]byte{
    "NOP": NOP,
    "NEG": NEG,
    "ADD": ADD,
    "SUB": SUB,
    "MUL": MUL,
    "DIV": DIV,
    "MOV": MOV,
    "JMP": JMP,
    "OUT": OUT,
}

type arg struct { //an arg, consisting out of prefixes and a numeric value
    Prefs string
    Val int
}

type cmd struct { //a single command with a name and two arguments
    Id byte
    Args []arg
}

type core struct { //one core, containing a program out of commands
    Code []cmd
    Pc int
    This int
    Up int
    Right int
    Down int
    Left int
    Color string
    Active bool
}

func (c *core) eval_arg(command_pos int, arg_pos int, b *board) (val_type string, arg int,
                        addr int, instr cmd, val int) { //evaluates prefixes
    val_type = "ARG"
    exp := c.Code[command_pos].Args[arg_pos]
    val = exp.Val
    addr = c.This
    arg = arg_pos
    for _, p := range exp.Prefs {
        if val_type == "ARG" {
            switch p {
            case AAT:
                addr = c.This
                instr = b.Cores[addr].Code[val]
                val_type = "CMD"
            case UPP:
                addr = c.Up
                instr = b.Cores[addr].Code[val]
                val_type = "CMD"
            case RGT:
                addr = c.Right
                instr = b.Cores[addr].Code[val]
                val_type = "CMD"
            case DWN:
                addr = c.Down
                instr = b.Cores[addr].Code[val]
                val_type = "CMD"
            case LFT:
                addr = c.Left
                instr = b.Cores[addr].Code[val]
                val_type = "CMD"
            default:
                panic("warning: pref"+strconv.QuoteRune(p)+"called with type ARG")
                return
            }
        } else if val_type == "CMD" {
            switch p {
            case NR1:
                if addr>=0 {
                    val = b.Cores[addr].Code[val].Args[0].Val
                    arg = 0
                    val_type = "ARG"
                } else {panic("addr under 0")}
            case NR2:
                if addr>=0 {
                    val = b.Cores[addr].Code[val].Args[1].Val
                    arg = 1
                    val_type = "ARG"
                } else {panic("addr under 0")}
            default:
                panic("warning: pref"+strconv.QuoteRune(p)+"called with false type ADDR")
                return
            }
        } else {
            panic("can't call any pref on"+strconv.QuoteRune(p))
        }
    }
    return
}

func (c core) tick(new_b *board, b *board) { //executes the current instruction
    r_val := c.Code[c.Pc].Args[1].Val
    switch instr := c.Code[c.Pc]; instr.Id {
    case NOP:
    case NEG:
        if val_type,arg_num,core_index,_,_ := c.eval_arg(c.Pc, 1, b); val_type=="ARG" {
            if b.Cores[core_index].Code[r_val].Args[arg_num].Val > 0 {
                new_b.Cores[core_index].Code[r_val].Args[arg_num].Val = 0
            } else {
                new_b.Cores[core_index].Code[r_val].Args[arg_num].Val = 1
            }
        } else {fmt.Println("NEG arg2 needs to be of type ARG in line",c.Pc); panic("")}
    case ADD:
        if val_type2,arg_num,core_index,_,_ := c.eval_arg(c.Pc, 1, b); val_type2=="ARG" {
            if val_type1,_,_,_,to_add := c.eval_arg(c.Pc, 0, b); val_type1=="ARG" {
                new_b.Cores[core_index].Code[r_val].Args[arg_num].Val += to_add
            } else {fmt.Println("ADD arg1 needs to be of type ARG in line",c.Pc); panic("")}
        } else {fmt.Println("ADD arg2 needs to be of type ARG in line",c.Pc); panic("")}
    case SUB:
        if val_type2,arg_num,core_index,_,_ := c.eval_arg(c.Pc, 1, b); val_type2=="ARG" {
            if val_type1,_,_,_,to_sub := c.eval_arg(c.Pc, 0, b); val_type1=="ARG" {
                new_b.Cores[core_index].Code[r_val].Args[arg_num].Val -= to_sub
            } else {fmt.Println("SUB arg1 needs to be of type ARG in line",c.Pc); panic("")}
        } else {fmt.Println("SUB arg2 needs to be of type ARG in line",c.Pc); panic("")}
    case MUL:
        if val_type2,arg_num,core_index,_,_ := c.eval_arg(c.Pc, 1, b); val_type2=="ARG" {
            if val_type1,_,_,_,to_mul := c.eval_arg(c.Pc, 0, b); val_type1=="ARG" {
                new_b.Cores[core_index].Code[r_val].Args[arg_num].Val *= to_mul
            } else {fmt.Println("MUL arg1 needs to be of type ARG in line",c.Pc); panic("")}
        } else {fmt.Println("MUL arg2 needs to be of type ARG in line",c.Pc); panic("")}
    case DIV:
        if val_type2,arg_num,core_index,_,_ := c.eval_arg(c.Pc, 1, b); val_type2=="ARG" {
            if val_type1,_,_,_,to_div := c.eval_arg(c.Pc, 0, b); val_type1=="ARG" {
                new_b.Cores[core_index].Code[r_val].Args[arg_num].Val /= to_div
            } else {fmt.Println("DIV arg1 needs to be of type ARG in line",c.Pc); panic("")}
        } else {fmt.Println("DIV arg2 needs to be of type ARG in line",c.Pc); panic("")}
    case MOV:
        if val_type2,arg_num,core_index,_,_ := c.eval_arg(c.Pc, 1, b); val_type2=="ARG" {
            if val_type1,_,_,_,to_mov:= c.eval_arg(c.Pc, 0, b); val_type1=="ARG" {
                new_b.Cores[core_index].Code[r_val].Args[arg_num].Val = to_mov
            } else {fmt.Println("Can't MOV command to arg in line",c.Pc); panic("")}
        } else {
            if val_type1,_,_,to_mov,_:= c.eval_arg(c.Pc, 0, b); val_type1=="CMD" {
                new_b.Cores[core_index].Code[r_val] = to_mov
                // fmt.Println(new_b.Cores[core_index].Code[r_val], to_mov)
                if !new_b.Cores[core_index].Active {
                    new_b.Cores[core_index].Active = true
                }
            } else {fmt.Println("Can't MOV arg to command in line",c.Pc); panic("")}
        }
    case JMP:
        if val_type1,arg_num,core_index,_,_ := c.eval_arg(c.Pc, 0, b); val_type1=="ARG" {
            if val_type2,_,_,_,to_jmp:= c.eval_arg(c.Pc, 1, b); val_type2=="ARG" {
                if b.Cores[core_index].Code[r_val].Args[arg_num].Val!=0 {c.Pc = int(to_jmp-1)}
            } else {fmt.Println("can't jump (set pc to) a command"); panic("")}
        } else {fmt.Println("can't read command as condition for jump"); panic("")}
    case OUT:
        if val_type1,arg_num,core_index,_,_ := c.eval_arg(c.Pc, 0, b); val_type1=="ARG" {
            if val_type2,_,_,_,to_print:= c.eval_arg(c.Pc, 1, b); val_type2=="ARG" {
                if b.Cores[core_index].Code[r_val].Args[arg_num].Val!=0 {
                    fmt.Println("OUT: t:",cycles,"n:",c.This,"m:",to_print)
                }
            } else {fmt.Println("can't jump (set pc to) a command"); panic("")}
        } else {fmt.Println("can't read command as condition for jump"); panic("")}
    default:
        panic("unknown command")
    }
    new_b.Cores[c.This].Pc = (c.Pc+1)%len(c.Code)
}

type board struct { //the whole board with multiple cores on it
    Cores []core
}

func (b *board) run(verbose bool) { //lets every core execute it's next instruction
    new_b := board{}
    deepcopy_board(b,&new_b) //bottleneck 1
    // fmt.Println(new_b,"\n",*b,"\n")
    for i, c := range b.Cores {
        if true {
            if verbose {
                fmt.Println(i,c.Color,c.Code[c.Pc])
            }
            c.tick(&new_b,b)
            // b.Cores[i].Pc = new_b.Cores[i].Pc
        }
    }
    if verbose {
        fmt.Println("###")
    }
    // time.Sleep(250*time.Millisecond)
    deepcopy_board(&new_b,b) //bottleneck 2
    cycles += 1
}

func code_to_codemap(source string, loc int) map[int][]cmd {
    code_map := make(map[int][]cmd)
    empty_core := make([]cmd, loc)
    for i:=0; i < loc; i++ {
        empty_core[i] = cmd{7,[]arg{arg{"",1},arg{"",0}}}
    }
    code_map[-1] = empty_core
    for j, nod_cod := range strings.Split(strings.TrimSpace(source), "*")[1:] {
        temp_code := make([]cmd, loc)
        for i:=0; i < loc; i++ {
            temp_code[i] = cmd{7,[]arg{arg{"",1},arg{"",0}}}
        }
        for i, val := range strings.Split(strings.TrimSpace(nod_cod[1:]),"\n") {
            if len(val) <= 1 {continue}
            if i>loc {
                panic("too many lines of code per node in source file")
            }
            cmd_string := strings.Split(val, " ")
            numeric := regexp.MustCompile(`[0-9]`)
            temp_cmd := cmd{cmd_ids[cmd_string[0]],[]arg{arg{"",0},arg{"",0}}}
            for k:=1; k<=2; k++ {
                arg_num_pos := numeric.FindIndex([]byte(cmd_string[k]))[0]
                arg_num, err := strconv.Atoi(cmd_string[k][arg_num_pos:])
                arg_prefs := []byte(cmd_string[k][:arg_num_pos])
                for i, j := 0, len(arg_prefs)-1; i < j; i, j = i+1, j-1 { //reverse arg_prefs
                    arg_prefs[i], arg_prefs[j] = arg_prefs[j], arg_prefs[i]
                }
                if err!=nil {panic("arg non-numeric")}
                temp_cmd.Args[k-1] = arg{string(arg_prefs),int(arg_num)}
            }
            temp_code[i] = temp_cmd
        }
        code_map[j] = temp_code
    }
    return code_map
}

func code_to_layout(lt string, code_map map[int][]cmd, verbose bool) (cores []core) {
    lt_split := strings.Split(strings.TrimSpace(lt), "\n")
    xwrap := false
    ywrap := false
    width := 1
    height := 1
    for _, property := range lt_split {
        option := strings.TrimSpace(strings.Split(property, "=")[0])
        value := strings.TrimSpace(strings.Split(property, "=")[1])
        switch option {
        case "wrap":
            if value == "x" {
                xwrap = true
            } else if value == "y" {
                ywrap = true
            } else if value == "xy" {
                xwrap = true
                ywrap = true
            }
        case "width":
            width, _ = strconv.Atoi(value)
        case "height":
            height, _ = strconv.Atoi(value)
        }
    }
    for y:=0; y<height; y++ {
        for x:=0; x<width; x++ {
            left := y*width+mod(x-1,width)
            right := y*width+mod(x+1,width)
            up := mod(y-1,height)*width+x
            down := mod(y+1,height)*width+x
            if !xwrap && x == 0 {
                left = -1
            } else if !xwrap && x == width {
                right = -1
            }
            if !ywrap && y == 0 {
                up = -1
            } else if !ywrap && y == height {
                down = -1
            }
            if len(code_map[y*width+x])>0 {
                cores = append(cores, core{code_map[y*width+x], 0, y*width+x, up, right, down, left, "#61d18d", true})
            } else {
                cores = append(cores, core{code_map[-1], 0, y*width+x, up, right, down, left, INACTIVE_COLOR, false})
            }
        }
    }
    if verbose {
        fmt.Println(width, height, INACTIVE_COLOR)
    }
    // for i,core := range cores {
    //     fmt.Println(i,core)
    // }
    return
}

func build_board(source string, verbose bool) (new_b board) { //builds a board from a string
    loc := 1
    source_split := strings.Split(strings.TrimSpace(source), "CODE:")
    properties := strings.Split(source_split[0],"LAYOUT:")[0]
    for _, property := range strings.Split(strings.TrimSpace(properties), "\n")[1:] {
        option := strings.TrimSpace(strings.Split(property, "=")[0])
        value := strings.TrimSpace(strings.Split(property, "=")[1])
        switch option {
        case "name":
            fmt.Println(value,"-")
        case "description":
            fmt.Println(value)
        case "loc":
            loc, _ = strconv.Atoi(value)
        }
    }
    fmt.Println("\nRENDER_INFO_BEGIN")
    code_map := code_to_codemap(source_split[1], loc)
    new_b.Cores = code_to_layout(strings.Split(source_split[0],"LAYOUT:")[1], code_map, verbose)
    // for _,val := range new_b.Cores { //DEBUGGING
    //     fmt.Println(val)
    // }
    return new_b
}

//FRONTEND
//FRONTEND
//FRONTEND
//FRONTEND
//FRONTEND
//FRONTEND

type emu struct {
	Pause bool
	Ended bool
    Speed int

    backboard board
}


func (e *emu) run(verbose bool) {
    e.backboard.run(verbose)
}

func main() {
    if len(os.Args)<=1 {
        panic("No input file given.")
    } else if len(os.Args) > 2 {
        if os.Args[2]=="-v" {

        }
    }
    source, err := ioutil.ReadFile(os.Args[1])
    if err != nil {
        panic(err)
    }
    // a := build_board(string(source), text_mode)
    // b := board{}
    // deepcopy_board(&a,&b)
    // fmt.Println(a,"\n",b,"\n")
    // b.Cores[0].Code[0].Id = 9
    // fmt.Println(a,"\n",b,"\n")
    // deepcopy_board(&b,&a)
    // fmt.Println(a,"\n",b,"\n")
    emulator := emu{false, false, 1, build_board(string(source), text_mode)}
    start := time.Now()
    for x:=0; x<10; x++ {
        emulator.run(text_mode)
    }
    elapsed := time.Since(start)
    fmt.Println("RENDER_INFO_END")
    fmt.Println(elapsed)

}
