#include <iostream>
#include <fstream>
#include <string>
#include <map>
#include <vector>
#include <algorithm>
#include <cctype>

using namespace std;

vector<string> string_rules;
vector<string> nonTerms;
vector<string> terms;

struct rightPart {
    int type; //1 - nterm, 2 - term
    string val; //A, B, ..., a, b, ...
};

struct Rule {
    string left;
    vector <rightPart> right;
};

vector<Rule> grammar;

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
    cout << endl;
    cout << "----second alt ----" << endl;
    cout << str << endl;
    cout << str.length() << endl;
    int k = 0;
	for (int i = 0; i < str.length(); i++) {
        if ( str[i] == '|'  && k == 1) {
            cout << "found: " << i << endl;
            cout << "---- ----" << endl;
			return i;
        } 
        if ( str[i] == '|'  && k < 1) {
            cout << i << endl;
			k++;
        }       
	}
    cout << "there is no second alt" << endl;
    cout << "---- ----" << endl;
    return -1;
}

vector<string> getListOfAltSubstrings(string str) {
    vector<string> s;
    cout << "ALT SUBS" << endl;
    //cout << "string: "<< str << endl;
    int k = str.length();
    cout << "k: "<< k << endl;
    int i = 0;
    while( i < k) {
        //cout << "i: "<< i << endl;
        cout << "string: "<< str << endl;
        int alt = getFirstAltIndex(str);
        cout << "alt index: "<< alt << endl;
        if (alt == -1) {
            s.push_back(str);
            return s;
        } else {
            //cout << "second alt index: " << getSecondAltIndex(str) << endl;
            if (alt == k - 1 || getSecondAltIndex(str)- alt == 1) {
                cout << "case last alt" << endl;
                s.push_back(str.substr(i, 1));
                s.push_back("");
                return s;
            } else {
                cout << "alternative substring push" << endl;
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
    cout << "CON" << endl;
    cout << "string: "<< str << endl;
    int i = 0;
    int len = str.length();
    if (str == "ε") {   // особенности заиписи этой буквы "Îµ"
        len = 1;
    }
    cout << "len: " << len << endl;
    while (i < len) {
        //cout << "i: " << i << endl;
        string s = str.substr(i);
        //cout << "str: " << s << endl;
        if (s[0] == '[') {
            string n = str.substr(i+1, 1);
            //cout << "n: " << n << endl;
            // if (find(nonTerms.begin(), nonTerms.end(), n) != nonTerms.end()) {
            //     nonTerms.push_back(n);
            // }    
            subs.push_back({1, n});
            i += 3;
        } else {
            string t = str.substr(i, 1);
            if (str == "ε") {   // особенности заиписи этой буквы 
                t = "";
            }
            cout << "t: " << t << endl;
            // if ( t != "" ) {
            //     terms.push_back(t);
            // } 
            subs.push_back({2, t});
            i++;
        }
    }
    return subs;
}

vector<Rule> parseRuleLine(string str) {
    vector<Rule> rules;
    string nt = str.substr(1,1);
    //nonTerms.push_back(nt);
    string leftPartOfRule = str.substr(5);
    //cout << getFirstAltIndex(str.substr(5)) << endl;
    if (getFirstAltIndex(leftPartOfRule) == -1) {
        cout << "i am here" << endl;
        vector<rightPart> subs;
        subs = parseCon(leftPartOfRule);
        Rule r;
        r.left = nt;
        r.right = subs;
        rules.push_back(r);
    } else {
        vector<string> subs = getListOfAltSubstrings(leftPartOfRule);
        cout << "alt subs" << endl;
        cout << "alt subs len " << subs.size() << endl;
        for (int i = 0;  i < subs.size(); i++) {
            if (subs[i] == "") {
                subs[i] = "ε";
                //cout << "ε" << endl;
            }
            cout << subs[i] << endl;
        }
        cout << "alt subs size again " << subs.size() << endl;
        for (int i = 0;  i < subs.size(); i++) {
            cout << "i " << i << endl;
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

int main(){
    int n;
    cout << "Enter test number" << endl;
    cin >> n;
    bool err = input(n);
    if (!err) {
        cout << "INCORRECT TEST FILE!";
        return 0;
    }
    for (int i = 0; i < grammar.size(); i++) {
        Rule r = grammar[i];
        cout << "-----------" << endl;
        cout << "RULE " << i + 1 << endl;
        cout << "left: " << r.left << ", ";
        cout << "right: ";
        for (int i = 0;  i < r.right.size(); i++) {
            cout << "{type:" << r.right[i].type << ", val:" << r.right[i].val << "}";
            if (i != r.right.size() - 1) cout << ", ";   
        }
        cout << endl;
    }
    /*cout << "nterms and terms" << endl;
    cout << "NTERMS: ";
    for (int i = 0;  i < nonTerms.size(); i++) cout << nonTerms[i] << " ";
    cout << endl;
    cout << "TERMS: ";
    for (int i = 0;  i < terms.size(); i++) cout << terms[i] << " ";*/
    return 0;
}