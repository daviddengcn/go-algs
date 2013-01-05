/*
ed package provides some types for edit-distance calculation.
*/
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

/*
String calculates the edit-distance between two strings. Input strings must be UTF-8 encoded.

The time complexity is O(mn) where m and n are lengths of a and b, and space complexity is O(n).
 */
func String(a, b string) int {
    f := make([]int, utf8.RuneCountInString(b) + 1)
    
    for j := range(f) {
        f[j] = j
    } // for i
    
    for _, ca := range(a) {
        j := 1
        fj1 := f[0] // fj1 is the value of f[j - 1] in last iteration
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
    
    return f[len(f) - 1]
}

/*
Interface defines a pair of list and the cost for each operation.
*/
type Interface interface {
    // LenA returns the lenght of the source list
    LenA() int
    
    // LenB returns the lenght of the destination list
    LenB() int
    
    // CostOfChange returns the change cost from an item in the source list at iA to an item in the destination list at iB
    CostOfChange(iA, iB int) int
    
    // CostOfDel returns the cost of deleting an item in the source list at iA
    CostOfDel(iA int) int
    
    // CostOfIns returns the cost of inserting an item in the destination list at iB
    CostOfIns(iB int) int
}

/*
EditDistance returns the edit-distance defined by Interface.

The time complexity is O(mn) where m and n are lengths of a and b, and space complexity is O(n).
*/
func EditDistance(in Interface) int {
    la, lb:= in.LenA(), in.LenB();
    
    f := make([]int, lb + 1)
    
    for j := 1; j <= lb; j ++ {
        f[j] = f[j - 1] + in.CostOfIns(j - 1)
    } // for i
    
    for i := 0; i < la; i ++ {
        fj1 := f[0] // fj1 is the value of f[j - 1] in last iteration
        f[0] += in.CostOfDel(i)
        for j := 1; j <= lb; j ++ {
            mn := min(f[j] + in.CostOfDel(i), f[j - 1] + in.CostOfIns(j - 1)) // delete & insert
            mn = min(mn, fj1 + in.CostOfChange(i, j - 1)) // change/matched
            
            fj1, f[j] = f[j], mn // save f[j] to fj1(j is about to increase), update f[j] to mn
        } // for cb
    } // for ca
    
    return f[lb]
}
