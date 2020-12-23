package main

import (
    "fmt"
    "io/ioutil"
    "strings"
    "bufio"
    "os"
    "strconv"
)

type Move struct {
  action string
  nb int
}

type Position struct {
  x int
  y int
  rotation int
  wx int
  wy int
}

func parseMove(data string) Move {
  var m Move
  m.action = data[:1]
  nb, err := strconv.Atoi(data[1:])
  if err != nil { panic("Can't parse the number of move") }
  m.nb = nb
  return m
}


func newPosition() Position {
  var p Position
  p.x = 0
  p.y = 0
  p.rotation = 0 // rotation (degree), 0 => facing east
  p.wx = 10
  p.wy = -1
  return p
}

func (p *Position) distance() int {
  dist := 0
  if p.x < 0 { 
    dist -= p.x 
  } else { 
    dist += p.x 
  }
  if p.y < 0 { 
    dist -= p.y 
  } else { 
    dist += p.y 
  }
  return dist
}

func (p *Position) applyMove(m Move) {
  switch m.action {
    case "N": p.y -= m.nb
    case "S": p.y += m.nb
    case "E": p.x += m.nb
    case "W": p.x -= m.nb
    case "L": p.rotation += m.nb
    case "R": p.rotation -= m.nb
    case "F":
      d := (p.rotation / 90) % 4
      if d < 0 { d += 4 }
      switch d {
        case 0: p.x += m.nb
        case 1: p.y -= m.nb
        case 2: p.x -= m.nb
        case 3: p.y += m.nb
        default: panic("unexpected rotation")
      }  
    default: panic("unexpected movement action")
  }
}

func (p *Position) applyMoveAlgo2(m Move) {
  rotation := 0
  switch m.action {
    case "N": p.wy -= m.nb
    case "S": p.wy += m.nb
    case "E": p.wx += m.nb
    case "W": p.wx -= m.nb
    case "L": rotation += m.nb/90
    case "R": rotation -= m.nb/90
    case "F":
      p.x += m.nb * p.wx
      p.y += m.nb * p.wy
    default: panic("unexpected movement action")
  }
  rotation = (rotation % 4)
  if rotation < 0 { rotation += 4 }
  switch rotation {
    case 1: 
      t := p.wx
      p.wx = p.wy
      p.wy = -t
    case 2:
      p.wx = -p.wx
      p.wy = -p.wy
    case 3:
      t := p.wx
      p.wx = -p.wy
      p.wy = t
  }
}

func parse(data string) []Move {
  var m []Move
  var lines = strings.Split(data, "\n")
  for _, line := range lines {
    m = append(m, parseMove(line))
  }
  return m
}

func parseFile(filename string) []Move {
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    panic(err)
  }
  return parse(string(data))
}

func pause(prompt string){
  fmt.Println(prompt)
  input := bufio.NewScanner(os.Stdin)
  input.Scan()
}

func algo1(moves []Move) int {
  p := newPosition()
  for _, move := range moves {
    p.applyMove(move)
  }
  return p.distance()
}

func algo2(moves []Move) int {
  p := newPosition()
  for _, move := range moves {
    p.applyMoveAlgo2(move)
  }
  return p.distance()
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
  assert_eq(algo1(parseFile("test1.txt")), 25, "1");
}

func question1() int {
  return algo1(parseFile("input.txt"));
}

func test2_1() {
  assert_eq(algo2(parseFile("test1.txt")), 286, "1");
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
