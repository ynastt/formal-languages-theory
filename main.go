package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

// удаление дупликатов среди терминальных форм для каждого нетерминала
func checkForDuplicatesForms(slice []string) []string {
	j := len(slice)
	for i := 1; i < len(slice); i++ {
		if slice[i] == slice[i-1] {
			j = 1
			//fmt.Println("dups:", slice[i], slice[i-1])
			slice[j] = slice[i]
			j++
		}
	}
	return slice[:j]
}

// построчное чтение входных данных из файла tests/test*.txt, * - номер теста
// заполнение массивов нетерминалов, терминалов и карты: нетерминал -> спсиок правил переписыванния
func parseTerms() ([]string, []string, map[string][]string) {
	var nonTerms, terms []string
	rules := make(map[string][]string)
	file, err := os.Open("tests/test2.txt")
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
	for _, r := range rules {
		sort.Slice(r, func(i, j int) bool {
			return len(r[i]) < len(r[j])
		})
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
			//sort.Slice(val, func(i, j int) bool {
			//	return len(val[i]) < len(val[j])
			//})
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
		} else {
			eqNotGenClass = append(eqNotGenClass, n)
			//fmt.Println("this nonterminal is not generating!")
		}
	}
	for k, f := range forms {
		f = checkForDuplicatesForms(f)
		forms[k] = f
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

// функция проверки гипотезы разделения на классы эквивалентности
// новое разбиение на классы эквивалентности в случае, когда в правилах
// нетерминалы не попадают в один класс эквивалентности
/*func checkEqClassDivision(firstEqClasses, rules map[string][]string, eqNotGenClass []string) map[string][]string {
	newEqGenClasses := make(map[string][]string)
	fmt.Println(len(firstEqClasses["Q"]))
	for nonTerm, eqs := range firstEqClasses {
		if len(eqs) != 0 {
			for _, e := range eqs {

			}
		}
	}
	return newEqGenClasses
}*/

func main() {
	terms, nonTerms, rules := parseTerms()
	fmt.Println("nonterms:", nonTerms[0:])
	fmt.Println("terms:", terms[0:])
	fmt.Println("rules:", rules)
	termForms, eqNotGenClass := makeListOfTermForms(nonTerms, rules)
	fmt.Println("term forms:", termForms)
	fmt.Println("notGenEqClass:", eqNotGenClass)
	eqGenClass := eqClassesDivision(termForms)
	fmt.Println("first eq classes:", eqGenClass)
	fmt.Println("=====test====")
	//eqGenClass = checkEqClassDivision(eqGenClass, rules, eqNotGenClass)
	//fmt.Println("new eq classes:", eqGenClass)
}
