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

func TestEditDistance(t *testing.T) {
    defer __(o_())
    
    cases := [][]interface{}{
        {"abcd", "bcde", 2},
        {"abcde", "", 5},
        {"", "abcde", 5},
        {"", "", 0},
        {"abcde", "dabce", 2}}
    
    for _, cs := range(cases) {
        d := EditDistance(cs[0].(string), cs[1].(string))
        fmt.Println(d, cs[2].(int))
    } // for case
}
