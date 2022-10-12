package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

type Node struct {
	Value    string
	Children []Node
	Label    string
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
		err := errors.New("regex has extra brackets")
		log.Fatalf("Error with regex: %v", err)
		os.Exit(1)
	}
	return parseAlt(regex)
}

func parseAlt(regex string) Node {
	//fmt.Println("TEST ALT:", regex)
	n := len(regex)
	closingParenthesis := func(str string) int {
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
	if regex[0] == '(' && regex[n-1] == ')' && closingParenthesis(regex) == n-1 {
		regex = regex[1 : n-1]
		//fmt.Println("TEST ALT: got rid of out brackets", regex)
	}
	if getFirstAltIndex(regex) == -1 {
		//fmt.Println("TEST ALT: there is no out alternatives in regex", regex)
		return parseCon(regex)
	}
	//fmt.Println("TEST ALT: get substrings of", regex)
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
	return Node{regex, childrenNodes, "Alt"}
}

func parseCon(regex string) Node {
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
	if regex == "" {
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
	return Node{regex, nil, "StarX"}
}

func printGraphNodes(start Node) {
	for i := range start.Children {
		fmt.Printf("\t%s [label = %s] -> ", start.Value, start.Label)
		fmt.Printf("%v [label = %s]\n", start.Children[i].Value, start.Children[i].Label)
		//fmt.Printf("%v\n", start.Children[i])
	}
	for i := range start.Children {
		printGraphNodes(start.Children[i])
	}
}

func printGraph(start Node) {
	fmt.Println("graph {")
	printGraphNodes(start)
	fmt.Println("}")
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
	//fmt.Println("TEST:", regex)
	//
	//start := regexParse(regex)
	////fmt.Println("TEST:", start)
	//printGraph(start)
	//brzozovskiAutomat(start)
	fmt.Println("=====tests for lambda()=====")

	tests := []string{
		"a",
		"", //это 'ε' в понимании парсера
		"a*",
		"ab",
		"abc",
		"a|b",
		"a||d",
		"a|b|c",
	}
	for _, t := range tests {
		if t == "" {
			fmt.Println("''")
		} else {
			fmt.Println(t)
		}
		r := regexParse(t)
		fmt.Println("parsed regex is:", r.Value)
		res := r.lambda()
		if res.Value == "" {
			fmt.Printf("lambda(%s)= %s\n", r.Value, "''")
		} else {
			fmt.Printf("lambda(%s)= %s\n", r.Value, r.lambda().Value)
		}
	}
}

/* Определим вспомогательную функцию lambda , такую, что
если аргумент принимает ε, то она возвращает ε, иначе – ""
(где "" - пустое множество, т.е. противоречие) */
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
			der.Value += children1[i].Value
		}
		der = Node{der.Value, children1, "Concat"}
	} else if regex.Label == "Star" {
		//regex.Value = "ε"
		der = Node{"ε", nil, "Empty"}
	} else if regex.Label == "Sym" {
		der = Node{"", nil, "Sym"}
	} else if regex.Label == "Empty" {
		der = Node{"ε", nil, "Empty"}
	}
	return der
}

/* использует вспом. функцию lambda()*/
func (regex Node) getDerivative(str string) Node {
	children1 := make([]Node, 0)
	var der Node
	if regex.Label == "Alt" {
		children1 = make([]Node, 0)
		for i := range regex.Children {
			children1 = append(children1, regex.Children[i].getDerivative(str))
		}
		for i := range children1 {
			der.Value += children1[i].Value + "|"
		}
		der = Node{der.Value, children1, "Alt"}
	} else if regex.Label == "Concat" {
		r, s := regex.Children[0], regex.Children[1]
		var c1, c2 Node
		c1 = Node{"conc", []Node{r.lambda(), s.getDerivative(str)}, "Concat"}
		c2 = Node{"conc", []Node{r.getDerivative(str), s}, "Concat"}
		regex = Node{regex.Value, []Node{c1, c2}, "Alt"}
	} else if regex.Label == "Star" {
		children1 = append(children1, regex.Children[0].getDerivative(str))
		for i := 1; i < len(regex.Children); i++ {
			children1 = append(children1, regex.Children[i])
		}
		der.Value = string(children1[0].Value) + regex.Value
		der = Node{der.Value, children1, "Concat"}
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

func deleteExtraAlt(str string) string {
	n := len(str)
	if n > 0 {
		if string(str[n-1]) == "|" {
			str = str[:n-1]
		}
	}
	return str
}
