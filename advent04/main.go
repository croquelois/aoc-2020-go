package main

import (
    "fmt"
    "io/ioutil"
    "strings"
    "regexp"
    "strconv"
)

type Passport struct {
  entries map[string]string
}

func NewPassport() Passport {
  var p Passport
  p.entries = make(map[string]string)
  return p
}

func (p Passport) addEntry(key string, value string) {
  p.entries[key] = value
}

func (p Passport) isValid() bool {
  required := [...]string{"byr","iyr","eyr","hgt","hcl","ecl","pid"} // ,"cid"
  for _, key := range required {
    _, ok := p.entries[key]
    if !ok {
      return false
    }
  }
  return true
}

func validateBirthYear(s string) bool {
  year, err := strconv.Atoi(s)
  if err != nil {
      return false
  }
  return (year >= 1920 && year <= 2002)
}

func validateIssueYear(s string) bool {
  year, err := strconv.Atoi(s)
  if err != nil {
      return false
  }
  return (year >= 2010 && year <= 2020)
}

func validateExpirationYear(s string) bool {
  year, err := strconv.Atoi(s)
  if err != nil {
      return false
  }
  return (year >= 2020 && year <= 2030)
}

var reHeight = regexp.MustCompile(`^([0-9]+)(in|cm)$`)
func validateHeight(s string) bool {
  match := reHeight.FindStringSubmatch(s)
  if match == nil { return false }
  height, _ := strconv.Atoi(match[1])
  if match[2] == "cm" {
    return (height >= 150 && height <= 193)
  }
  return (height >= 59 && height <= 76)
}

var reHairColor = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)
func validateHairColor(s string) bool { return reHairColor.MatchString(s) }

var reEyeColor = regexp.MustCompile(`^(amb|blu|brn|gry|grn|hzl|oth)$`)
func validateEyeColor(s string) bool { return reEyeColor.MatchString(s) }

var rePassportID = regexp.MustCompile(`^[0-9]{9}$`)
func validatePassportID(s string) bool { return rePassportID.MatchString(s) }

func (p Passport) isValidExtended() bool {
  if !p.isValid() { return false }
  if !validateBirthYear(p.entries["byr"]) { return false }
  if !validateIssueYear(p.entries["iyr"]) { return false }
  if !validateExpirationYear(p.entries["eyr"]) { return false }
  if !validateHeight(p.entries["hgt"]) { return false }
  if !validateHairColor(p.entries["hcl"]) { return false }
  if !validateEyeColor(p.entries["ecl"]) { return false }
  if !validatePassportID(p.entries["pid"]) { return false }
  return true
}

func parse(data string) []Passport {
  var arr []Passport
  p := NewPassport()
  lines := strings.Split(data, "\n")
  for _, line := range lines {
    if len(line) == 0 {
      arr = append(arr, p)
      p = NewPassport()
      continue
    }
    pairs := strings.Split(line, " ")
    for _, pair := range pairs {
      tmp := strings.Split(pair, ":")
      p.addEntry(tmp[0], tmp[1])
    }
  }
  arr = append(arr, p)
  return arr
}

func parseFile(filename string) []Passport {
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    panic(err)
  }
  return parse(string(data))
}

func nbValidPassport(passports []Passport) int {
  count := 0
  for _, p := range passports {
    if p.isValid() {
      count++
    }
  }
  return count
}

func nbValidExtendedPassport(passports []Passport) int {
  count := 0
  for _, p := range passports {
    if p.isValidExtended() {
      count++
    }
  }
  return count
}

func algo1(passports []Passport) int {
  return nbValidPassport(passports)
}

func algo2(passports []Passport) int {
  return nbValidExtendedPassport(passports)
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
  assert_eq(algo1(parseFile("test1.txt")), 2, "test 1");
}

func test2_1(){
  assert(validateBirthYear("2002"), "2002 should be a valid birthday")
  assert(!validateBirthYear("2003"), "2003 should not be a valid birthday")
  
  assert(validateHeight("60in"), "60in should be a valid height")
  assert(validateHeight("190cm"), "190cm should be a valid height")
  assert(!validateHeight("190in"), "190in should not be a valid height")
  assert(!validateHeight("190"), "190 should not be a valid height")
  
  assert(validateHairColor("#123abc"), "#123abc should be a valid hair color")
  assert(!validateHairColor("#123abz"), "#123abz should not be a valid hair color")
  assert(!validateHairColor("123abc"), "123abc should not be a valid hair color")
  
  assert(validateEyeColor("brn"), "brn should be a valid eye color")
  assert(!validateEyeColor("wat"), "wat should not be a valid eye color")
  
  assert(validatePassportID("000000001"), "000000001 should be a valid passport id")
  assert(!validatePassportID("0123456789"), "0123456789 should not be a valid passport id")
}

func test2_2(){
  passports := parseFile("test2_invalid.txt")
  for _, p := range passports {
    assert(!p.isValidExtended(), "all the passports in the test2_invalid.txt file should be invalid");
  }
}

func test2_3(){
  passports := parseFile("test2_valid.txt")
  for _, p := range passports {
    assert(p.isValidExtended(), "all the passports in the test2_valid.txt file should be valid");
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
  test2_2()
  test2_3()
  fmt.Printf("Question2: %d\n", question2())
}
