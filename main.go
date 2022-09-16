package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// построчное чтение из файла
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
			tmp := strings.TrimSpace(strings.Split(fileScanner.Text(), "=")[1])
			nonTerms = strings.Split(tmp, ",")
			//fmt.Println(nonTerms[0:])
		} else if strings.Contains(fileScanner.Text(), "terminals") {
			tmp := strings.TrimSpace(strings.Split(fileScanner.Text(), "=")[1])
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

func makeListOfTermForms(nonTerms []string, rules map[string][]string) map[string][]string {
	forms := make(map[string][]string)
	//var curForm string
	nt := strings.ReplaceAll(strings.Join(nonTerms, ""), " ", "")
	//fmt.Println(nt)
	/*for _, n := range nonTerms {

		//for v := 0; v < len(val); v++ {
		fmt.Printf("%s- %s\n", key, v)
		curForm = strings.ReplaceAll(v, n, "")
		//fmt.Printf("term form: %s\n", curForm)
		forms[key] = append(forms[key], curForm)

	}*/
	for _, val := range rules {
		for _, v := range val {
			if strings.ContainsAny(v, nt) {

			}
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
	//termForms := makeListOfTermForms(nonTerms, rules)
	//fmt.Println(termForms)
}
