package main

import (
    "fmt"
    "io/ioutil"
    "strings"
    "bufio"
    "os"
    "strconv"
)

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

type Range struct {
  min int
  max int
}

func parseRange(s string) Range {
  var r Range
  tmp := splitTrimInt(s, "-")
  r.min = tmp[0]
  r.max = tmp[1]
  return r
}

func (r Range) isInside(n int) bool {
  return n >= r.min && n <= r.max 
}

type Rule struct {
  name string
  ranges []Range
  pos int
  validPos []bool
}

func parseRule(s string) Rule {
  var r Rule
  tmp := splitTrim(s, ":")
  r.name = tmp[0]
  rangesStr := splitTrim(tmp[1],"or")
  for _, s := range rangesStr {
    r.ranges = append(r.ranges, parseRange(s))
  }
  r.pos = -1
  return r
}

func (rule Rule) isValid(n int) bool {
  for _, r := range rule.ranges {
    if r.isInside(n) {
      return true
    }
  }
  return false
}

type Ticket struct {
  numbers []int
}

func parseTicket(s string) Ticket {
  var t Ticket
  t.numbers = splitTrimInt(s, ",")
  return t
}

func (t Ticket) isValid(rules []Rule) bool {
  for i, n := range t.numbers {
    isValid := false
    posFound := false
    for _, rule := range rules {
      if rule.pos == i {
        posFound = true
        if rule.isValid(n) {
          isValid = true
        } else {
          return false
        }
        break
      }
    }
    if !posFound {
      for _, rule := range rules {
        if rule.pos == -1 && rule.isValid(n) {
          isValid = true
          break
        }
      }
    }
    if !isValid { 
      return false
    }
  }
  return true
}

func (t Ticket) errorRate(rules []Rule) int {
  err := 0
  for i, n := range t.numbers {
    isValid := false
    posFound := false
    for _, rule := range rules {
      if rule.pos == i {
        posFound = true
        if rule.isValid(n) {
          isValid = true
        } else {
          isValid = false
        }
        break
      }
    }
    if !posFound {
      for _, rule := range rules {
        if rule.pos == -1 && rule.isValid(n) {
          isValid = true
          break
        }
      }
    }
    if !isValid { err += n }
  }
  return err
}

type Problem struct {
  rules []Rule
  myTicket Ticket
  nearbyTickets []Ticket
}

func parseProblem(s string) Problem {
  var p Problem
  lines := splitTrim(s, "\n")
  mode := "rules:"
  for _, line := range lines {
    if len(line) == 0 { 
      mode = ""
      continue
    }
    switch mode {
      case "":
        mode = line
      case "rules:":
        p.rules = append(p.rules, parseRule(line))
      case "your ticket:":
        p.myTicket = parseTicket(line)
      case "nearby tickets:":
        p.nearbyTickets = append(p.nearbyTickets, parseTicket(line))
      default:
        panic("failed parsing, unexpected mode")
    }
  }
  return p
}

func (p Problem) errorRate() int {
  err := 0
  for _, t := range p.nearbyTickets {
    err += t.errorRate(p.rules)
  }
  return err
}

func (p *Problem) solve() {
  var validTickets []Ticket
  //fmt.Println("Keep only valid tickets")
  for _, t := range p.nearbyTickets {
    if t.isValid(p.rules) {
      validTickets = append(validTickets, t)
    }
  }
  //fmt.Println("nb ticket initial", len(p.nearbyTickets), "nb ticket valid", len(validTickets))
  
  
  //fmt.Println("initialise rules valid position range")
  for rIdx := range p.rules {
    r := &p.rules[rIdx]
    for _, n := range p.myTicket.numbers {
      r.validPos = append(r.validPos, r.isValid(n))
    }
  }
  for rIdx := range p.rules {
    r := &p.rules[rIdx]
    for _, t := range validTickets {
      for i, n := range t.numbers {
        if !r.isValid(n) { r.validPos[i] = false }
      }
    }
  }
  
  for {
    //fmt.Println("Resolution loop")
    hasChanged := false
    for rIdx := range p.rules {
      r := &p.rules[rIdx]
      if r.pos > -1 { continue }
      //fmt.Println("Rule", r)
      count := 0
      pos := 0
      for i, v := range r.validPos {
        if v {
          count += 1 
          pos = i
        }
      }
      if count == 0 { panic("can't solve the problem") }
      if count > 1 { continue }
      hasChanged = true
      r.pos = pos
            
      for rIdx2 := range p.rules {
        r2 := &p.rules[rIdx2]
        r2.validPos[pos] = false
      }
    }
    if !hasChanged {
      for _, r := range p.rules {
        if r.pos == -1 { panic("can't solve the problem") }
      }
      return
    }
  }  
}

func (p Problem) get(name string) int {
  for _, r := range p.rules {
    if r.name == name {
      return p.myTicket.numbers[r.pos]
    }
  }
  panic("rule not found")
}

func parseFile(filename string) Problem {
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    panic(err)
  }
  return parseProblem(string(data))
}

func pause(prompt string){
  fmt.Println(prompt)
  input := bufio.NewScanner(os.Stdin)
  input.Scan()
}

func algo1(p Problem) int {
  return p.errorRate()
}

func algo2(p Problem) int {
  p.solve()
  mul := 1
  for _, r := range p.rules {
    if strings.HasPrefix(r.name, "departure") {
      mul *= p.get(r.name)
    }
  }
  return mul
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
  assert_eq(algo1(parseFile("test1.txt")), 71, "1");
}

func question1() int {
  return algo1(parseFile("input.txt"));
}

func test2_1() {
  p := parseFile("test2.txt")
  p.solve()
  assert_eq(p.get("class"), 12, "class");
  assert_eq(p.get("row"), 11, "row");
  assert_eq(p.get("seat"), 13, "seat");
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