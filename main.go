package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func addUnderscore(v, nt string) (string, string) {
	s, curForm := "", ""
	for _, sym := range v {
		if strings.Contains(nt, string(sym)) {
			s = string(sym)
			curForm += "_"
		} else {
			curForm += string(sym)
		}
	}
	return curForm, s
}

// удаление дупликатов среди терминальных форм для каждого нетерминала
func checkForDuplicates(slice []string) []string {
	j := len(slice)
	for i := 1; i < j; i++ {
		if slice[i] == slice[i-1] {
			//fmt.Println("duplicates:", slice[i], slice[i-1])
			slice = append(slice[0:i-1], slice[i:]...)
			i--
			j--
		}
	}
	return slice
}

func findClass(nonTerm string, firstEqClasses []string) string {
	var class string
	for _, classNonTerm := range firstEqClasses {
		if strings.ContainsAny(classNonTerm, nonTerm) {
			//fmt.Println("nonTerm has class: ", classNonTerm)
			class = classNonTerm
		}
	}
	return class
}

func removeClass(firstEqClasses []string, class1 string, nonTerm string) []string {
	for ind, class := range firstEqClasses {
		if class == class1 {
			class = strings.ReplaceAll(class, nonTerm, "")
			firstEqClasses[ind] = class
			firstEqClasses = append(firstEqClasses, nonTerm)
			sort.Strings(firstEqClasses)
			fmt.Println("classes:", firstEqClasses)
			fmt.Printf("%s was removed from class %s\n", nonTerm, class)
		}
	}
	return firstEqClasses
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
func makeListOfTermForms(nonTerms []string, rules map[string][]string, nt string) (map[string][]string, []string) {
	var eqNotGenClass []string
	forms := make(map[string][]string)
	for _, n := range nonTerms {
		//fmt.Println(n)
		val, ok := rules[n]
		if ok {
			for _, v := range val {
				curForm, _ := addUnderscore(v, nt)
				forms[n] = append(forms[n], curForm)
			}
		} else {
			eqNotGenClass = append(eqNotGenClass, n)
		}
	}
	for k, f := range forms {
		f = checkForDuplicates(f)
		forms[k] = f
	}
	return forms, eqNotGenClass
}

// гипотеза разделения на классы эквивалентности на основе сравнения списка терминальных форм
func eqClassesDivision(termForms map[string][]string) []string {
	var eqGenClasses []string
	var arr []byte
	for nonTerm, _ := range termForms {
		str := nonTerm
		for key, _ := range termForms {
			if nonTerm != key {
				if strings.Join(termForms[nonTerm], "") == strings.Join(termForms[key], "") {
					str += key
				}
				arr = []byte(str)
				sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
			}
		}
		eqGenClasses = append(eqGenClasses, string(arr))
	}

	sort.Slice(eqGenClasses, func(i, j int) bool {
		return len(eqGenClasses[i]) < len(eqGenClasses[j])
	})
	eqGenClasses = checkForDuplicates(eqGenClasses)
	return eqGenClasses
}

// функция проверки гипотезы разделения на классы эквивалентности
// новое разбиение на классы эквивалентности в случае, когда в правилах
// нетерминалы не попадают в один класс эквивалентности
func checkEqClassDivision(firstEqClasses []string, rules, termForms map[string][]string, notGenEqClass, t, nt string) (bool, []string) {
	flag := false
	var nonTerms []string
	for k := range termForms {
		nonTerms = append(nonTerms, k)
	}
	sort.Strings(nonTerms)
	for _, nonTerm := range nonTerms {
		fmt.Println("nonterm now:", nonTerm)
		for _, rule1 := range rules[nonTerm] {
			fmt.Printf("nonterm rule %s->%s\n", nonTerm, rule1)
			if len(rule1) == 1 && strings.ContainsAny(t, rule1) {
				//fmt.Println("-> term", rule1)
				continue
			}
			f1, c1 := addUnderscore(rule1, nt)
			for _, key := range nonTerms {
				if nonTerm != key {
					fmt.Printf("\tnt is: %s\n", key)
					for _, rule2 := range rules[key] {
						fmt.Printf("\t\tnt rule %s->%s\n", key, rule2)
						if len(rule2) == 1 && strings.ContainsAny(t, rule2) {
							//fmt.Println("-> term", rule2)
							continue
						}
						f2, c2 := addUnderscore(rule2, nt)
						fmt.Println("HERE F1 F2, c1, c2", f1, f2, c1, c2)
						fmt.Printf("%s->%s %s->%s\n", nonTerm, rule1, key, rule2)
						class1 := findClass(nonTerm, firstEqClasses)
						class2 := findClass(key, firstEqClasses)
						fmt.Printf("class nonterm: %s, class nt: %s\n", class1, class2)
						if class1 == class2 && f1 == f2 {
							fmt.Println("CLASSES and termforms ARE SAME")
							class3 := findClass(c1, firstEqClasses)
							class4 := findClass(c2, firstEqClasses)
							fmt.Printf("class nonterm _ : %s, class nt _ : %s\n", class3, class4)
							if class3 != class4 {
								fmt.Println("несовпали классы экв у Ni' Nj'")
								if class3 != class1 {
									firstEqClasses = removeClass(firstEqClasses, class1, nonTerm)
									flag = true
									fmt.Println("now classes are:", firstEqClasses)
								} else if class4 != class1 {
									firstEqClasses = removeClass(firstEqClasses, class1, key)
									flag = true
									fmt.Println("now classes are:", firstEqClasses)
								}
							}
						}
					}
				}
			}
		}
	}
	return flag, firstEqClasses
}

/*func outputNewRules() {

}*/

func main() {
	terms, nonTerms, rules := parseTerms()
	nt := strings.Join(nonTerms, "")
	fmt.Println("nonterms:", nonTerms[0:])
	fmt.Println("terms:", terms[0:])
	fmt.Println("rules:", rules)
	termForms, eqNotGenClass := makeListOfTermForms(nonTerms, rules, nt)
	fmt.Println("term forms:", termForms)
	fmt.Println("notGenEqClass:", eqNotGenClass)
	eqGenClass := eqClassesDivision(termForms)
	fmt.Println("first eq classes:", eqGenClass)

	fmt.Println("====================test====================")
	var flag bool
	//fmt.Println(len(eqGenClass))
	//fmt.Println(len(termForms))
	//fmt.Println(len(eqNotGenClass))
	t := strings.Join(terms, "")
	notGenEqClass := strings.Join(eqNotGenClass, "")
	flag, eqGenClass = checkEqClassDivision(eqGenClass, rules, termForms, notGenEqClass, t, nt)
	fmt.Println(flag)
	for flag {
		fmt.Println("again")
		flag, eqGenClass = checkEqClassDivision(eqGenClass, rules, termForms, notGenEqClass, t, nt)
	}
	fmt.Println("============================================")
	fmt.Println("new eq classes:", eqGenClass)

}
