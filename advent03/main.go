package main

import (
    "fmt"
    "io/ioutil"
    "strings"
)

type Map struct {
  w,h int
  tree []bool
}

func (m Map) hasTree(x int, y int) bool {
  return m.tree[(x % m.w) + (y * m.w)]
}

func parse(data string) Map {
  var m Map
  var lines = strings.Split(data, "\n")
  m.h = len(lines)
  m.w = len(lines[0])
  for _, line := range lines {
    for _, rune := range line {
      m.tree = append(m.tree, rune == '#')
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

func nbArborealStop(m Map, dx int, dy int) int {
  x := 0
  count := 0
  for y:=0;y<m.h;y+=dy {
    if m.hasTree(x, y) {
      count++
    }
    x += dx
  }
  return count
}

func algo1(m Map) int {
  return nbArborealStop(m, 3, 1)
}

func algo2(m Map) int {
  total := 1
  total *= nbArborealStop(m, 1, 1)
  total *= nbArborealStop(m, 3, 1)
  total *= nbArborealStop(m, 5, 1)
  total *= nbArborealStop(m, 7, 1)
  total *= nbArborealStop(m, 1, 2)
  return total
}

func test1_1() {
  var expected = 7
  var v = algo1(parseFile("test1.txt"));
  if v != expected {
    fmt.Printf("the test give %d instead of %d\n", v, expected)
    panic("test failed !")
  }
}

func test2_1() {
  var expected = 336
  var v = algo2(parseFile("test1.txt"));
  if v != expected {
    fmt.Printf("the test give %d instead of %d\n", v, expected)
    panic("test failed !")
  }
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
  test2_1()
  fmt.Printf("Question2: %d\n", question2())
}
