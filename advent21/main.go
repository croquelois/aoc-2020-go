package main

import (
  "fmt"
  "io/ioutil"
  "strings"
  "bufio"
  "os"
  "strconv"
  "runtime"
  "sort"
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

type Item struct {
  ingredients []string
  allergens []string
}

type Menu struct {
  items []Item
}

func parseItem(line string) Item {
  var item Item
  tmp := splitTrim(line, "(contains")
  item.ingredients = splitTrim(tmp[0], " ")
  item.allergens = splitTrim(splitTrim(tmp[1], ")")[0], ",")
  return item
}

func parseItemWorker(lineChan chan string, itemChan chan Item){
  for {
    line, ok := <-lineChan
    if !ok { return }
    itemChan <- parseItem(line)
  }
}

func parseMenu(data string) Menu {
  var m Menu
  lines := splitTrim(data, "\n")
  n := len(lines)
  lineChan := make(chan string, n)
  itemChan := make(chan Item, n)
  for i:=0;i<nbWorkers;i++ { go parseItemWorker(lineChan, itemChan) }
  for _, line := range lines { lineChan <- line }
  close(lineChan)
  m.items = make([]Item, n)
  for i:=0;i<n;i++ { m.items[i] = <-itemChan }
  return m
}

func parseFile(filename string) Menu {
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    panic(err)
  }
  return parseMenu(string(data))
}

func pause(prompt string){
  fmt.Println(prompt)
  input := bufio.NewScanner(os.Stdin)
  input.Scan()
}

type Problem struct {
  menu                Menu
  allergensSet        map[string]bool
  ingredientsSet      map[string]bool
  allAllergens        []string
  allIngredients      []string
  possibleAllergens   map[string](map[string]bool)
  allergen2ingredient map[string]string
  ingredient2allergen map[string]string
}

func createProblem(m Menu) Problem {
  var p Problem
  p.menu = m
  p.allergensSet          = make(map[string]bool)
  p.ingredientsSet        = make(map[string]bool)
  p.possibleAllergens   = make(map[string](map[string]bool))
  p.allergen2ingredient = make(map[string]string)
  p.ingredient2allergen = make(map[string]string)
  for _, item := range m.items {
    for _, allergen := range item.allergens {
      p.allergensSet[allergen] = true
    }
    for _, ingredient := range item.ingredients {
      p.ingredientsSet[ingredient] = true
    }
  }
  p.allAllergens = make([]string, len(p.allergensSet))
  i := 0
  for allergen, _ := range p.allergensSet { 
    p.allAllergens[i] = allergen
    i++
  }
  p.allIngredients = make([]string, len(p.ingredientsSet))
  i = 0
  for ingredient, _ := range p.ingredientsSet { 
    p.allIngredients[i] = ingredient
    i++
  }
  return p
}

func (p Problem) print() {
  fmt.Println("possible allergens")
  
  for ingredient, possibleAllergens := range p.possibleAllergens {
    fmt.Println(ingredient)
    fmt.Println(possibleAllergens)
  }
  
  fmt.Println("allergen => ingredient")
  fmt.Println(p.allergen2ingredient)
  
  fmt.Println("ingredient => allergen")
  fmt.Println(p.ingredient2allergen)
}

func (p Problem) getAllIngredients() []string {
  return p.allIngredients
}

func (p Problem) getAllAllergens() []string {
  return p.allAllergens
}

func (p *Problem) setPossibleAllergens(ingredient string, allergens []string){
  p.possibleAllergens[ingredient] = make(map[string]bool)
  for _, allergen := range allergens {
    p.possibleAllergens[ingredient][allergen] = true
  }
}

func (p *Problem) removeAllergens(ingredients []string, allergens []string){
  for _, ingredient := range ingredients {
    for _, allergen := range allergens {
      p.possibleAllergens[ingredient][allergen] = false
    }
  }
}

func (p Problem) getIngredientsNotInside(ingredients []string) []string {
  ingredientsSet := make(map[string]bool)
  for _, ingredient := range ingredients { ingredientsSet[ingredient] = true }
  ret := []string{}
  for _, ingredient := range p.allIngredients {
    _, p := ingredientsSet[ingredient]
    if !p { ret = append(ret, ingredient) }
  }
  return ret
}

func (p Problem) getPossibleIngredients(allergen string) []string {
  ret := []string{}
  for ingredient, possibleAllergens := range p.possibleAllergens {
    _, present := p.ingredient2allergen[ingredient]
    if present { continue }
    contain, presentAllergen := possibleAllergens[allergen] 
    if contain && presentAllergen { ret = append(ret, ingredient) }
  }
  return ret
}

func (p *Problem) setAllergens(ingredient string, allergen string){
  p.allergen2ingredient[allergen] = ingredient
  p.ingredient2allergen[ingredient] = allergen
}

func (p Problem) getAllergensFreeIngredients() []string {
  var ingredients []string
  for _, ingredient := range p.allIngredients {
    _, present := p.ingredient2allergen[ingredient]
    if !present {
      ingredients = append(ingredients, ingredient)
    }
  }
  return ingredients
}

func (p Problem) countOccurence(ingredient string) int {
  count := 0
  for _, item := range p.menu.items {
    for _, i := range item.ingredients {
      if i == ingredient { count++ }
    }
  }
  return count
}

func (p *Problem) solve(){
  for _, ingredients := range p.getAllIngredients() {
    p.setPossibleAllergens(ingredients, p.getAllAllergens())
  }
  for _, item := range p.menu.items {
    ingrediensNotIn := p.getIngredientsNotInside(item.ingredients)
    p.removeAllergens(ingrediensNotIn, item.allergens)
  }
  
  for {
    change := false
    for _, allergen := range p.getAllAllergens() {
      ingredients := p.getPossibleIngredients(allergen)
      if len(ingredients) == 1 {
        p.setAllergens(ingredients[0], allergen)
        change = true
      }
    }
    if !change {
      return
    }
  }
}

func algo1(m Menu) int {
  p := createProblem(m)
  p.solve()
  ingredients := p.getAllergensFreeIngredients()
  count := 0
  for _, ingredient := range ingredients {
    count += p.countOccurence(ingredient)
  }
  return count
}

func algo2(m Menu) string {
  p := createProblem(m)
  p.solve()
  allergens := p.getAllAllergens()
  s := ""
  sort.Strings(allergens)
  for _, allergen := range allergens {
    if len(s) > 0 { s +="," }
    s += p.allergen2ingredient[allergen]
  }
  return s
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
  assert_eq(algo1(parseFile("test1.txt")), 5, "1");
}

func question1() int {
  return algo1(parseFile("input.txt"));
}

func test1_2() {
  assert_eqStr(algo2(parseFile("test1.txt")), "mxmxvkd,sqjhc,fvjkl", "2");
}

func question2() string {
  return algo2(parseFile("input.txt"));
}

func main() {
  nbWorkers = runtime.NumCPU()
  test1_1()
  fmt.Printf("Question1: %d\n", question1())
  test1_2()
  fmt.Printf("Question2: %s\n", question2())
}