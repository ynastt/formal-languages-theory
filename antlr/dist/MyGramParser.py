# Generated from MyGram.g4 by ANTLR 4.11.1
# encoding: utf-8
from antlr4 import *
from io import StringIO
import sys
if sys.version_info[1] > 5:
	from typing import TextIO
else:
	from typing import TextIO

def serializedATN():
    return [
        4,1,3,16,2,0,7,0,2,1,7,1,1,0,1,0,1,0,1,0,3,0,9,8,0,1,1,1,1,1,1,3,
        1,14,8,1,1,1,0,0,2,0,2,0,0,15,0,8,1,0,0,0,2,13,1,0,0,0,4,5,5,1,0,
        0,5,9,3,0,0,0,6,7,5,2,0,0,7,9,3,2,1,0,8,4,1,0,0,0,8,6,1,0,0,0,9,
        1,1,0,0,0,10,11,5,3,0,0,11,14,3,2,1,0,12,14,5,0,0,1,13,10,1,0,0,
        0,13,12,1,0,0,0,14,3,1,0,0,0,2,8,13
    ]

class MyGramParser ( Parser ):

    grammarFileName = "MyGram.g4"

    atn = ATNDeserializer().deserialize(serializedATN())

    decisionsToDFA = [ DFA(ds, i) for i, ds in enumerate(atn.decisionToState) ]

    sharedContextCache = PredictionContextCache()

    literalNames = [ "<INVALID>", "'a'", "'b'", "'c'" ]

    symbolicNames = [  ]

    RULE_s = 0
    RULE_t = 1

    ruleNames =  [ "s", "t" ]

    EOF = Token.EOF
    T__0=1
    T__1=2
    T__2=3

    def __init__(self, input:TokenStream, output:TextIO = sys.stdout):
        super().__init__(input, output)
        self.checkVersion("4.11.1")
        self._interp = ParserATNSimulator(self, self.atn, self.decisionsToDFA, self.sharedContextCache)
        self._predicates = None




    class SContext(ParserRuleContext):
        __slots__ = 'parser'

        def __init__(self, parser, parent:ParserRuleContext=None, invokingState:int=-1):
            super().__init__(parent, invokingState)
            self.parser = parser

        def s(self):
            return self.getTypedRuleContext(MyGramParser.SContext,0)


        def t(self):
            return self.getTypedRuleContext(MyGramParser.TContext,0)


        def getRuleIndex(self):
            return MyGramParser.RULE_s

        def enterRule(self, listener:ParseTreeListener):
            if hasattr( listener, "enterS" ):
                listener.enterS(self)

        def exitRule(self, listener:ParseTreeListener):
            if hasattr( listener, "exitS" ):
                listener.exitS(self)




    def s(self):

        localctx = MyGramParser.SContext(self, self._ctx, self.state)
        self.enterRule(localctx, 0, self.RULE_s)
        try:
            self.state = 8
            self._errHandler.sync(self)
            token = self._input.LA(1)
            if token in [1]:
                self.enterOuterAlt(localctx, 1)
                self.state = 4
                self.match(MyGramParser.T__0)
                self.state = 5
                self.s()
                pass
            elif token in [2]:
                self.enterOuterAlt(localctx, 2)
                self.state = 6
                self.match(MyGramParser.T__1)
                self.state = 7
                self.t()
                pass
            else:
                raise NoViableAltException(self)

        except RecognitionException as re:
            localctx.exception = re
            self._errHandler.reportError(self, re)
            self._errHandler.recover(self, re)
        finally:
            self.exitRule()
        return localctx


    class TContext(ParserRuleContext):
        __slots__ = 'parser'

        def __init__(self, parser, parent:ParserRuleContext=None, invokingState:int=-1):
            super().__init__(parent, invokingState)
            self.parser = parser

        def t(self):
            return self.getTypedRuleContext(MyGramParser.TContext,0)


        def EOF(self):
            return self.getToken(MyGramParser.EOF, 0)

        def getRuleIndex(self):
            return MyGramParser.RULE_t

        def enterRule(self, listener:ParseTreeListener):
            if hasattr( listener, "enterT" ):
                listener.enterT(self)

        def exitRule(self, listener:ParseTreeListener):
            if hasattr( listener, "exitT" ):
                listener.exitT(self)




    def t(self):

        localctx = MyGramParser.TContext(self, self._ctx, self.state)
        self.enterRule(localctx, 2, self.RULE_t)
        try:
            self.state = 13
            self._errHandler.sync(self)
            token = self._input.LA(1)
            if token in [3]:
                self.enterOuterAlt(localctx, 1)
                self.state = 10
                self.match(MyGramParser.T__2)
                self.state = 11
                self.t()
                pass
            elif token in [-1]:
                self.enterOuterAlt(localctx, 2)
                self.state = 12
                self.match(MyGramParser.EOF)
                pass
            else:
                raise NoViableAltException(self)

        except RecognitionException as re:
            localctx.exception = re
            self._errHandler.reportError(self, re)
            self._errHandler.recover(self, re)
        finally:
            self.exitRule()
        return localctx





