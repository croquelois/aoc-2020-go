package main

import (
    "fmt"
    "io/ioutil"
    "strings"
)

type Group struct {
  questions map[rune]int
  members int
}

func NewGroup() Group {
  var g Group
  g.questions = make(map[rune]int)
  g.members = 0
  return g
}

func (g *Group) addQuestion(r rune) {
  val, ok := g.questions[r]
  if !ok {
    g.questions[r] = 1
  } else {
    g.questions[r] = val + 1
  }
}

func (g *Group) addMember() {
  g.members += 1
}

func (g *Group) nbQuestions() int {
  return len(g.questions)
}

func (g *Group) nbQuestionsFullYes() int {
  count := 0
  for _,v := range g.questions {
    if v == g.members {
      count += 1
    }
  }
  return count
}


func parse(data string) []Group {
  var arr []Group
  g := NewGroup()
  lines := strings.Split(data, "\n")
  for _, line := range lines {
    if len(line) == 0 {
      arr = append(arr, g)
      g = NewGroup()
      continue
    }
    g.addMember()
    for _, r := range line {
      g.addQuestion(r)
    }
  }
  arr = append(arr, g)
  return arr
}

func parseFile(filename string) []Group {
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    panic(err)
  }
  return parse(string(data))
}

func algo1(groups []Group) int {
  count := 0
  for _, g := range groups {
    count = count + g.nbQuestions()
  }
  return count
}

func algo2(groups []Group) int {
  count := 0
  for _, g := range groups {
    count = count + g.nbQuestionsFullYes()
  }
  return count
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
  assert_eq(algo1(parseFile("test1.txt")), 11, "1");
}

func test2_1() {
  assert_eq(algo2(parseFile("test1.txt")), 6, "2");
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
