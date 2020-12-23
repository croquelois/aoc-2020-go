package main

import (
  "fmt"
  "io/ioutil"
  "strings"
  "bufio"
  "os"
  "strconv"
)

func itoa(i int) string {
  return strconv.Itoa(i)
}

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

type Deck struct {
  cards []int
}

func parseDeck(data string) Deck {
  var d Deck
  tmp := splitTrim(data, "\n")
  for _, card := range tmp[1:] {
    d.cards = append(d.cards, atoi(card))
  }
  return d
}

func (d Deck) newDeck(n int) Deck {
  var nd Deck
  nd.cards = make([]int, n)
  copy(nd.cards, d.cards[:n])
  return nd
}

func (d *Deck) pop() int {
  card := d.cards[0]
  d.cards = d.cards[1:]
  return card
}

func (d *Deck) push(card int) {
  d.cards = append(d.cards, card)
}

func (d Deck) len() int {
  return len(d.cards)
}

func (d Deck) score() int {
  n := len(d.cards)
  score := 0
  for i, card := range d.cards {
    score += (n-i) * card
  }
  return score
}

type Game struct {
  deck1 Deck
  deck2 Deck
}

func parseGame(data string) Game {
  var game Game
  tmp := splitTrim(data, "\n\n")
  game.deck1 = parseDeck(tmp[0])
  game.deck2 = parseDeck(tmp[1])
  return game
}

func (g *Game) playCombat() int {
  for {
    card1 := g.deck1.pop()
    card2 := g.deck2.pop()
    if card1 > card2 {
      g.deck1.push(card1)
      g.deck1.push(card2)
    }else{
      g.deck2.push(card2)
      g.deck2.push(card1)
    }
    if len(g.deck1.cards) == 0 { return g.deck2.score() }
    if len(g.deck2.cards) == 0 { return g.deck1.score() }
  }
}

func hash(deck1 Deck, deck2 Deck) string {
  s := ""
  for _, card := range deck1.cards { s += itoa(card) }
  s = "-"
  for _, card := range deck2.cards { s += itoa(card) }
  return s
}

func playRecursiveCombat(deck1 Deck, deck2 Deck) (int, Deck) {
  loopCheck := make(map[string]bool)
  for {
    h := hash(deck1, deck2)
    _,p := loopCheck[h]
    if p { return 1, deck1 }
    loopCheck[h] = true
    card1 := deck1.pop()
    card2 := deck2.pop()
    victory := 1
    if deck1.len() >= card1 && deck2.len() >= card2 {
      newDeck1 := deck1.newDeck(card1)
      newDeck2 := deck2.newDeck(card2)
      victory, _ = playRecursiveCombat(newDeck1, newDeck2)
    }else{
      if card1 < card2 { victory = 2 }
    }
    if victory == 1 {
      deck1.push(card1)
      deck1.push(card2)
    }else{
      deck2.push(card2)
      deck2.push(card1)
    }
    if len(deck1.cards) == 0 { return 2, deck2 }
    if len(deck2.cards) == 0 { return 1, deck1 }
  }
}

func (g *Game) playRecursiveCombat() int {
  _, deck := playRecursiveCombat(g.deck1, g.deck2)
  return deck.score()
}

func parseFile(filename string) Game {
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    panic(err)
  }
  return parseGame(string(data))
}

func pause(prompt string){
  fmt.Println(prompt)
  input := bufio.NewScanner(os.Stdin)
  input.Scan()
}

func algo1(game Game) int {
  return game.playCombat()
}

func algo2(game Game) int {
  return game.playRecursiveCombat()
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
  assert_eq(algo1(parseFile("test1.txt")), 306, "1");
}

func question1() int {
  return algo1(parseFile("input.txt"));
}

func test2_1() {
  assert_eq(algo2(parseFile("test1.txt")), 291, "1");
}

func question2() int {
  return algo2(parseFile("input.txt"));
}

func main() {
  test1_1()
  fmt.Printf("Question1: %d\n", question1())
  test2_1()
  fmt.Printf("Question1: %d\n", question2())
}