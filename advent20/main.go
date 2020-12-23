package main

import (
  "fmt"
  "io/ioutil"
  "strings"
  "bufio"
  "os"
  "strconv"
  "runtime"
  "math"
)

var (
  nbWorkers int
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

type Image struct {
  size int
  pattern []bool
}

func (img Image) print() {
  size := img.size
  for y:=0;y<size;y++ {
    s := ""
    for x:=0;x<size;x++ {
      if img.pattern[y*size+x] {
        s += "#"
      }else{
        s += "."
      }
    }
    fmt.Println(s)
  }  
}

func (img Image) clone() Image{
  var nImg Image
  nImg.pattern = make([]bool,len(img.pattern))
  nImg.size = img.size
  copy(nImg.pattern, img.pattern)
  return nImg
}

func (img Image) flip() Image{
  var nImg Image
  size := img.size
  nImg.pattern = make([]bool,len(img.pattern))
  nImg.size = size
  for x:=0;x<size;x++ {
    for y:=0;y<size;y++ {
      nImg.pattern[y*size+x] = img.pattern[y*size+(size-x-1)]
    }
  }
  return nImg
}

func (img Image) rotate(n int) Image{
  if n == 0 { return img }
  var nImg Image
  size := img.size
  nImg.pattern = make([]bool,len(img.pattern))
  nImg.size = size
  if size % 2 == 1 {
    c := (size-1)/2
    for x:=-c;x<=c;x++ {
      for y:=-c;y<=c;y++ {
        nImg.pattern[(c+y)*size+(c+x)] = img.pattern[(c-x)*size+(c+y)]
      }
    }
  } else {
    cf := float64(size-1)/2.0
    for x:=0;x<size;x++ {
      for y:=0;y<size;y++ {
        xf := float64(x)-cf
        yf := float64(y)-cf
        rxf := yf
        ryf := -xf
        rx := int(rxf+cf)
        ry := int(ryf+cf)
        nImg.pattern[y*size+x] = img.pattern[ry*size+rx]
      }
    }
  }
  return nImg.rotate(n-1)
}

func parseImage(lines []string) Image {
  var img Image
  img.size = len(lines[0])
  img.pattern = make([]bool, img.size*img.size)
  i := 0
  for _, line := range lines {
    for _, rune := range line {
      img.pattern[i] = (rune=='#')
      i++
    }
  }
  return img
}

type Tile struct {
  id int
  flipped bool
  rotation int
  img Image
}

func parseTile(dataChan chan string, tileChan chan Tile) {
  for {
    data, ok := <-dataChan
    if !ok { return }
    var t Tile
    lines := splitTrim(data, "\n")
    t.id = atoi(splitTrim(splitTrim(lines[0], ":")[0]," ")[1])
    t.img = parseImage(lines[1:])
    tileChan <- t
    tileChan <- t.flip()
    tileChan <- t.rotate(1)
    tileChan <- t.rotate(1).flip()
    tileChan <- t.rotate(2)
    tileChan <- t.rotate(2).flip()
    tileChan <- t.rotate(3)
    tileChan <- t.rotate(3).flip()
  }
}

func (t Tile) hash() string{
  s := strconv.Itoa(t.id)
  if t.flipped {
    s += "-Flipped"
  }
  if t.rotation > 0 {
    s += "-Rot" + strconv.Itoa(t.rotation)
  }
  return s
}

func (t Tile) print() {
  t.img.print()
}

func (t Tile) flip() Tile{
  var nt Tile
  nt.id = t.id
  nt.flipped = !t.flipped
  nt.rotation = t.rotation
  nt.img = t.img.flip()
  return nt
}

func (t Tile) rotate(n int) Tile{
  if n%4 == 0 { return t }
  var nt Tile
  nt.id = t.id
  nt.flipped = t.flipped
  nt.rotation = (t.rotation+n)%4
  nt.img = t.img.rotate(n)
  return nt
}

type Problem struct {
  tiles map[string]Tile
  n int // Size of a side
  n2 int // Number of tiles
}

func parse(data string) Problem {
  var p Problem
  tileData := splitTrim(data, "\n\n")
  n := len(tileData)
  dataChan := make(chan string, n)
  tileChan := make(chan Tile, 8*n)
  for i:=0;i<nbWorkers;i++ { go parseTile(dataChan, tileChan) }
  for _, s := range tileData { dataChan <- s }
  close(dataChan)
  p.tiles = make(map[string]Tile)
  for i:=0;i<8*n;i++ { 
    tile := <-tileChan
    p.tiles[tile.hash()] = tile
  }
  p.n2 = n
  p.n = int(math.Round(math.Sqrt(float64(n))))
  return p
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

func (p *Problem) checkBorderUpDown(uId string, dId string) bool {
  imgUp := p.tiles[uId].img
  imgDown := p.tiles[dId].img
  n := imgUp.size
  for i:=0;i<n;i++ {
    if imgUp.pattern[i + n*(n-1)] != imgDown.pattern[i] {
      return false
    }
  }  
  return true
}

func (p *Problem) checkBorderLeftRight(lId string, rId string) bool {
  imgLeft := p.tiles[lId].img
  imgRight := p.tiles[rId].img
  n := imgLeft.size
  for i:=0;i<n;i++ {
    if imgLeft.pattern[i*n + (n-1)] != imgRight.pattern[i*n] {
      return false
    }
  }  
  return true
}

func (p *Problem) checkNeighbors(puzzle []string, i int) bool {
  if (i % p.n) > 0 && len(puzzle[i-1]) > 0 {
    if !p.checkBorderLeftRight(puzzle[i-1], puzzle[i]) { 
      return false 
    }
  }
  if (i % p.n) < (p.n-1) && len(puzzle[i+1]) > 0 {
    if !p.checkBorderLeftRight(puzzle[i], puzzle[i+1]) { 
      return false 
    }
  }
  if i >= p.n && len(puzzle[i-p.n]) > 0 {
    if !p.checkBorderUpDown(puzzle[i-p.n], puzzle[i]) { 
      return false 
    }
  }
  if i < p.n*(p.n-1) && len(puzzle[i+p.n]) > 0 {
    if !p.checkBorderLeftRight(puzzle[i], puzzle[i+p.n]) { 
      return false 
    }
  }  
  return true
}

func (p *Problem) search(argPuzzle []string, i int, argUsed map[int]bool) ([]string, bool) {
  if i == p.n2 { return argPuzzle, true }
  puzzle := make([]string, len(argPuzzle))
  copy(puzzle, argPuzzle)
  used := make(map[int]bool)
  for k,v := range argUsed { if v { used[k] = v } }
  for _, tile := range p.tiles {
    isUsed := used[tile.id]
    if isUsed { continue }
    puzzle[i] = tile.hash()
    if !p.checkNeighbors(puzzle, i) { continue }
    used[tile.id] = true
    retPuzzle, ok := p.search(puzzle, i+1, used)
    if ok { return retPuzzle, true }
    used[tile.id] = false
  }
  return []string{}, false
}

func (p *Problem) rebuildImage() Image {
  puzzle := make([]string, p.n2)
  used := make(map[int]bool)
  retPuzzle, ok := p.search(puzzle, 0, used)
  if !ok { panic("can't rebuild the image") }
  var img Image
  tSize := p.tiles[retPuzzle[0]].img.size
  img.size = p.n*(tSize-2)
  img.pattern = make([]bool,img.size*img.size)
  for x:=0;x<p.n;x++ {
    for y:=0;y<p.n;y++ {
      tImg := p.tiles[retPuzzle[x + y*p.n]].img
      for tx:=1;tx<tSize-1;tx++ {
        for ty:=1;ty<tSize-1;ty++ {
          ix := x*(tSize-2)+(tx-1)
          iy := y*(tSize-2)+(ty-1)
          img.pattern[ix + iy*img.size] = tImg.pattern[tx + ty*tSize]
        }
      }
    }
  }
  return img
}

func algo1(p Problem) int {
  puzzle := make([]string, p.n2)
  used := make(map[int]bool)
  retPuzzle, ok := p.search(puzzle, 0, used)
  if !ok { panic("unable to found a solution") }    
  corners := 1
  corners *= p.tiles[retPuzzle[0]].id
  corners *= p.tiles[retPuzzle[p.n-1]].id
  corners *= p.tiles[retPuzzle[p.n*(p.n-1)]].id
  corners *= p.tiles[retPuzzle[p.n*p.n-1]].id
  return corners
}

func getSeaMonster() chan[2]int {
//  000000000011111111112
//  012345678901234567890
// 0                  #
// 1#    ##    ##    ###
// 2 #  #  #  #  #  #   
  c := make(chan[2]int, 15)
  c <- [2]int{18,0}
  c <- [2]int{ 0,1}
  c <- [2]int{ 5,1}
  c <- [2]int{ 6,1}
  c <- [2]int{11,1}
  c <- [2]int{12,1}
  c <- [2]int{17,1}
  c <- [2]int{18,1}
  c <- [2]int{19,1}
  c <- [2]int{ 1,2}
  c <- [2]int{ 4,2}
  c <- [2]int{ 7,2}
  c <- [2]int{10,2}
  c <- [2]int{13,2}
  c <- [2]int{16,2}
  close(c)
  return c
}

func (img Image) isSeaMonster(x int, y int) bool {
  for p := range getSeaMonster() {
    px := p[0]
    py := p[1]
    if x+px >= img.size { return false }
    if y+py >= img.size { return false }
    if !img.pattern[x+px + (y+py)*img.size] { return false }
  }
  return true
}

func (img Image) hasSeaMonster() bool {
  for x:=0;x<img.size;x++ {
    for y:=0;y<img.size;y++ {
      if img.isSeaMonster(x,y) {
        return true
      }
    }
  }
  return false
}

func (img *Image) removeSeaMonster(x int, y int) {
  for p := range getSeaMonster() {
    img.pattern[x+p[0] + (y+p[1])*img.size] = false
  }
}

func (img Image) removeSeaMonsters() Image {
  newImg := img.clone()
  for x:=0;x<img.size;x++ {
    for y:=0;y<img.size;y++ {
      if img.isSeaMonster(x,y) {
        newImg.removeSeaMonster(x,y)
      }
    }
  }
  return newImg
}

func (img Image) nbWave() int {
  count := 0
  for x:=0;x<img.size;x++ {
    for y:=0;y<img.size;y++ {
      if img.pattern[x + y*img.size] {
        count++
      }
    }
  }
  return count
}

func searchSeaMonster(imageChan chan Image, roughnessChan chan int){
  for {
    image, ok := <- imageChan
    if !ok { return }
    if !image.hasSeaMonster() {
      roughnessChan <- 0
      continue
    }
    roughnessChan <- image.removeSeaMonsters().nbWave()
  }
}

func algo2(p Problem) int {
  n := 8
  imageChan := make(chan Image, n)
  roughnessChan := make(chan int, n)
  for i:=0;i<nbWorkers;i++ { go searchSeaMonster(imageChan, roughnessChan) }
  
  image := p.rebuildImage()
  imageChan <- image
  imageChan <- image.flip()
  imageChan <- image.rotate(1)
  imageChan <- image.rotate(1).flip()
  imageChan <- image.rotate(2)
  imageChan <- image.rotate(2).flip()
  imageChan <- image.rotate(3)
  imageChan <- image.rotate(3).flip()
  close(imageChan)
  
  roughness := 0
  for i:=0;i<n;i++ {
    roughness += <- roughnessChan 
  }
  return roughness
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
  p := parseFile("test1.txt")
  assert(p.checkBorderLeftRight("1951-Flipped-Rot2", "2311-Flipped-Rot2"), "both should match");
}

func test1_2() {
  assert_eq(algo1(parseFile("test1.txt")), 20899048083289, "rebuild image");
}

func question1() int {
  return algo1(parseFile("input.txt"));
}

func test2_1() {
  assert_eq(algo2(parseFile("test1.txt")), 273, "sea monster");
}

func question2() int {
  return algo2(parseFile("input.txt"));
}

func main() {
  nbWorkers = runtime.NumCPU()
  test1_1()
  test1_2()
  fmt.Printf("Question1: %d\n", question1())
  test2_1()
  fmt.Printf("Question2: %d\n", question2())
}