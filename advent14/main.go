package main

import (
    "fmt"
    "io/ioutil"
    "strings"
    "bufio"
    "os"
    "strconv"
)

func splitTrim(data string, sep string) []string {
  tmp := strings.Split(data, sep)
  for i := range tmp {
    tmp[i] = strings.TrimSpace(tmp[i]) 
  }
  return tmp;
}

type Instruction struct {
  mask string
  addr int
  value int
}

func (i Instruction) nbFloating() int{
  return strings.Count(i.mask,"X")
}

type Memory struct {
  mem map[int]int
}

func newMemory() Memory {
  var m Memory
  m.mem = make(map[int]int)
  return m
}

func parseMask(data string) (isMask bool, maskAnd int, maskOr int){
  tmp := splitTrim(data, "=")
  isMask = (tmp[0] == "mask")
  if(!isMask){ return }
  return
}

func (m *Memory) apply(instr Instruction) {
  var err error
  var mask int64
  
  v := instr.value
  mask, err = strconv.ParseInt(strings.ReplaceAll(instr.mask, "X", "1"), 2, 0)
  if err != nil { panic("can't parse mask !") }
  v &= int(mask)
  mask, err = strconv.ParseInt(strings.ReplaceAll(instr.mask, "X", "0"), 2, 0)
  if err != nil { panic("can't parse mask !") }
  v |= int(mask)
  m.mem[instr.addr] = v
}

func buildPossibleAddresses(addresses *[]int, addr int, mask string){
  if len(mask) == 1 {
    
  }
}

func allPossibleAddresses(addr int, mask string) []int {
  var addresses []int
  addresses = append(addresses,0)
  n := len(mask)
  p := 1
  for i:=(n-1);i>=0;i-- {
    c := mask[i] 
    var addressesIn []int
    if c == '0' {
      for _,addrIn := range addresses {
        if (addr & p) != 0 { addrIn += p }
        addressesIn = append(addressesIn, addrIn)
      }
    } else if c == '1' {
      for _,addrIn := range addresses {
        addrIn += p
        addressesIn = append(addressesIn, addrIn)
      }
    } else {
      for _,addrIn := range addresses {
        addressesIn = append(addressesIn, addrIn)
        addrIn += p
        addressesIn = append(addressesIn, addrIn)
      }
    }
    addresses = addressesIn
    p *= 2
  }
  return addresses
}

func (m *Memory) set(instr Instruction) {
  for _, addr := range allPossibleAddresses(instr.addr, instr.mask) {
    m.mem[addr] = instr.value
  }
}


func (m Memory) print() {
  for k, v := range m.mem { fmt.Println(k, v) }
}


func (m Memory) sum() int {
  r := 0
  for _, v := range m.mem { r += v }
  return r
}

func atoi(str string) int {
  value, err := strconv.Atoi(str)
  if err != nil { panic("Can't parse the number") }
  return value
}

func parseInstruction(data string, mask string) Instruction {
  var instr Instruction
  tmp := splitTrim(data, "=")
  if(tmp[0] == "mask") { panic("can't be a mask !") }
  instr.mask = mask
  instr.addr = atoi(splitTrim(splitTrim(tmp[0], "[")[1], "]")[0])
  instr.value = atoi(tmp[1])
  return instr
}

func parse(data string) []Instruction {
  var instr []Instruction
  lines := strings.Split(data, "\n")
  mask := ""
  for _, line := range lines {
    tmp := splitTrim(line, "=")
    if tmp[0] == "mask" {
      mask = tmp[1]
      continue
    }
    instr = append(instr, parseInstruction(line, mask))
  }
  return instr
}

func parseFile(filename string) []Instruction {
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    panic(err)
  }
  return parse(string(data))
}

func power2(n int) int {
  r := 1
  for n > 0 {
    r *= 2
    n -= 1
  }
  return r
}

func pause(prompt string){
  fmt.Println(prompt)
  input := bufio.NewScanner(os.Stdin)
  input.Scan()
}

func algo1(instr []Instruction) int {
  m := newMemory()
  for _, i := range instr {
    m.apply(i)
  }
  return m.sum()
}
  
func algo2(instr []Instruction) int {
  m := newMemory()
  for _, i := range instr {
    m.set(i)
  }
  return m.sum()
}
/*
func investigate() {
  instr := parseFile("input.txt")
  addrMax := 0
  nbAddr := 0
  for _,i := range instr {
    if addrMax < i.addr { addrMax = i.addr }
    nbAddr += power2(i.nbFloating())
  }
  fmt.Println(addrMax)
  fmt.Println(nbAddr)
}*/

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
  assert_eq(algo1(parseFile("test1.txt")), 165, "1");
}

func question1() int {
  return algo1(parseFile("input.txt"));
}

func test2_1() {
  assert_eq(algo2(parseFile("test2.txt")), 208, "2");
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