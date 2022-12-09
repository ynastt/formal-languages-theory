import sys
from antlr4 import *
from dist.MyGramLexer import MyGramLexer
from dist.MyGramParser import MyGramParser

def main(argv):
    data = InputStream(input(">>> "))
    #input_stream = FileStream(data)
    lexer = MyGramLexer(data)
    stream = CommonTokenStream(lexer)
    parser = MyGramParser(stream)

    parser.s()
    print('DONE!')

if __name__ == "__main__":
    main(sys.argv)


# генератор ANTLR
# ограничения на входную грамматику
# жаловался no fragment rules когда нетерминалы были с заглавной буквы
# написала грамматику для академической регулярки (a*bc*)


