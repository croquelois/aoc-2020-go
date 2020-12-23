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

func isNumberOk(data []int, nb int) bool{
  for idx, n1 := range data {
    for _, n2 := range data[idx+1:] {
      if n1+n2 ==  nb {
        return true
      }
    }
  }
  return false
}

func algo1(data []int, depth int) int {
  for idx, n1 := range data[depth:] {
    if !isNumberOk(data[idx:idx+depth],n1) {
      return n1
    }
  }
  panic("One of the number should not match");
}

func Min(v []int) int {
  m := v[0]
  for i := 1; i < len(v); i++ {
      if v[i] < m {
          m = v[i]
      }
  }
  return m
}

func Max(v []int) int {
  m := v[0]
  for i := 1; i < len(v); i++ {
      if v[i] > m {
          m = v[i]
      }
  }
  return m
}

func algo2(data []int, depth int) int {
  nb := algo1(data, depth)
  val := 0
  low := 0
  high := 0
  for {
    if val < nb {
      val += data[high]
      high += 1
    } else if val > nb {
      val -= data[low]
      low += 1
    } else {
      return Min(data[low:high]) + Max(data[low:high])
    }
  }
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
  assert_eq(algo1(parseFile("test1.txt"),5), 127, "1");
}

func question1() int {
  return algo1(parseFile("input.txt"),25);
}

func test2_1() {
  assert_eq(algo2(parseFile("test1.txt"),5), 62, "2");
}

func question2() int {
  return algo2(parseFile("input.txt"),25);
}

func main() {
  test1_1()
  fmt.Printf("Question1: %d\n", question1())
  test2_1()
  fmt.Printf("Question2: %d\n", question2())
}
