import sys
from antlr4 import *
from dist.MyGramLexer import MyGramLexer
from dist.MyGramParser import MyGramParser

def main(argv):
    data = InputStream(input(">>> "))
    lexer = MyGramLexer(data)
    stream = CommonTokenStream(lexer)
    parser = MyGramParser(stream)

    parser.s()
    print('DONE!')

if __name__ == "__main__":
    main(sys.argv)



