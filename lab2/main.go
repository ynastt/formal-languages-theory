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

type BrzozovskiAuto struct {
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
		//fmt.Println("str is:", str, "s is:", s, "k is:", k, "i:", i)
		alt := getFirstAltIndex(str)
		//fmt.Println("alt", alt)
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

func getFirstConcatIndex(str string) int { //index of next after ')'
	outBracketsCount := 0
	for i := range str {
		if i > 0 {
			if str[i-1] == '(' {
				outBracketsCount++
				//fmt.Println("i after (:", i)
			} else if str[i-1] == ')' {
				outBracketsCount--
				//fmt.Println("i after ):", i)
			}
			if outBracketsCount == 0 {
				if str[i] == '(' || unicode.IsLetter(rune(str[i])) {
					//fmt.Println("concatenation")
					return i
				}
			}
		}
	}
	return -1
}

func getListOfSubstrings(str string) []string {
	s := make([]string, 0)
	//fmt.Println("initial s:", s)
	k := len(str)
	i := 0
	for i <= k {
		//fmt.Println("str is:", str, "s is:", s, "k is:", k, "i:", i)
		concAfter := getFirstConcatIndex(str)
		//fmt.Println("concAfter", concAfter)
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
			//fmt.Println(": got rid of out brackets", regex)
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
	//fmt.Println("TEST ALT:", regex)
	regex = removeExtraParenthesis(regex)
	//fmt.Println(": got rid of out brackets", regex)
	if getFirstAltIndex(regex) == -1 {
		//fmt.Println("TEST ALT: there is no out alternatives in regex", regex)
		return parseCon(regex)
	}
	children := getListOfAltSubstrings(regex)
	childrenNodes := make([]Node, 0, len(regex))
	//fmt.Printf("TEST ALT: substrings of '%s' are: ", regex)
	//for i := range children {
	//	fmt.Printf("%s ", children[i])
	//}
	//fmt.Println()
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
	//fmt.Println("TEST CONCAT:", regex)
	children := getListOfSubstrings(regex)
	if len(children) == 1 {
		//fmt.Println("TEST CONCAT: there is no out concat in regex")
		return parseStar(regex)
	}

	childrenNodes := make([]Node, 0, len(regex))
	//fmt.Println("TEST CONCAT: get substrings of", regex)
	//fmt.Printf("TEST CONCAT: substrings of '%s' are: ", regex)
	//for i := range children {
	//	fmt.Printf("%s ", children[i])
	//}
	//fmt.Println()
	for i := 0; i < len(children); i++ {
		child := parseStar(children[i])
		childrenNodes = append(childrenNodes, child)
		//fmt.Println("i`m here")
	}
	return Node{regex, childrenNodes, "Concat"}
}

func parseStar(regex string) Node {
	//fmt.Println("TEST STAR:", regex)
	n := len(regex)
	if regex == "" || regex == "ε" {
		return Node{"ε", nil, "Empty"}
	}
	if len(regex) == 1 && unicode.IsLetter(rune(regex[0])) {
		//fmt.Println("TEST STAR: letter", regex)
		return Node{string(regex), nil, "Sym"}
	} else if regex[n-1] == '*' {
		if len(regex) == 2 && unicode.IsLetter(rune(regex[0])) {
			//fmt.Println("TEST STAR: letter*", regex)
			child := Node{string(regex[0]), nil, "Sym"}
			return Node{regex, []Node{child}, "Star"}
		} else if regex[0] == '(' {
			//fmt.Println("TEST STAR: construction (regex)*", regex)
			child := parseAlt(regex[:n-1])
			return Node{regex, []Node{child}, "Star"}
		}
	} else if regex[0] == '(' && regex[n-1] == ')' {
		//fmt.Println("TEST STAR: construction (regex)", regex)
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

func printAutomata(automata BrzozovskiAuto) {
	fmt.Println("\nBrzozovski automat {")
	for _, e := range automata.Edges {
		fmt.Printf("\t%s  --- %s ---> %s\n", e.Parent, e.EdgeLabel, e.Child)
	}
	fmt.Printf("}\n")
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
	//fmt.Println("\nnow states of auto are: ")
	//for _, c := range arr {
	//	fmt.Print(c, ",")
	//}
	//fmt.Println()
	//fmt.Println(el)
	for _, a := range arr {
		//fmt.Println("a:", a)
		if a == el {
			//fmt.Println("SAME a:", a, el)
			return true
		}
	}
	return false
}

func areEdgesDuplicates(el Edge, arr []Edge) bool {
	//fmt.Println("\nnow edges of auto are: ")
	//for _, c := range arr {
	//	fmt.Printf("\t%s  --- %s ---> %s\n", c.Parent, c.EdgeLabel, c.Child)
	//}
	//fmt.Println()
	//fmt.Printf("el edge is: %s  --- %s ---> %s\n", el.Parent, el.EdgeLabel, el.Child)
	for _, a := range arr {
		//fmt.Printf("a edge is: %s  --- %s ---> %s\n", a.Parent, a.EdgeLabel, a.Child)
		if a == el {
			//fmt.Println("SAME")
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
			//fmt.Println("DER OF ALT CHILD IS:", child)
			children1 = append(children1, child)
		}
		children1 = simplifyAlt(children1)
		for i := range children1 {
			if children1[i].Value == "" /*|| children1[i].Value == "ε"*/ {
				continue
			} else {
				der.Value += children1[i].Value + "|"
			}
		}
		der.Value = deleteExtraAlt(der.Value)
		if len(notNullNotEmpty(children1)) > 1 {
			der.Value = "(" + der.Value[:] + ")"
		}
		//fmt.Println("Alt derivative:", der.Value)
		der = Node{der.Value, children1, "Alt"}
	} else if regex.Label == "Concat" {
		//fmt.Println("CONC input regex is; need to be checked:", regex)
		correctRegex := regexParse(regex.Value)
		//fmt.Println("correct regex:", correctRegex.Value)
		//fmt.Println("childred of correwct regex:", correctRegex.Children)
		//fmt.Println("label of correct regex should be same:", correctRegex.Label)
		//fmt.Println(regex.Children)
		correctRegex.Children = notNullNotEmpty(correctRegex.Children)
		//fmt.Print("CONC not empty children: ")
		//fmt.Println(correctRegex.Children)
		if len(correctRegex.Children) == 1 {
			//fmt.Println("SINGLE child")
			phi := correctRegex.Children[0]
			phi = regexParse(phi.Value)
			//fmt.Println("Single child is:", phi)
			phiDer := phi.derivative(str)
			//fmt.Println("der of child is :", phiDer)
			der = Node{phiDer.Value, phiDer.Children, phiDer.Label}
		} else if len(correctRegex.Children) != 0 {
			//fmt.Print("Concat of >1 children; Children are :")
			//for _, c := range correctRegex.Children {
			//	fmt.Print(c.Value, ", ")
			//}
			//fmt.Println()
			phi, psi := correctRegex.Children[0], simplifyConcat(correctRegex.Children[1:])
			//fmt.Println("Left and Right parts of concatination:", phi.Value, psi.Value)
			//fmt.Println("FIND DERIVATIVE OF left part of concat - der phi + psi")
			//phiRegex := regexParse(phi.Value)
			phiDer := phi.derivative(str)
			//fmt.Println("Phi derivative:", phiDer.Value)
			//leftAlt := make([]Node, 0)
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
			//fmt.Println("LEFT alt of concat derivative is:", leftAlt.Value)
			if leftAlt.Value != "" {
				children1 = append(children1, leftAlt)
				//der.Value += leftAlt.Value
			}
			var rightAlt Node
			//fmt.Println("FIND DERIVATIVE OF right part of concat - lambda phi + der psi")
			phiLambda := phi.lambda()
			//fmt.Println("Lambda of Phi", phiLambda)
			if phiLambda.Value != "" { // lambda().Value = "" or "ε" <- definition; "ε" is concatenated => rightAlt is psiDer; "" => rightAlt is ""
				//fmt.Println("HERE WHAT")
				psiRegex := regexParse(psi.Value)
				psiDer := psiRegex.derivative(str)
				children := make([]Node, 0)
				//children = append(children, phiLambda)
				children = append(children, psiDer)
				rightAlt = Node{psiDer.Value, children, "Concat"}
				//rightAlt = Node{"(" + phiLambda.Value + psiDer.Value + ")", children, "Concat"}
				//fmt.Println("RIGHT alt of concat derivative is:", rightAlt.Value)
				if rightAlt.Value != "" {
					if leftAlt.Value != "" {
						leftAlt.Value = "(" + leftAlt.Value + ")" + "|"
						children1[0] = leftAlt
						//der.Value += leftAlt.Value
						//der.Value += "|"
					}
					children1 = append(children1, rightAlt)
					//der.Value += rightAlt.Value
				}
			}
			for _, c := range children1 {
				der.Value += c.Value
			}
			//fmt.Println("RIGHT alt of concat derivative is:", rightAlt.Value)
			der = Node{der.Value, children1, "Alt"}
		}
	} else if regex.Label == "Star" {
		//fmt.Println("STAR")
		childRegex := regexParse(regex.Children[0].Value)
		children1 = append(children1, childRegex.derivative(str))
		for i := 1; i < len(regex.Children); i++ {
			children1 = append(children1, regex.Children[i])
		}
		//fmt.Println("children of Star node:")
		//for _, c := range children1 {
		//	fmt.Print(c.Value, ", ")
		//}
		//fmt.Println()
		if children1[0].Value != "" {
			if children1[0].Value == "ε" { //derivative(a*) = εa* => remove ε from concat
				//fmt.Println("производная = конкатенация с ε, значит значение это сама ()*")
				//fmt.Println("производаня от ()*:", regex.Value)
				der = Node{regex.Value, children1, "Star"}
			} else {
				der.Value = children1[0].Value + regex.Value
				//fmt.Println("производная не пустая, значение = конкатенация производной и ()*: ", der.Value)
				der = Node{der.Value, children1, "Concat"}
			}
		} else {
			//fmt.Println("производная от ()* дала пустое множество")
			der = Node{"", children1, "Null"}
		}
	} else if regex.Label == "Sym" {
		if regex.Value == str {
			//fmt.Println("производная по s от s = ε")
			der = Node{"ε", nil, "Empty"}
		} else {
			//fmt.Println("производная по s от НЕ s = '' ")
			der = Node{"", nil, "Null"}
		}
	} else if regex.Label == "Empty" {
		//fmt.Println("производная по s от '' = '' ")
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
	//fmt.Println("TEST children:", start.Children)
	//fmt.Println("value:", start.Value)
	printGraph(start)
	//fmt.Println("TEST alphabet derivatives for regex")
	//alp := getAlphabetForRegex(start.Value)
	//fmt.Println("alphabet of regex:", alp)
	//fmt.Println("parsed regex is:", start.Value)
	//fmt.Println("label of parsed regex:", start.Label)
	//for _, s := range alp {
	//	fmt.Println("derivative with respect to", s)
	//	res := start.derivative(s)
	//	if res.Value == "" {
	//		fmt.Printf("derivative(%s)= %s\n", start.Value, "''")
	//	} else {
	//		fmt.Printf("derivative(%s)= %s\n", start.Value, res.Value)
	//	}
	//}
	//fmt.Println()
	automata := brzozowskiAutomata(start)
	printAutomata(automata)
	//fmt.Println("States of R auto:", automataStates(automata.Nodes))
	res := make([]string, 0)
	for _, s := range automata.Nodes {
		s = reverse(s)
		rsi := regexParse(s)
		qij := automataStates(brzozowskiAutomata(rsi).Nodes)
		//fmt.Println("States of Sr auto:", qij)
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

	//fmt.Println("\n\n\ntest derivative func")
	//str := "ba(aba)*|(baa)*ba"
	//str = "(baa)*ba"
	//str = "(ab)*"
	//str = "a(aba)*|(aa(baa)*ba)|a"
	//str = "((aba)*|a(baa)*ba|ε)"
	//str = "(ab)*b(a(a(ab)*)b)*"  //der by b
	//str = "b(ab)*b(a(a(ab)*)b)*" //der by b
	//str = "(bb|ab|(aaa)*)*"
	//str = "(b|aa(aaa)*)(bb|ab|(aaa)*)*"
	//str = "a(aaa)*(bb|ab|(aaa)*)*"                                //der by a
	//str = "(aaa)*(bb|ab|(aaa)*)*"                                 //der by a
	//str = "(aa(aaa)*(bb|ab|(aaa)*)*)|(b|aa(aaa)*)(bb|ab|(aaa)*)*" //der by b
	//str = "b(bb|ab|(aaa)*)*"
	//node := regexParse(str)
	//fmt.Println("node:", node)
	//fmt.Println("\tStart Derivative process")
	//der := node.derivative("a")
	//fmt.Println(der)
	//fmt.Println(der.Value)
	//fmt.Println()
	//der = node.derivative("b")
	//fmt.Println(der)
	//fmt.Println(der.Value)

	//fmt.Println("\n\ntest reverse func")
	//str := "(baa)*ba"
	//str = "a(aba)*|(aa(baa)*ba)|a"
	//str = "(aa(aaa)*(bb|ab|(aaa)*)*)|(b|aa(aaa)*)(bb|ab|(aaa)*)*"
	//str = "b(bb|ab|(aaa)*)*"
	//str = "a|b|b|ab|b*|(((ab(a(|ab*|)b))*)*)*|c"
	//str = "((|a(b)*c|b)*|(ba)*)*"
	//fmt.Println("str is:", str)
	//r := reverse(str)
	//fmt.Println("reversed str is:", r)
	//fmt.Println("list of subs of str:", getListOfAltSubstrings(str), len(getListOfAltSubstrings(str)))
	//fmt.Println("list of subs of rev:", getListOfAltSubstrings(r), len(getListOfAltSubstrings(r)))
	//parsed := regexParse(r)
	//fmt.Println("parsed:", parsed.Value)
	//fmt.Println("children of parsed:")
	//for _, c := range parsed.Children {
	//	fmt.Print(c.Value, ", ")
	//}
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
		//fmt.Println("\talt")
		ch := node.Children
		//fmt.Println("alt ch:", ch)
		for _, c := range ch {
			//fmt.Println("child:", c.Value)
			if c.Label == "Alt" {
				revChild = "(" + reverse(c.Value) + ")"
				//fmt.Println("revChild:", revChild)
			} else {
				revChild = reverse(c.Value)
				//fmt.Println("revChild:", revChild)
			}
			revCh = append(revCh, revChild)
		}
		//fmt.Println("reversed children:", revCh)
		for j, s := range revCh {
			if j == len(revCh)-1 {
				rev += s
			} else {
				rev += s
				rev += "|"
			}
		}
		//fmt.Println("reversed alt:", rev)
	} else if node.Label == "Concat" {
		//fmt.Println("\tconcat")
		ch := node.Children
		//fmt.Println("concat ch:", ch)
		for i := len(ch) - 1; i >= 0; i-- {
			//fmt.Println("child:", ch[i].Value)
			if ch[i].Label == "Alt" {
				revChild = "(" + reverse(ch[i].Value) + ")"
				//fmt.Println("revChild:", revChild)
			} else {
				revChild = reverse(ch[i].Value)
				//fmt.Println("revChild:", revChild)
			}
			revCh = append(revCh, revChild)
		}
		//fmt.Println("reversed children:", revCh)
		for _, s := range revCh {
			rev += s
		}
		//fmt.Println("reversed concat:", rev)
	} else if node.Label == "Star" {
		//fmt.Println("\tstar")
		ch := node.Children
		//fmt.Println("star ch:", ch)
		for _, c := range ch {
			//fmt.Println("child:", c.Value)
			revCh = append(revCh, reverse(c.Value))
		}
		rev = "("
		for _, s := range revCh {
			rev += s
		}
		rev += ")*"
		//fmt.Println("reversed star:", rev)
	} else if node.Label == "Sym" {
		if node.Value != "" {
			rev = node.Value
		}
		//fmt.Println("reversed sym:", rev)
	} else if node.Label == "Empty" {
		rev = "ε"
		//fmt.Println("reversed empty:", rev)
	} else if node.Label == "Null" {
		rev = ""
		//fmt.Println("reversed '':", rev)
	}
	return rev
}

func brzozowskiAutomata(start Node) BrzozovskiAuto {
	//var fa BrzozovskiAuto
	fa := BrzozovskiAuto{
		start.Value,
		getAlphabetForRegex(start.Value),
		[]string{start.Value},
		make([]Edge, 0),
	}
	//fmt.Println("alphabet of regex:", fa.Alphabet)
	//fmt.Println("parsed regex is:", start.Value)
	//fmt.Println("label of parsed regex:", start.Label)
	condition := func(a BrzozovskiAuto) (bool, BrzozovskiAuto) {
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
			a = BrzozovskiAuto{a.Regex, a.Alphabet, a.Nodes, a.Edges}
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

func addStateDerivative(a BrzozovskiAuto, w string) ([]string, []Edge) {
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

/* то, как я сама себе все усложнила это уже своего рода анекдот */
