package ed

import(
    "runtime"
    "strings"
    "fmt"
    "testing"
)

func o_() string {
    pc, _, _, _ := runtime.Caller(1)
    name := runtime.FuncForPC(pc).Name()
    if p := strings.LastIndexAny(name, `./\`); p >= 0 {
        name = name[p+1:]
    } // if
    fmt.Println("== BEGIN", name, "===")
    return name
}

func __(name string) {
    fmt.Println("== END", name, "===")
}

func TestString(t *testing.T) {
    defer __(o_())
    
    cases := [][]interface{}{
        { "abcd",  "bcde", 2},
        {"abcde",      "", 5},
        {     "", "abcde", 5},
        {     "",      "", 0},
        {"abcde", "abcde", 0},
        {"abcde", "dabce", 2}}

    for _, cs := range(cases) {
        d := String(cs[0].(string), cs[1].(string))
        if d != cs[2].(int) {
            t.Errorf("Edit-distance between %s and %s is expected to be %d, but %d got!", cs[0], cs[1], cs[2], d)
        } // if
    } // for case
}

type stringInterface struct {
    a, b []rune
}

func (in *stringInterface) LenA() int {
    return len(in.a)
}

func (in *stringInterface) LenB() int {
    return len(in.b)
}

func (in *stringInterface) CostOfChange(iA, iB int) int {
    if in.a[iA] == in.b[iB] {
        return 0
    } // if
    
    return 100
}

func (in *stringInterface) CostOfDel(iA int) int {
    return 100 + iA
}

func (in *stringInterface) CostOfIns(iB int) int {
    return 110 + iB
}

func TestEditDistance(t *testing.T) {
    defer __(o_())

    cases := [][]interface{}{
        { "abcd",  "bcde", 213},
        {"abcde",      "", 510},
        {     "", "abcde", 560},
        {     "",      "", 0},
        {"abcde", "abcde", 0},
        {"abcde", "dabce", 213}}

    for _, cs := range(cases) {
        d := EditDistance(&stringInterface{[]rune(cs[0].(string)), []rune(cs[1].(string))})
        if d != cs[2].(int) {
            t.Errorf("Edit-distance between %s and %s is expected to be %d, but %d got!", cs[0], cs[1], cs[2].(int), d)
        } // if
    } // for case
}
