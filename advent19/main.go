package main

import (
    "fmt"
    "io/ioutil"
    "strings"
    "bufio"
    "os"
    "strconv"
)

type Rule struct {
  key int
  rules [][]int
  char string
}

func atoi(str string) int {
  value, err := strconv.Atoi(str)
  if err != nil { panic("Can't parse the number") }
  return value
}

func splitTrim(data string, sep string) []string {
  tmp := strings.Split(data, sep)
  for i := range tmp {
    tmp[i] = strings.TrimSpace(tmp[i]) 
  }
  return tmp;
}

func splitTrimInt(data string, sep string) []int {
  var arr = []int{}
  tmp := splitTrim(data, sep)
  for i := range tmp {
    arr = append(arr, atoi(tmp[i]))
  }
  return arr;
}

func parseRule(data string) Rule {
  var r Rule
  tmp := splitTrim(data,":")
  r.key = atoi(tmp[0])
  tmp2 := splitTrim(tmp[1],"\x22")
  if len(tmp2) == 3 {
    r.char = tmp2[1]
    return r
  }
  tmp2 = splitTrim(tmp[1],"|")
  for _,v := range tmp2 {
    r.rules = append(r.rules, splitTrimInt(v," "))
  }
  return r
}

type Rules map[int]Rule

func parseRules(lines []string) Rules {
  ret := make(Rules)
  for _, line := range lines {
    rule := parseRule(line)
    ret[rule.key] = rule
  }
  return ret
}

type Problem struct {
  rules Rules
  messages []string
}

func parse(data string) Problem {
  var p Problem
  parts := splitTrim(data, "\n\n")
  p.rules = parseRules(splitTrim(parts[0], "\n"))
  p.messages = splitTrim(parts[1], "\n")
  return p
}

func (problem Problem) matchRuleFwd(s string, rules []int, p int) []int {
  //fmt.Println("IN matchRuleFwd", s, rules, p)
  possibles := problem.matchRule(s, rules[0], p)
  if len(rules) == 1 { 
    //fmt.Println("OUT matchRuleFwd", s, rules, p, possibles)
    return possibles 
  }
  dedup := make(map[int]bool)
  for _, p2 := range possibles {
    possiblesIn := problem.matchRuleFwd(s, rules[1:], p2)
    for _, p3 := range possiblesIn {
      dedup[p3] = true
    }
  }
  keys := make([]int, len(dedup))
  i := 0
  for k := range dedup {
      keys[i] = k
      i++
  }
  //fmt.Println("OUT matchRuleFwd", s, rules, p, keys)
  return keys
}

func (problem Problem) matchRule(s string, n int, p int) []int {
  //fmt.Println("IN matchRule", s, n, p)
  if p >= len(s) {
    //fmt.Println("OUT matchRule", s, n, p, []int{})
    return []int{}
  }
  r := problem.rules[n]
  if len(r.rules) == 0 {
    if s[p:p+1] == r.char {
      //fmt.Println("OUT matchRule", s, n, p, []int{p+1})
      return []int{p+1}
    }
    //fmt.Println("OUT matchRule", s, n, p, []int{})
    return []int{}
  }
  dedup := make(map[int]bool)
  for _, m := range r.rules {
    possiblesIn := problem.matchRuleFwd(s, m, p)
    for _, p3 := range possiblesIn {
      dedup[p3] = true
    }
  }
  keys := make([]int, len(dedup))
  i := 0
  for k := range dedup {
      keys[i] = k
      i++
  }
  //fmt.Println("OUT matchRule", s, n, p, keys)
  return keys
}

func (p Problem) isCorrectMsg(msg string, c chan bool) {
  for _, p := range p.matchRule(msg, 0, 0) {
    if p == len(msg) { 
      c <- true
      return
    }
  }
  c <- false
}

func (p Problem) nbCorrectMsg() int {
  count := 0
  c := make(chan bool)
  for _,m := range p.messages { go p.isCorrectMsg(m, c) }
  for _,_ = range p.messages { if <-c { count += 1 } }
  return count
}

func parseFile(filename string) Problem {
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

func algo1(p Problem) int {
  return p.nbCorrectMsg()
}

func algo2(p Problem) int {
  r := p.rules[8]
  r.rules = append(r.rules, []int{42,8})
  p.rules[8] = r
  
  r = p.rules[11]
  r.rules = append(r.rules, []int{42,11,31})
  p.rules[11] = r
  return p.nbCorrectMsg()
}

func assert(v bool, msg string){
  if !v {
    fmt.Printf("the test fail: %s\n", msg)
    panic("test failed !")
  }
}

func assert_eq(v int, e int, msg string){
  if v != e {
    fmt.Printf("the test '%s' fail, it give %d instead of %d\n", msg, v, e)
    panic("test failed !")
  }
}

func test1_1() {
  assert_eq(algo1(parseFile("test1.txt")), 2, "1");
}

func test1_2() {
  assert_eq(algo1(parseFile("test2.txt")), 3, "2");
}

func test2_1() {
  assert_eq(algo2(parseFile("test2.txt")), 12, "2");
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
  test1_2()
  test2_1()
  fmt.Printf("Question2: %d\n", question2())
}