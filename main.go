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
	file, err := os.Open("tests/test6.txt")
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
			//fmt.Printf("%s - %s\n", tmp[0], rules[tmp[0]])
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

// заполнение карты: нетерминал -> список терминальных форм (например, S-нетерминал, a-терминал, S->aSa => S->a_a)
// добавление непорождающих нетерминалов в отдельный класс эквивалентности.
// изначально созданы 2 класса эквивалентности: для порождающих и для непорождающих нетерминалов
// далее при необходимости создаются новые классы эквивалентности.
func makeListOfTermForms(nonTerms []string, rules map[string][]string) (map[string][]string, []string) {
	var eqNotGenClass []string
	nt := strings.Join(nonTerms, "")
	forms := make(map[string][]string)
	for _, n := range nonTerms {
		//fmt.Println(n)
		val, ok := rules[n]
		if ok {
			//fmt.Println(val)
			sort.Slice(val, func(i, j int) bool {
				return len(val[i]) < len(val[j])
			})
			//fmt.Println(val)
			for _, v := range val {
				curForm := ""
				for _, sym := range v {
					if strings.Contains(nt, string(sym)) {
						curForm += "_"
					} else {
						curForm += string(sym)
					}
				}
				forms[n] = append(forms[n], curForm)
			}
			//forms[n] = append(forms[n], curForm)
		} else {
			eqNotGenClass = append(eqNotGenClass, n)
			//fmt.Println("this nonterminal is not generating!")
		}
	}
	return forms, eqNotGenClass
}

// гипотеза разделения на классы эквивалентности на основе сравнения списка терминальных форм
func eqClassesDivision(termForms map[string][]string) map[string][]string {
	eqGenClasses := make(map[string][]string)
	var single []string
	for nonTerm, _ := range termForms {
		eqGenClasses[nonTerm] = single
		//fmt.Println()
		//fmt.Println(nonTerm)
		//fmt.Printf("joined nonterm: %s", strings.Join(termForms[nonTerm], ""))
		for key, _ := range termForms {
			if nonTerm != key {
				//fmt.Printf("\n %s", key)
				//fmt.Printf("\njoined eifq nonterm: %s,", strings.Join(termForms[key], ""))
				if strings.Compare(strings.Join(termForms[nonTerm], ""),
					strings.Join(termForms[key], "")) == 0 {
					eqGenClasses[nonTerm] = append(eqGenClasses[nonTerm], key)
				}
			}
		}
	}
	return eqGenClasses
}

func main() {
	terms, nonTerms, rules := parseTerms()
	fmt.Println("nonterms:", nonTerms[0:])
	fmt.Println("terms:", terms[0:])
	fmt.Println("rules:", rules)
	termForms, eqNotGenClass := makeListOfTermForms(nonTerms, rules)
	fmt.Println("term forms:", termForms)
	fmt.Println("notGenEqClass:", eqNotGenClass)
	fmt.Println("=====test====")
	eqGenClass := eqClassesDivision(termForms)
	fmt.Println("first eq classes:", eqGenClass)
}
