package main

import (
    "fmt"
    "io/ioutil"
    "strings"
    "regexp"
    "strconv"
)

func splitTrim(data string, sep string) []string {
  tmp := strings.Split(data, sep)
  for i := range tmp {
    tmp[i] = strings.TrimSpace(tmp[i]) 
  }
  return tmp;
}

type Rule struct {
  color string
  contains map[string]int
}

var reLine = regexp.MustCompile(`^(.+) bags contain (.+)\.$`)
var reContain = regexp.MustCompile(`^(\d+) (.+) bags?$`)
func ParseRule(line string) Rule {
  match := reLine.FindStringSubmatch(line)
  if match == nil { 
    panic("Can't parse the rule") 
  }
  var r Rule
  r.color = strings.TrimSpace(match[1])
  r.contains = make(map[string]int)
  if match[2] == "no other bags" { 
    return r 
  }
  for _,v := range splitTrim(match[2],",") {
    matchContain := reContain.FindStringSubmatch(v)
    if matchContain == nil { 
      panic("Can't parse the right side of the rule") 
    }
    nb, err := strconv.Atoi(matchContain[1])
    if err != nil {
      panic("Can't parse the number of bags") 
    }
    r.contains[matchContain[2]] = nb
  }
  return r
}

func parse(data string) map[string]Rule {
  ret := make(map[string]Rule)
  for _, line := range strings.Split(data, "\n") {
    r := ParseRule(line)
    ret[r.color] = r
  }
  return ret
}

func parseFile(filename string) map[string]Rule {
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    panic(err)
  }
  return parse(string(data))
}

func algo1(rules map[string]Rule, color string) int {
  graph := make(map[string][]string)
  for _, r := range rules {
    for k, _ := range r.contains {
      graph[k] = append(graph[k],r.color)
    }
  }
  
  visited := make(map[string]bool)
  visited[color] = true
  var open []string
  open = append(open, color)
  pos := 0
  
  for ;pos<len(open); {
    cur := open[pos]
    pos++
    for _, c := range graph[cur] {
      _, present := visited[c]
      if !present {
        open = append(open, c)
        visited[c] = true
      }
    }
  }
  
  return len(visited)-1
}

func algo2(rules *map[string]Rule, color string) int {
  total := 0
  for k, v := range ((*rules)[color].contains) {
    total += v*(1+algo2(rules, k))
  }
  return total
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
  assert_eq(algo1(parseFile("test1.txt"), "shiny gold"), 4, "1");
}

func test2_1() {
  rules := parseFile("test1.txt")
  assert_eq(algo2(&rules, "dark olive"), 7, "dark olive");
  assert_eq(algo2(&rules, "vibrant plum"), 11, "vibrant plum");
  assert_eq(algo2(&rules, "shiny gold"), 32, "shiny gold");
}

func test2_2() {
  rules := parseFile("test2.txt")
  assert_eq(algo2(&rules, "shiny gold"), 126, "test 2 - second file");
}

func question1() int {
  return algo1(parseFile("input.txt"), "shiny gold");
}

func question2() int {
  rules := parseFile("input.txt")
  return algo2(&rules, "shiny gold");
}

func main() {
  test1_1()
  fmt.Printf("Question1: %d\n", question1())
  test2_1()
  test2_2()
  fmt.Printf("Question2: %d\n", question2())
}
