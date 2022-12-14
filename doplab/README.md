### Реализация множеств ```First_1``` , ```Follow_1``` , ```First_k``` для КС грамматики

#### Информация
- [First1 Follow1 множества](https://neerc.ifmo.ru/wiki/index.php?title=%D0%9F%D0%BE%D1%81%D1%82%D1%80%D0%BE%D0%B5%D0%BD%D0%B8%D0%B5_FIRST_%D0%B8_FOLLOW#lemmafirst1)
- [First_k](https://github.com/TonitaN/FormalLanguageTheory/blob/main/2022/lect_tfl_8.pdf)
- [Удаление непорождающих и недостижимых нетерминалов грамматики](https://neerc.ifmo.ru/wiki/index.php?title=%D0%A3%D0%B4%D0%B0%D0%BB%D0%B5%D0%BD%D0%B8%D0%B5_%D0%B1%D0%B5%D1%81%D0%BF%D0%BE%D0%BB%D0%B5%D0%B7%D0%BD%D1%8B%D1%85_%D1%81%D0%B8%D0%BC%D0%B2%D0%BE%D0%BB%D0%BE%D0%B2_%D0%B8%D0%B7_%D0%B3%D1%80%D0%B0%D0%BC%D0%BC%D0%B0%D1%82%D0%B8%D0%BA%D0%B8)

#### Формат входных данных грамматики
- Грамматика записана в файле с названием ```test<номер>.txt``` (При запуске main.cpp после 'Enter test name' надо ввести этот <номер>)
- Нетереминалы должны быть записаны в квадратных скобках ```[S], [A], ...```
- Стартовый нетерминал грамматики ```[S]```
- Можно использовать знак альтернативы ``` [S] -> [S] b [B] | [B] c | [B] ```
- Если есть правило перевода нетерминала в пустую строку, оно должно быть записано с помощью альтернативы ```[S]->[A]a |```

#### Формат выходных данных 
- ```First1(N) = { t, eps }```
- ```Follow(N) = { t, $ }```
- ```Firstk(N) = { t1...tk, t1...tj, eps }``` (где j < k и N -> t1...tj)
##### Дополнительно можно вывести распарсенную грамматику с помощью:
``` printGrammar(); ```
Она будет иметь вид: 
``` 
RULE 1
left: N, right: {type:1, val:N}, {type:2, val:t}
.
.
.
RULE n
left: N, right: {type:1, val:N}, {type:2, val:t}
```
##### Дополнительно можно вывести грамматику без непорождающих нетерминалов:
```
removeNonGeneratingNterms();   
updateGrammar();
printGrammar(); 
```
формат выхода аналогичный, но убраны правила с непорождающими нетерминалами 

##### Дополнительно можно вывести грамматику без непорождающих и недостижимых нетерминалов:
```
removeNonGeneratingNterms(); 
removeUnreachableNterms(); 
updateGrammar();
printGrammar(); 
```
Важен именно такой порядок выоплнения функций
формат выхода аналогичный, но убраны правила с непорождающими  и недостижимыми нетерминалами 

#### Пример грамматики и множеств ```First_1``` , ```Follow_1``` , ```First_k```

```
[S]->[A]a |
[S]->[A][B][C]d
[A]->abc |
[B]->g |
[C]->l
[C]->ki |

> FIRST 1 sets for nonterminals <
FIRST1(A) = {a, eps}
FIRST1(B) = {eps, g}
FIRST1(C) = {eps, k, l}
FIRST1(S) = {A, a, eps, g, k, l}

> FOLLOW sets for nonterminals <
FOLLOW(A) = {a, d, g, k, l}
FOLLOW(B) = {d, k, l}
FOLLOW(C) = {d}
FOLLOW(S) = {$}

> FIRST k sets for nonterminals <
FIRST2(A) = {ab, eps}
FIRST2(B) = {eps, g}
FIRST2(C) = {eps, ki, l}
FIRST2(S) = {a, ab, eps, gd, gk, gl}
```

###### TODO
нормальные логи в файл
