/*
ed package provides some functions for generalized edit-distance calculation.

String calculates the standard edit-distance, and StringFull returns an extra longest-common-string.

EditDistance calculates the generalized edit-distance, and EditDistanceFull returns extra matching infomation.
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
String calculates the edit-distance and longest-common-string between two strings. Input strings must be UTF-8 encoded.

The time and space complexity are all O(mn) where m and n are lengths of a and b.
 */
func StringFull(a, b string) (dist int, lcs string) {
    la, lb := utf8.RuneCountInString(a), utf8.RuneCountInString(b)
    f := make([]int, lb + 1)
    ops := make([]byte, la * lb)
    
    const(
        DEL byte = iota
        INS
        CHANGE
        MATCH
    )
    
    for j := range(f) {
        f[j] = j
    } // for i
    
    p := 0; // the index to ops
    
    for _, ca := range(a) {
        j := 1
        fj1 := f[0] // fj1 is the value of f[j - 1] in last iteration
        f[0] ++
        for _, cb := range(b) {
            mn, op := f[j] + 1, DEL
            if f[j - 1] + 1 < mn {
                mn, op = f[j - 1] + 1, INS
            }// if

            if cb != ca {
                if fj1 + 1 < mn {
                    mn, op = fj1 + 1, CHANGE
                } // if
            } else {
                    mn, op = fj1, MATCH
            } // else
            
            fj1, f[j], ops[p] = f[j], mn, op // save f[j] to fj1(j is about to increase), update f[j] to mn
            p ++
            j ++
        } // for cb
    } // for ca
    
    // Calculate longest-common-string
    lcsi := make([]int, 0, la)
    for i, j := la, lb; i > 0 || j > 0; {
        var op byte
        if i == 0 {
            op = INS
        } else if j == 0 {
            op = DEL
        } else {
            op = ops[(i - 1)*lb + j - 1]
        } // else
        
        switch op {
            case INS:
                j --
            case DEL:
                i --
            case CHANGE, MATCH:
                i --
                j --
                if op == MATCH {
                    lcsi = append(lcsi, i)
                } // if
        } // switch
    } // for i, j

    if i, l := 0, len(lcsi) - 1; l > 0 {
        for _, ca := range(a) {
            if i == lcsi[l] {
                lcs = lcs + string(ca)
                l --
                if l < 0 {
                    break
                } // if
            } // if
            i ++
        } // for ca
    } // for  l
    
    return f[len(f) - 1], lcs
}

/*
Interface defines a pair of lists and the cost for each operation.
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

/*
EditDistanceFull returns the edit-distance and corresponding match indexes defined by Interface.
Each element in matA and matB is the index in the other list, if it is equal to or greater than zero; or -1 meaning a deleting or inserting in matA or matB respectively.

The time and space complexity are all O(mn) where m and n are lengths of a and b.
*/
func EditDistanceFull(in Interface) (dist int, matA, matB []int) {
    la, lb:= in.LenA(), in.LenB();
    
    f := make([]int, lb + 1)
    ops := make([]byte, la * lb)
    
    for j := 1; j <= lb; j ++ {
        f[j] = f[j - 1] + in.CostOfIns(j - 1)
    } // for i
    
    const(
        DEL byte = iota
        INS
        CHANGE
    )
    
    // Matching with dynamic programming
    p := 0
    for i := 0; i < la; i ++ {
        fj1 := f[0] // fj1 is the value of f[j - 1] in last iteration
        f[0] += in.CostOfDel(i)
        for j := 1; j <= lb; j ++ {
            mn, op := f[j] + in.CostOfDel(i), DEL // delete
            
            if v := f[j - 1] + in.CostOfIns(j - 1); v < mn {
                // insert
                mn, op = v, INS
            } // if
            
             // change/matched
            if v := fj1 + in.CostOfChange(i, j - 1); v < mn {
                // insert
                mn, op = v, CHANGE
            } // if
            
            fj1, f[j], ops[p] = f[j], mn, op // save f[j] to fj1(j is about to increase), update f[j] to mn
            p ++
        } // for cb
    } // for ca
    // Reversely find the match info
    matA, matB = make([]int, la), make([]int, lb)
    for i, j := la, lb; i > 0 || j > 0; {
        var op byte
        if i == 0 {
            op = INS
        } else if j == 0 {
            op = DEL
        } else {
            op = ops[(i - 1)*lb + j - 1]
        } // else
        
        switch op {
            case INS:
                j --
                matB[j] = -1
            case DEL:
                i --
                matA[i] = -1
            case CHANGE:
                i --
                j --
                matA[i], matB[j] = j, i
        } // switch
    } // for i, j
    
    return f[lb], matA, matB
}
