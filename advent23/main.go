package main

import (
  "fmt"
  "bufio"
  "os"
  "strconv"
  "runtime/pprof"
  "flag"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var verbose = false

func itoa(i int) string {
  return strconv.Itoa(i)
}

func atoi(str string) int {
  value, err := strconv.Atoi(str)
  if err != nil { panic("Can't parse the number") }
  return value
}

func pause(prompt string){
  fmt.Println(prompt)
  input := bufio.NewScanner(os.Stdin)
  input.Scan()
}


type Cup struct {
  label int
  next *Cup
  m1 *Cup
}

func newCup(label int) *Cup {
  var c Cup
  c.label = label
  c.next = nil
  return &c
}

func (c *Cup) moveNext() *Cup {
  return c.next
}

func (c *Cup) removeNext() *Cup {
  r := c.next
  c.next = r.next
  return r
}

func (c *Cup) insertNext(cup *Cup) {
  cup.next = c.next
  c.next = cup
}

func (c *Cup) search(label int) *Cup {
  cur := c
  for {
    if cur.label == label { return cur }
    cur = cur.moveNext()
  }
}

func (c *Cup) print() {
  start := c
  cur := c
  s := ""
  for {
    s += itoa(cur.label)
    cur = cur.moveNext()
    if(start == cur){
      fmt.Println(s)
      return
    }
  }  
}

type Problem struct {
  cur *Cup
  n int
}

func parseProblem(data string, n int) Problem {
  var p Problem
  var root *Cup
  var cup *Cup
  var prev *Cup
  p.n = 0
  for i, _ := range data {
    cup = newCup(atoi(data[i:i+1]))
    p.n += 1
    if root == nil { root = cup }
    if prev != nil { prev.next = cup }
    prev = cup
  }
  cup.next = root
  for i:=2;i<=p.n;i++ {
    root.search(i).m1 = root.search(i-1)
  }
  biggestOne := root.search(p.n)
  first := true
  for p.n < n {
    p.n += 1
    cup = newCup(p.n)
    prev.next = cup
    if first {
      first = false
      cup.m1 = biggestOne
    }else{
      cup.m1 = prev
    }
    biggestOne = cup
    prev = cup
  }
  root.search(1).m1 = biggestOne
  cup.next = root
  p.cur = root
  return p
}

func (p *Problem) doMoves(n int){
  if verbose && n > 1000 {
    nbStepPct1 := n / 100
    for i:=0;i<n;i++ { 
      if i % nbStepPct1 == 0 {
        fmt.Printf("pct: %.2f\n", 100.0 * float64(i) / float64(n))
      }
      p.doOneMove() 
    }
  } else {
    for i:=0;i<n;i++ { 
      
      p.doOneMove() 
    }
  }
}

func (p *Problem) doOneMove(){
  //p.cur.print()
  cup1 := p.cur.removeNext()
  cup2 := p.cur.removeNext()
  cup3 := p.cur.removeNext()
  dest := p.cur.m1
  for dest == cup1 || dest == cup2 || dest == cup3 { 
    dest = dest.m1
  }
  dest.insertNext(cup1)
  cup1.insertNext(cup2)
  cup2.insertNext(cup3)
  p.cur = p.cur.moveNext()
}

func (p *Problem) result() string {
  start := p.cur.search(1)
  cur := start
  s := ""
  for {
    cur = cur.moveNext()
    if cur == start { return s }
    s += itoa(cur.label)
  }
}

func algo1(p Problem, n int) string {
  p.doMoves(n)
  return p.result()
}

func algo2(p Problem, n int) (int, int) {
  p.doMoves(n)
  start := p.cur.search(1)
  cupN := start.moveNext()
  cupN2 := cupN.moveNext()
  return cupN.label, cupN2.label
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

func assert_eqStr(v string, e string, msg string){
  if v != e {
    fmt.Printf("the test '%s' fail, it give %s instead of %s\n", msg, v, e)
    panic("test failed !")
  }
}

func test1_1() {
  assert_eqStr(algo1(parseProblem("389125467",9),10), "92658374", "1");
  assert_eqStr(algo1(parseProblem("389125467",9),100), "67384529", "1");
}

func question1() string {
  return algo1(parseProblem("974618352",9),100);
}

func test2_1() {
  star1, star2 := algo2(parseProblem("389125467",1000000),10000000)
  assert_eq(star1, 934001, "first star")
  assert_eq(star2, 159792, "second star")
}

func question2() int {
  star1, star2 := algo2(parseProblem("974618352",1000000),10000000)
  return star1 * star2
}

func main() {
  flag.Parse()
  if *cpuprofile != "" {
      f, err := os.Create(*cpuprofile)
      if err != nil { panic(err) }
      pprof.StartCPUProfile(f)
      defer pprof.StopCPUProfile()
  }
  test1_1()
  fmt.Printf("Question1: %s\n", question1())
  test2_1()
  fmt.Printf("Question1: %d\n", question2())
}