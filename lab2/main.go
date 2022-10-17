package main

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/exp/slices"
	"log"
	"os"
	"sort"
	"strings"
	"unicode"
)

type Node struct {
	Value    string
	Children []Node
	Label    string
}

type Edge struct {
	Parent    string
	Child     string
	EdgeLabel string // derivative with respect to
}

type BrzozowskiAuto struct {
	Regex    string
	Alphabet []string
	Nodes    []string
	Edges    []Edge
}

func getFirstAltIndex(str string) int {
	outBracketsCount := 0
	for i := range str {
		if outBracketsCount == 0 && str[i] == '|' {
			return i
		} else if str[i] == '(' {
			outBracketsCount++
		} else if str[i] == ')' {
			outBracketsCount--
		}
	}
	return -1
}

func getListOfAltSubstrings(str string) []string {
	s := make([]string, 0)
	k := len(str)
	i := 0
	for i <= k {
		alt := getFirstAltIndex(str)
		if alt == -1 {
			s = append(s, str)
			return s
		} else {
			s = append(s, string([]byte(str)[:alt]))
			str = string([]byte(str)[alt+1:])
			k = len(str)
			i = 0
		}
		i++
	}
	return append(s, "")
}

func getFirstConcatIndex(str string) int {
	outBracketsCount := 0
	for i := range str {
		if i > 0 {
			if str[i-1] == '(' {
				outBracketsCount++
			} else if str[i-1] == ')' {
				outBracketsCount--
			}
			if outBracketsCount == 0 {
				if str[i] == '(' || unicode.IsLetter(rune(str[i])) {
					return i
				}
			}
		}
	}
	return -1
}

func getListOfSubstrings(str string) []string {
	s := make([]string, 0)
	k := len(str)
	i := 0
	for i <= k {
		concAfter := getFirstConcatIndex(str)
		if concAfter == -1 {
			s = append(s, str)
			return s
		} else {
			s = append(s, string([]byte(str)[:concAfter]))
			str = string([]byte(str)[concAfter:])
			k = len(str)
			i = 0
		}
		i++
	}
	return s
}

func closingParenthesis(str string) int {
	c := 0
	for i := range str {
		if str[i] == '(' {
			c++
		} else if str[i] == ')' {
			c--
			if c == 0 {
				return i
			}
		}
	}
	return -1
}

func removeExtraParenthesis(regex string) string {
	n := len(regex)
	if n > 0 {
		if regex[0] == '(' && regex[n-1] == ')' && closingParenthesis(regex) == n-1 {
			regex = regex[1 : n-1]
		}
	}
	return regex
}

func regexParse(regex string) Node {
	if regex == "" {
		return Node{"ε", nil, "Empty"}
	}
	pairBracketsCount := 0
	for i := range regex {
		if regex[i] == '(' {
			pairBracketsCount++
		}
		if regex[i] == ')' {
			pairBracketsCount--
		}
	}
	if pairBracketsCount != 0 {
		err := errors.New("regex has extra no-pair brackets")
		log.Fatalf("Error with regex: %v", err)
		os.Exit(1)
	}
	return parseAlt(regex)
}

func parseAlt(regex string) Node {
	value := ""
	regex = removeExtraParenthesis(regex)
	if getFirstAltIndex(regex) == -1 {
		return parseCon(regex)
	}
	children := getListOfAltSubstrings(regex)
	childrenNodes := make([]Node, 0, len(regex))
	for i := range children {
		child := parseCon(children[i])
		childrenNodes = append(childrenNodes, child)
	}
	for i := range childrenNodes {
		value += childrenNodes[i].Value + "|"
	}
	return Node{deleteExtraAlt(value), childrenNodes, "Alt"}
}

func parseCon(regex string) Node {
	regex = removeExtraParenthesis(regex)
	children := getListOfSubstrings(regex)
	if len(children) == 1 {
		return parseStar(regex)
	}

	childrenNodes := make([]Node, 0, len(regex))
	for i := 0; i < len(children); i++ {
		child := parseStar(children[i])
		childrenNodes = append(childrenNodes, child)
	}
	return Node{regex, childrenNodes, "Concat"}
}

func parseStar(regex string) Node {
	n := len(regex)
	if regex == "" || regex == "ε" {
		return Node{"ε", nil, "Empty"}
	}
	if len(regex) == 1 && unicode.IsLetter(rune(regex[0])) {
		return Node{regex, nil, "Sym"}
	} else if regex[n-1] == '*' {
		if len(regex) == 2 && unicode.IsLetter(rune(regex[0])) {
			child := Node{string(regex[0]), nil, "Sym"}
			return Node{regex, []Node{child}, "Star"}
		} else if regex[0] == '(' {
			child := parseAlt(regex[:n-1])
			return Node{regex, []Node{child}, "Star"}
		}
	} else if regex[0] == '(' && regex[n-1] == ')' {
		return parseAlt(regex)
	}
	return Node{regex, nil, "Null"}
}

func printNodes(start Node) {
	for i := range start.Children {
		fmt.Printf("\t%s [label = %s] -> ", start.Value, start.Label)
		fmt.Printf("%s [label = %s]\n", start.Children[i].Value, start.Children[i].Label)
	}
	for i := range start.Children {
		printNodes(start.Children[i])
	}
}

func printGraph(start Node) {
	fmt.Println("\ngraph {")
	printNodes(start)
	fmt.Printf("}\n")
}

func printAutomata(automata BrzozowskiAuto) {
	fmt.Println("\nBrzozovski automat {")
	for _, e := range automata.Edges {
		fmt.Printf("\t%s  --- %s ---> %s\n", e.Parent, e.EdgeLabel, e.Child)
	}
	fmt.Printf("}\n\n")
}

func getAlphabetForRegex(t string) []string {
	alp := make([]string, 0)
	for i := range t {
		if unicode.IsLetter(rune(t[i])) && t[i] != 206 && !slices.Contains(alp, string(t[i])) { //'ε' = 206
			alp = append(alp, string(t[i]))
		}
	}
	return alp
}

func notNullNotEmpty(arr []Node) []Node {
	upd := make([]Node, 0)
	for i := range arr {
		if arr[i].Value == "" || arr[i].Value == "ε" {
			continue
		} else {
			upd = append(upd, arr[i])
		}
	}
	return upd
}

func deleteExtraAlt(str string) string {
	n := len(str)
	if n > 0 {
		if string(str[n-1]) == "|" {
			str = str[:n-1]
		}
	}
	return str
}

func areNodesDuplicates(el string, arr []string) bool {
	for _, a := range arr {
		if a == el {
			return true
		}
	}
	return false
}

func areEdgesDuplicates(el Edge, arr []Edge) bool {
	for _, a := range arr {
		if a == el {
			return true
		}
	}
	return false
}

// Определим вспомогательную функцию lambda , такую, что
// если аргумент принимает ε, то она возвращает ε, иначе – ""
// (где "" - пустое множество, т.е. противоречие)
func (regex Node) lambda() Node {
	children1 := make([]Node, 0)
	var der Node
	der.Value = ""
	if regex.Label == "Alt" {
		for i := range regex.Children {
			children1 = append(children1, regex.Children[i].lambda())
		}
		for i := range children1 {
			if children1[i].Value == "" {
				continue
			} else {
				der.Value += children1[i].Value + "|"
			}
		}
		der = Node{deleteExtraAlt(der.Value), children1, "Alt"}
	} else if regex.Label == "Concat" {
		for i := range regex.Children {
			children1 = append(children1, regex.Children[i].lambda())
		}
		for i := range children1 {
			if children1[i].Value == "" {
				der.Value = ""
				break
			} else {
				der.Value += children1[i].Value
			}
		}
		der = Node{der.Value, children1, "Concat"}
	} else if regex.Label == "Star" {
		der = Node{"ε", nil, "Empty"}
	} else if regex.Label == "Sym" {
		der = Node{"", nil, "Null"}
	} else if regex.Label == "Empty" {
		der = Node{"ε", nil, "Empty"}
	}
	return der
}

// Функция производной по строке использует вспом. функцию lambda()
func (regex Node) derivative(str string) Node {
	children1 := make([]Node, 0)
	var der Node
	der.Value = ""
	if regex.Label == "Alt" {
		children1 = make([]Node, 0)
		for _, c := range regex.Children {
			childRegex := regexParse(c.Value)
			child := childRegex.derivative(str)
			children1 = append(children1, child)
		}
		children1 = simplifyAlt(children1)
		for i := range children1 {
			if children1[i].Value == "" {
				continue
			} else {
				der.Value += children1[i].Value + "|"
			}
		}
		der.Value = deleteExtraAlt(der.Value)
		if len(notNullNotEmpty(children1)) > 1 {
			der.Value = "(" + der.Value[:] + ")"
		}
		der = Node{der.Value, children1, "Alt"}
	} else if regex.Label == "Concat" {
		correctRegex := regexParse(regex.Value)
		correctRegex.Children = notNullNotEmpty(correctRegex.Children)
		if len(correctRegex.Children) == 1 {
			phi := correctRegex.Children[0]
			phi = regexParse(phi.Value)
			phiDer := phi.derivative(str)
			der = Node{phiDer.Value, phiDer.Children, phiDer.Label}
		} else if len(correctRegex.Children) != 0 {
			phi, psi := correctRegex.Children[0], simplifyConcat(correctRegex.Children[1:])
			phiDer := phi.derivative(str)
			var leftAlt Node
			if phiDer.Value == "ε" {
				leftAlt = psi
			} else if phiDer.Value == "" {
				leftAlt = Node{"", nil, "Null"}
			} else {
				children := make([]Node, 0)
				children = append(children, phiDer)
				children = append(children, psi)
				leftAlt = Node{phiDer.Value + psi.Value, children, "Concat"}
			}
			if leftAlt.Value != "" {
				children1 = append(children1, leftAlt)
			}
			var rightAlt Node
			phiLambda := phi.lambda()
			if phiLambda.Value != "" { // lambda().Value = "" or "ε" <- definition; "ε" is concatenated => rightAlt is psiDer; "" => rightAlt is ""
				psiRegex := regexParse(psi.Value)
				psiDer := psiRegex.derivative(str)
				children := make([]Node, 0)
				children = append(children, psiDer)
				rightAlt = Node{psiDer.Value, children, "Concat"}
				if rightAlt.Value != "" {
					if leftAlt.Value != "" {
						leftAlt.Value = "(" + leftAlt.Value + ")" + "|"
						children1[0] = leftAlt
					}
					children1 = append(children1, rightAlt)
				}
			}
			for _, c := range children1 {
				der.Value += c.Value
			}
			der = Node{der.Value, children1, "Alt"}
		}
	} else if regex.Label == "Star" {
		childRegex := regexParse(regex.Children[0].Value)
		children1 = append(children1, childRegex.derivative(str))
		for i := 1; i < len(regex.Children); i++ {
			children1 = append(children1, regex.Children[i])
		}
		if children1[0].Value != "" {
			if children1[0].Value == "ε" { //derivative(a*) = εa* => remove ε from concat
				der = Node{regex.Value, children1, "Star"}
			} else {
				der.Value = children1[0].Value + regex.Value
				der = Node{der.Value, children1, "Concat"}
			}
		} else {
			der = Node{"", children1, "Null"}
		}
	} else if regex.Label == "Sym" {
		if regex.Value == str {
			der = Node{"ε", nil, "Empty"}
		} else {
			der = Node{"", nil, "Null"}
		}
	} else if regex.Label == "Empty" {
		der = Node{"", nil, "Null"}
	}
	return der
}

func simplifyConcat(nodes []Node) Node {
	nodeVal := ""
	nodeChildren := make([]Node, 0)
	for _, n := range nodes {
		nodeVal += n.Value
		nodeChildren = append(nodeChildren, n)
	}
	return Node{nodeVal, nodeChildren, "Concat"}
}

func simplifyAlt(children []Node) []Node {
	unique := make(map[string]bool, len(children))
	nodeChildren := make([]Node, 0)
	for _, n := range children {
		if !unique[n.Value] {
			nodeChildren = append(nodeChildren, n)
			unique[n.Value] = true
		}
	}
	return nodeChildren
}

func automataStates(nodes []string) []string {
	states := make([]string, 0)
	for _, n := range nodes {
		reg := regexParse(n)
		if reg.Label == "Alt" {
			for _, c := range reg.Children {
				states = append(states, c.Value)
			}
		} else {
			states = append(states, n)
		}
	}
	return states
}

func reverse(nod string) string {
	node := regexParse(nod)
	rev := ""
	var revChild string
	revCh := make([]string, 0)
	if node.Label == "Alt" {
		ch := node.Children
		for _, c := range ch {
			if c.Label == "Alt" {
				revChild = "(" + reverse(c.Value) + ")"
			} else {
				revChild = reverse(c.Value)
			}
			revCh = append(revCh, revChild)
		}
		for j, s := range revCh {
			if j == len(revCh)-1 {
				rev += s
			} else {
				rev += s
				rev += "|"
			}
		}
	} else if node.Label == "Concat" {
		ch := node.Children
		for i := len(ch) - 1; i >= 0; i-- {
			if ch[i].Label == "Alt" {
				revChild = "(" + reverse(ch[i].Value) + ")"
			} else {
				revChild = reverse(ch[i].Value)
			}
			revCh = append(revCh, revChild)
		}
		for _, s := range revCh {
			rev += s
		}
	} else if node.Label == "Star" {
		ch := node.Children
		for _, c := range ch {
			revCh = append(revCh, reverse(c.Value))
		}
		rev = "("
		for _, s := range revCh {
			rev += s
		}
		rev += ")*"
	} else if node.Label == "Sym" {
		if node.Value != "" {
			rev = node.Value
		}
	} else if node.Label == "Empty" {
		rev = "ε"
	} else if node.Label == "Null" {
		rev = ""
	}
	return rev
}

func brzozowskiAutomata(start Node) BrzozowskiAuto {
	//var fa BrzozowskiAuto
	fa := BrzozowskiAuto{
		start.Value,
		getAlphabetForRegex(start.Value),
		[]string{start.Value},
		make([]Edge, 0),
	}
	//fmt.Println("alphabet of regex:", fa.Alphabet)
	//fmt.Println("parsed regex is:", start.Value)
	//fmt.Println("label of parsed regex:", start.Label)
	condition := func(a BrzozowskiAuto) (bool, BrzozowskiAuto) {
		//fmt.Print("\nstates of fa: ")
		//for _, c := range a.Nodes {
		//	fmt.Print(c, " ")
		//}
		oldStates := a.Nodes
		//fmt.Print("\nOLDSTATES: ")
		//for _, c := range oldStates {
		//	fmt.Print(c, " ")
		//}
		for _, w := range a.Alphabet {
			a.Nodes, a.Edges = addStateDerivative(a, w)
			a = BrzozowskiAuto{a.Regex, a.Alphabet, a.Nodes, a.Edges}
			//fmt.Print("\nHERE states: ")
			//fmt.Print("{")
			//for _, c := range a.Nodes {
			//	fmt.Print(c, ", ")
			//}
			//fmt.Println("}")
		}
		newStatesLength := len(a.Nodes)
		//fmt.Println("lengths:", len(oldStates), newStatesLength)
		return newStatesLength > len(oldStates), a
	}
	c, a := condition(fa)
	for c {
		//fmt.Println("ПОКА НЕ ПОСТРОИЛИ")
		fa = a
		c, a = condition(fa)
	}
	return a
}

func addStateDerivative(a BrzozowskiAuto, w string) ([]string, []Edge) {
	nodesCopy := a.Nodes
	for _, r := range nodesCopy {
		//fmt.Println("\nwhat state regex is derivated?:", r)
		n := regexParse(r)
		//fmt.Println(n)
		node := n.derivative(w)
		//fmt.Println(node)
		//fmt.Println("node:", node.Value)
		//fmt.Printf("==\nder is: %s for str '%s'", node.Value, w)
		if node.Value == "" {
			continue
		}
		node.Value = removeExtraParenthesis(node.Value)

		//fmt.Print("\nBEFORE DUPS: ")
		//for _, c := range a.Nodes {
		//	fmt.Print(c, " ")
		//}
		if !areNodesDuplicates(node.Value, a.Nodes) {
			//fmt.Println("dups")
			//fmt.Println("i am going to add node")
			a.Nodes = append(a.Nodes, node.Value)
			//fmt.Print("added: ")
			//for _, c := range a.Nodes {
			//	fmt.Print(c, ", ")
			//}
			a.Edges = append(a.Edges, Edge{r, node.Value, w})
			//fmt.Print("\nEDGES ARE: ")
			//for _, c := range a.Edges {
			//	fmt.Println("\t", c.Parent, "--", c.EdgeLabel, "->", c.Child)
			//}
		} else {
			if !areEdgesDuplicates(Edge{r, node.Value, w}, a.Edges) {
				a.Edges = append(a.Edges, Edge{r, node.Value, w})
			}
		}
	}
	return a.Nodes, a.Edges
}

func main() {
	var n int
	var regex string
	fmt.Println("Input test number:")
	fmt.Scan(&n)
	file, err := os.Open(fmt.Sprintf("tests/test%d.txt", n))
	if err != nil {
		log.Fatalf("Error with openning file: %s", err)
		os.Exit(1)
	}
	defer file.Close()
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		regex += strings.ReplaceAll(fileScanner.Text(), " ", "")
	}
	if err = fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}
	ok := file.Close()
	if ok != nil {
		log.Fatalf("Error with closing file: %s", ok)
	}
	fmt.Println("regex:", regex)
	start := regexParse(regex)
	printGraph(start)
	automata := brzozowskiAutomata(start)
	printAutomata(automata)
	res := make([]string, 0)
	for _, s := range automata.Nodes {
		s = reverse(s)
		rsi := regexParse(s)
		qij := automataStates(brzozowskiAutomata(rsi).Nodes)
		for _, q := range qij {
			res = append(res, reverse(q))
		}
	}
	unique := make(map[string]bool, len(res))
	result := make([]string, 0)
	for _, r := range res {
		if r == "" || r == "ε" {
			continue
		}
		if !unique[r] {
			result = append(result, r)
			unique[r] = true
		}
	}
	fmt.Println("Regex infixes:")
	sort.Sort(sort.StringSlice(result))
	for _, r := range result {
		fmt.Printf("\t%s\n", r)
	}
}

// здесь мог быть анекдот, но лучше будет котик
//          ／＞　 フ
//　　　　　| 　_　 _|
//　 　　　／`ミ _x 彡
//　　 　 /　　　 　 |
//　　　 /　 ヽ　　 ﾉ
//　／￣|　　 |　|　|
//　| (￣ヽ＿_ヽ_)_)
//　＼二つ
//
// он не грустный, он просто спит
