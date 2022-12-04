#include <iostream>
#include <fstream>
#include <string>
#include <map>
#include <vector>
#include <algorithm>

using namespace std;

vector<string> nonTerms;
vector<string> terms;
vector<string> genNterms;
vector<string> reachNterms;

// struct First1Set {
//     string nterm;
//     vector<string> first1;
// };

struct rightPart {
    int type; //1 - nterm, 2 - term
    string val; //A, B, ..., a, b, ...
};

struct Rule {
    string left;
    vector <rightPart> right;
};

vector<Rule> grammar;
map <string, vector<string>> first_one_set;

bool err = false;

int getFirstAltIndex(string str) {
    //cout << str << endl;
    //cout << str.length() << endl;
	for (int i = 0; i < str.length(); i++) {
		if ( str[i] == '|' )
			return i;
	}
    return -1;
}

int getSecondAltIndex(string str) {
    //cout << endl;
    //cout << "----second alt ----" << endl;
    //cout << str << endl;
    //cout << str.length() << endl;
    int k = 0;
	for (int i = 0; i < str.length(); i++) {
        if ( str[i] == '|'  && k == 1) {
            //cout << "found: " << i << endl;
            //cout << "---- ----" << endl;
			return i;
        } 
        if ( str[i] == '|'  && k < 1) {
            //cout << i << endl;
			k++;
        }       
	}
    //cout << "there is no second alt" << endl;
    //cout << "---- ----" << endl;
    return -1;
}

vector<string> getListOfAltSubstrings(string str) {
    vector<string> s;
    //cout << "ALT SUBS" << endl;
    //cout << "string: "<< str << endl;
    int k = str.length();
    //cout << "k: "<< k << endl;
    int i = 0;
    while( i < k) {
        //cout << "i: "<< i << endl;
        //cout << "string: "<< str << endl;
        int alt = getFirstAltIndex(str);
        //cout << "alt index: "<< alt << endl;
        if (alt == -1) {
            s.push_back(str);
            return s;
        } else {
            //cout << "second alt index: " << getSecondAltIndex(str) << endl;
            if (alt == k - 1 || getSecondAltIndex(str)- alt == 1) {
                //cout << "case last alt" << endl;
                s.push_back(str.substr(i, 1));
                s.push_back("");
                return s;
            } else {
                //cout << "alternative substring push" << endl;
                s.push_back(str.substr(0, alt));
                str = str.substr(alt + 1);
                k = str.length();
                //i = 0;
            }
        }
        //i++;
    }
    return s;
}

vector<rightPart> parseCon(string str) {
    vector<rightPart> subs;
    //cout << "CON" << endl;
    //cout << "string: "<< str << endl;
    int i = 0;
    int len = str.length();
    if (str == "ε") {   // особенности заиписи этой буквы "Îµ"
        len = 1;
    }
    //cout << "len: " << len << endl;
    while (i < len) {
        //cout << "i: " << i << endl;
        string s = str.substr(i);
        //cout << "str: " << s << endl;
        if (s[0] == '[') {
            string n = str.substr(i+1, 1);
            //cout << "n: " << n << endl;
            nonTerms.push_back(n);  
            subs.push_back({1, n});
            i += 3;
        } else {
            string t = str.substr(i, 1);
            if (str == "ε") {   // особенности заиписи этой буквы 
                t = "";
            }
            //cout << "t: " << t << endl;
            if ( t != "" ) {
                terms.push_back(t);
            } 
            subs.push_back({2, t});
            i++;
        }
    }
    return subs;
}

vector<Rule> parseRuleLine(string str) {
    vector<Rule> rules;
    string nt = str.substr(1,1);
    nonTerms.push_back(nt);
    string leftPartOfRule = str.substr(5);
    //cout << getFirstAltIndex(str.substr(5)) << endl;
    if (getFirstAltIndex(leftPartOfRule) == -1) {
        //cout << "i am here" << endl;
        vector<rightPart> subs;
        subs = parseCon(leftPartOfRule);
        Rule r;
        r.left = nt;
        r.right = subs;
        rules.push_back(r);
    } else {
        vector<string> subs = getListOfAltSubstrings(leftPartOfRule);
        //cout << "alt subs" << endl;
        //cout << "alt subs len " << subs.size() << endl;
        for (int i = 0;  i < subs.size(); i++) {
            if (subs[i] == "") {
                subs[i] = "ε";
            }
            //cout << subs[i] << endl;
        }
        //cout << "alt subs size again " << subs.size() << endl;
        for (int i = 0;  i < subs.size(); i++) {
            //cout << "i " << i << endl;
            //cout << "len of subs{i}: " << subs[i].length() << endl;
            vector<rightPart> s;
            s = parseCon(subs[i]);
            Rule r;
            r.left = nt;
            r.right = s;
            rules.push_back(r);
        }
    }
    return rules;
}

string removeSpaces(string input) {
  input.erase(std::remove(input.begin(),input.end(),' '),input.end());
  return input;
}

int input(int n){
    string testName = "tests\\test" + to_string(n) + ".txt";
    //cout << testName << endl;
    ifstream file(testName);
    if (file.is_open()) {
        string str;
        while (getline(file,str)) {
            if (str.size()) {
                vector<Rule> rules;
                rules = parseRuleLine(removeSpaces(str));
                for (int i = 0; i < rules.size(); i++) {
                    grammar.push_back(rules[i]);
                }
            }
        }
        file.close();
    }
    if (grammar.size() == 0)
        return false;
    else 
        return true;
}

void printTerms() {
    sort(nonTerms.begin(), nonTerms.end());
    auto last = unique(nonTerms.begin(), nonTerms.end());
    nonTerms.erase(last, nonTerms.end());
    cout << "NTERMS: ";
    for (int i = 0;  i < nonTerms.size(); i++) cout << nonTerms[i] << " ";
    cout << endl;
    // sort(terms.begin(), terms.end());
    // auto lastUniqueTerm = unique(terms.begin(), terms.end());
    // terms.erase(lastUniqueTerm, terms.end());
    // cout << "TERMS: ";
    // for (int i = 0;  i < terms.size(); i++) cout << terms[i] << " ";
    // cout << endl;
}

void printGrammar() {
    cout << "grammar size: " << grammar.size() << endl;
    for (int i = 0; i < grammar.size(); i++) {
        Rule r = grammar[i];
        cout << endl;
        cout << "RULE " << i + 1 << endl;
        cout << "left: " << r.left << ", ";
        cout << "right: ";
        for (int j = 0;  j < r.right.size(); j++) {
            cout << "{type:" << r.right[j].type << ", val:" << r.right[j].val << "}";
            if (j != r.right.size() - 1) cout << ", ";   
        }
        cout << endl;
    }
    cout << endl;
}

// если есть нетерминал в правой части, то вернет его "индекс", иначе -1
int FirstIndexNterminRightPart(vector<rightPart> rights) {
    for(int i = 0; i < rights.size(); i++) {
        if (rights[i].type == 1) return i;
    }
    return -1;
}

// принадлежит ли нетерминал вектору (множеству порождающих нетерминалов)
bool isInGenNterms(string nterm, vector<string> gen) {
    for (auto el : gen) {
        if (el == nterm) return true;
    }
    // if (binary_search(gen.begin(), gen.end(), nterm)) {
    //     return true;
    // }
    return false;
}

bool isInReachNterms(string nterm, vector<string> reach) {
    for (auto el : reach) {
        if (el == nterm) return true;
    }
    return false;
}

// количество нетерминалов в правой части правила
int getQuantityOfNterms(vector<rightPart> rights) {
    int s = 0;
    for(int i = 0; i < rights.size(); i++) {
        if (rights[i].type == 1) s++;
    }
    return s;
}

void updateGrammar() {
    int len = grammar.size();
    for (int i = 0; i < len; i++) {
        //cout << "left: " << grammar[i].left << endl;
        if (!isInGenNterms(grammar[i].left, nonTerms)) {
            //cout << "it is not gen" << endl;
            grammar.erase(grammar.begin() + i);
            len = grammar.size();
            //cout << "grammar len: " << len << endl;
            i = 0;
            //printGrammar();
        } else {
            for (int j = 0;  j < grammar[i].right.size(); j++) {
                if (grammar[i].right[j].type == 2) {
                    //cout << "terminal" << endl;
                    continue;
                }
                //cout << "{type:" << grammar[i].right[j].type << ", val:" << grammar[i].right[j].val << "}" << endl;;
                if (!isInGenNterms(grammar[i].right[j].val, nonTerms) && grammar[i].right[j].type == 1) {
                    //cout << "it is not gen" << endl;
                    grammar.erase(grammar.begin() + i);
                    j++;
                    len = grammar.size();
                    //cout << "grammar len: " << len << endl;
                    i = 0;
                    //printGrammar();
                }
            }
        }        
    }
}

void removeNonGeneratingNterms() {
    int setSize = 0;
    int r = 0;
    //cout << "grammar size:" << grammar.size() << endl;
    //cout << "==ШАГ 1==" <<endl;
    while (r < grammar.size()) {
        string nterm = grammar[r].left;
        vector<rightPart> rights = grammar[r].right;
        //cout << "NTERM: " << nterm << endl;
        //cout << "RIGHTPART: ";
        // for (int j = 0;  j < rights.size(); j++) {
        //     cout << "{type:" << rights[j].type << ", val:" << rights[j].val << "}";
        //     if (j != rights.size() - 1) cout << ", ";   
        // }
        // cout << endl;
        // шаг 1 находим правила не содерж нетерминалов в правой части
        if (FirstIndexNterminRightPart(rights) != -1)  {
            r++;
        } else {
            //cout << "не содержит нетерминалов в правой части" << endl;
            if (!isInGenNterms(nterm, genNterms) && nterm != "") {
                genNterms.push_back(nterm);
            }
            r++;
        }      
    }
    // cout << "===" << endl;   
    // for (auto n: genNterms) {
    //     cout << n << " ";
    // }
    // cout << endl;
    // cout << "===" << endl; 
    // cout << endl;
    // cout << "==ШАГ 2==" <<endl;

    // шаг 2 если найдено правило, все нетерминалы правой части которого уже
    // входят в множество, то добавляем левый нетерминал 
    // если множество порождающих нетерминалов изменилось, повторяем шаг 2
    while (genNterms.size() > setSize) {
        r = 0;
        setSize = genNterms.size();
        while (r < grammar.size()) {   
            string nterm = grammar[r].left;
            vector<rightPart> rights = grammar[r].right; 
            // cout << "NTERM: " << nterm << endl;
            // cout << "RIGHTPART: ";
            // for (int j = 0;  j < rights.size(); j++) {
            //     cout << "{type:" << rights[j].type << ", val:" << rights[j].val << "}";
            //     if (j != rights.size() - 1) cout << ", ";   
            // }
            // cout << endl;
            int col = getQuantityOfNterms(rights);
            int k = 0;
            //cout << "количество нетерминалов в правой части: " << col << endl;
            for (int j = 0;  j < rights.size(); j++) {
                if (col == 0) {
                    //cout << "ZERO NTERMS AT RIGHT" << endl;
                    r++;
                }
                if (col > 0) {    
                    string cur = rights[j].val;
                    //cout << "поиск " << cur << endl;
                    if (isInGenNterms(cur,  genNterms)) {
                        k++;
                        //cout << "нетерминал уже есть в множестве порождающих k="<< k << endl;
                        if (k == col) {
                            //cout << "!правило где все нетерминалы справа в множестве порождающих!" << endl;
                            if (!isInGenNterms(nterm, genNterms)) {
                                genNterms.push_back(nterm);
                                //cout << "size of gen nterms: " << genNterms.size() << endl;
                                // cout << "===" << endl;   
                                // for (auto n: genNterms) {
                                //     cout << n << " ";
                                // }
                                // cout << endl;
                                // cout << "===" << endl; 
                                // cout << endl;
                            }
                            r++;
                        }
                    } else {
                        if (binary_search(nonTerms.begin(), nonTerms.end(), cur)) {
                            //cout << "не принадлежит порождающим, следующее правило" << endl;
                            r++;
                        }
                    }
                }
                
            }
        }
    }  
    /*cout << "===" << endl;   
    for (auto n: genNterms) {
        cout << n << " ";
    }
    cout << endl;
    cout << "===" << endl; */
    cout << endl;
    nonTerms = genNterms;   
}

void removeUnreachableNterms() {
    reachNterms.push_back("S"); 
    int setSize = 0;
    int r = 0;
    //cout << "grammar size:" << grammar.size() << endl;
    while (reachNterms.size() > setSize) {
        setSize = reachNterms.size();
        while (r < grammar.size()) {
            string nterm = grammar[r].left;
            vector<rightPart> rights = grammar[r].right;
            // cout << "NTERM: " << nterm << endl;
            // cout << "RIGHTPART: " << endl;
            // добавляем нетерминалы достижимые из данного
            if (isInReachNterms(nterm, reachNterms)) {
                for (int j = 0;  j < rights.size(); j++) {
                    //cout << "{type:" << rights[j].type << ", val:" << rights[j].val << "}" << endl;
                    if (rights[j].type == 1) {
                        reachNterms.push_back(rights[j].val);
                    }
                }    
                sort(reachNterms.begin(), reachNterms.end());
                auto l = unique(reachNterms.begin(), reachNterms.end());
                reachNterms.erase(l, reachNterms.end());
                
                // cout << "===" << endl;   
                // for (auto n: reachNterms) {
                //     cout << n << " ";
                // }
                // cout << endl;
                // cout << "===" << endl;
                // обновляем правила
                nonTerms = reachNterms;
            }
            r++;     
        }
    }
    // cout << "===" << endl;   
    // for (auto n: reachNterms) {
    //     cout << n << " ";
    // }
    // cout << endl;
    // cout << "===" << endl; 
    // cout << endl;
    nonTerms = reachNterms; 
}

void printFirst1Set() {
    for (auto n : nonTerms) {
        cout << "FIRST1(" << n << ") = {";
        for(int i = 0; i < first_one_set[n].size(); i++) {
            if ( i == first_one_set[n].size() - 1) cout << first_one_set[n][i];
            else cout << first_one_set[n][i] << ", ";
        }
        cout << "}" << endl;
    }
}

vector<string> findFirst1(string nterm) {
    vector<string> f;
    for (auto rule : grammar) {
        if (rule.left == nterm) {
            if (first_one_set[nterm].size() == 1 && first_one_set[nterm][0] == "") {
                cout << "remove first empty set" << endl;
                first_one_set[nterm].pop_back();
            }
            vector<rightPart> rt = rule.right;
            if (rt.size() == 1) {
                cout << "nonterm -> term" << endl;
                if (rt[0].val == "") {
                    cout << "term is eps" << endl;
                    f.push_back("eps");
                } else {
                    if (rt[0].type == 2) {
                        cout << "term is letter" << endl;
                        f.push_back(rt[0].val);
                    }
                    if (rt[0].type == 1) {
                        cout << "term is another nonterm" << endl;
                        vector<string> vec = findFirst1(rt[0].val);
                        for (auto v : vec) {
                            f.push_back(v);
                        }  
                    }      
                }
            } else {
                cout << "nterm -> term1 term2..." << endl;
                if(rt[0].type == 2) {
                    cout << "term1 is letter" << endl;
                    f.push_back(rt[0].val);
                }
                if(rt[0].type == 1) {
                    cout << "term1 is nonterm" << endl;\
                    if (rt[0].val[0] == nterm.c_str()[0]) {
                        cout << "левая рекурсия" << endl;
                        break;
                    }
                    vector<string> vec2 = findFirst1(rt[0].val);
                    cout << "find first1 set of term1" << endl;
                    for(auto v : vec2) {
                        f.push_back(v);
                    }
                    if (isInGenNterms("eps", vec2)) {
                        cout << "first1 set of term1 might has eps, do the union with next first1"<< endl;
                        for(int i = 1; i < rt.size(); i++) {
                            vector<string> vec3 = findFirst1(rt[i].val);
                            if (isInGenNterms("eps", vec3)) {
                                for(auto v : vec3) {
                                    f.push_back(v);
                                }
                            }
                            if (rt[rt.size() - 1].type == 2 && isInGenNterms("eps", vec3)) {
                                f.push_back(rt[rt.size() - 1].val);
                            }
                        }
                    }      
                }
            }
        }
    }
    sort(f.begin(), f.end());
    auto lt = unique(f.begin(), f.end());
    f.erase(lt, f.end());
    return f;
}

void constructFirst1() {
    int setSize = 0;
    for (auto n : nonTerms) {
        first_one_set[n].push_back("");
    }
    for (auto n : nonTerms) {
        bool changed = true;
        cout << "NETERM: " << n << endl;
        while(changed) {
            changed = false;
            setSize = first_one_set[n].size();
            cout << "setsize: " << setSize << endl;
            vector<string> res = findFirst1(n);
            changed = res.size() > setSize;
            first_one_set[n] = res;
            printFirst1Set();
        }
        cout << endl;
    }
}

int main() {
    int n;
    cout << "Enter test number" << endl;
    cin >> n;
    bool err = input(n);
    if (!err) {
        cout << "INCORRECT TEST FILE!";
        return 0;
    }
    cout << "> Parsed grammar <" << endl;
    printTerms();
    printGrammar();
    //убираем непорождающие нетерминалы
    removeNonGeneratingNterms();
    cout << "> Grammar with removed non-generating nonterminals <" << endl;
    printTerms();
    updateGrammar();
    printGrammar();
    //убираем недостижимые нетерминалы
    cout << "> Grammar with removed unreachable nonterminals <" << endl;
    removeUnreachableNterms();
    printTerms();
    updateGrammar();
    printGrammar();
    cout << "> FIRST 1 sets for nonterminals <" << endl;
    constructFirst1();
    printFirst1Set();
    return 0;
}

// ааааааааа