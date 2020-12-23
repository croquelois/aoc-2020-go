package main

import (
    "fmt"
    "io/ioutil"
    "strings"
    "bufio"
    "os"
)

type Map struct {
  w,h int
  seat []bool
  occupied []bool
  willChange []bool
}

func (m *Map) isInside(x int, y int) bool {
  if x < 0 || x >=  m.w { return false }
  if y < 0 || y >=  m.h { return false }
  return true
}

func (m *Map) hasSeat(x int, y int) bool {
  return m.isInside(x,y) && m.seat[(x % m.w) + (y * m.w)]
}

func (m *Map) applyChange() bool {
  hasChanged := false
  for idx, v := range m.willChange {
    if v { 
      m.occupied[idx] = !m.occupied[idx]
      m.willChange[idx] = false
      hasChanged = true 
    }
  }
  return hasChanged
}

func (m *Map) flipSeat(x int, y int) {
  m.willChange[(x % m.w) + (y * m.w)] = true
}

func (m *Map) isOccupied(x int, y int) bool {
  return m.isInside(x,y) && m.occupied[(x % m.w) + (y * m.w)]
}

func (m *Map) seeOccupied(x int, y int, dx int, dy int) bool {
  x += dx
  y += dy
  for m.isInside(x, y) {
    if m.hasSeat(x, y) {
      return m.isOccupied(x, y)
    }
    x += dx
    y += dy
  }
  return false
}

func (m *Map) peopleSeenArround(x int, y int) int {
  count := 0
  if m.seeOccupied(x,y,-1,-1) { count += 1 }
  if m.seeOccupied(x,y,-1, 0) { count += 1 }
  if m.seeOccupied(x,y,-1,+1) { count += 1 }
  if m.seeOccupied(x,y, 0,-1) { count += 1 }
  if m.seeOccupied(x,y, 0,+1) { count += 1 }
  if m.seeOccupied(x,y,+1,-1) { count += 1 }
  if m.seeOccupied(x,y,+1, 0) { count += 1 }
  if m.seeOccupied(x,y,+1,+1) { count += 1 }
  return count
}

func (m *Map) peopleArround(x int, y int) int {
  count := 0
  if m.isOccupied(x-1,y-1) { count += 1 }
  if m.isOccupied(x-1,y  ) { count += 1 }
  if m.isOccupied(x-1,y+1) { count += 1 }
  if m.isOccupied(x  ,y-1) { count += 1 }
  if m.isOccupied(x  ,y+1) { count += 1 }
  if m.isOccupied(x+1,y-1) { count += 1 }
  if m.isOccupied(x+1,y  ) { count += 1 }
  if m.isOccupied(x+1,y+1) { count += 1 }
  return count
}

func (m *Map) doOneTurn() bool {
  for x:=0;x<m.w;x++ {
    for y:=0;y<m.h;y++ {
      if !m.hasSeat(x, y) { continue }
      if m.isOccupied(x,y) {
        if m.peopleArround(x,y) >= 4 {
          m.flipSeat(x, y)
        }
      }else{
        if m.peopleArround(x,y) == 0 {
          m.flipSeat(x, y)
        }
      }
    }
  }
  return m.applyChange()
}

func (m *Map) doOneTurnAlgo2() bool {
  for x:=0;x<m.w;x++ {
    for y:=0;y<m.h;y++ {
      if !m.hasSeat(x, y) { continue }
      if m.isOccupied(x,y) {
        if m.peopleSeenArround(x,y) >= 5 {
          m.flipSeat(x, y)
        }
      }else{
        if m.peopleSeenArround(x,y) == 0 {
          m.flipSeat(x, y)
        }
      }
    }
  }
  return m.applyChange()
}

func (m *Map) countOccupied() int {
  count := 0
  for _, v := range m.occupied {
    if v { count += 1 }
  }
  return count
}

func (m *Map) toString() string {
  str := ""
  for y:=0;y<m.h;y++ {
    for x:=0;x<m.w;x++ {
      if !m.hasSeat(x, y) { 
        str += "."
        continue
      }
      if m.isOccupied(x,y) {
        str += "#"
      }else{
        str += "L"
      }
    }
    str += "\n"
  }
  return str
}

func parse(data string) Map {
  var m Map
  var lines = strings.Split(data, "\n")
  m.h = len(lines)
  m.w = len(lines[0])
  for _, line := range lines {
    for _, rune := range line {
      m.seat = append(m.seat, rune == 'L')
      m.occupied = append(m.occupied, false)
      m.willChange = append(m.willChange, false)
    }
  }
  return m
}

func parseFile(filename string) Map {
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

func algo1(m Map) int {
  for m.doOneTurn() {}
  return m.countOccupied()
}

func algo2(m Map) int {
  for m.doOneTurnAlgo2() {}
  return m.countOccupied()
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
  assert_eq(algo1(parseFile("test1.txt")), 37, "1");
}

func question1() int {
  return algo1(parseFile("input.txt"));
}

func test2_1() {
  assert_eq(algo2(parseFile("test1.txt")), 26, "2");
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
