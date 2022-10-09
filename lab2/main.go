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
	Label    string //Alt or Concat or Star or Sym
}

func getFirstAltIndex(raw string) int {
	outBracketsCount := 0
	for i := range raw {
		if outBracketsCount == 0 && raw[i] == '|' {
			return i
		} else if raw[i] == '(' {
			outBracketsCount++
		} else if raw[i] == ')' {
			outBracketsCount--
		}
	}
	return -1
}

//func getAltIndexes(raw string) []int {
//	outBracketsCount := 0
//	altIndexes := make([]int, 0)
//	for i := range raw {
//		if outBracketsCount == 0 && raw[i] == '|' {
//			altIndexes = append(altIndexes, i)
//		} else if raw[i] == '(' {
//			outBracketsCount++
//		} else if raw[i] == ')' {
//			outBracketsCount--
//		}
//	}
//	return altIndexes
//}

//func getNumberOfSubstrings(raw string) int {
//	if getFirstAltIndex(raw) == -1 {
//		return 1
//	}
//	return len(getAltIndexes(raw)) + 1
//}

func getListOfAltSubstrings(raw string) []string {
	s := make([]string, 0)
	k := len(raw)
	i := 0
	for i <= k {
		//fmt.Println("raw is:", raw, "s is:", s, "k is:", k, "i:", i)
		alt := getFirstAltIndex(raw)
		//fmt.Println("alt", alt)
		if alt == -1 {
			s = append(s, raw)
			return s
		} else {
			s = append(s, string([]byte(raw)[:alt]))
			raw = string([]byte(raw)[alt+1:])
			k = len(raw)
			i = 0
		}
		i++
	}
	return append(s, "")
}

func getFirstConcIndex(raw string) (int, error) { //index of next after ')'
	outBracketsCount := 0
	for i := range raw {
		if i > 0 {
			if raw[i-1] == '(' {
				outBracketsCount++
				//fmt.Println("i after (:", i)
			} else if raw[i-1] == ')' {
				outBracketsCount--
				//fmt.Println("i after ):", i)
			}
			if outBracketsCount == 0 {
				if raw[i] == '(' || unicode.IsLetter(rune(raw[i])) {
					//fmt.Println("concatenation")
					return i, nil
				}
				if raw[i] == raw[i-1] && raw[i] == '*' {
					return -1, errors.New("error due to ** in regex")
				}
			}
		}
	}
	return -1, nil
}

func getListOfSubstrings(raw string) []string {
	s := make([]string, 0)
	//fmt.Println("initial s:", s)
	k := len(raw)
	i := 0
	for i <= k {
		//fmt.Println("raw is:", raw, "s is:", s, "k is:", k, "i:", i)
		concAfter, err := getFirstConcIndex(raw)
		if err != nil {
			fmt.Println(err)
			s = make([]string, 0)
			return s
		}
		//fmt.Println("concAfter", concAfter)
		if concAfter == -1 {
			s = append(s, raw)
			return s
		} else {
			s = append(s, string([]byte(raw)[:concAfter]))
			raw = string([]byte(raw)[concAfter:])
			k = len(raw)
			i = 0
		}
		i++
	}
	return s
}

func regexParse(regex string) Node {
	if regex == "" {
		err := errors.New("regex is an empty string")
		log.Fatalf("Error with regex: %v", err)
		os.Exit(1)
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
	fmt.Println("TEST ALT:", regex)
	n := len(regex)
	if regex[0] == '(' && regex[n-1] == ')' {
		regex = regex[1 : n-1]
		fmt.Println("TEST ALT: got rid of out brackets", regex)
	}
	if getFirstAltIndex(regex) == -1 {
		fmt.Println("TEST ALT: there is no out alternatives in regex", regex)
		return parseCon(regex)
	}
	fmt.Println("TEST ALT: get substrings of", regex)
	children := getListOfAltSubstrings(regex)
	childrenNodes := make([]Node, 0, len(regex))
	fmt.Printf("TEST ALT: substrings of '%s' are: ", regex)
	for i := range children {
		fmt.Printf("%s ", children[i])
	}
	fmt.Println()
	for i := range children {
		child := parseCon(children[i])
		childrenNodes = append(childrenNodes, child)
	}
	return Node{regex, childrenNodes, "Alt"}
}

func parseCon(regex string) Node {
	fmt.Println("TEST CONCAT:", regex)
	children := getListOfSubstrings(regex)
	if len(children) == 1 {
		fmt.Println("TEST CONCAT: there is no out concat in regex")
		return parseStar(regex)
	}
	childrenNodes := make([]Node, 0, len(regex))
	fmt.Println("TEST CONCAT: get substrings of", regex)
	fmt.Printf("TEST CONCAT: substrings of '%s' are: ", regex)
	for i := range children {
		fmt.Printf("%s ", children[i])
	}
	fmt.Println()
	for i := 0; i < len(children); i++ {
		child := parseStar(children[i])
		childrenNodes = append(childrenNodes, child)
		//fmt.Println("i`m here")
	}
	return Node{regex, childrenNodes, "Concat"}
}

func parseStar(regex string) Node {
	fmt.Println("TEST STAR:", regex)
	n := len(regex)
	if regex == "" {
		return Node{regex, nil, "Sym"}
	}
	if len(regex) == 1 && unicode.IsLetter(rune(regex[0])) {
		fmt.Println("TEST STAR: letter", regex)
		return Node{string(regex), nil, "Sym"}
	} else if regex[n-1] == '*' {
		if len(regex) == 2 && unicode.IsLetter(rune(regex[0])) {
			fmt.Println("TEST STAR: letter*", regex)
			child := Node{string(regex[0]), nil, "Sym"}
			return Node{regex, []Node{child}, "Star"}
		} else if regex[0] == '(' {
			fmt.Println("TEST STAR: construction (regex)*", regex)
			child := parseAlt(regex[:n-1])
			return Node{regex, []Node{child}, "Star"}
		}
	} else if regex[0] == '(' && regex[n-1] == ')' {
		fmt.Println("TEST STAR: construction (regex)", regex)
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
	fmt.Println("TEST:", regex)

	start := regexParse(regex)
	//fmt.Println("TEST:", start)
	printGraph(start)
	//testAlt := []string{
	//	"abc",
	//	"abc|",
	//	"a|b|c",
	//	"|a|b|c",
	//	"(ab)|cd",
	//	"(a||b)|cd|",
	//	"c|a",
	//	"a*|(ab)",
	//}
	//fmt.Println("test the alt parser")
	//for _, t := range testAlt {
	//	fmt.Println("=====")
	//	fmt.Println(t)
	//	s := getListOfAltSubstrings(t)
	//	for i, v := range s {
	//		if i == len(s)-1 {
	//			fmt.Printf("'%s'\n", v)
	//		} else {
	//			fmt.Printf("'%s', ", v)
	//		}
	//	}
	//}
	//testConc := []string{
	//	"abcdef",
	//	"abcdef*",
	//	"a*bcdef",
	//	//"*abcdef",
	//	"abc*def",
	//	"a(bc)def",
	//	"ab(cd)*ef*",
	//	//"ab**c",
	//	//"",
	//	"a*",
	//}
	//fmt.Println("test the conc parser")
	//for _, t := range testConc {
	//	fmt.Println("=====")
	//	fmt.Println(t)
	//	s := getListOfSubstrings(t)
	//	//fmt.Println(s)
	//	for i, v := range s {
	//		if i == len(s)-1 {
	//			fmt.Printf("'%s'\n", v)
	//		} else {
	//			fmt.Printf("'%s', ", v)
	//		}
	//	}
	//}
}
