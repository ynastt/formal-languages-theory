package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

// построчное чтение входных данных из файла tests/test*.txt, * - номер теста
// заполнение массивов нетерминалов, терминалов и карты: нетерминал -> спсиок правил переписыванния
func parseTerms() ([]string, []string, map[string][]string) {
	var nonTerms, terms []string
	rules := make(map[string][]string)
	file, err := os.Open("tests/test3.txt")
	if err != nil {
		log.Fatalf("Error with openning file: %s", err)
	}
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		if strings.Contains(fileScanner.Text(), "nonterminals") {
			tmp := strings.ReplaceAll(strings.Split(fileScanner.Text(), "=")[1], " ", "")
			nonTerms = strings.Split(tmp, ",")
			//fmt.Println(nonTerms[0:])
		} else if strings.Contains(fileScanner.Text(), "terminals") {
			tmp := strings.ReplaceAll(strings.Split(fileScanner.Text(), "=")[1], " ", "")
			terms = strings.Split(tmp, ",")
			//fmt.Println(terms[0:])
		} else if strings.Contains(fileScanner.Text(), "->") {
			tmp := strings.Split(fileScanner.Text(), " -> ")
			rules[tmp[0]] = append(rules[tmp[0]], strings.ReplaceAll(tmp[1], " ", ""))
			fmt.Printf("%s - %s\n", tmp[0], rules[tmp[0]])
		}
	}
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}
	ok := file.Close()
	if ok != nil {
		log.Fatalf("Error with closing file: %s", ok)
	}
	return terms, nonTerms, rules
}

// заполнение карты: нетерминал -> список терминальных форм (например, aSa преобразуется в a_a)
// добавить сразу добавление непорождающих нетерминалов в отдельные классы эквивалентности
func makeListOfTermForms(nonTerms []string, rules map[string][]string) map[string][]string {
	nt := strings.Join(nonTerms, "")
	forms := make(map[string][]string)
	for _, n := range nonTerms {
		//fmt.Println(n)
		curForm := ""
		val, ok := rules[n]
		if ok {
			//fmt.Println(val)
			sort.Slice(val, func(i, j int) bool {
				return len(val[i]) < len(val[j])
			})
			//fmt.Println(val)
			for _, v := range val {
				for _, sym := range v {
					if strings.Contains(nt, string(sym)) {
						curForm += "_"
					} else {
						curForm += string(sym)
					}
				}
				curForm += " "
			}
			forms[n] = append(forms[n], curForm)
		} else {
			// add to special eq class
			fmt.Println("this nonterminal is not generating!")
		}
	}
	return forms
}

//func EqClassesDivision() {
//}

func main() {
	terms, nonTerms, rules := parseTerms()
	fmt.Println("nonterms:")
	fmt.Println(nonTerms[0:])
	fmt.Println("terms:")
	fmt.Println(terms[0:])
	fmt.Println(rules)
	fmt.Println("==========terms forms===========")
	termForms := makeListOfTermForms(nonTerms, rules)
	fmt.Println(termForms)
}
