### Дополнительное задание 
- Установить генератор парсеров (не YACC, не Bison);
- Разобраться, какие ограничения этот генератор накладывает на входную грамматику, написать отчет;
- Написать на языке генератора грамматику для языка академической регулярки.

### Реализация
- Установлен генератор ANTLR
- Отчет об ограничениях, накладываемых генератором на грамматику [здесь](https://github.com/ynastt/formal-languages-theory/blob/main/antlr/otchet.pdf)
- На языке генератора написана грамматика для академической регулярки (a\*bc\*)
- Если слово, введенное во входной поток, не принадлежит языку данной грамматики, то появится строка с предупреждением о нарушении грамматики.

Например:
```
>>> aaaaaac
line 1:5 mismatched input 'c' expecting {'a', 'b'}
DONE!
```

### Полезные ссылочки
- [ANTLR v4](https://github.com/pboyer/antlr4)
- [Introduction to ANTLR using python](https://faun.pub/introduction-to-antlr-python-af8a3c603d23)
- [How to create a Python lexer or parser](https://github.com/antlr/antlr4/blob/master/doc/python-target.md)
- [ANTLR Tool Command Line Options](https://github.com/antlr/antlr4/blob/master/doc/tool-options.md)

 ~~не нашла туториал на Go :(~~
