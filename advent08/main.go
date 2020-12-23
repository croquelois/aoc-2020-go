package main

import (
    "fmt"
    "io/ioutil"
    "strings"
    "strconv"
)

func splitTrim(data string, sep string) []string {
  tmp := strings.Split(data, sep)
  for i := range tmp {
    tmp[i] = strings.TrimSpace(tmp[i]) 
  }
  return tmp;
}

type Command struct {
  operation string
  value int
}

func ParseCommand(line string) Command {
  var c Command
  v := splitTrim(line, " ")
  c.operation = v[0]
  nb, err := strconv.Atoi(v[1])
  if err != nil {
    panic("Can't parse the number") 
  }
  c.value = nb
  return c
}

func parse(data string) []Command {
  var ret []Command
  for _, line := range strings.Split(data, "\n") {
    ret = append(ret, ParseCommand(line))
  }
  return ret
}

func parseFile(filename string) []Command {
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    panic(err)
  }
  return parse(string(data))
}

func runProgram(program []Command) (int, bool) {
  acc := 0
  pos := 0
  var visited []bool
  for _, _ = range program {
    visited = append(visited, false)
  }
  for pos >= 0 && pos < len(program) && !visited[pos] {
    cmd := program[pos]
    visited[pos] = true
    switch cmd.operation {
      case "acc":
        acc += cmd.value
        pos += 1        
      case "jmp":
        pos += cmd.value
      case "nop":
        pos += 1
      default:
        panic("Unknow operation")
    }
  }
  return acc, pos == len(program)
}

func algo1(program []Command) int {
  acc, _ := runProgram(program)
  return acc
}

func algo2(program []Command) int {
  for idx, v := range program {
    if v.operation == "acc" {
      continue
    }
    if v.operation == "jmp" {
      program[idx].operation = "nop"
      acc, ok := runProgram(program)
      if ok {
        return acc
      }
      program[idx].operation = "jmp"
    }
    if v.operation == "nop" {
      program[idx].operation = "jmp"
      acc, ok := runProgram(program)
      if ok {
        return acc
      }
      program[idx].operation = "nop"
    }
  }
  panic("Can't found a way to repair the program")
}

func assert(v bool, msg string){
  if !v {
    fmt.Printf("the test fail: %s\n", msg)
    panic("test failed !")
  }
}

func assert_eq(v int, e int, msg string){
  if v != e {
    fmt.Printf("the test %s fail, it give %d instead of %d\n", msg, v, e)
    panic("test failed !")
  }
}

func test1_1() {
  assert_eq(algo1(parseFile("test1.txt")), 5, "1");
}

func question1() int {
  return algo1(parseFile("input.txt"));
}


func test2_1() {
  assert_eq(algo2(parseFile("test1.txt")), 8, "1");
}

func question2() int {
  return algo2(parseFile("input.txt"));
}

func main() {
  test1_1()
  fmt.Printf("Question1: %d\n", question1())
  test2_1()
  fmt.Printf("Question2: %d\n", question2())
}
