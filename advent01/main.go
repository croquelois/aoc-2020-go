package main

import (
    "fmt"
    "io/ioutil"
    "strings"
    "strconv"
)

func parse(data string) []int {
  var arrStr = strings.Split(data, "\n")
  var arrInt = []int{}
  
  for _, s := range arrStr {
      i, err := strconv.Atoi(s)
      if err != nil {
          panic(err)
      }
      arrInt = append(arrInt, i)
  }
  return arrInt
}

func parseFile(filename string) []int {
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    panic(err)
  }
  return parse(string(data))
}

func algo1(data []int) int {
  for idx, n1 := range data {
    for _, n2 := range data[idx:] {
      if n1 + n2 == 2020 {
        return n1*n2;
      }
    }
  }
  panic("can't found a correct pair")
}

func algo2(data []int) int {
  for idx1, n1 := range data {
    for idx2, n2 := range data[idx1:] {
      for _, n3 := range data[idx2:] {
        if n1 + n2 + n3 == 2020 {
          return n1*n2*n3;
        }
      }
    }
  }
  panic("can't found a correct tuple")
}

func test1_1() {
  var expected = 514579
  var v = algo1(parseFile("test1.txt"));
  if v != expected {
    fmt.Printf("the test give %d instead of %d\n", v, expected)
    panic("test failed !")
  }
}

func test2_1() {
  var expected = 241861950
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
