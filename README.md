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

## `rnreverse`

```bash
$ rnreverse --help          
Usage: rnreverse [options] <file1> [<file2> ...]

Options:
  -algorithm string
        Choose the reverse algorithm to use, supported: searking. (default "searking")
  -dry-run
        Display the intended changes without renaming.
  -interactive
        Ask for confirmation before renaming each file.

```

## `rntocamel`

```bash
$ rntocamel -h
Usage: rntocamel [options] <file1> [<file2> ...]

Options:
  -acronym value
        Words to consider acronyms and not to assume they are words, ie ID => ID rather than ID => Id
  -acronym-from-file value
        Words to consider acronyms and not to assume they are words, ie ID => ID rather than ID => Id
  -algorithm string
        Choose the camel algorithm to use, supported: searking,iancoleman,lower-iancoleman,lower-searking. (default "iancoleman")
  -dry-run
        Display the intended changes without renaming.
  -interactive
        Ask for confirmation before renaming each file.

```


## `rntodarwin`

```bash
$ rntocamel -h
Usage: rntocamel [options] <file1> [<file2> ...]

Options:
  -acronym value
        Words to consider acronyms and not to assume they are words, ie ID => ID rather than ID => Id
  -acronym-from-file value
        Words to consider acronyms and not to assume they are words, ie ID => ID rather than ID => Id
  -algorithm string
        Choose the camel algorithm to use, supported: searking,iancoleman,lower-iancoleman,lower-searking. (default "iancoleman")
  -dry-run
        Display the intended changes without renaming.
  -interactive
        Ask for confirmation before renaming each file.
$ rntodarwin -h
Usage: rntodarwin [options] <file1> [<file2> ...]

Options:
  -algorithm string
        Choose the darwincase algorithm to use, supported: searking. (default "searking")
  -dry-run
        Display the intended changes without renaming.
  -interactive
        Ask for confirmation before renaming each file.

```

## `rntodelimited`

```bash
$ rntodelimited -h
Usage: rntodelimited [options] <file1> [<file2> ...]

Options:
  -acronym value
        Words to consider acronyms and not to assume they are words, ie ID => ID rather than ID => Id
  -acronym-from-file value
        Words to consider acronyms and not to assume they are words, ie ID => ID rather than ID => Id
  -algorithm string
        Choose the delimited algorithm to use, supported: iancoleman,screaming-iancoleman,searking. (default "iancoleman")
  -delimiter string
        Delimiter, default '.' but can be any single ascii character. (default ".")
  -dry-run
        Display the intended changes without renaming.
  -ignore string
        Other delimiter characters to ignore when parsing
  -interactive
        Ask for confirmation before renaming each file.

```


## `rntokebab`

```bash
$ rntokebab -h
Usage: rntokebab [options] <file1> [<file2> ...]

Options:
  -acronym value
        Words to consider acronyms and not to assume they are words, ie ID => ID rather than ID => Id
  -acronym-from-file value
        Words to consider acronyms and not to assume they are words, ie ID => ID rather than ID => Id
  -algorithm string
        Choose the kebab algorithm to use, supported: searking,screaming-iancoleman,iancoleman. (default "iancoleman")
  -dry-run
        Display the intended changes without renaming.
  -interactive
        Ask for confirmation before renaming each file.

```

## `rntolowerleading`

```bash
$ rntolowerleading -h
Usage: rntolowerleading [options] <file1> [<file2> ...]

Options:
  -algorithm string
        Choose the lowerleading algorithm to use, supported: searking. (default "searking")
  -dry-run
        Display the intended changes without renaming.
  -interactive
        Ask for confirmation before renaming each file.

```

## `rntosnake`

```bash
$ rntosnake -h
Usage: rntosnake [options] <file1> [<file2> ...]

Options:
  -acronym value
        Words to consider acronyms and not to assume they are words, ie ID => ID rather than ID => Id
  -acronym-from-file value
        Words to consider acronyms and not to assume they are words, ie ID => ID rather than ID => Id
  -algorithm string
        Choose the snake algorithm to use, supported: iancoleman,screaming-iancoleman,lowercase-searking,searking. (default "iancoleman")
  -dry-run
        Display the intended changes without renaming.
  -interactive
        Ask for confirmation before renaming each file.

```

## `rntotitle`

```bash
$ rntotitle -h
Usage: rntotitle [options] <file1> [<file2> ...]

Options:
  -algorithm string
        Choose the titlecase algorithm to use, supported: revett,searking. (default "revett")
  -dry-run
        Display the intended changes without renaming.
  -interactive
        Ask for confirmation before renaming each file.

```

## `rntoupperleading`

```bash
$ rntoupperleading -h
Usage: rntoupperleading [options] <file1> [<file2> ...]

Options:
  -algorithm string
        Choose the upperleading algorithm to use, supported: searking. (default "searking")
  -dry-run
        Display the intended changes without renaming.
  -interactive
        Ask for confirmation before renaming each file.

```

## `rntrim`

```bash
$ rntrim -h
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
