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

type Key struct {
  x []int
}

func newKeyFrom2d(x int, y int, n int) Key {
  var k Key
  k.x = append(k.x, x)
  k.x = append(k.x, y)
  for i:=2;i<n;i++ {
    k.x = append(k.x, 0)
  }
  return k
}

func newKeyFrom2dComplete(x int, y int, s []int) Key {
  var k Key
  k.x = append(k.x, x)
  k.x = append(k.x, y)
  for _,v := range s { k.x = append(k.x, v) }
  return k
}

func newKeyExtend(k1 Key, x int) Key {
  var k Key
  k.x = make([]int, len(k1.x))
  copy(k.x,k1.x)
  k.x = append(k.x, x)
  return k
}

func newKey1d(x int) Key {
  var k Key
  k.x = append(k.x, x)
  return k
}

func eachCellGenerator(minX []int, maxX []int) func() (Key, bool) {
  pos := 0
  n := len(minX)
  //fmt.Println("eachCellGenerator", n)
  if n == 1 {
    return func() (Key, bool) {
      if pos > (maxX[0] - minX[0]) {
        return Key{}, false
      }
      pos++
      return newKey1d(minX[0]+(pos-1)), true
    }
  }
  next := eachCellGenerator(minX[:(n-1)], maxX[:(n-1)])
  cur, ok := next()
  return func() (Key, bool) {
    if pos > (maxX[n-1] - minX[n-1]) {
      if !ok { return Key{}, false }
      cur, ok = next()
      if !ok { return Key{}, false }
      pos = 0
    }
    pos++
    return newKeyExtend(cur, minX[n-1]+(pos-1)), true
  }
}

func (k Key) hash() string {
  s := ""
  for i,v := range k.x { 
    if i != 0 {
      s += "," 
    }
    s += strconv.Itoa(v)
  }
  return s
}


func sum(k1 Key, k2 Key) Key {
  var k Key
  for i, _ := range k1.x { 
    k.x = append(k.x, k1.x[i] + k2.x[i])
  }
  return k
}

func (k Key) isZero() bool {
  for _, v := range k.x {
    if v != 0 { 
      return false 
    }
  }
  return true
}

func (k Key) neighbors() []Key {
  var ret []Key
  minX := make([]int, len(k.x))
  maxX := make([]int, len(k.x))
  for i, _ := range k.x { 
    minX[i] = -1
    maxX[i] = 1
  }
  next := eachCellGenerator(minX, maxX)
  for {
    d, ok := next()
    if !ok { return ret }
    //fmt.Println("check neighbor", d.hash())
    if !d.isZero() {
      ret = append(ret, sum(k,d))
    }
  }
  return ret
}

type Conway struct {
  cell map[string]Key
}

func parseConway(data string, n int) Conway {
  var c Conway
  c.cell = make(map[string]Key)
  lines := splitTrim(data, "\n")
  for y, line := range lines {
    for x, r := range line {
      if r == '#' { 
        k := newKeyFrom2d(x,y,n)
        c.cell[k.hash()] = k
      }
    }
  }
  return c
}

func checkHash(cell map[string]Key) () {
  for hash, k := range cell {
    if hash != k.hash() {
      fmt.Println(hash, k, k.hash())
      panic("incorrect hash")
    }
  }
}

func (c *Conway) maxMin() (minX []int, maxX []int, empty bool) {
  init := false
  for _, k := range c.cell {
    if !init {
      minX = make([]int, len(k.x))
      maxX = make([]int, len(k.x))
      copy(minX, k.x)
      copy(maxX, k.x)
      init = true
    }else{
      for i, v := range k.x {
        if minX[i] > v { minX[i] = v }
        if maxX[i] < v { maxX[i] = v }
      }
    }
  }
  empty = !init
  return
}

func (c *Conway) print(slice []int) {
  minX, maxX, empty := c.maxMin()
  if empty { 
    fmt.Println("<<EMPTY>>") 
    return
  }
  for y:=minX[1];y<=maxX[1];y++ {
    s := ""
    for x:=minX[0];x<=maxX[0];x++ {
      k := newKeyFrom2dComplete(x,y,slice)
      _, p := c.cell[k.hash()] 
      if p {
        s += "#"
      } else {
        s += "."
      }
    }
    fmt.Println(s) 
  }
  return
}

func (c *Conway) step() {
  newCell := make(map[string]Key)
  minX, maxX, empty := c.maxMin()
  if empty { return }
  for i, _ := range minX { minX[i] -= 1 }
  for i, _ := range maxX { maxX[i] += 1 }
  //fmt.Println(minX, maxX)
  next := eachCellGenerator(minX, maxX)
  for {
    cur, ok := next()
    if !ok {
      c.cell = newCell
      return
    }
    //fmt.Println("check cell", cur.hash())
    count := 0
    for _, k := range cur.neighbors() {
      _, p := c.cell[k.hash()]
      if p { count++ }
    }
    //fmt.Println("nb neighbors", count)
    _, p := c.cell[cur.hash()]
    if p { 
      if count == 2 || count == 3 {
        newCell[cur.hash()] = cur
      }
    }else{
      if count == 3 {
        newCell[cur.hash()] = cur
      }
    }
  }
}

func (c *Conway) count() int {
  return len(c.cell)
}

func parseFile(filename string, n int) Conway {
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    panic(err)
  }
  return parseConway(string(data), n)
}

func pause(prompt string){
  fmt.Println(prompt)
  input := bufio.NewScanner(os.Stdin)
  input.Scan()
}

func algo1(c Conway) int {
  for i:=0;i<6;i++ { 
    /*fmt.Println("After", i,"cycle:")
    for z:=-i;z<=i;z++ {
      for w:=-i;w<=i;w++ {
        fmt.Println("z=", z,"w=",w)
        c.print([]int{z,w})
      }
    }*/
    c.step() 
  }
  /*fmt.Println("Final step")
  c.print([]int{0,0})*/
  return c.count()
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
  assert_eq(algo1(parseFile("test1.txt",3)), 112, "3dim");
}

func question1() int {
  return algo1(parseFile("input.txt",3));
}

func test2_1() {
  assert_eq(algo1(parseFile("test1.txt",4)), 848, "4dim");
}

func question2() int {
  return algo1(parseFile("input.txt",4));
}

func main() {
  test1_1()
  fmt.Printf("Question1: %d\n", question1())
  test2_1()
  fmt.Printf("Question2: %d\n", question2())
}