package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

// замена нетерминалов в правой части правил на "_" для формирования терминальных форм
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

// поиск класса эквивалентности для нетерминала
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

// удаление неверного класса эквивалентности и замена на новые классы
func removeClass(firstEqClasses []string, class1 string, nonTerm string) []string {
	for ind, class := range firstEqClasses {
		if class == class1 {
			class = strings.ReplaceAll(class, nonTerm, "")
			firstEqClasses[ind] = class
			firstEqClasses = append(firstEqClasses, nonTerm)
			sort.Strings(firstEqClasses)
			//fmt.Println("classes:", firstEqClasses)
			//fmt.Printf("%s was removed from class %s\n", nonTerm, class)
		}
	}
	return firstEqClasses
}

// построчное чтение входных данных из файла tests/test*.txt, * - номер теста
// заполнение массивов нетерминалов, терминалов и карты: нетерминал -> спсиок правил переписыванния
func parseTerms() ([]string, []string, map[string][]string) {
	var nonTerms, terms []string
	rules := make(map[string][]string)
	file, err := os.Open("tests/test7.txt")
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
func checkEqClassDivision(firstEqClasses []string, rules, termForms map[string][]string, t, nt, gnt string) (bool, []string) {
	flag := false
	var nonTerms []string
	for k := range termForms {
		nonTerms = append(nonTerms, k)
	}
	sort.Strings(nonTerms)
	for _, nonTerm := range nonTerms {
		//fmt.Println("nonterm now:", nonTerm)
		for _, rule1 := range rules[nonTerm] {
			//fmt.Printf("nonterm rule %s->%s\n", nonTerm, rule1)
			if len(rule1) == 1 && (strings.ContainsAny(t, rule1) || strings.ContainsAny(gnt, rule1)) {
				//fmt.Println("-> term or not generating neterm", rule1)
				continue
			}
			f1, c1 := addUnderscore(rule1, nt)
			for _, key := range nonTerms {
				if nonTerm != key {
					//fmt.Printf("\tnt is: %s\n", key)
					for _, rule2 := range rules[key] {
						//fmt.Printf("\t\tnt rule %s->%s\n", key, rule2)
						if len(rule2) == 1 && (strings.ContainsAny(t, rule2) || strings.ContainsAny(gnt, rule2)) {
							//fmt.Println("-> term or not generating neterm", rule2)
							continue
						}
						f2, c2 := addUnderscore(rule2, nt)
						//fmt.Println("HERE F1 F2, c1, c2", f1, f2, c1, c2)
						//fmt.Printf("%s->%s %s->%s\n", nonTerm, rule1, key, rule2)
						class1 := findClass(nonTerm, firstEqClasses)
						class2 := findClass(key, firstEqClasses)
						//fmt.Printf("class nonterm: %s, class nt: %s\n", class1, class2)
						if class1 == class2 && f1 == f2 {
							//fmt.Println("CLASSES and termforms ARE SAME")
							class3 := findClass(c1, firstEqClasses)
							class4 := findClass(c2, firstEqClasses)
							//fmt.Printf("class nonterm _ : %s, class nt _ : %s\n", class3, class4)
							if class3 != class4 {
								//fmt.Println("несовпали классы экв у Ni' Nj'")
								if class3 != class1 {
									firstEqClasses = removeClass(firstEqClasses, class1, nonTerm)
									flag = true
									//fmt.Println("now classes are:", firstEqClasses)
								} else if class4 != class1 {
									firstEqClasses = removeClass(firstEqClasses, class1, key)
									flag = true
									//fmt.Println("now classes are:", firstEqClasses)
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

// ответ: Классы эквивалентности нетерминалов + новая грамматика, где в терминальные формы подставлены
// соответствующие представители классов эквивалентности
func outputNewRules(termForms, rules map[string][]string, eqGenClass []string, eqNotGenClass []string, t, nt, gnt string) {
	var eqGenClassOld []string
	for i, _ := range eqGenClass {
		e := eqGenClass[i]
		eqGenClassOld = append(eqGenClassOld, e)
	}
	//fmt.Println("old classes", eqGenClassOld)
	for i, c := range eqGenClass {
		if len(c) > 1 {
			fmt.Print("{")
			for j, sym := range c {
				if j == len(c)-1 {
					fmt.Printf("%s}\n", string(sym))
				} else {
					fmt.Printf("%s,", string(sym))
				}
			}
			c = strings.ReplaceAll(c, c[1:], "")
			eqGenClass[i] = c
		} else {
			fmt.Printf("{%s}\n", c)
		}
	}
	if len(eqNotGenClass) != 0 {
		fmt.Print("{")
		for i, c := range eqNotGenClass {
			if i == len(eqNotGenClass)-1 {
				fmt.Printf("%s}\n", c)
			} else {
				fmt.Printf("%s,", c)
			}
		}
	}
	for nonTerm, form := range termForms {
		for _, f := range form {
			classOld := findClass(nonTerm, eqGenClassOld)
			class := findClass(nonTerm, eqGenClass)
			if len(class) != 0 {
				if len(f) == 1 && (strings.ContainsAny(t, f) || strings.ContainsAny(gnt, f)) {
					continue
				} else {
					for j, r := range rules[nonTerm] {
						if len(f) == len(r) {
							for _, sym := range r {
								if strings.ContainsAny(t, string(sym)) {
									continue
								}
								if strings.ContainsAny(classOld, string(sym)) {
									r = strings.ReplaceAll(r, string(sym), class)
									rules[nonTerm][j] = r
									rules[nonTerm] = checkForDuplicates(rules[nonTerm])
								}
							}
						}
					}
				}
			}
		}
	}
}

func main() {
	terms, nonTerms, rules := parseTerms()
	nt := strings.Join(nonTerms, "")
	//fmt.Println("nonterms:", nonTerms[0:])
	//fmt.Println("terms:", terms[0:])
	//fmt.Println("rules:", rules)
	termForms, eqNotGenClass := makeListOfTermForms(nonTerms, rules, nt)
	//fmt.Println("term forms:", termForms)
	//fmt.Println("notGenEqClass:", eqNotGenClass)
	eqGenClass := eqClassesDivision(termForms)
	//fmt.Println("first eq classes:", eqGenClass)
	var flag bool
	t := strings.Join(terms, "")
	gnt := strings.Join(eqNotGenClass, "")
	flag, eqGenClass = checkEqClassDivision(eqGenClass, rules, termForms, t, nt, gnt)
	//fmt.Println(flag)
	for flag {
		//fmt.Println("again")
		flag, eqGenClass = checkEqClassDivision(eqGenClass, rules, termForms, t, nt, gnt)
	}
	//fmt.Println("new eq classes:", eqGenClass)
	outputNewRules(termForms, rules, eqGenClass, eqNotGenClass, t, nt, gnt)
	for nonTerm, _ := range termForms {
		for j, _ := range rules[nonTerm] {
			class := findClass(nonTerm, eqGenClass)
			if nonTerm == class {
				fmt.Printf("%s -> %s\n", nonTerm, rules[nonTerm][j])
			}
		}
	}
}
