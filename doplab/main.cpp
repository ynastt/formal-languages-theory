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
vector<string> hasEpsInGenTerms;

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
map <string, vector<string>> first_k_set;
map <string, vector<string>> follow_set;

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
    if(str[k-1] == '|') k += 1;
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
            if ( k == 1  && str == "" ) {
                //cout << "case last alt" << endl;
                s.push_back(str.substr(i, 1));
                s.push_back("");
                return s;
            } else {
                //cout << "alternative substring push" << endl;
                s.push_back(str.substr(0, alt));
                str = str.substr(alt + 1);
                k = str.length();
                //cout << "now k is" << k << endl;
                if (k == 0) {
                    k = 1;
                }
            }
        }
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

// принадлежит ли нетерминал вектору (например множеству порождающих нетерминалов)
bool isInGenNterms(string nterm, vector<string> gen) {
    for (auto el : gen) {
        if (el == nterm) return true;
    }
    // if (binary_search(gen.begin(), gen.end(), nterm)) {
    //     return true;
    // }
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
    // cout << "grammar size:" << grammar.size() << endl;
    // cout << "==ШАГ 1==" <<endl;
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

        // шаг 1 находим правила не содерж нетерминалов в правой части
        if (FirstIndexNterminRightPart(rights) != -1)  {
            // cout << "есть нетерминалы, идем дальше" << endl;
            r++;
        } else {
            // cout << "не содержит нетерминалов в правой части" << endl;
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
    // cout << "grammar size" << grammar.size() << endl;
    while (genNterms.size() != setSize) {
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
            // cout << "количество нетерминалов в правой части: " << col << endl;
            for (int j = 0;  j < rights.size(); j++) {
                if (col == 0) {
                    // cout << "ZERO NTERMS AT RIGHT" << endl;
                    continue;
                }
                if (col > 0) {    
                    string cur = rights[j].val;
                    // cout << "поиск " << cur << endl;
                    if (isInGenNterms(cur,  genNterms)) {
                        k++;
                        // cout << "нетерминал уже есть в множестве порождающих k="<< k << endl;
                        if (k == col) {
                            // cout << "!правило где все нетерминалы справа в множестве порождающих!" << endl;
                            if (!isInGenNterms(nterm, genNterms)) {
                                genNterms.push_back(nterm);
                                // cout << "size of gen nterms: " << genNterms.size() << endl;
                                // cout << "===" << endl;   
                                // for (auto n: genNterms) {
                                //     cout << n << " ";
                                // }
                                // cout << endl;
                                // cout << "===" << endl; 
                                // cout << endl;
                            }
                        }
                    } else {
                        //if (binary_search(nonTerms.begin(), nonTerms.end(), cur)) {
                            // cout << "не принадлежит порождающим, следующее правило" << endl;
                            continue;
                            //r++;
                        //}
                    }
                    
                }
                
            }
            r++;
        }
    }  
    cout << "===" << endl;   
    for (auto n: genNterms) {
        cout << n << " ";
    }
    cout << endl;
    cout << "===" << endl; 
    cout << endl;
    nonTerms = genNterms;   
}

void removeUnreachableNterms() {
    reachNterms.push_back("S"); 
    int setSize = 0;
    int r = 0;
    //cout << "grammar size:" << grammar.size() << endl;
    while (reachNterms.size() != setSize) {
        setSize = reachNterms.size();
        while (r < grammar.size()) {
            string nterm = grammar[r].left;
            vector<rightPart> rights = grammar[r].right;
            // cout << "NTERM: " << nterm << endl;
            // cout << "RIGHTPART: " << endl;
            // добавляем нетерминалы достижимые из данного
            if (isInGenNterms(nterm, reachNterms)) {
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
    cout << endl;
}

void printFirstkSet(int k)  {
    for (auto n : nonTerms) {
        cout << "FIRST" << k << "(" << n << ") = {";
        for(int i = 0; i < first_k_set[n].size(); i++) {
            if ( i == first_k_set[n].size() - 1) cout << first_k_set[n][i];
            else cout << first_k_set[n][i] << ", ";
        }
        cout << "}" << endl;
    }
    cout << endl;
}

void printFollowSet() {
    for (auto n : nonTerms) {
        cout << "FOLLOW(" << n << ") = {";
        for(int i = 0; i < follow_set[n].size(); i++) {
            if ( i == follow_set[n].size() - 1) cout << follow_set[n][i];
            else cout << follow_set[n][i] << ", ";
        }
        cout << "}" << endl;
    }
    cout << endl;
}

void constructFirst1() {
    int setSize = 0;
    for (auto n : nonTerms) {
        first_one_set[n].push_back("");
    }
    bool changed = true;
    while(changed) {
        changed = false;
        for (auto n : nonTerms) {
            //cout << "NETERM: " << n << endl;
            setSize = first_one_set[n].size();
            //cout << "setsize: " << setSize << endl;
            for (auto rule : grammar) {
                if (rule.left == n) {
                    if (first_one_set[n].size() == 1 && first_one_set[n][0] == "") {
                        //cout << "remove first empty set" << endl;
                        first_one_set[n].pop_back();
                        setSize = 0;
                    }
                    vector<rightPart> rt = rule.right;
                    for (int z = 0; z < rt.size(); z++) {
                        if (rt[z].type == 2) {
                            if (rt[0].val == "") {
                                //cout << "term is eps" << endl;
                                first_one_set[n].push_back("eps");
                            } else {
                                first_one_set[n].push_back(rt[0].val);
                            }
                            break;
                        }
                        if (rt[z].type == 1) {
                            vector<string> vec = first_one_set[rt[z].val];
                            for(auto v : vec) {
                                first_one_set[n].push_back(v);
                            }
                            if (!isInGenNterms("eps", first_one_set[n]))
                                break;
                        }
                    }
                    
                }
            }
            for (int i = 0; i < first_one_set[n].size(); i++) {
                if (first_one_set[n][i] == "") 
                    first_one_set[n].erase(first_one_set[n].begin()+ i);
            }
            sort(first_one_set[n].begin(), first_one_set[n].end());
            auto lt = unique(first_one_set[n].begin(), first_one_set[n].end());
            first_one_set[n].erase(lt, first_one_set[n].end());

            
            changed = first_one_set[n].size() != setSize;
            //printFirst1Set();
            //cout << endl;
        }
    }
}

void constructFollow() {
    int setSize = 0;
    for (auto n : nonTerms) {
        if (n == "S") follow_set["S"].push_back("$");
        else follow_set[n].push_back("");
    }
    //cout << "initial follow" << endl;
    //printFollowSet();
    bool changed = true;
    while(changed) {
        changed = false;
        for (auto rule : grammar) {
            
            // cout << "==rule: ";
            // cout << "left: " << rule.left << ", ";
            // cout << "right: ";
            // for (int j = 0;  j < rule.right.size(); j++) {
            //     cout << "{type:" << rule.right[j].type << ", val:" << rule.right[j].val << "}";
            //     if (j != rule.right.size() - 1) cout << ", ";   
            // }
            // cout << endl;    
            string nterm = rule.left;
            vector<rightPart> rt = rule.right;
            for(int i = 0; i < rt.size(); i++) {
                //cout << "i " << i << " rt size: " << rt.size() << endl;
                rightPart right = rt[i];
                if (right.type == 1) {
                    //cout << "nonterm right " << rt[i].val << endl;
                    setSize = follow_set[right.val].size();
                    if (setSize == 1 && follow_set[right.val][0] == "") {
                        setSize = 0;
                        //follow_set[right.val].pop_back();
                    }    
                    //cout << "setsize is: "<< setSize << endl;
                    for(int j = i + 1; j < rt.size() + 1; j++) {
                        
                        //cout << "searching for follow" << endl;
                        
                        if (j == rt.size()) {
                            //cout << " last sym" << endl;
                            vector<string> vec2 = follow_set[nterm];
                            if (follow_set[right.val][0] == "") {
                                follow_set[right.val].erase(follow_set[right.val].begin() + 0);
                            } 
                            for (auto v : vec2) {
                                if (v == "") v ="eps";
                                follow_set[right.val].push_back(v);
                            }
                            if (j == rt.size()) break;
                        } else {
                            rightPart next_sym = rt[j];
                            //cout << "next sym is: " << next_sym.val << endl;
                            if (follow_set[right.val][0] == "") {
                                follow_set[right.val].erase(follow_set[right.val].begin() + 0);
                            } 
                            if (next_sym.type == 2) {
                                follow_set[right.val].push_back(next_sym.val);
                                break;
                            }    
                            if (next_sym.type == 1) {
                                //cout << "rule: nterm -> nterm nterm" << endl;
                                vector<string> vecf = first_one_set[next_sym.val];
                                //cout << "first set is" << endl;
                                //for (auto v : vecf) {
                                //    cout << v << " ";
                                //}
                                // cout << endl;
                                
                                for (auto v : vecf) {
                                    if(v != "eps") {
                                        follow_set[right.val].push_back(v);
                                    }
                                }
                                if (!isInGenNterms("eps", first_one_set[next_sym.val])) {
                                    break;
                                }
                            }
                            //|| isInGenNterms("eps", first_one_set[rt[j].val])

                        }
                            
                    }
                    //unique and sort
                    sort(follow_set[right.val].begin(), follow_set[right.val].end());
                    auto last = unique(follow_set[right.val].begin(), follow_set[right.val].end());
                    follow_set[right.val].erase(last, follow_set[right.val].end());

                    int newsetSize = follow_set[right.val].size();
                    //cout << "newsize: " << newsetSize << endl;
                    //printFollowSet();
                    if (newsetSize != setSize) {
                        changed = true;
                    }
                }
            }   
        }
    }
}

void findNtermsWithEpsinGenTerm() {
    for(auto n : nonTerms) {
        for (auto rule : grammar) {
            if (rule.left == n) {
                vector<rightPart> rt = rule.right;
                for (auto r : rt) {
                    if(r.val == "") {
                        hasEpsInGenTerms.push_back(n);
                    }
                }
            }
        }
    }
    sort(hasEpsInGenTerms.begin(), hasEpsInGenTerms.end());
    auto end = unique(hasEpsInGenTerms.begin(), hasEpsInGenTerms.end());
    hasEpsInGenTerms.erase(end, hasEpsInGenTerms.end());
}

void constructFirstk(int k) {
    int setSize = 0;
    // for (auto n : nonTerms) {
    //     first_k_set[n].push_back("");
    // }
    bool changed = true;
    while(changed) {
        changed = false;
        for(auto n : nonTerms) {
            // cout << "nonterm: " << n << endl;
            setSize = first_k_set[n].size();
            // cout << "setsize: " << setSize << endl;
            if(isInGenNterms(n, hasEpsInGenTerms)) {
                first_k_set[n].push_back("eps");
            }
            for (auto rule : grammar) {
                if (rule.left == n) { 
                    vector<string> mightbe;
                    mightbe.push_back("");
                    vector<rightPart> rt = rule.right;
                    for(int i = 0; i < rt.size(); i++) {
                        sort(mightbe.begin(), mightbe.end());
                        auto end = unique(mightbe.begin(), mightbe.end());
                        mightbe.erase(end, mightbe.end());
                        rightPart right = rt[i];
                        // cout << "rule right is: " << right.val << endl;
                        if (right.type == 1) {
                            // cout << "neterm" << endl;
                            vector<string> buf = mightbe;
                            // cout << "buf size: " << buf.size() << endl;                   
                            for(int i = 0; i < buf.size(); i++) {
                                // если следующий нетерминал может перевести по правилу в непустую строку, 
                                // удаляем элемент, добавленный на предыдущем шаге в возможный к-префикс
                                findNtermsWithEpsinGenTerm();
                                if(!isInGenNterms(right.val, hasEpsInGenTerms)) {  
                                    mightbe.erase(mightbe.begin() + i);         
                                }
                                // cout << "updated mightbe: ";
                                // for (auto b : mightbe) {
                                //     cout << b << " ";
                                // }
                                // cout << endl;
                                vector<string> bufFirstK = first_k_set[right.val];
                                // cout << "copied first K set: ";
                                // for(int i = 0; i < bufFirstK.size(); i++) {
                                //     cout << bufFirstK[i] << " ";
                                // }
                                // cout << endl;
                                for ( int j = 0; j < bufFirstK.size(); j++) {
                                    // cout << "EL n:"<< bufFirstK[j] << endl;
                                    if (bufFirstK[j] == "eps") continue;
                                    string test = buf[i] + bufFirstK[j];
                                    // cout << "test: " << test << endl;
                                    if (test.length() == k) {
                                        first_k_set[n].push_back(test);
                                    } else if (test.length() > k) {
                                        first_k_set[n].push_back(test.substr(0,k));
                                    } else {
                                        mightbe.push_back(test);
                                        for (int s = 0; s < mightbe.size(); s++) {
                                            if (mightbe[s] == "") 
                                                mightbe.erase(mightbe.begin()+ s);
                                        }
                                    }
                                    
                                }
                            }
                        }
                        if (right.type == 2) {
                            // cout << "term" << endl;
                            vector<string> buf = mightbe;
                            // cout << "buf now ";
                            // for(int i = 0; i < buf.size(); i++) {
                            //     cout << buf[i] << " ";
                            // }
                            // cout << endl;
                            for(int i = 0; i < buf.size(); i++) {
                                // чистим возможный к-префикс
                                mightbe.erase(mightbe.begin()+ i);
                                // cout << "EL t:"<< buf[i] << endl;
                                string test = buf[i] + right.val;
                                // cout << "test is now: " << test << endl;
                                if (test.length() == k) {
                                    first_k_set[n].push_back(test);
                                } else {
                                    mightbe.push_back(test);
                                    for (int s = 0; s < mightbe.size(); s++) {
                                        if (mightbe[s] == "") 
                                            mightbe.erase(mightbe.begin()+ s);
                                    }
                                }
                            }
                        }    
                        if (mightbe.size() == 0) break;
                    }
                    // cout << "MIGHTBE K PREFIX: ";
                    // for (auto c : mightbe) {
                    //     cout << c << " ";
                    // }
                    // cout << endl;
                    for (auto c : mightbe) {
                        first_k_set[n].push_back(c);
                    }
                    sort(first_k_set[n].begin(), first_k_set[n].end());
                    auto lt = unique(first_k_set[n].begin(), first_k_set[n].end());
                    first_k_set[n].erase(lt, first_k_set[n].end());
                }     
            }
            // cout << "set now" << endl;
            // printFirstkSet(k);
            if (first_k_set[n].size() != setSize)
            changed = true;    
        }
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
    //printTerms();
    updateGrammar();
    printGrammar();

    //убираем недостижимые нетерминалы
    cout << "> Grammar with removed unreachable nonterminals <" << endl;
    removeUnreachableNterms();
    //printTerms();
    updateGrammar();
    printGrammar();
    cout << "> FIRST 1 sets for nonterminals <" << endl;
    constructFirst1();
    printFirst1Set();
    cout << "> FOLLOW sets for nonterminals <" << endl;
    constructFollow();
    printFollowSet();
    cout << "> FIRST k sets for nonterminals <" << endl;
    constructFirstk(2);
    printFirstkSet(2);
    return 0;
}

// ааааааааа (крик)
// хотела бы я сказать, что сиплюсплюснулась настолько, чтобы понять 3 лабораторную, но пока нет...
