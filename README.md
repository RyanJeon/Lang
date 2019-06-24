# RyanLang

##### Warning: RyanLang is in a very crude stage. Also I am also accepting a good name for this language. As the language is in a development stage, keywords will be added as the language becomes more useful!

# Info

RyanLang is a programming langauge that supports both Korean and English keywords to allow flexibility in terms of spoken language when writing code. 

#### Example .rlang code
```
Int foo => Int c , Int b
    Int grub = 10
    print ( c + grub )
End

Int execute =>
    Int a = 15 + 32
    foo ( a + 2 , 5 )
End
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
- Korean support for delclaration

### Some Problems
- What about a conflict between variable and keywords? For example, someone who wants to write in Korean decides to declare a variable in English, does that person have to be aware of English keywords as well?:
    - In the future, the user will specify the code language at compile time
    to avoid such conflict.