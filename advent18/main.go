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

type Node struct {
  op string
  v int
  l *Node
  r *Node
}

func parseLine(line string) []string {
  var tokens []string
  var buf = ""
  for _, r := range line {
    switch r {
      case ' ': continue
      case '+': fallthrough
      case '*': fallthrough
      case '(': fallthrough
      case ')': 
        if len(buf) > 0 {
          tokens = append(tokens, buf)
          buf = ""
        }
        tokens = append(tokens, string(r))
      default: buf += string(r)
    }
  }
  if len(buf) > 0 {
    tokens = append(tokens, buf)
    buf = ""
  }
  return tokens
}

func createGraph(tokens []string, idx int) (*Node, int) {
  var cur *Node
  nbTokens := len(tokens)
  for i:=idx;i<nbTokens;i++ {
    token := tokens[i]
    switch token {
      case "+": fallthrough
      case "*":
        var newNode Node
        newNode.op = token
        newNode.l = cur
        cur = &newNode
      case "(":
        var newNodePtr *Node
        newNodePtr, i = createGraph(tokens, i+1)
        if cur != nil {
          cur.r = newNodePtr
        } else {
          cur = newNodePtr
        }
      case ")": 
        return cur, i
      default:
        var newNode Node
        newNode.op = "value"
        newNode.v = atoi(token)
        if cur != nil {
          cur.r = &newNode
        } else {
          cur = &newNode
        }
    }
  }
  if cur == nil { 
    panic("empty tree") 
  }
  return cur, nbTokens
}

func createGraphAtom(tokens []string, idx int) (*Node, int) {
  switch tokens[idx] {
    case "(":
      newNode, idx := createGraphExpr(tokens, idx+1)
      if tokens[idx] != ")" { panic("parsing error") }
      return newNode, idx+1
    default:
      var newNode Node
      newNode.op = "value"
      newNode.v = atoi(tokens[idx])
      return &newNode, idx+1
  }
}

func createGraphAddExpr(tokens []string, idx int) (*Node, int) {
  var rootNode *Node
  rootNode, idx = createGraphAtom(tokens, idx)
  for idx < len(tokens) && tokens[idx] == "+" {
    idx += 1
    var newNode Node
    newNode.op = "+"
    newNode.l = rootNode
    newNode.r, idx = createGraphAtom(tokens, idx)
    rootNode = &newNode
  }
  return rootNode, idx
}

func createGraphExpr(tokens []string, idx int) (*Node, int) {
  var rootNode *Node
  rootNode, idx = createGraphAddExpr(tokens, idx)
  for idx < len(tokens) && tokens[idx] == "*" {
    idx += 1
    var newNode Node
    newNode.op = "*"
    newNode.l = rootNode
    newNode.r, idx = createGraphAddExpr(tokens, idx)
    rootNode = &newNode
  }
  return rootNode, idx
}

func createGraph1(tokens []string) *Node {
  root, _ := createGraph(tokens, 0)
  return root
}

func createGraph2(tokens []string) *Node {
  root, _ := createGraphExpr(tokens, 0)
  return root
}

func (n *Node) compute() int {
  switch n.op {
    case "+": return n.l.compute() + n.r.compute()
    case "*": return n.l.compute() * n.r.compute()
    case "value": return n.v
  }
  panic("unexpected operation")
}

func (n *Node) print(depth string) {
  if n.op == "value" { 
    println(depth + strconv.Itoa(n.v))
  } else {
    println(depth + n.op)
    n.l.print(depth + " ")
    n.r.print(depth + " ")
  }
}

func parse(data string) [][]string {
  lines := splitTrim(data, "\n")
  ret := make([][]string, len(lines))
  for i, line := range lines {
    ret[i] = parseLine(line)
  }
  return ret
}

func parseFile(filename string) [][]string {
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

func algo1(lines [][]string) int {
  sum := 0
  for _, line := range lines {
    sum += createGraph1(line).compute()
  }
  return sum
}

func algo2(lines [][]string) int {
  sum := 0
  for _, line := range lines {
    sum += createGraph2(line).compute()
  }
  return sum
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

func testLine(line string, expected int) {
  assert_eq(createGraph1(parseLine(line)).compute(), expected, line)
}

func testLineAlgo2(line string, expected int) {
  assert_eq(createGraph2(parseLine(line)).compute(), expected, line)
}

func test1_1() {
  testLine("2 * 3", 2*3);
}

func test1_2() {
  testLine("2 * 3 + (4 * 5)", 2*3+4*5);
}

func test1_3() {
  testLine("5 + (8 * 3 + 9 + 3 * 4 * 3)", 5+(8*3+9+3)*4*3);
}

func test1_4() {
  testLine("5 * 9 * (7 * 3 * 3 + 9 * 3 + (8 + 6 * 4))", 12240);
}

func test1_5() {
  assert_eq(algo1(parseFile("test1.txt")), 26+437+12240+13632, "total");
}

func question1() int {
  return algo1(parseFile("input.txt"));
}

func test2_1() {
  testLineAlgo2("8 * 3 + 9 + 3 * 4 * 3", 8 * (3 + 9 + 3) * 4 * 3);
}

func test2_2() {
  testLineAlgo2("2 * 3 + (4 * 5)", 2*(3+4*5));
}

func test2_3() {
  testLineAlgo2("5 + (8 * 3 + 9 + 3 * 4 * 3)", 1445);
}

func test2_4() {
  testLineAlgo2("5 * 9 * (7 * 3 * 3 + 9 * 3 + (8 + 6 * 4))", 669060);
}

func test2_5() {
  testLineAlgo2("6 + (4 * 4 + 8 + 2) * 2 + 4 * 5", 1860);
}

func test2_6() {
  testLineAlgo2("((2 + 4 * 9) * (6 + 9 * 8 + 6) + 6) + 2 + 4 * 2", 23340);
}

func test2_7() {
  assert_eq(algo2(parseFile("test1.txt")), 46+1445+669060+23340, "total");
}

func question2() int {
  return algo2(parseFile("input.txt"));
}

func main() {
  test1_1()
  test1_2()
  test1_3()
  test1_4()
  fmt.Printf("Question1: %d\n", question1())
  test2_1()
  test2_2()
  test2_3()
  test2_4()
  test2_5()
  test2_6()
  test2_7()
  fmt.Printf("Question2: %d\n", question2())
}