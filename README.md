# Jennifer

Jennifer is a code generator for Go:

```go
package main

import (
	"fmt"

	. "github.com/davelondon/jennifer/jen"
)

func main() {
	f := NewFile("main")
    f.Func().Id("main").Params().Block(
        Id("fmt.Println").Call(
            Lit("Hello, world"),
        ),
    )
    fmt.Printf("%#v", f)
}
```

Output:

```go
package main

import fmt "fmt"

func main() {
    fmt.Println("Hello, world")
}
```

# Imports

Jennifer manages your imports and aliases:

```go
package main

import (
	"fmt"

	. "github.com/davelondon/jennifer/jen"
)

func main() {
	f := NewFilePath("a.b/c")
    f.Func().Id("init").Params().Block(
        Id("a.b/c.Local").Call(),
        Id("d.e/f.Remote").Call(),
        Id("g.h/f.Collision").Call(),
    )
    fmt.Printf("%#v", f)
}
```

Output:

```go
package c

import (
    f "d.e/f"
    f1 "g.h/f"
)

func init() {
    Local()
    f.Remote()
    f1.Collision()
}
```

# Examples
The tests are written mostly as examples - [see godoc.org](https://godoc.org/github.com/davelondon/jennifer/jen#pkg-examples) for an index.

Most of the code is generated using jennifer itself, see the [genjen package](https://github.com/davelondon/jennifer/tree/master/genjen) for a real-world example of usage - it generates [generated.go](https://github.com/davelondon/jennifer/blob/master/jen/generated.go).

# Rendering
For testing, a `File` or `Statement` can be rendered with the `fmt` package:

```go
c := Id("a").Call(Lit("b"))
fmt.Printf("%#v", c)
// Output: a("b")
```

This is not recommended for use in production because any error will cause a 
panic. For production use, `File.Render` or `File.Save` are preferred.

# Id
`Id` renders an identifier. For a local identifier, simply use a string:
 
```go
c := Id("a")
fmt.Printf("%#v", c)
// Output: a
```

For a remote identifier, prefix with the full package path:

```go
c := Id("encoding/gob.NewEncoder").Call()
fmt.Printf("%#v", c)
// Output: gob.NewEncoder()
```

The imports are automatically handled when used with a `File`, and if the path
matches the local path, the package alias is omitted from the rendered code.

To access fields, more items may be added to the `Id` method:

```go
c := Id("a", "b", "c")
fmt.Printf("%#v", c)
// Output: a.b.c
```

This can be combined with the remote syntax:

```go
c := Id("a.b/c.Foo", "Bar", "Baz")
fmt.Printf("%#v", c)
// Output: c.Foo.Bar.Baz
```

More complex chains can be formed by using Code items instead of strings:

```go
c := Id("a.b/c.Foo", Id("Bar").Call(), "Baz")
fmt.Printf("%#v", c)
// Output: c.Foo.Bar().Baz
```

More control over the package import can be gained by using the `Alias` method 
to specify the remote package:
 
```go
c := Id(Alias("a.b/c"), Id("Foo").Call(), "Bar")
fmt.Printf("%#v", c)
// Output: c.Foo().Bar
```

# Op
`Op` renders the provided string. Use for operators and tokens:

```go
c := Id("a").Op(":=").Id("b").Call()
fmt.Printf("%#v", c)
// Output: a := b()
```

```go
c := Op("*").Id("a")
fmt.Printf("%#v", c)
// Output: *a
```

```go
c := Id("a").Call(Id("b").Op("..."))
fmt.Printf("%#v", c)
// Output: a(b...)
```

# Identifiers 

Identifiers are simple methods with no parameters. They render as the 
identifier token:

```go
c := Break()
fmt.Printf("%#v", c)
// Output: break
```

Keywords: `Break`, `Default`, `Func`, `Select`, `Defer`, `Go`, `Struct`, `Chan`, `Else`, `Goto`, `Const`, `Fallthrough`, `Range`, `Type`, `Continue`, `Var`

Built-in types: `Bool`, `Byte`, `Complex64`, `Complex128`, `Error`, `Float32`, `Float64`, `Int`, `Int8`, `Int16`, `Int32`, `Int64`, `Rune`, `String`, `Uint`, `Uint8`, `Uint16`, `Uint32`, `Uint64`, `Uintptr`

Constants: `True`, `False`, `Iota`, `Nil`

Also included is `Err` for the commonly used `err` variable.

Note: `Interface`, `Map`, `Return`, `Switch`, `For`, `Case` and `If` are special cases, 
and treated as groups - see below.

Note: The `import` and `package` keywords are always rendered automatically, so 
not included.

# Built-in functions

Built in functions take a variadic list of code items, and render the function
name followed by the items as a comma seperated list of parameters in 
parenthesis:

```go
c := Append(Id("a"), Id("b"))
fmt.Printf("%#v", c)
// Output: append(a, b)
```

Functions: `Append`, `Cap`, `Close`, `Complex`, `Copy`, `Delete`, `Imag`, `Len`, `Make`, `New`, `Panic`, `Print`, `Println`, `Real`, `Recover`

# Groups

Groups take either a single code item or a varidic list of code items. The 
items are rendered between open and closing tokens. Multiple items are 
seperated by a separator token.

### Groups accepting a list of items:

| Group     | Opening       | Separator | Closing | Usage                             |
| --------- | ------------- | --------- | ------- | --------------------------------- |
| List      |               | `,`       |         | `a, b := c()`                     |
| Call      | `(`           | `,`       | `)`     | `fmt.Println(b, c)`               |
| Params    | `(`           | `,`       | `)`     | `func (a *A) Foo(i int) { ... }`  |
| Values    | `{`           | `,`       | `}`     | `[]int{1, 2}` or `interface{}`    |
| Index     | `[`           | `:`       | `]`     | `a[1:2]` or `[]int{}`             |
| Block     | `{`           | `\n`      | `}`     | `func a() { ... }`                |
| Case      | `case`        | `,`       |         | `switch a {case "b", "c": ... }`  |
| CaseBlock | `:`           | `\n`      |         | `switch i {case 1: ... }`         |
| Return    | `return`      | `,`       |         | `return a, b`                     |
| If        | `if`          | `;`       |         | `if a, ok := b(); ok { ... }`     |
| For       | `for`         | `;`       |         | `for i := 0; i < 10; i++ { ... }` |
| Switch    | `switch`      | `;`       |         | `switch a { ... }`                |
| Interface | `interface {` | `\n`      | `}`     | `interface { ... }`               |

### Groups accepting a single item:

| Group  | Opening  | Closing | Usage                        |
| ------ | -------- | ------- | ---------------------------- |
| Parens | `(`      | `)`     | `[]byte(s)` or `a / (b + c)` |
| Assert | `.(`     | `)`     | `s, ok := i.(string)`        |
| Map    | `map[`   | `]`     | `map[int]string`             |

### List
`List` renders a comma seperated list with no open or closing tokens. Use for 
multiple return functions:

```go
c := List(Id("a"), Id("b")).Op(":=").Id("c").Call()
fmt.Printf("%#v", c)
// Output: a, b := c()
```

### Values
`Values` renders a comma seperated list enclosed by curly braces. Use for slice 
literals:

```go
c := Index().String().Values(Lit("a"), Lit("b"))
fmt.Printf("%#v", c)
// Output: []string{"a", "b"}
```

### Call
`Call` renders a comma seperated list enclosed by parenthesis. Use for function
calls:

```go
c := Id("a").Call(Id("b"), Id("c"))
fmt.Printf("%#v", c)
// Output: a(b, c)
```

### Params
`Params` renders a comma seperated list enclosed by parenthesis. Use for 
function parameters and method receivers:

```go
c := Func().Params(Id("a").Id("A")).Id("foo").Params(Id("b").String()).String().Block()
fmt.Printf("%#v", c)
// Output: func (a A) foo(b string) string {}
```

### Index
`Params` renders a colon seperated list enclosed by square brackets. Use for 
array / slice indexes and definitions:

```go
c := Var().Id("a").Index().String()
fmt.Printf("%#v", c)
// Output: var a []string
```

```go
c := Id("a").Op(":=").Id("b").Index(Lit(0), Lit(1))
fmt.Printf("%#v", c)
// Output: a := b[0:1]
```

```go
c := Id("a").Op(":=").Id("b").Index(Lit(1), Empty())
fmt.Printf("%#v", c)
// Output: a := b[1:]
```

### Block
`Block` renders a list of statements enclosed by curly braces. Use for all 
blocks:

```go
c := Func().Id("main").Params().Block(
    Id("a").Op("++"),
    Id("b").Op("--"),
)
fmt.Printf("%#v", c)
// Output: func main() {
//  a++
//  b--
// }
}
```

```go
c := If(Id("a").Op(">").Lit(10)).Block(
    Id("a").Op("=").Id("a").Op("/").Lit(2),
)
fmt.Printf("%#v", c)
// Output: if a > 10 {
//  a = a / 2
// }
}
```

### Switch, Case, CaseBlock
`Switch`, `Case` and `CaseBlock` can be used to build `switch` statements:

```go
c := Switch(Id("a")).Block(
    Case(Lit("1")).CaseBlock(
        Return(Lit(1)),
    ),
    Case(Lit("2"), Lit("3")).CaseBlock(
        Return(Lit(2)),
    ),
    Case(Lit("4")).CaseBlock(
        Fallthrough(),
    ),
    Default().CaseBlock(
        Return(Lit(3)),
    ),
)
fmt.Printf("%#v", c)
// Output: switch a {
// case "1":
// 	return 1
// case "2", "3":
// 	return 2
// case "4":
// 	fallthrough
// default:
// 	return 3
// }
```

### Return
`Return` renders the `return` keyword followed by a comma separated list of 
Code items:

```go
c := Return(Id("a"), Id("b"))
fmt.Printf("%#v", c)
// Output: return a, b
```

### If
`If` renders the `if` keyword followed by a semicolon separated list of Code 
items:

```go
c := If(Err().Op(":=").Id("a").Call(), Err().Op("!=").Nil()).Block(
    Return(Err()),
)
fmt.Printf("%#v", c)
// Output: if err := a(); err != nil {
//  return err
// }
```

### For
`For` renders the `for` keyword followed by a semicolon separated list of Code 
items:

```go
c := For(Id("i").Op(":=").Lit(0), Id("i").Op("<").Lit(10), Id("i").Op("++")).Block(
    Id("fmt.Println").Call(Id("i")),
)
fmt.Printf("%#v", c)
// Output: for i := 0; i < 10; i++ {
//  fmt.Println(i)
// }
```

### Interface
`Interface` renders the interface keyword followed by a statement block:

```go
c := Var().Id("a").Interface()
fmt.Printf("%#v", c)
// Output: var a interface{}
```

```go
c := Type().Id("a").Interface(
    Id("b").Params().String(),
)
fmt.Printf("%#v", c)
// Output: type a interface {
// 	b() string
// }
```

### Parens
`Parens` renders a single code item in parenthesis. Use for type conversion or 
to specify evaluation order:

```go
c := Id("b").Op(":=").Index().Byte().Parens(Id("s"))
fmt.Printf("%#v", c)
// Output: b := []byte(s)
```

```go
c := Id("a").Op("/").Parens(Id("b").Op("+").Id("c"))
fmt.Printf("%#v", c)
// Output: a / (b + c)
```

### Assert
`Assert` renders a period followed by a single Code item enclosed by 
parenthesis. Use for type assertions:

```go
c := Id("a").Op("=").Id("b").Assert(String())
fmt.Printf("%#v", c)
// Output: a = b.(string)
```

### Map
`Map` renders the `map` keyword followed by a single Code item enclosed by 
square brackets. Use for map definitions:

```go
c := Id("a").Op(":=").Map(String()).String().Values()
fmt.Printf("%#v", c)
// Output: a := map[string]string{}
```

### Alternate GroupFunc methods
All the Group functions have GroupFunc functions that accept a `func(*Group)`. 
Use this for embedding logic:

```go
increment = true
c := Func().Id("a").Params().BlockFunc(func(g *Group) {
    if increment {
        g.Id("a").Op("++")
    } else {
        g.Id("a").Op("--")
    }
})
fmt.Printf("%#v", c)
// Output: func a() {
// 	a++
// }
```

# Add
`Add` adds the provided Code to the Statement.
 
```go
ptr := Op("*")
c := Id("a").Op("=").Add(ptr).Id("b")
fmt.Printf("%#v", c)
// Output: a = *b
```

# Do
`Do` takes a `func(*Statement)` and executes it on the current statement. This 
is useful for embedding logic:

```go
f := func(name string, isMap bool) *Statement {
    return Id(name).Op(":=").Do(func(s *Statement) {
        if isMap {
            s.Map(String()).String()
        } else {
            s.Index().String()
        }
    }).Values()
}
fmt.Printf("%#v\n%#v", f("a", true), f("b", false))
// Output: a := map[string]string{}
// b := []string{}
```

# Lit, LitFunc
`Lit` renders a literal, using the format provided by `fmt.Sprintf("%#v", ...)`.

TODO: This probably isn't good enough for all cases. 

```go
c := Id("a").Op(":=").Lit("a")
fmt.Printf("%#v", c)
// Output: a := "a"
```

```go
c := Id("a").Op(":=").Lit(1.5)
fmt.Printf("%#v", c)
// Output: a := 1.5
```

# Dict, DictFunc
`Dict` takes a `map[Code]Code` and renders a list of colon separated key value 
pairs, enclosed in curly braces. Use for map literals:

```go
c := Id("a").Op(":=").Map(String()).String().Dict(map[Code]Code{
    Lit("a"): Lit("b"),
    Lit("c"): Lit("b"),
})
fmt.Printf("%#v", c)
// Output: a := map[string]string{
// 	"a": "b",
// 	"c": "d",
// }
```

`DictFunc` does the same by executing the provided `func(map[Code]Code)`:

```go
c := Id("a").Op(":=").Map(String()).String().DictFunc(func(m map[Code]Code) {
    m[Lit("a")] = Lit("b")
    m[Lit("c")] = Lit("d")
})
fmt.Printf("%#v", c)
// Output: a := map[string]string{
// 	"a": "b",
// 	"c": "d",
// }
```

# Tag
`Tag` renders a struct tag:

```go
c := Type().Id("foo").Struct().Block(
    Id("A").String().Tag(map[string]string{"json": "a"}),
    Id("B").Int().Tag(map[string]string{"json": "b", "bar": "baz"}),
)
fmt.Printf("%#v", c)
// Output: type foo struct {
// 	A string `json:"a"`
// 	B int    `json:"b" bar:"baz"`
// }
```

# Null, Empty
`Null` adds a null item. Null items render nothing and are not followed by a 
separator in lists.
`Empty` adds an empty item. Empty items render nuothing but are followed by a 
separator in lists.

```go
c := Id("a").Op(":=").Id("b").Index(Null(), Lit(1))
fmt.Printf("%#v", c)
// Output: a := b[1]
```

```go
c := Id("a").Op(":=").Id("b").Index(Empty(), Lit(1))
fmt.Printf("%#v", c)
// Output: a := b[:1]
```

# Line
`Line` inserts a blank line.

# Comment, Commentf
`Comment` adds a comment. If the provided string contains a newline, the 
comment is formatted in multiline style:

```go
c := Comment("a")
fmt.Printf("%#v", c)
// Output: // a
```

```go
c := Comment("a\nb")
fmt.Printf("%#v", c)
// Output: /*
// a
// b
// */
```

```go
c := Id("a").Call().Comment("b")
fmt.Printf("%#v", c)
// Output: a() // b
```

`Commentf` accepts a format string and a list of parameters:

```go
c := Commentf("a %d", 1)
fmt.Printf("%#v", c)
// Output: // a 1
```

```go
c := Id("a").Call().Commentf("b %d", 1)
fmt.Printf("%#v", c)
// Output: a() // b 1
```

# File

### NewFile
`NewFile` Creates a new file, with the specified package name. 

### NewFilePath
`NewFilePath` creates a new file while specifying 
the package path - the package name is inferred from the path.

### NewFilePathName
`NewFilePathName` 
additionally specifies the package name.

```go
f := NewFilePathName("a.b/c", "main")
f.Func().Id("main").Params().Block(
    Id("a.b/c.Foo").Call(),
)
fmt.Printf("%#v", f)
// Output: package main
//
// func main() {
// 	Foo()
// }
```

### PackageComment
`PackageComment` adds a comment to the very top of the file, above the 
`package` keyword:

```go
f := NewFile("c")
f.PackageComment("a")
f.PackageComment("b")
f.Func().Id("init").Params().Block()
fmt.Printf("%#v", f)
// Output: // a
// // b
// package c
//
// func init() {}
```

### Anon
`Anon` adds an anonymous import:

```go
f := NewFile("c")
f.Anon("a")
f.Func().Id("init").Params().Block()
fmt.Printf("%#v", f)
// Output:
// package c
//
// import _ "a"
//
// func init() {}
```

### PackagePrefix
If you're worried about package aliases conflicting with local variable names, 
you can set a prefix here:

```go
f := NewFile("c")
f.PackagePrefix("pkg")
f.Func().Id("main").Params().Block(
    Id("fmt.Println").Call(),
)
fmt.Printf("%#v", f)
// Output:
// package c
//
// import pkg_fmt "fmt"
//
// func main() {
// 	pkg_fmt.Println()
// }
```

### Save, Render
`Save` renders the file and saves to the filename provided. `Render` renders 
the file to the provided writer:
 
```go
f := NewFile("a")
f.Func().Id("main").Params().Block()
buf := &bytes.Buffer{}
err := f.Render(buf)
if err != nil {
    fmt.Println(err.Error())
} else {
    fmt.Println(buf.String())
}
// Output: package a
//
// func main() {}
```

# Pointers
Be careful when passing `*Statement` around. Consider the following example:

```go
caller := func(s *Statement) *Statement {
    return s.Call()
}
a := Id("a")
c := Block(
    caller(a),
    caller(a),
)
fmt.Printf("%#v", c)
// Output: {
// 	a()()
// 	a()()
// }
```

`Id("a")` returns a `*Statement`, which the `Call()` method appends to twice. To
avoid this, pass `Statement` instead of `*Statement`:

```go
caller := func(s Statement) *Statement {
    return s.Call()
}
a := *Id("a")
c := Block(
    caller(a),
    caller(a),
)
fmt.Printf("%#v", c)
// Output: {
// 	a()
// 	a()
// }
```

Here is another variation, which can't be solved by pointer indirection:

```go
a := Id("a")
c := Block(
    a.Call(),
    a.Call(),
)
fmt.Printf("%#v", c)
// Output: {
// 	a()()
// 	a()()
// }
```

Here we can prevent the double call by using `Clone` to create a new `*Statement`:  

```go
a := Id("a")
c := Block(
    a.Clone().Call(),
    a.Clone().Call(),
)
fmt.Printf("%#v", c)
// Output: {
// 	a()
// 	a()
// }
```