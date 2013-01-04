package ed

import(
    "unicode/utf8"
//    "fmt"
)

func min(a, b int) int {
    if a < b {
        return a
    } // if
    
    return b
}

func EditDistance(a, b string) int {
    lb := utf8.RuneCountInString(b)

    f := make([]int, lb + 1)
    
    for j := 0; j <= lb; j ++ {
        f[j] = j
    } // for i
    
    for _, ca := range(a) {
        j := 1
        fj1 := f[0] // the value of f[j - 1] in last iteration
        f[0] ++
        for _, cb := range(b) {
            mn := min(f[j] + 1, f[j - 1] + 1) // delete & insert
            if cb != ca {
                mn = min(mn, fj1 + 1) // change
            } else {
                mn = min(mn, fj1) // matched
            } // else
            
            fj1, f[j] = f[j], mn // save f[j] to fj1(j is about to increase), update f[j] to mn
            j ++
        } // for cb
    } // for ca
    
    return f[lb]
}
