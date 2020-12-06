package main

import (
    "fmt"
    "io/ioutil"
    "strings"
    "strconv"
)

type Rule struct {
  min int
  max int
  chr rune
}

type RulePassword struct {
  rule Rule
  password string
}

func firstRune(str string) rune {
  for _, r := range str {
      return r
  }
  panic("empty string !")
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
    value, err := strconv.Atoi(tmp[i])
    if err != nil {
        panic(err)
    }
    arr = append(arr, value)
  }
  return arr;
}

func parseRule(data string) Rule {
  tmp := splitTrim(data, " ")
  rng := splitTrimInt(tmp[0], "-")
  var r Rule
  r.min = rng[0]
  r.max = rng[1]
  r.chr = firstRune(tmp[1])
  return r
}

func parseLine(data string) RulePassword{
  var rulePassword RulePassword
  tmp := splitTrim(data, ":")
  rulePassword.rule = parseRule(tmp[0])
  rulePassword.password = tmp[1]
  return rulePassword
}

func parse(data string) []RulePassword {
  var arrStr = strings.Split(data, "\n")
  var arr = []RulePassword{}
  for _, s := range arrStr {
      arr = append(arr, parseLine(s))
  }
  return arr
}

func parseFile(filename string) []RulePassword {
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    panic(err)
  }
  return parse(string(data))
}

type isValidAlgo func(rulePassword RulePassword) bool

func isValidAlgo1(rulePassword RulePassword) bool {
  var count = 0
  for _, s := range rulePassword.password {
    if s == rulePassword.rule.chr {
      count++
    }
  }
  return rulePassword.rule.min <= count && count <= rulePassword.rule.max;
}

func isValidAlgo2(rulePassword RulePassword) bool {
  var count = 0
  for i, s := range rulePassword.password {
    if (i == (rulePassword.rule.min-1) || i == (rulePassword.rule.max-1)) && s == rulePassword.rule.chr {
      count++
    }
  }
  return count == 1;
}

func nbValid(rulePasswords []RulePassword, isValid isValidAlgo) int {
  var count = 0
  for _, rulePassword := range rulePasswords {
    if isValid(rulePassword) {
      count++;
    }
  }
  return count
}

func test1_1() {
  var expected = 2
  var v = nbValid(parseFile("test1.txt"), isValidAlgo1);
  if v != expected {
    fmt.Printf("the test give %d instead of %d\n", v, expected)
    panic("test failed !")
  }
}

func test2_1() {
  var expected = 1
  var v = nbValid(parseFile("test1.txt"), isValidAlgo2);
  if v != expected {
    fmt.Printf("the test give %d instead of %d\n", v, expected)
    panic("test failed !")
  }
}

func question1() int {
  return nbValid(parseFile("input.txt"), isValidAlgo1);
}

func question2() int {
  return nbValid(parseFile("input.txt"), isValidAlgo2);
}

func main() {
  test1_1()
  fmt.Printf("Question1: %d\n", question1())
  test2_1()
  fmt.Printf("Question2: %d\n", question2())
}
