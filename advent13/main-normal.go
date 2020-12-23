package main

import (
    "fmt"
    "io/ioutil"
    "strings"
    "bufio"
    "os"
    "strconv"
)

type Problem struct {
  time int
  bus []int
  offset map[int]int
}

func (p Problem) busHere(t int) int {
  for _, n := range p.bus {
    if t % n == 0 { return n }
  }
  return -1
}

func (p Problem) nextBus(t int) (bus int, dt int) {
  dt = 0
  for {
    bus = p.busHere(t + dt)
    if bus > 0 {
      return
    }
    dt += 1
  }
}

func atoi(str string) int {
  value, err := strconv.Atoi(str)
  if err != nil { panic("Can't parse the number") }
  return value
}

func parseProblem(data string) Problem {
  var p Problem
  tmp := strings.Split(data, "\n")
  p.time = atoi(tmp[0])
  p.offset = make(map[int]int)
  busSchedule := strings.Split(tmp[1], ",")
  for idx, busId := range busSchedule {
    if busId != "x" {
      b := atoi(busId)
      p.bus = append(p.bus, b)
      p.offset[b] = idx
    }
  }
  return p
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
  bus, dt := p.nextBus(p.time)
  return bus*dt
}

func gcdExtended(a int, b int) (g int, x int, y int) {
  if a == 0 {
    g = b
    x = 0
    y = 1
  }else{
    g1, x1, y1 := gcdExtended(b%a, a)
    g = g1
    x = y1 - (b/a) * x1
    y = x1
  }
  return
}

func step(a int, b int, o1 int, o2 int) (t int, m int) {
  //fmt.Printf("step: %d %d %d %d\n", a,b,o1,o2)
  d := o1-o2
  g, x, y := gcdExtended(a,b)
  if(g != 1) { panic("all number should be coprime") }
  x *= d
  y *= d
  t = x*a - o1
  t = t % (a*b)
  if t < 0 { t += a*b }
  m = a*b
  //fmt.Printf("return: %d %d\n", t,m)
  return
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

func algo2(p Problem) int {
  n := len(p.bus)
  t := p.offset[p.bus[0]]
  m := p.bus[0]
  for i:=1;i<n;i++ {
    t2, m2 := step(m,p.bus[i],-t,p.offset[p.bus[i]])
    t = t2
    m = m2
  }
  for idx, b := range p.bus {
    fmt.Println(t, b, p.offset[p.bus[idx]])
    assert_eq((t + p.offset[p.bus[idx]])% b, 0, "bus is not here")
  }
  return t
}

func test1_1() {
  assert_eq(algo1(parseFile("test1.txt")), 295, "1");
}

func question1() int {
  return algo1(parseFile("input.txt"));
}

func test2_1() {
  assert_eq(algo2(parseProblem("0\n17,x,13,19")), 3417, "2-1"); // period 4'199
}

func test2_2() {
  assert_eq(algo2(parseProblem("0\n67,7,59,61")), 754018, "2-2"); // period 1'612'352
}

func test2_3() {
  assert_eq(algo2(parseProblem("0\n67,x,7,59,61")), 779210, "2-3"); // period 1'612'352
}

func test2_4() {
  assert_eq(algo2(parseProblem("0\n67,7,x,59,61")), 1261476, "2-4"); // period 1'612'352
}

func test2_5() {
  assert_eq(algo2(parseProblem("0\n1789,37,47,1889")), 1202161486, "2-5"); // period 5'876'813'119
}

func test2_6() {
  assert_eq(algo2(parseFile("test1.txt")), 1068781, "2-final");
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
  test2_4()
  test2_5()
  test2_6()
  fmt.Printf("Question2: %d\n", question2())
}

// 610963397446931 <= too low