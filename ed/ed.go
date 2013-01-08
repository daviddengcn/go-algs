/*
ed package provides some functions for generalized edit-distance calculation.

String calculates the standard edit-distance, and StringFull returns an extra longest-common-string.

EditDistance calculates the generalized edit-distance, and EditDistanceFull returns extra matching infomation. Base is a helper type using as a base type for Interface implementations.

EditDistanceF and EditDistanceFFull are similar to EditDistance and EditDistanceFull, but using parameter and functions instead of an interface. This is sometimes more easy to use. ConstCost is a helper function.
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
    } // for j
    
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

// Constants of operations
const(
    opDEL byte = iota
    opINS
    opCHANGE
    opMATCH
)
    
/*
String calculates the edit-distance and longest-common-string between two strings. Input strings must be UTF-8 encoded.

The time and space complexity are all O(mn) where m and n are lengths of a and b.

NOTE if detailed matching information is not necessary, call String instead because it needs much less memories.
 */
func StringFull(a, b string) (dist int, lcs string) {
    la, lb := utf8.RuneCountInString(a), utf8.RuneCountInString(b)
    f := make([]int, lb + 1)
    ops := make([]byte, la * lb)
    
    for j := range(f) {
        f[j] = j
    } // for j
    
    p := 0; // the index to ops
    
    for _, ca := range(a) {
        j := 1
        fj1 := f[0] // fj1 is the value of f[j - 1] in last iteration
        f[0] ++
        for _, cb := range(b) {
            mn, op := f[j] + 1, opDEL
            if f[j - 1] + 1 < mn {
                mn, op = f[j - 1] + 1, opINS
            }// if

            if cb != ca {
                if fj1 + 1 < mn {
                    mn, op = fj1 + 1, opCHANGE
                } // if
            } else {
                mn, op = fj1, opMATCH
            } // else
            
            fj1, f[j], ops[p] = f[j], mn, op // save f[j] to fj1(j is about to increase), update f[j] to mn
            p ++;  j ++
        } // for cb
    } // for ca
    
    // Calculate longest-common-string
    lcsi := make([]int, 0, la)
    for i, j := la, lb; i > 0 || j > 0; {
        var op byte
        switch {
            case i == 0:
                op = opINS
            case j == 0:
                op = opDEL
            default:
                op = ops[(i - 1)*lb + j - 1]
        } // switch
        
        switch op {
            case opINS:
                j --
            case opDEL:
                i --
            case opCHANGE, opMATCH:
                i --;  j --
                if op == opMATCH {
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
    } // if
    
    return f[len(f) - 1], lcs
}

// Ternary returns vT if cond equals true, or vF otherwise.
func Ternary(cond bool, vT, vF int) int {
    if cond {
        return vT
    } // if
    
    return vF
}

// ConstCost returns a const-function whose return value is the specified cost. This is useful for EditDistanceF and EditDistanceFFull if the del/ins cost is constant to positions.
func ConstCost(cost int) func(int) int {
    return func(int) int {
        return cost
    }
}

/*
Base is a helper type which defines LenA, LenB, CostOfDel and CostOfIns methods.
*/
type Base struct {
    LA, LB, Cost int
}

// Interface.LenA
func (b* Base) LenA() int {
    return b.LA
}

// Interface.LenB
func (b* Base) LenB() int {
    return b.LB
}

// Interface.CostOfDel
func (b* Base) CostOfDel(iA int) int {
    return b.Cost
}

// Interface.CostOfIns
func (b* Base) CostOfIns(iB int) int {
    return b.Cost
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
    } // for j

    for i := 0; i < la; i ++ {
        fj1 := f[0] // fj1 is the value of f[j - 1] in last iteration
        f[0] += in.CostOfDel(i)
        for j := 1; j <= lb; j ++ {
            mn := min(f[j] + in.CostOfDel(i), f[j - 1] + in.CostOfIns(j - 1)) // delete & insert
            mn = min(mn, fj1 + in.CostOfChange(i, j - 1)) // change/matched
            
            fj1, f[j] = f[j], mn // save f[j] to fj1(j is about to increase), update f[j] to mn
        } // for j
    } // for i
    
    return f[lb]
}

func matchingFromOps(la, lb int, ops []byte) (matA, matB []int) {
    matA, matB = make([]int, la), make([]int, lb)
    for i, j := la, lb; i > 0 || j > 0; {
        var op byte
        switch {
            case i == 0:
                op = opINS
            case j == 0:
                op = opDEL
            default:
                op = ops[(i - 1)*lb + j - 1]
        } // switch
        
        switch op {
            case opINS:
                j --
                matB[j] = -1
            case opDEL:
                i --
                matA[i] = -1
            case opCHANGE:
                i --;  j --
                matA[i], matB[j] = j, i
        } // switch
    } // for i, j
    
    return matA, matB
}

/*
EditDistanceFull returns the edit-distance and corresponding match indexes defined by Interface.
Each element in matA and matB is the index in the other list, if it is equal to or greater than zero; or -1 meaning a deleting or inserting in matA or matB, respectively.

The time and space complexity are all O(mn) where m and n are lengths of a and b.

NOTE if detailed matching information is not necessary, call EditDistance instead because it needs much less memories.
*/
func EditDistanceFull(in Interface) (dist int, matA, matB []int) {
    la, lb:= in.LenA(), in.LenB();
    
    f := make([]int, lb + 1)
    ops := make([]byte, la * lb)
    
    for j := 1; j <= lb; j ++ {
        f[j] = f[j - 1] + in.CostOfIns(j - 1)
    } // for j
    
    // Matching with dynamic programming
    p := 0
    for i := 0; i < la; i ++ {
        fj1 := f[0] // fj1 is the value of f[j - 1] in last iteration
        f[0] += in.CostOfDel(i)
        for j := 1; j <= lb; j ++ {
            mn, op := f[j] + in.CostOfDel(i), opDEL // delete
            
            if v := f[j - 1] + in.CostOfIns(j - 1); v < mn {
                // insert
                mn, op = v, opINS
            } // if
            
             // change/matched
            if v := fj1 + in.CostOfChange(i, j - 1); v < mn {
                // insert
                mn, op = v, opCHANGE
            } // if
            
            fj1, f[j], ops[p] = f[j], mn, op // save f[j] to fj1(j is about to increase), update f[j] to mn
            p ++
        } // for j
    } // for i
    // Reversely find the match info
    matA, matB = matchingFromOps(la, lb, ops)
    
    return f[lb], matA, matB
}

/*
EditDistanceFunc returns the edit-distance defined by parameters and functions.

The time complexity is O(mn) where m and n are lengths of a and b, and space complexity is O(n).
*/
func EditDistanceF(lenA, lenB int, costOfChange func(iA, iB int) int, costOfDel func(iA int) int, costOfIns func(iB int) int)int {
    la, lb:= lenA, lenB;
    
    f := make([]int, lb + 1)
    
    for j := 1; j <= lb; j ++ {
        f[j] = f[j - 1] + costOfIns(j - 1)
    } // for j

    for i := 0; i < la; i ++ {
        fj1 := f[0] // fj1 is the value of f[j - 1] in last iteration
        f[0] += costOfDel(i)
        for j := 1; j <= lb; j ++ {
            mn := min(f[j] + costOfDel(i), f[j - 1] + costOfIns(j - 1)) // delete & insert
            mn = min(mn, fj1 + costOfChange(i, j - 1)) // change/matched
            
            fj1, f[j] = f[j], mn // save f[j] to fj1(j is about to increase), update f[j] to mn
        } // for j
    } // for i
    
    return f[lb]
}

/*
EditDistanceFFull returns the edit-distance and corresponding match indexes defined by parameters and functions.
Each element in matA and matB is the index in the other list, if it is equal to or greater than zero; or -1 meaning a deleting or inserting in matA or matB, respectively.

The time and space complexity are all O(mn) where m and n are lengths of a and b.

NOTE if detailed matching information is not necessary, call EditDistance instead because it needs much less memories.
*/
func EditDistanceFFull(lenA, lenB int, costOfChange func(iA, iB int) int, costOfDel func(iA int) int, costOfIns func(iB int) int) (dist int, matA, matB []int) {
    la, lb:= lenA, lenB;
    
    f := make([]int, lb + 1)
    ops := make([]byte, la * lb)
    
    for j := 1; j <= lb; j ++ {
        f[j] = f[j - 1] + costOfIns(j - 1)
    } // for j
    
    // Matching with dynamic programming
    p := 0
    for i := 0; i < la; i ++ {
        fj1 := f[0] // fj1 is the value of f[j - 1] in last iteration
        f[0] += costOfDel(i)
        for j := 1; j <= lb; j ++ {
            mn, op := f[j] + costOfDel(i), opDEL // delete
            
            if v := f[j - 1] + costOfIns(j - 1); v < mn {
                // insert
                mn, op = v, opINS
            } // if
            
             // change/matched
            if v := fj1 + costOfChange(i, j - 1); v < mn {
                // insert
                mn, op = v, opCHANGE
            } // if
            
            fj1, f[j], ops[p] = f[j], mn, op // save f[j] to fj1(j is about to increase), update f[j] to mn
            p ++
        } // for j
    } // for i
    // Reversely find the match info
    matA, matB = matchingFromOps(la, lb, ops)
    
    return f[lb], matA, matB
}
