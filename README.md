# Banana language


This repository contains an interpreter for the "Banana" programming language, as a variation of "Monkey" programming language.


### 0. Also Inspired By
* python
* javascript
* golang
* java



## 1. Installation

you'll need to compile this project with go version 1.16 or higher.

You can install from source like:

```markdown
git clone https://github.com/proost/banana
cd banana
go install
```



## 2. Syntax
### 2.1 Variables
Variables are defined using the `let` keyword.
```markdown
let a = 1
let b = 2.3
```

Using `;`, multiple variables can be defined in one-line.
```markdown
let a = 1; let b = 2.3
```

Variables may be integers, floats, strings, boolean, null or arrays/hashes (which are discussed later).

Variables may be updated `let`. 
```markdown
let welcome = "Hello";
welcome = "Hello, world";
print(welcome);
```



### 2.2 Arithmetic operations
`banana` supports all the basic arithmetic operation of `int` and `float` types.

The `int` type is represented by `int64` and `float` type is represented by `float64`.

```markdown
let a = 2
let b = 3.5

print(a+b) // 5.5
print(a-b) // -1.5
print(a*b) // 7.0
print(b/a) // 1.75
print(b%a) // 1.5
```



### 2.3 Bitwise operations
`banana` supports all the basic bitwise operation of a `int` type.

```markdown
let a = 4
let b = 7

print(a&b) // 4
print(a|b) // 7
print(a^b) // 3
print(a >> 2) // 1
print(a << 2) // 64
```



### 2.4 Builtin containers
#### 2.4.1 Array
An array is a list which organizes items by linear sequence. An array can hold multiple types.

```markdown
let a = [1, 2.3, "array", true]
```

Using `append` function, adding a new element to array is done.
```markdown
let a = [1, 2.3, "array", true]
let b = append(a, false)
print(a)  //  [1, 2.3, "array", true]
print(b) // [1, 2.3, "array", true, false]
```

Using `range` function, you can generate array
```markdown
let a = range(1,5) // [1,2,3,4]
let b = range(1,7,2)  // [1,3,5]
let c = range(5,1,-1) // [5,4,3,2]
```

You can iterate over the elements of an array like:
```markdown
let a = [1, 2.3, "array"];
for el in a { print("Array contains ", el, "\n") } 

shows:
// Array contains 1
// Array contains 2.300000
// Array contains array
```

Also, You can iterate over the elements of an array with index(Like, `enumerate` in python or `for` in golang:
```markdown
let a = [1, 2.3, "array"];
for i,el in a { print("Array contains ", el, " at index ", i, "\n") }

shows:
// Array contains 1 at index 0
// Array contains 2.300000 at index 1
// Array contains array at index 2
```



#### 2.4.2 Hash
A hash is a key/value container, every object can be a key.
```markdown
let a = {"name": "banana", true: 1, 2: "two", {}: "hash"}

print(a) // {name: banana, true: 1, 2: two, {}: hash}
```

Updating a hash is done via indexing, changing in-place:
```markdown
let a = {"name": "banana", true: 1, 2: "two", {}: "hash"}

a[[1]] = "arr"

print(a) // {true: 1, 2: two, [1]: arr, {}: hash, name: banana}
```

You can iterate over the hash, using `for` loop. 
```markdown
let a = {"name": "banana", true: 1, 2: "two", {}: "hash"}

for k,v in a { print(k, ": ", v) }

shows: 
// name: banana
// 2: two
// {}: hash
// true: 1
```

If you want iterate over the hash's keys, use one variable for `for`.
for example,
```markdown
let a = {"name": "banana", true: 1, 2: "two", {}: "hash"}

for k in a { print(k) }

shows: 
// {}
// name
// true
// 2
```

Using `delete`, you can delete key from hash.
for example,
```markdown
let a = {"name": "banana", true: 1, 2: "two", {}: "hash"}

delete(a, "name")

print(a) // {true: 1, 2: two, {}: hash}
```



### 2.5 Builtin Functions
* `len`: return the length of builtin containers, string
    > len("abc") // 3
* `first`: return the first element of array
    >  first([1,2,3]) // 1
* `last`: return the last element of array
    > last([1,2,3]) // 3
* `append`: return a new array, after add an item to end of array
    > append([1,2,3], 4) // [1,2,3,4]
* `print`: print out argument to stdout
    > print("abc", "d") // abcd
* `type`: return type object of argument and print out type
    > type(1) // Type: INTEGER
* `range`: return array [start, end), 3rd argument is interval
    > range(0, -10, -3) // [0, -3, -6, -9]
* `delete`: remove key from hash



### 2.6 Function



