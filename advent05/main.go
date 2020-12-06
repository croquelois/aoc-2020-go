package main

import (
    "fmt"
    "io/ioutil"
    "strings"
)

type BoardingPass struct {
  position string
}

func (b BoardingPass) getRow() int {
  row := 0
  m := 128
  for _, rune := range b.position[:7] {
    m /= 2
    if rune == 'B' {
      row += m
    }
  }
  return row
}
func (b BoardingPass) getColumn() int {
  col := 0
  m := 8
  for _, rune := range b.position[7:10] {
    m /= 2
    if rune == 'R' {
      col += m
    }
  }
  return col
}
func (b BoardingPass) getSeatId() int {
  return b.getRow() * 8 + b.getColumn()
}

func parseBoardingPass(data string) BoardingPass {
  return BoardingPass{data}
}

type Plane struct {
  seats [128*8]bool
}
func newPlane() Plane {
  var p Plane
  for i, _ := range p.seats {
    p.seats[i] = false
  }
  return p
}
func (p *Plane) addPassenger(b BoardingPass) {
  p.seats[b.getSeatId()] = true
}
func (p *Plane) addPassengers(bs []BoardingPass) {
  for _, b := range bs {
    p.addPassenger(b)
  }
}
func (p *Plane) firstFreeSeat() int {
  i := 1
  for !p.seats[i] { i++ }
  for p.seats[i] { i++ }
  return i
}

func parse(data string) []BoardingPass {
  var arr []BoardingPass
  lines := strings.Split(data, "\n")
  for _, line := range lines {
    arr = append(arr, parseBoardingPass(line))
  }
  return arr
}

func parseFile(filename string) []BoardingPass {
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    panic(err)
  }
  return parse(string(data))
}

func algo1(boardingPasses []BoardingPass) int {
  biggestSeatID := -1
  for _, b := range boardingPasses {
    seatId := b.getSeatId()
    if seatId > biggestSeatID {
      biggestSeatID = seatId
    }
  }
  return biggestSeatID
}

func algo2(boardingPasses []BoardingPass) int {
  p := newPlane()
  p.addPassengers(boardingPasses)
  return p.firstFreeSeat()
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
  assert_eq(parseBoardingPass("BFFFBBFRRR").getRow(), 70, "BFFFBBF");
  assert_eq(parseBoardingPass("BFFFBBFRRR").getColumn(), 7, "RRR");
  assert_eq(parseBoardingPass("BFFFBBFRRR").getSeatId(), 567, "BFFFBBFRRR");
  assert_eq(parseBoardingPass("FFFBBBFRRR").getSeatId(), 119, "FFFBBBFRRR");
  assert_eq(parseBoardingPass("BBFFBBFRLL").getSeatId(), 820, "BBFFBBFRLL");
}

func question1() int {
  return algo1(parseFile("input.txt"));
}

func question2() int {
  return algo2(parseFile("input.txt"));
}

func main() {
  test1_1()
  fmt.Printf("Question1: %d\n", question1())
  fmt.Printf("Question2: %d\n", question2())
}
