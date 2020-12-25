package main

import (
  "fmt"
  "io/ioutil"
  "strings"
  "os"
  "runtime/pprof"
  "runtime"
  "flag"
)

var nbWorkers = 1 + 0*runtime.NumCPU()
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func splitTrim(data string, sep string) []string {
  tmp := strings.Split(data, sep)
  for i := range tmp {
    tmp[i] = strings.TrimSpace(tmp[i]) 
  }
  return tmp;
}

type HexCoord struct {
  x int
  y int
}

func zeroHexCoord() HexCoord {
  var h HexCoord
  h.x = 0
  h.y = 0
  return h
}

func newHexCoord(x int,y int) HexCoord {
  var h HexCoord
  h.x = x
  h.y = y
  return h
}

func (h HexCoord) print() {
  fmt.Printf("(%d,%d)\n", h.x, h.y)
}

func (h HexCoord) move(dir string) HexCoord {
  switch dir {
    case "e": h.x += 1
    case "w": h.x -= 1
    case "se": 
      if h.y % 2 == 0 { h.x += 1 }
      h.y += 1
    case "sw": 
      if h.y % 2 != 0 { h.x -= 1 }
      h.y += 1
    case "ne":
      if h.y % 2 == 0 { h.x += 1 }
      h.y -= 1
    case "nw":
      if h.y % 2 != 0 { h.x -= 1 }
      h.y -= 1
    default: 
      panic("unsupported direction: " + dir)
  }
  return h
}

func parseDirections(line string) []string {
  var directions []string
  push := func(d string){ directions = append(directions, d) }
  buf := ' '
  for _, r := range line {
    switch r {
      case 'n': fallthrough
      case 's':
        switch buf {
          case ' ':
          default:
            panic("should not happen")
        }
        buf = r
      case 'w':
        switch buf {
          case ' ':
            push("w")
          case 'n':
            push("nw")
          case 's':
            push("sw")
          default:
            panic("should not happen")
        }
        buf = ' '
      case 'e':
        switch buf {
          case ' ':
            push("e")
          case 'n':
            push("ne")
          case 's':
            push("se")
          default:
            panic("should not happen")
        }
        buf = ' '
      default:
        panic("should not happen")
    }
  }
  if buf != ' ' { panic("should not happen") }
  return directions
}

func parseDirectionsWorker(lineChan chan string, dirChan chan []string){
  for {
    line, ok := <-lineChan
    if !ok { return }
    dirChan <- parseDirections(line)
  }
}

func parse(data string) [][]string {
  var dirs [][]string
  lines := splitTrim(data, "\n")
  n := len(lines)
  lineChan := make(chan string, n)
  dirChan := make(chan []string, n)
  for i:=0;i<nbWorkers;i++ { go parseDirectionsWorker(lineChan, dirChan) }
  for _, line := range lines { lineChan <- line }
  close(lineChan)
  dirs = make([][]string, n)
  for i:=0;i<n;i++ { dirs[i] = <-dirChan }
  return dirs
}

func parseFile(filename string) [][]string {
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    panic(err)
  }
  return parse(string(data))
}

func applyMove(directions []string) HexCoord{
  cur := zeroHexCoord()
  for _, dir := range directions {
    cur = cur.move(dir)
  }
  return cur
}

func moveWorker(dirChan chan []string, coordChan chan HexCoord){
  for {
    dirs, ok := <-dirChan
    if !ok { return }
    coordChan <- applyMove(dirs)
  }
}

type HexConway struct {
  floor map[HexCoord]bool
}

func initConway(directions [][]string) HexConway {
  var conway HexConway
  conway.floor = make(map[HexCoord]bool)
  n := len(directions)
  dirChan := make(chan []string, n)
  coordChan := make(chan HexCoord, n)
  for i:=0;i<nbWorkers;i++ { go moveWorker(dirChan, coordChan) }
  for _, dirs := range directions { dirChan <- dirs }
  close(dirChan)
  for i:=0;i<n;i++ { 
    coord := <-coordChan 
    v, p := conway.floor[coord]
    if !p {
      conway.floor[coord] = true
    }else{
      conway.floor[coord] = !v
    }
  }
  return conway
}

func (conway HexConway) count() int {
  count := 0
  for _, v := range conway.floor {
    if v { count++ }
  }
  return count
}

func (conway *HexConway) isBlack(coord HexCoord) bool {
  v, p := conway.floor[coord]
  return p && v
}

func (conway *HexConway) minMax() (minX int, maxX int, minY int, maxY int) {
  init := false
  for k, v := range conway.floor {
    if !v { continue }
    if !init {
      minX = k.x
      maxX = k.x
      minY = k.y
      maxY = k.y
      init = true
    } else {
      if minX > k.x { minX = k.x }
      if maxX < k.x { maxX = k.x }
      if minY > k.y { minY = k.y }
      if maxY < k.y { maxY = k.y }
    }
  }
  return
}

func (conway *HexConway) oneStep() {
  newFloor := make(map[HexCoord]bool)
  minX, maxX, minY, maxY := conway.minMax()
  minX -= 1
  maxX += 1
  minY -= 1
  maxY += 1
  for x:=minX;x<=maxX;x++ {
    for y:=minY;y<=maxY;y++ {
      coord := newHexCoord(x,y)
      iAmBlack := conway.isBlack(coord)
      count := 0
      if conway.isBlack(coord.move("e")) { count++ }
      if conway.isBlack(coord.move("w")) { count++ }
      if conway.isBlack(coord.move("ne")) { count++ }
      if conway.isBlack(coord.move("nw")) { count++ }
      if conway.isBlack(coord.move("se")) { count++ }
      if conway.isBlack(coord.move("sw")) { count++ }
      if iAmBlack {
        if count == 1 || count == 2 {
          newFloor[coord] = true
        }
      } else {
        if count == 2 {
          newFloor[coord] = true
        }
      }
    }
  }
  conway.floor = newFloor
}

func algo1(directions [][]string) int {
  return initConway(directions).count()
}

func algo2(directions [][]string, step int) int {
  conway := initConway(directions)
  for i:=0;i<step;i++ {
    conway.oneStep()
  }
  return conway.count()
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
  assert_eq(algo1(parseFile("test1.txt")), 10, "1")
}

func question1() int {
  return algo1(parseFile("input.txt"))
}

func test2_1() {
  directions := parseFile("test1.txt")
  assert_eq(algo2(directions,1), 15, "1")
  assert_eq(algo2(directions,2), 12, "1")
  assert_eq(algo2(directions,3), 25, "1")
  assert_eq(algo2(directions,4), 14, "1")
  assert_eq(algo2(directions,5), 23, "1")
  assert_eq(algo2(directions,6), 28, "1")
  assert_eq(algo2(directions,7), 41, "1")
  assert_eq(algo2(directions,8), 37, "1")
  assert_eq(algo2(directions,9), 49, "1")
  assert_eq(algo2(directions,10), 37, "1")
  assert_eq(algo2(directions,20), 132, "1")
  assert_eq(algo2(directions,30), 259, "1")
  assert_eq(algo2(directions,40), 406, "1")
  assert_eq(algo2(directions,50), 566, "1")
  assert_eq(algo2(directions,60), 788, "1")
  assert_eq(algo2(directions,70), 1106, "1")
  assert_eq(algo2(directions,80), 1373, "1")
  assert_eq(algo2(directions,90), 1844, "1")
  assert_eq(algo2(directions,100), 2208, "1")
}

func question2() int {
  return algo2(parseFile("input.txt"),100)
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
  fmt.Printf("Question1: %d\n", question1())
  test2_1()
  fmt.Printf("Question2: %d\n", question2())
}