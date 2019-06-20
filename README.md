# RyanLang

##### Warning: RyanLang is in a very crude stage. Also I am also accepting a good name for this language.

# Info

RyanLang is a programming langauge that supports both Korean and English keywords to allow flexibility in terms of spoken language when writing code. 

#### Example .rlang code
```
Int execute ( ) {
    Int a = 10
    출력 ( a + 1 )
}
```
The same can be done through this code
```
Int execute ( ) {
    Int a = 10
    print ( a + 1 )
}
```


#### Compiling / Running .rlang file
```sh
$ ./rlang example.rlang
$ ./example.rlang.out
```

### Current Update
- Supports Functions
- Local variables
- Basic Arithmetic

### Need to be done this week
- Functions should support parameters
- Korean support for delclaration

### Some Problems
- What about a conflict between variable and keywords? For example, someone who wants to write in Korean decides to declare a variable in English, does that person have to be aware of English keywords as well?:
    - In the future, the user will specify the code language at compile time
    to avoid such conflict.