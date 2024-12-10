# rntocase

A series of files to help you /rename/ files from one type to another.

| Command              | What it does                                                            | Example                                              | Note |
|----------------------|-------------------------------------------------------------------------|------------------------------------------------------|------|
| `rnreverse`          | Reverses a file's name (excluding extension)                            | `Hello World.txt` -> `dlroW olleH.txt`               |      |
| `rntocamel`          | Renames a file to camel case                                            | `Hello World.txt` -> `HelloWorld.txt`                |      |
| `rntodarwin`         | Renames a file to darwin case                                           | `Hello World.txt` -> `Hello_ _World.txt`             |      |
| `rntodelimited`      | Renames a file changing the delimiter (think . or some other character) | `Hello World.txt` -> `hello.world.txt`               |      |
| `rntokebab`          | Renames a file to kebab case                                            | `Hello World.txt` -> `hello-world.txt`               |      |
| `rntolowerleading`   | Renames a file to lower leading letter case                             | `Hello WORLD.txt` -> `hello WORLD.txt`               |      |
| `rntosnake`          | Renames a file to snake case                                            | `Hello World.txt` -> `hello_world.txt`               |      |
| `rntotitle`          | Renames a file to title case (English)                                  | `hello of THE world.txt` -> `Hello of THE World.txt` |      |
| `rntoupperleading`   | Renames a file to upper leading letter case                             | `hello world.txt` -> `Hello world.txt`               |      |
| `rntrim`             | Renames a file removing white spaces from either side of the name       | ` Hello world .txt` -> `Hello world.txt`             |      |

* Note: this is a work in progress. I will happily change how things work so please do not use in a script (or without verification) unless you're fixing it to a particular version.

# Usage

## `rnacronym`

```bash
$ rnacronym
Error: No files provided.
Usage: rnacronym [options] <file1> [<file2> ...]

Options:
  -algorithm string
    	Choose the acronym algorithm to use, supported: gobeam. (default "gobeam")
  -dry-run
    	Display the intended changes without renaming.
  -interactive
    	Ask for confirmation before renaming each file.

Conversion Examples:
+-----------+-------------+------------+------------+-------------+
| ALGORITHM |             |            |            |             |
+-----------+-------------+------------+------------+-------------+
|           | Hello World | helloWorld | HelloWorld | Hello_WORLD |
| gobeam    | HW          | h          | H          | H           |
+-----------+-------------+------------+------------+-------------+
exit status 1
```


## `rnreverse`

```bash
$ rnreverse
Error: No files provided.
Usage: rnreverse [options] <file1> [<file2> ...]

Options:
  -algorithm string
    	Choose the reverse algorithm to use, supported: searking. (default "searking")
  -dry-run
    	Display the intended changes without renaming.
  -interactive
    	Ask for confirmation before renaming each file.

Conversion Examples:
+-----------+-------------+------------+------------+-------------+
| ALGORITHM |             |            |            |             |
+-----------+-------------+------------+------------+-------------+
|           | Hello World | helloWorld | HelloWorld | Hello_WORLD |
| searking  | dlroW olleH | dlroWolleh | dlroWolleH | DLROW_olleH |
+-----------+-------------+------------+------------+-------------+
exit status 1
```


## `rntocamel`

```bash
$ rntocamel
Error: No files provided.
Usage: rntocamel [options] <file1> [<file2> ...]

Options:
  -acronym value
    	Words to consider acronyms and not to assume they are words, ie ID => ID rather than ID => Id
  -acronym-from-file value
    	Words to consider acronyms and not to assume they are words, ie ID => ID rather than ID => Id
  -algorithm string
    	Choose the camel algorithm to use, supported: tomeh,golang-cz,lower-iancoleman,lower-searking,ettle,resenje,resenje-snake,gobeam,iancoleman,searking,go-ettle,resenje-kebab. (default "iancoleman")
  -dry-run
    	Display the intended changes without renaming.
  -interactive
    	Ask for confirmation before renaming each file.
  -word-seperators value
    	Word separators. (gobeam only) Default: "_-. "

Conversion Examples:
+------------------+-------------+-------------+-------------+-------------+
| ALGORITHM        |             |             |             |             |
+------------------+-------------+-------------+-------------+-------------+
|                  | Hello World | helloWorld  | HelloWorld  | Hello_WORLD |
| ettle            | helloWorld  | helloWorld  | helloWorld  | helloWorld  |
| go-ettle         | helloWorld  | helloWorld  | helloWorld  | helloWorld  |
| gobeam           | helloWorld  | helloworld  | helloworld  | helloWORLD  |
| golang-cz        | helloWorld  | helloWorld  | helloWorld  | helloWorld  |
| iancoleman       | HelloWorld  | HelloWorld  | HelloWorld  | HelloWorld  |
| lower-iancoleman | helloWorld  | helloWorld  | helloWorld  | helloWorld  |
| lower-searking   | hello World | helloWorld  | helloWorld  | hello_WORLD |
| resenje          | helloWorld  | helloWorld  | helloWorld  | helloWorld  |
| resenje-kebab    | Hello-World | Hello-World | Hello-World | Hello-World |
| resenje-snake    | Hello_World | Hello_World | Hello_World | Hello_World |
| searking         | Hello World | HelloWorld  | HelloWorld  | Hello_WORLD |
| tomeh            | HelloWorld  | helloWorld  | HelloWorld  | HelloWORLD  |
+------------------+-------------+-------------+-------------+-------------+
exit status 1
```


## `rntodarwin`

```bash
$ rntodarwin
Error: No files provided.
Usage: rntodarwin [options] <file1> [<file2> ...]

Options:
  -algorithm string
    	Choose the darwincase algorithm to use, supported: searking. (default "searking")
  -dry-run
    	Display the intended changes without renaming.
  -interactive
    	Ask for confirmation before renaming each file.

Conversion Examples:
+-----------+---------------+-------------+-------------+-------------------+
| ALGORITHM |               |             |             |                   |
+-----------+---------------+-------------+-------------+-------------------+
|           | Hello World   | helloWorld  | HelloWorld  | Hello_WORLD       |
| searking  | Hello_ _World | Hello_World | Hello_World | Hello___W_O_R_L_D |
+-----------+---------------+-------------+-------------+-------------------+
exit status 1
```


## `rntodelimited`

```bash
$ rntodelimited
Error: No files provided.
Usage: rntodelimited [options] <file1> [<file2> ...]

Options:
  -acronym value
    	Words to consider acronyms and not to assume they are words, ie ID => ID rather than ID => Id (iancoleman only)
  -acronym-from-file value
    	Words to consider acronyms and not to assume they are words, ie ID => ID rather than ID => Id (iancoleman only)
  -algorithm string
    	Choose the delimited algorithm to use, supported: iancoleman,screaming-iancoleman,searking,gobeam. (default "iancoleman")
  -delimiter string
    	Delimiter, default '.' but can be any single ascii character. (default ".")
  -dry-run
    	Display the intended changes without renaming.
  -ignore string
    	Other delimiter characters to ignore when parsing
  -input-delimiters string
    	Input's Delimiters, default ' ' but can be multiple ascii character. (searking only) (default ".")
  -interactive
    	Ask for confirmation before renaming each file.
  -word-seperators value
    	Word separators. (gobeam only) Default: "_-. "

Conversion Examples:
+----------------------+---------------+-------------+-------------+-------------------+
| ALGORITHM            |               |             |             |                   |
+----------------------+---------------+-------------+-------------+-------------------+
|                      | Hello World   | helloWorld  | HelloWorld  | Hello_WORLD       |
| gobeam               | Hello.World   | hello.World | Hello.World | Hello.WORLD       |
| iancoleman           | hello.world   | hello.world | hello.world | hello.world       |
| screaming-iancoleman | HELLO.WORLD   | HELLO.WORLD | HELLO.WORLD | HELLO.WORLD       |
| searking             | hello. .world | hello.world | hello.world | hello._.w.o.r.l.d |
+----------------------+---------------+-------------+-------------+-------------------+
exit status 1
```


## `rntokebab`

```bash
$ rntokebab
Error: No files provided.
Usage: rntokebab [options] <file1> [<file2> ...]

Options:
  -acronym value
    	Words to consider acronyms and not to assume they are words, ie ID => ID rather than ID => Id (iancoleman only)
  -acronym-from-file value
    	Words to consider acronyms and not to assume they are words, ie ID => ID rather than ID => Id (iancoleman only)
  -algorithm string
    	Choose the kebab algorithm to use, supported: gobeam,iancoleman,ettle,ETTLE,screaming-resenje,resenje-snake,golang-cz,searking,screaming-iancoleman,go-ettle,resenje,resenje-camel,tomeh. (default "iancoleman")
  -dry-run
    	Display the intended changes without renaming.
  -interactive
    	Ask for confirmation before renaming each file.
  -word-seperators value
    	Word separators. (gobeam only) Default: "_-. "

Conversion Examples:
+----------------------+---------------+-------------+-------------+-------------------+
| ALGORITHM            |               |             |             |                   |
+----------------------+---------------+-------------+-------------+-------------------+
|                      | Hello World   | helloWorld  | HelloWorld  | Hello_WORLD       |
| ETTLE                | HELLO-WORLD   | HELLO-WORLD | HELLO-WORLD | HELLO-WORLD       |
| ettle                | hello-world   | hello-world | hello-world | hello-world       |
| go-ettle             | hello-world   | hello-world | hello-world | hello-world       |
| gobeam               | Hello-World   | hello-World | Hello-World | Hello-WORLD       |
| golang-cz            | hello-world   | hello-world | hello-world | hello-world       |
| iancoleman           | hello-world   | hello-world | hello-world | hello-world       |
| resenje              | hello-world   | hello-world | hello-world | hello-world       |
| resenje-camel        | Hello-World   | Hello-World | Hello-World | Hello-World       |
| resenje-snake        | Hello_World   | Hello_World | Hello_World | Hello_World       |
| screaming-iancoleman | HELLO-WORLD   | HELLO-WORLD | HELLO-WORLD | HELLO-WORLD       |
| screaming-resenje    | HELLO-WORLD   | HELLO-WORLD | HELLO-WORLD | HELLO-WORLD       |
| searking             | hello- -world | hello-world | hello-world | hello-_-w-o-r-l-d |
| tomeh                | hello-world   | hello-world | hello-world | hello-w-o-r-l-d   |
+----------------------+---------------+-------------+-------------+-------------------+
exit status 1
```


## `rntolower`

```bash
$ rntolower
Error: No files provided.
Usage: rntolower [options] <file1> [<file2> ...]

Options:
  -algorithm string
    	Choose the lowercase algorithm to use, supported: searking,go,resenje,gobeam. (default "resenje")
  -dry-run
    	Display the intended changes without renaming.
  -interactive
    	Ask for confirmation before renaming each file.

Conversion Examples:
+-----------+-------------+-------------+-------------+-------------+
| ALGORITHM |             |             |             |             |
+-----------+-------------+-------------+-------------+-------------+
|           | Hello World | helloWorld  | HelloWorld  | Hello_WORLD |
| go        | hello world | helloworld  | helloworld  | hello_world |
| gobeam    | hello world | helloworld  | helloworld  | hello_world |
| resenje   | hello world | hello world | hello world | hello world |
| searking  | hello world | helloworld  | helloworld  | hello_world |
+-----------+-------------+-------------+-------------+-------------+
exit status 1
```


## `rntolowerleading`

```bash
$ rntolowerleading
Error: No files provided.
Usage: rntolowerleading [options] <file1> [<file2> ...]

Options:
  -algorithm string
    	Choose the lowerleading algorithm to use, supported: searking,gobeam. (default "searking")
  -dry-run
    	Display the intended changes without renaming.
  -interactive
    	Ask for confirmation before renaming each file.

Conversion Examples:
+-----------+-------------+------------+------------+-------------+
| ALGORITHM |             |            |            |             |
+-----------+-------------+------------+------------+-------------+
|           | Hello World | helloWorld | HelloWorld | Hello_WORLD |
| gobeam    | hello World | helloWorld | helloWorld | hello_WORLD |
| searking  | hello World | helloWorld | helloWorld | hello_WORLD |
+-----------+-------------+------------+------------+-------------+
exit status 1
```


## `rntopascal`

```bash
$ rntopascal
Error: No files provided.
Usage: rntopascal [options] <file1> [<file2> ...]

Options:
  -acronym value
    	Words to consider acronyms and not to assume they are words, ie ID => ID rather than ID => Id
  -acronym-from-file value
    	Words to consider acronyms and not to assume they are words, ie ID => ID rather than ID => Id
  -algorithm string
    	Choose the pascal algorithm to use, supported: resenje,tomeh,golang-cz,gobeam,searKing,ettle,go-ettle. (default "ettle")
  -dry-run
    	Display the intended changes without renaming.
  -interactive
    	Ask for confirmation before renaming each file.
  -word-seperators value
    	Word separators. (gobeam only) Default: "_-. "

Conversion Examples:
+-----------+-------------+------------+------------+-------------+
| ALGORITHM |             |            |            |             |
+-----------+-------------+------------+------------+-------------+
|           | Hello World | helloWorld | HelloWorld | Hello_WORLD |
| ettle     | HelloWorld  | HelloWorld | HelloWorld | HelloWorld  |
| go-ettle  | HelloWorld  | HelloWorld | HelloWorld | HelloWorld  |
| gobeam    | HelloWorld  | HelloWorld | HelloWorld | HelloWORLD  |
| golang-cz | HelloWorld  | HelloWorld | HelloWorld | HelloWorld  |
| resenje   | HelloWorld  | HelloWorld | HelloWorld | HelloWorld  |
| searKing  | Hello World | HelloWorld | HelloWorld | Hello_WORLD |
| tomeh     | HelloWorld  | HelloWorld | HelloWorld | HelloWORLD  |
+-----------+-------------+------------+------------+-------------+
exit status 1
```


## `rntosnake`

```bash
$ rntosnake
Error: No files provided.
Usage: rntosnake [options] <file1> [<file2> ...]

Options:
  -acronym value
    	Words to consider acronyms and not to assume they are words, ie ID => ID rather than ID => Id
  -acronym-from-file value
    	Words to consider acronyms and not to assume they are words, ie ID => ID rather than ID => Id
  -algorithm string
    	Choose the snake algorithm to use, supported: iancoleman,screaming-iancoleman,resenje,golang-cz,ETTLE,screaming-resenje,tomeh,gobeam,lowercase-searking,searking,ettle,go-ettle. (default "iancoleman")
  -dry-run
    	Display the intended changes without renaming.
  -interactive
    	Ask for confirmation before renaming each file.
  -word-seperators value
    	Word separators. (gobeam only) Default: "_-. "

Conversion Examples:
+----------------------+---------------+-------------+-------------+-----------------+
| ALGORITHM            |               |             |             |                 |
+----------------------+---------------+-------------+-------------+-----------------+
|                      | Hello World   | helloWorld  | HelloWorld  | Hello_WORLD     |
| ETTLE                | HELLO_WORLD   | HELLO_WORLD | HELLO_WORLD | HELLO_WORLD     |
| ettle                | hello_world   | hello_world | hello_world | hello_world     |
| go-ettle             | hello_world   | hello_world | hello_world | hello_world     |
| gobeam               | Hello_World   | hello_World | Hello_World | Hello_WORLD     |
| golang-cz            | hello_world   | hello_world | hello_world | hello_world     |
| iancoleman           | hello_world   | hello_world | hello_world | hello_world     |
| lowercase-searking   | hello_ _world | hello_world | hello_world | hello_w_o_r_l_d |
| resenje              | hello_world   | hello_world | hello_world | hello_world     |
| screaming-iancoleman | HELLO_WORLD   | HELLO_WORLD | HELLO_WORLD | HELLO_WORLD     |
| screaming-resenje    | HELLO_WORLD   | HELLO_WORLD | HELLO_WORLD | HELLO_WORLD     |
| searking             | hello_ _world | hello_world | hello_world | hello_w_o_r_l_d |
| tomeh                | hello_world   | hello_world | hello_world | hello_w_o_r_l_d |
+----------------------+---------------+-------------+-------------+-----------------+
exit status 1
```


## `rntotitle`

```bash
$ rntotitle
Error: No files provided.
Usage: rntotitle [options] <file1> [<file2> ...]

Options:
  -algorithm string
    	Choose the titlecase algorithm to use, supported: searking,go,resenje,gobeam,revett. (default "revett")
  -dry-run
    	Display the intended changes without renaming.
  -interactive
    	Ask for confirmation before renaming each file.

Conversion Examples:
+-----------+-------------+-------------+-------------+-------------+
| ALGORITHM |             |             |             |             |
+-----------+-------------+-------------+-------------+-------------+
|           | Hello World | helloWorld  | HelloWorld  | Hello_WORLD |
| go        | HELLO WORLD | HELLOWORLD  | HELLOWORLD  | HELLO_WORLD |
| gobeam    | Hello World | Helloworld  | Helloworld  | Hello_world |
| resenje   | Hello World | Hello World | Hello World | Hello World |
| revett    | Hello World | Helloworld  | Helloworld  | Hello_world |
| searking  | HELLO WORLD | HELLOWORLD  | HELLOWORLD  | HELLO_WORLD |
+-----------+-------------+-------------+-------------+-------------+
exit status 1
```


## `rntoupper`

```bash
$ rntoupper
Error: No files provided.
Usage: rntoupper [options] <file1> [<file2> ...]

Options:
  -algorithm string
    	Choose the uppercase algorithm to use, supported: searking,go,resenje,gobeam. (default "searking")
  -dry-run
    	Display the intended changes without renaming.
  -interactive
    	Ask for confirmation before renaming each file.

Conversion Examples:
+-----------+-------------+-------------+-------------+-------------+
| ALGORITHM |             |             |             |             |
+-----------+-------------+-------------+-------------+-------------+
|           | Hello World | helloWorld  | HelloWorld  | Hello_WORLD |
| go        | HELLO WORLD | HELLOWORLD  | HELLOWORLD  | HELLO_WORLD |
| gobeam    | HELLO WORLD | HELLOWORLD  | HELLOWORLD  | HELLO_WORLD |
| resenje   | HELLO WORLD | HELLO WORLD | HELLO WORLD | HELLO WORLD |
| searking  | Hello World | HelloWorld  | HelloWorld  | Hello_WORLD |
+-----------+-------------+-------------+-------------+-------------+
exit status 1
```


## `rntoupperleading`

```bash
$ rntoupperleading
Error: No files provided.
Usage: rntoupperleading [options] <file1> [<file2> ...]

Options:
  -algorithm string
    	Choose the upperleading algorithm to use, supported: searking,gobeam. (default "searking")
  -dry-run
    	Display the intended changes without renaming.
  -interactive
    	Ask for confirmation before renaming each file.

Conversion Examples:
+-----------+-------------+------------+------------+-------------+
| ALGORITHM |             |            |            |             |
+-----------+-------------+------------+------------+-------------+
|           | Hello World | helloWorld | HelloWorld | Hello_WORLD |
| gobeam    | Hello World | HelloWorld | HelloWorld | Hello_WORLD |
| searking  | Hello World | HelloWorld | HelloWorld | Hello_WORLD |
+-----------+-------------+------------+------------+-------------+
exit status 1
```


## `rntrim`

```bash
$ rntrim
Error: No files provided.
Usage: rntrim [options] <file1> [<file2> ...]

Options:
  -algorithm string
    	Choose the trim algorithm to use, supported: go. (default "go")
  -dry-run
    	Display the intended changes without renaming.
  -interactive
    	Ask for confirmation before renaming each file.
  -trim string
    	Characters to trim off end and start of name white space if not set.

Conversion Examples:
+-----------+--------------+-------------+----------------+-------------+
| ALGORITHM |              |             |                |             |
+-----------+--------------+-------------+----------------+-------------+
|           |  Hello World |  helloWorld |  HelloWorld    | Hello_WORLD |
| go        | Hello World  | helloWorld  | HelloWorld     | Hello_WORLD |
+-----------+--------------+-------------+----------------+-------------+
exit status 1
```

# Examples

## Normal run

```bash
$ touch "Hello World"
$ touch "Hello World.txt"
$ rnreverse Hello\ World
Rename: 'Hello World' -> 'dlroW olleH'
Renamed successfully.
$ ls
'dlroW olleH'  'Hello World.txt'
$ rnreverse Hello\ World.txt 
Rename: 'Hello World.txt' -> 'dlroW olleH.txt'
Renamed successfully.
$ ls
'dlroW olleH'  'dlroW olleH.txt'

```

## Interactive mode

```bash
$ ls                    
'Hello World'  'Hello World.txt'
$ rnreverse -interactive * 
Rename: 'Hello World' -> 'dlroW olleH'
Rename 'Hello World' to 'dlroW olleH'? [y/n]: y
Renamed successfully.
Rename: 'Hello World.txt' -> 'dlroW olleH.txt'
Rename 'Hello World.txt' to 'dlroW olleH.txt'? [y/n]: n 
Skipped.
$ ls
'dlroW olleH'  'Hello World.txt'
```

## Sample outputs

```bash
$ rnreverse *             
Rename: 'Hello World' -> 'dlroW olleH'
Renamed successfully.
Rename: 'Hello World.txt' -> 'dlroW olleH.txt'
Renamed successfully.

$ touch Hello\ World Hello\ World.txt
$ rntocamel *
Rename: 'Hello World' -> 'HelloWorld'
Renamed successfully.
Rename: 'Hello World.txt' -> 'HelloWorld.txt'
Renamed successfully.

$ touch Hello\ World Hello\ World.txt
$ ls
'Hello World'  'Hello World.txt'
$ rntocamel *                        
Rename: 'Hello World' -> 'HelloWorld'
Renamed successfully.
Rename: 'Hello World.txt' -> 'HelloWorld.txt'
Renamed successfully.
$ ls
HelloWorld  HelloWorld.txt

$ touch Hello\ World Hello\ World.txt
$ rntodarwin *                       
Rename: 'Hello World' -> 'Hello_ _World'
Renamed successfully.
Rename: 'Hello World.txt' -> 'Hello_ _World.txt'
Renamed successfully.
$ ls
'Hello_ _World'  'Hello_ _World.txt'

$ touch Hello\ World Hello\ World.txt
$ ls
'Hello World'  'Hello World.txt'
$ rntodelimited *
Rename: 'Hello World' -> 'hello.world'
Renamed successfully.
Rename: 'Hello World.txt' -> 'hello.world.txt'
Renamed successfully.
$ ls
hello.world  hello.world.txt


$ touch Hello\ World Hello\ World.txt
$ rntokebab *
Rename: 'Hello World' -> 'hello-world'
Renamed successfully.
Rename: 'Hello World.txt' -> 'hello-world.txt'
Renamed successfully.
$ ls    
hello-world  hello-world.txt


$ touch Hello\ World Hello\ World.txt
$ touch Hello\ WORLD Hello\ WORLD.txt
$ ls
'Hello World'  'Hello WORLD'  'Hello World.txt'  'Hello WORLD.txt'
$ rntolowerleading *                 
Rename: 'Hello World' -> 'hello World'
Renamed successfully.
Rename: 'Hello WORLD' -> 'hello WORLD'
Renamed successfully.
Rename: 'Hello World.txt' -> 'hello World.txt'
Renamed successfully.
Rename: 'Hello WORLD.txt' -> 'hello WORLD.txt'
Renamed successfully.

$ touch Hello\ World Hello\ World.txt
$ rntosnake *
Rename: 'Hello World' -> 'hello_world'
Renamed successfully.
Rename: 'Hello World.txt' -> 'hello_world.txt'
Renamed successfully.
$ ls
hello_world  hello_world.txt


$ touch Hello\ World Hello\ World.txt
$ touch "hello of THE world.txt"
$ ls                            
'hello of THE world.txt'  'Hello World'  'Hello World.txt'
$ rntotitle *                   
Rename: 'hello of THE world.txt' -> 'Hello of THE World.txt'
Renamed successfully.
Skipping 'Hello World' (already matches desired format).
Skipping 'Hello World.txt' (already matches desired format).

$ touch hello\ world.txt
$ ls
'hello world.txt'
$ rntoupperleading *
Rename: 'hello world.txt' -> 'Hello world.txt'
Renamed successfully.


$ touch Hello.txt ' Hello world .txt' 'hello     .txt'
$ ls                                                  
'hello     .txt'   Hello.txt  ' Hello world .txt'
$ rntrim *
Rename: 'hello     .txt' -> 'hello.txt'
Renamed successfully.
Skipping 'Hello.txt' (already matches desired format).
Rename: ' Hello world .txt' -> 'Hello world.txt'
Renamed successfully.
$ ls  
 hello.txt   Hello.txt  'Hello world.txt'


```

# Download

See Github releases here: https://github.com/arran4/rntocase/releases

## Gentoo

See: https://github.com/arran4/rntocase

```bash
$ eselect repository enable arrans-overlay
$ emerge -va app-misc/rntocase-bin 
```
