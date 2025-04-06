package main

import Al.core.fmt as fmt

main := proc() {
    fmt.Println("Hello World!")

    fmt.Println(f"2 + 3.1415 = {add(2, 3.1415)}")
    
}

add := func(a, b: int) int {
    return a + b
}
