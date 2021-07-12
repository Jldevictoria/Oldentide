# Oldentide Style Guide

## Go

### USE "go fmt *.go" for all your go code.

* **Boilerplate** Every code-containing file should have the same boilerplate at the top of the file.  
  The boilerplate should include the file name, creation date, author's name, and purpose (in that order).
```Go
    // File:     user.go
    // Created:  January 1, 2016
    // Author:   John Smith
    // Purpose:  A simple package used for manipulating user information.
```

* **Files** should be all lowercase, have the .go suffix, and be stored in an appropriate src/folder\_name location. (named after their namespace and purpose)
  underscores are acceptable!
```Go
    server.go should be found in Oldentide/Server/src/server/server.go
```

* **Functions** should be lowerCamelCase.
```Go
    func addTwo(x int, y int) int {
        return x + y
    }
```

* **Packages** should be all lowercase, have the .go suffix, and be stored in an appropriate src/folder\_name location. (named after their namespace and purpose)
* **Package Methods** should be UpperCamelCase.
```Go
    // in src/foo/foo.go
    package foo

    func AddTwo(x int, y int) int {
        return x + y
    }
```

* **Structs** should typically be all lowercase.
* **Comments** should all be implemented using line comments ("//").
```Go
    // This is a description for the foo struct.
    // This is the way to do a multi-line comment.
    type foo struct {
        bar int    // This is a description for the bar variable.
        baz byte   // This is a description for the baz variable.
    }
```

* **Variables** should be named using lowerCamelCase or under\_score notation.
```Go
    phoneNumber := 9251234567
    phone_number := 9251234567
    fullHomeAddress := "123 Sesame Street, Sacramento, CA 94203"
    full_home_address := "123 Sesame Street, Sacramento, CA 94203"
```

* **Braces** should always open on the same line as the declaration, with a space between any other character.  
  Closing braces should be separated from other characters by a space.  If multi-lined, they should be aligned with the start of the declaration on their own line.
```Go
    func foo(bar int) int {
        return bar
    }

    baz := []string { "a", "b", "c" }
    
    baz := []string {
        "a",
        "b",
        "c",
    }
```

* **Brackets and Parenthases** should touch at least one other character on the outside and have no spaces touching their insides.
```Go
    names := []int { 1, 2, 3, 4, 5 }
    fmt.Println(names[3])
```

* **if-else statements** are required to be minimize bracket space by sharing declarations with brackets.
```Go
    if (foo) {
        return bar
    } else if (qux){
        return baz
    } else {
        return quux
    }
```

## C# 

* **To be implemented**

## Python

* **To be implemented**