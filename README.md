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
+-----------+--------------------+-----------------------+-------------------------+--------------------------------+-------------------------+
| ALGORITHM |                    |                       |                         |                                |                         |
+-----------+--------------------+-----------------------+-------------------------+--------------------------------+-------------------------+
|           | Basic Cases        | Hello World           | helloWorld              | HelloWorld                     | Hello_WORLD             |
| gobeam    |                    | HW                    | h                       | H                              | H                       |
|           | Underscore Cases   | __camel_snake_kebab__ | _leading_snake_case_    | __trailing__underscore__       |                         |
| gobeam    |                    | _                     | _                       | _                              |                         |
|           | Hyphen Cases       | --camel-snake-kebab   | -leading-kebab-case-    | --trailing--hyphen--           |                         |
| gobeam    |                    | -                     | -                       | -                              |                         |
|           | Mixed Delimiters   | hello_world-and-kebab | Mixed_Snake-Kebab--Case | UPPER_snake-KEBAB_Case         |                         |
| gobeam    |                    | h                     | M                       | U                              |                         |
|           | Acronym Handling   | ID_NUMBER             | HTTP_Response_Code      | XML_HTTP_REQUEST               | API_Version_2           |
| gobeam    |                    | I                     | H                       | X                              | A                       |
|           | Spaces and Words   |    leading spaces     | trailing spaces         |    both sides                  | This is a test sentence |
| gobeam    |                    | ls                    | ts                      | bs                             | Tiats                   |
|           | Symbols and Random | 123_ABC-xyz--789      | !!special$$characters** | ___!!!___weird___CASE___!!!___ |                         |
| gobeam    |                    | 1                     | !                       | _                              |                         |
+-----------+--------------------+-----------------------+-------------------------+--------------------------------+-------------------------+
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
+-----------+--------------------+-----------------------+-------------------------+--------------------------------+-------------------------+
| ALGORITHM |                    |                       |                         |                                |                         |
+-----------+--------------------+-----------------------+-------------------------+--------------------------------+-------------------------+
|           | Basic Cases        | Hello World           | helloWorld              | HelloWorld                     | Hello_WORLD             |
| searking  |                    | dlroW olleH           | dlroWolleh              | dlroWolleH                     | DLROW_olleH             |
|           | Underscore Cases   | __camel_snake_kebab__ | _leading_snake_case_    | __trailing__underscore__       |                         |
| searking  |                    | __babek_ekans_lemac__ | _esac_ekans_gnidael_    | __erocsrednu__gniliart__       |                         |
|           | Hyphen Cases       | --camel-snake-kebab   | -leading-kebab-case-    | --trailing--hyphen--           |                         |
| searking  |                    | babek-ekans-lemac--   | -esac-babek-gnidael-    | --nehpyh--gniliart--           |                         |
|           | Mixed Delimiters   | hello_world-and-kebab | Mixed_Snake-Kebab--Case | UPPER_snake-KEBAB_Case         |                         |
| searking  |                    | babek-dna-dlrow_olleh | esaC--babeK-ekanS_dexiM | esaC_BABEK-ekans_REPPU         |                         |
|           | Acronym Handling   | ID_NUMBER             | HTTP_Response_Code      | XML_HTTP_REQUEST               | API_Version_2           |
| searking  |                    | REBMUN_DI             | edoC_esnopseR_PTTH      | TSEUQER_PTTH_LMX               | 2_noisreV_IPA           |
|           | Spaces and Words   |    leading spaces     | trailing spaces         |    both sides                  | This is a test sentence |
| searking  |                    | secaps gnidael        |    secaps gniliart      |    sedis htob                  | ecnetnes tset a si sihT |
|           | Symbols and Random | 123_ABC-xyz--789      | !!special$$characters** | ___!!!___weird___CASE___!!!___ |                         |
| searking  |                    | 987--zyx-CBA_321      | **sretcarahc$$laiceps!! | ___!!!___ESAC___driew___!!!___ |                         |
+-----------+--------------------+-----------------------+-------------------------+--------------------------------+-------------------------+
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
    	Choose the camel algorithm to use, supported: iancoleman,go-ettle,resenje-snake,tomeh,gobeam,resenje-kebab,golang-cz,lower-iancoleman,lower-searking,searking,ettle,resenje. (default "iancoleman")
  -dry-run
    	Display the intended changes without renaming.
  -interactive
    	Ask for confirmation before renaming each file.
  -word-seperators value
    	Word separators. (gobeam only) Default: "_-. "

Conversion Examples:
+------------------+--------------------+-----------------------+-----------------------------+--------------------------------+-------------------------+
| ALGORITHM        |                    |                       |                             |                                |                         |
+------------------+--------------------+-----------------------+-----------------------------+--------------------------------+-------------------------+
|                  | Basic Cases        | Hello World           | helloWorld                  | HelloWorld                     | Hello_WORLD             |
| ettle            |                    | helloWorld            | helloWorld                  | helloWorld                     | helloWorld              |
| go-ettle         |                    | helloWorld            | helloWorld                  | helloWorld                     | helloWorld              |
| gobeam           |                    | helloWorld            | helloworld                  | helloworld                     | helloWORLD              |
| golang-cz        |                    | helloWorld            | helloWorld                  | helloWorld                     | helloWorld              |
| iancoleman       |                    | HelloWorld            | HelloWorld                  | HelloWorld                     | HelloWorld              |
| lower-iancoleman |                    | helloWorld            | helloWorld                  | helloWorld                     | helloWorld              |
| lower-searking   |                    | hello World           | helloWorld                  | helloWorld                     | hello_WORLD             |
| resenje          |                    | helloWorld            | helloWorld                  | helloWorld                     | helloWorld              |
| resenje-kebab    |                    | Hello-World           | Hello-World                 | Hello-World                    | Hello-World             |
| resenje-snake    |                    | Hello_World           | Hello_World                 | Hello_World                    | Hello_World             |
| searking         |                    | Hello World           | HelloWorld                  | HelloWorld                     | Hello_WORLD             |
| tomeh            |                    | HelloWorld            | helloWorld                  | HelloWorld                     | HelloWORLD              |
|                  | Underscore Cases   | __camel_snake_kebab__ | _leading_snake_case_        | __trailing__underscore__       |                         |
| ettle            |                    | camelSnakeKebab       | leadingSnakeCase            | trailingUnderscore             |                         |
| go-ettle         |                    | camelSnakeKebab       | leadingSnakeCase            | trailingUnderscore             |                         |
| gobeam           |                    | camelSnakeKebab       | leadingSnakeCase            | trailingUnderscore             |                         |
| golang-cz        |                    | camelSnakeKebab       | leadingSnakeCase            | trailingUnderscore             |                         |
| iancoleman       |                    | CamelSnakeKebab       | LeadingSnakeCase            | TrailingUnderscore             |                         |
| lower-iancoleman |                    | CamelSnakeKebab       | LeadingSnakeCase            | TrailingUnderscore             |                         |
| lower-searking   |                    | __camel_snake_kebab__ | _leading_snake_case_        | __trailing__underscore__       |                         |
| resenje          |                    | camelSnakeKebab       | leadingSnakeCase            | trailingUnderscore             |                         |
| resenje-kebab    |                    | Camel-Snake-Kebab     | Leading-Snake-Case          | Trailing-Underscore            |                         |
| resenje-snake    |                    | __Camel_Snake_Kebab__ | _Leading_Snake_Case_        | __Trailing_Underscore__        |                         |
| searking         |                    | __camel_snake_kebab__ | _leading_snake_case_        | __trailing__underscore__       |                         |
| tomeh            |                    | CamelSnakeKebab       | LeadingSnakeCase            | TrailingUnderscore             |                         |
|                  | Hyphen Cases       | --camel-snake-kebab   | -leading-kebab-case-        | --trailing--hyphen--           |                         |
| ettle            |                    | camelSnakeKebab       | leadingKebabCase            | trailingHyphen                 |                         |
| go-ettle         |                    | camelSnakeKebab       | leadingKebabCase            | trailingHyphen                 |                         |
| gobeam           |                    | camelSnakeKebab       | leadingKebabCase            | trailingHyphen                 |                         |
| golang-cz        |                    | camelSnakeKebab       | leadingKebabCase            | trailingHyphen                 |                         |
| iancoleman       |                    | CamelSnakeKebab       | LeadingKebabCase            | TrailingHyphen                 |                         |
| lower-iancoleman |                    | CamelSnakeKebab       | LeadingKebabCase            | TrailingHyphen                 |                         |
| lower-searking   |                    | --camel-snake-kebab   | -leading-kebab-case-        | --trailing--hyphen--           |                         |
| resenje          |                    | camelSnakeKebab       | leadingKebabCase            | trailingHyphen                 |                         |
| resenje-kebab    |                    | --Camel-Snake-Kebab   | -Leading-Kebab-Case-        | --Trailing-Hyphen--            |                         |
| resenje-snake    |                    | Camel_Snake_Kebab     | Leading_Kebab_Case          | Trailing_Hyphen                |                         |
| searking         |                    | --camel-snake-kebab   | -leading-kebab-case-        | --trailing--hyphen--           |                         |
| tomeh            |                    | CamelSnakeKebab       | LeadingKebabCase            | TrailingHyphen                 |                         |
|                  | Mixed Delimiters   | hello_world-and-kebab | Mixed_Snake-Kebab--Case     | UPPER_snake-KEBAB_Case         |                         |
| ettle            |                    | helloWorldAndKebab    | mixedSnakeKebabCase         | upperSnakeKebabCase            |                         |
| go-ettle         |                    | helloWorldAndKebab    | mixedSnakeKebabCase         | upperSnakeKebabCase            |                         |
| gobeam           |                    | helloWorldAndKebab    | mixedSnakeKebabCase         | upperSnakeKEBABCase            |                         |
| golang-cz        |                    | helloWorldAndKebab    | mixedSnakeKebabCase         | upperSnakeKebabCase            |                         |
| iancoleman       |                    | HelloWorldAndKebab    | MixedSnakeKebabCase         | UpperSnakeKebabCase            |                         |
| lower-iancoleman |                    | helloWorldAndKebab    | mixedSnakeKebabCase         | upperSnakeKebabCase            |                         |
| lower-searking   |                    | hello_world-and-kebab | mixed_Snake-Kebab--Case     | uPPER_snake-KEBAB_Case         |                         |
| resenje          |                    | helloWorldAndKebab    | mixedSnakeKebabCase         | upperSnakeKebabCase            |                         |
| resenje-kebab    |                    | Hello-World-And-Kebab | Mixed-Snake-Kebab-Case      | Upper-Snake-Kebab-Case         |                         |
| resenje-snake    |                    | Hello_World_And_Kebab | Mixed_Snake_Kebab_Case      | Upper_Snake_Kebab_Case         |                         |
| searking         |                    | Hello_world-and-kebab | Mixed_Snake-Kebab--Case     | UPPER_snake-KEBAB_Case         |                         |
| tomeh            |                    | helloWorldAndKebab    | MixedSnakeKebabCase         | UPPERSnakeKEBABCase            |                         |
|                  | Acronym Handling   | ID_NUMBER             | HTTP_Response_Code          | XML_HTTP_REQUEST               | API_Version_2           |
| ettle            |                    | idNumber              | httpResponseCode            | xmlHttpRequest                 | apiVersion2             |
| go-ettle         |                    | idNumber              | httpResponseCode            | xmlHTTPRequest                 | apiVersion2             |
| gobeam           |                    | idNUMBER              | httpResponseCode            | xmlHTTPREQUEST                 | apiVersion2             |
| golang-cz        |                    | idNumber              | httpResponseCode            | xmlHttpRequest                 | apiVersion2             |
| iancoleman       |                    | IdNumber              | HttpResponseCode            | XmlHttpRequest                 | ApiVersion2             |
| lower-iancoleman |                    | idNumber              | httpResponseCode            | xmlHttpRequest                 | apiVersion2             |
| lower-searking   |                    | iD_NUMBER             | hTTP_Response_Code          | xML_HTTP_REQUEST               | aPI_Version_2           |
| resenje          |                    | idNumber              | httpResponseCode            | xmlHttpRequest                 | apiVersion2             |
| resenje-kebab    |                    | Id-Number             | Http-Response-Code          | Xml-Http-Request               | Api-Version-2           |
| resenje-snake    |                    | Id_Number             | Http_Response_Code          | Xml_Http_Request               | Api_Version_2           |
| searking         |                    | ID_NUMBER             | HTTP_Response_Code          | XML_HTTP_REQUEST               | API_Version_2           |
| tomeh            |                    | IDNUMBER              | HTTPResponseCode            | XMLHTTPREQUEST                 | APIVersion2             |
|                  | Spaces and Words   |    leading spaces     | trailing spaces             |    both sides                  | This is a test sentence |
| ettle            |                    | leadingSpaces         | trailingSpaces              | bothSides                      | thisIsATestSentence     |
| go-ettle         |                    | leadingSpaces         | trailingSpaces              | bothSides                      | thisIsATestSentence     |
| gobeam           |                    | leadingSpaces         | trailingSpaces              | bothSides                      | thisIsATestSentence     |
| golang-cz        |                    | leadingSpaces         | trailingSpaces              | bothSides                      | thisIsATestSentence     |
| iancoleman       |                    | LeadingSpaces         | TrailingSpaces              | BothSides                      | ThisIsATestSentence     |
| lower-iancoleman |                    | leadingSpaces         | trailingSpaces              | bothSides                      | thisIsATestSentence     |
| lower-searking   |                    |    leading spaces     | trailing spaces             |    both sides                  | this is a test sentence |
| resenje          |                    | leadingSpaces         | trailingSpaces              | bothSides                      | thisIsATestSentence     |
| resenje-kebab    |                    | Leading-Spaces        | Trailing-Spaces             | Both-Sides                     | This-Is-A-Test-Sentence |
| resenje-snake    |                    | Leading_Spaces        | Trailing_Spaces             | Both_Sides                     | This_Is_A_Test_Sentence |
| searking         |                    |    leading spaces     | Trailing spaces             |    both sides                  | This is a test sentence |
| tomeh            |                    | LeadingSpaces         | trailingSpaces              | BothSides                      | ThisIsATestSentence     |
|                  | Symbols and Random | 123_ABC-xyz--789      | !!special$$characters**     | ___!!!___weird___CASE___!!!___ |                         |
| ettle            |                    | 123AbcXyz789          | !!special$$characters**     | !!!WeirdCase!!!                |                         |
| go-ettle         |                    | 123AbcXyz789          | !!special$$characters**     | !!!WeirdCase!!!                |                         |
| gobeam           |                    | 123ABCXyz789          | !!special$$characters**     | !!!WeirdCASE!!!                |                         |
| golang-cz        |                    | 123AbcXyz789          | specialCharacters           | weirdCase                      |                         |
| iancoleman       |                    | 123AbcXyz789          | specialcharacters           | WeirdCase                      |                         |
| lower-iancoleman |                    | 123AbcXyz789          | specialcharacters           | WeirdCase                      |                         |
| lower-searking   |                    | 123_ABC-xyz--789      | !!special$$characters**     | ___!!!___weird___CASE___!!!___ |                         |
| resenje          |                    | 123AbcXyz789          | !!special$$characters**     | !!!WeirdCase!!!                |                         |
| resenje-kebab    |                    | 123-Abc-Xyz-789       | !-!special-$-$characters-** | !!!-Weird-Case-!!!             |                         |
| resenje-snake    |                    | 123_Abc_Xyz_789       | !_!special_$_$characters_** | ___!!!_Weird_Case_!!!___       |                         |
| searking         |                    | 123_ABC-xyz--789      | !!special$$characters**     | ___!!!___weird___CASE___!!!___ |                         |
| tomeh            |                    | 123ABCXyz789          | !!specialCharacters**       | !!!WeirdCASE!!!                |                         |
+------------------+--------------------+-----------------------+-----------------------------+--------------------------------+-------------------------+
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
+-----------+--------------------+------------------------------+--------------------------------+--------------------------------------------------------+-----------------------------+
| ALGORITHM |                    |                              |                                |                                                        |                             |
+-----------+--------------------+------------------------------+--------------------------------+--------------------------------------------------------+-----------------------------+
|           | Basic Cases        | Hello World                  | helloWorld                     | HelloWorld                                             | Hello_WORLD                 |
| searking  |                    | Hello_ _World                | Hello_World                    | Hello_World                                            | Hello___W_O_R_L_D           |
|           | Underscore Cases   | __camel_snake_kebab__        | _leading_snake_case_           | __trailing__underscore__                               |                             |
| searking  |                    | ___camel__snake__kebab____   | _leading__snake__case__        | ___trailing____underscore____                          |                             |
|           | Hyphen Cases       | --camel-snake-kebab          | -leading-kebab-case-           | --trailing--hyphen--                                   |                             |
| searking  |                    | -_-Camel_-Snake_-Kebab       | -Leading_-Kebab_-Case_-        | -_-Trailing_-_-Hyphen_-_-                              |                             |
|           | Mixed Delimiters   | hello_world-and-kebab        | Mixed_Snake-Kebab--Case        | UPPER_snake-KEBAB_Case                                 |                             |
| searking  |                    | Hello__world_-And_-Kebab     | Mixed___Snake_-_Kebab_-_-_Case | U_P_P_E_R__snake_-_K_E_B_A_B___Case                    |                             |
|           | Acronym Handling   | ID_NUMBER                    | HTTP_Response_Code             | XML_HTTP_REQUEST                                       | API_Version_2               |
| searking  |                    | I_D___N_U_M_B_E_R            | H_T_T_P___Response___Code      | X_M_L___H_T_T_P___R_E_Q_U_E_S_T                        | A_P_I___Version___2         |
|           | Spaces and Words   |    leading spaces            | trailing spaces                |    both sides                                          | This is a test sentence     |
| searking  |                    |  _ _ Leading_ Spaces         | Trailing_ Spaces_ _ _          |  _ _ Both_ Sides_ _ _                                  | This_ Is_ A_ Test_ Sentence |
|           | Symbols and Random | 123_ABC-xyz--789             | !!special$$characters**        | ___!!!___weird___CASE___!!!___                         |                             |
| searking  |                    | 1_2_3___A_B_C_-Xyz_-_-_7_8_9 | !_!Special_$_$Characters_*_*   | ______!_!_!______weird_______C_A_S_E_______!_!_!______ |                             |
+-----------+--------------------+------------------------------+--------------------------------+--------------------------------------------------------+-----------------------------+
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
    	Choose the delimited algorithm to use, supported: searking,gobeam,iancoleman,screaming-iancoleman. (default "iancoleman")
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
+----------------------+--------------------+------------------------------+--------------------------------+--------------------------------------------------------+-----------------------------+
| ALGORITHM            |                    |                              |                                |                                                        |                             |
+----------------------+--------------------+------------------------------+--------------------------------+--------------------------------------------------------+-----------------------------+
|                      | Basic Cases        | Hello World                  | helloWorld                     | HelloWorld                                             | Hello_WORLD                 |
| gobeam               |                    | Hello.World                  | hello.World                    | Hello.World                                            | Hello.WORLD                 |
| iancoleman           |                    | hello.world                  | hello.world                    | hello.world                                            | hello.world                 |
| screaming-iancoleman |                    | HELLO.WORLD                  | HELLO.WORLD                    | HELLO.WORLD                                            | HELLO.WORLD                 |
| searking             |                    | hello. .world                | hello.world                    | hello.world                                            | hello._.w.o.r.l.d           |
|                      | Underscore Cases   | __camel_snake_kebab__        | _leading_snake_case_           | __trailing__underscore__                               |                             |
| gobeam               |                    | camel.snake.kebab            | leading.snake.case             | trailing.underscore                                    |                             |
| iancoleman           |                    | ..camel.snake.kebab..        | .leading.snake.case.           | ..trailing..underscore..                               |                             |
| screaming-iancoleman |                    | ..CAMEL.SNAKE.KEBAB..        | .LEADING.SNAKE.CASE.           | ..TRAILING..UNDERSCORE..                               |                             |
| searking             |                    | _._camel._snake._kebab._._   | _leading._snake._case._        | _._trailing._._underscore._._                          |                             |
|                      | Hyphen Cases       | --camel-snake-kebab          | -leading-kebab-case-           | --trailing--hyphen--                                   |                             |
| gobeam               |                    | camel.snake.kebab            | leading.kebab.case             | trailing.hyphen                                        |                             |
| iancoleman           |                    | ..camel.snake.kebab          | .leading.kebab.case.           | ..trailing..hyphen..                                   |                             |
| screaming-iancoleman |                    | ..CAMEL.SNAKE.KEBAB          | .LEADING.KEBAB.CASE.           | ..TRAILING..HYPHEN..                                   |                             |
| searking             |                    | -.-camel.-snake.-kebab       | -leading.-kebab.-case.-        | -.-trailing.-.-hyphen.-.-                              |                             |
|                      | Mixed Delimiters   | hello_world-and-kebab        | Mixed_Snake-Kebab--Case        | UPPER_snake-KEBAB_Case                                 |                             |
| gobeam               |                    | hello.world.and.kebab        | Mixed.Snake.Kebab.Case         | UPPER.snake.KEBAB.Case                                 |                             |
| iancoleman           |                    | hello.world.and.kebab        | mixed.snake.kebab..case        | upper.snake.kebab.case                                 |                             |
| screaming-iancoleman |                    | HELLO.WORLD.AND.KEBAB        | MIXED.SNAKE.KEBAB..CASE        | UPPER.SNAKE.KEBAB.CASE                                 |                             |
| searking             |                    | hello._world.-and.-kebab     | mixed._.snake.-.kebab.-.-.case | u.p.p.e.r._snake.-.k.e.b.a.b._.case                    |                             |
|                      | Acronym Handling   | ID_NUMBER                    | HTTP_Response_Code             | XML_HTTP_REQUEST                                       | API_Version_2               |
| gobeam               |                    | ID.NUMBER                    | HTTP.Response.Code             | XML.HTTP.REQUEST                                       | API.Version.2               |
| iancoleman           |                    | id.number                    | http.response.code             | xml.http.request                                       | api.version.2               |
| screaming-iancoleman |                    | ID.NUMBER                    | HTTP.RESPONSE.CODE             | XML.HTTP.REQUEST                                       | API.VERSION.2               |
| searking             |                    | i.d._.n.u.m.b.e.r            | h.t.t.p._.response._.code      | x.m.l._.h.t.t.p._.r.e.q.u.e.s.t                        | a.p.i._.version._.2         |
|                      | Spaces and Words   |    leading spaces            | trailing spaces                |    both sides                                          | This is a test sentence     |
| gobeam               |                    | leading.spaces               | trailing.spaces                | both.sides                                             | This.is.a.test.sentence     |
| iancoleman           |                    | leading.spaces               | trailing.spaces                | both.sides                                             | this.is.a.test.sentence     |
| screaming-iancoleman |                    | LEADING.SPACES               | TRAILING.SPACES                | BOTH.SIDES                                             | THIS.IS.A.TEST.SENTENCE     |
| searking             |                    |  . . leading. spaces         | trailing. spaces. . .          |  . . both. sides. . .                                  | this. is. a. test. sentence |
|                      | Symbols and Random | 123_ABC-xyz--789             | !!special$$characters**        | ___!!!___weird___CASE___!!!___                         |                             |
| gobeam               |                    | 123.ABC.xyz.789              | !!special$$characters**        | !!!.weird.CASE.!!!                                     |                             |
| iancoleman           |                    | 123.abc.xyz..789             | !!special$$characters**        | ...!!!...weird...case...!!!...                         |                             |
| screaming-iancoleman |                    | 123.ABC.XYZ..789             | !!SPECIAL$$CHARACTERS**        | ...!!!...WEIRD...CASE...!!!...                         |                             |
| searking             |                    | 1.2.3._.a.b.c.-xyz.-.-.7.8.9 | !.!special.$.$characters.*.*   | _._._.!.!.!._._._weird._._._.c.a.s.e._._._.!.!.!._._._ |                             |
+----------------------+--------------------+------------------------------+--------------------------------+--------------------------------------------------------+-----------------------------+
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
    	Choose the kebab algorithm to use, supported: iancoleman,screaming-iancoleman,go-ettle,resenje-camel,searking,ettle,ETTLE,resenje,screaming-resenje,resenje-snake,tomeh,golang-cz,gobeam. (default "iancoleman")
  -dry-run
    	Display the intended changes without renaming.
  -interactive
    	Ask for confirmation before renaming each file.
  -word-seperators value
    	Word separators. (gobeam only) Default: "_-. "

Conversion Examples:
+----------------------+--------------------+------------------------------+--------------------------------+--------------------------------------------------------+-----------------------------+
| ALGORITHM            |                    |                              |                                |                                                        |                             |
+----------------------+--------------------+------------------------------+--------------------------------+--------------------------------------------------------+-----------------------------+
|                      | Basic Cases        | Hello World                  | helloWorld                     | HelloWorld                                             | Hello_WORLD                 |
| ETTLE                |                    | HELLO-WORLD                  | HELLO-WORLD                    | HELLO-WORLD                                            | HELLO-WORLD                 |
| ettle                |                    | hello-world                  | hello-world                    | hello-world                                            | hello-world                 |
| go-ettle             |                    | hello-world                  | hello-world                    | hello-world                                            | hello-world                 |
| gobeam               |                    | Hello-World                  | hello-World                    | Hello-World                                            | Hello-WORLD                 |
| golang-cz            |                    | hello-world                  | hello-world                    | hello-world                                            | hello-world                 |
| iancoleman           |                    | hello-world                  | hello-world                    | hello-world                                            | hello-world                 |
| resenje              |                    | hello-world                  | hello-world                    | hello-world                                            | hello-world                 |
| resenje-camel        |                    | Hello-World                  | Hello-World                    | Hello-World                                            | Hello-World                 |
| resenje-snake        |                    | Hello_World                  | Hello_World                    | Hello_World                                            | Hello_World                 |
| screaming-iancoleman |                    | HELLO-WORLD                  | HELLO-WORLD                    | HELLO-WORLD                                            | HELLO-WORLD                 |
| screaming-resenje    |                    | HELLO-WORLD                  | HELLO-WORLD                    | HELLO-WORLD                                            | HELLO-WORLD                 |
| searking             |                    | hello- -world                | hello-world                    | hello-world                                            | hello-_-w-o-r-l-d           |
| tomeh                |                    | hello-world                  | hello-world                    | hello-world                                            | hello-w-o-r-l-d             |
|                      | Underscore Cases   | __camel_snake_kebab__        | _leading_snake_case_           | __trailing__underscore__                               |                             |
| ETTLE                |                    | CAMEL-SNAKE-KEBAB-           | LEADING-SNAKE-CASE-            | TRAILING-UNDERSCORE-                                   |                             |
| ettle                |                    | camel-snake-kebab-           | leading-snake-case-            | trailing-underscore-                                   |                             |
| go-ettle             |                    | camel-snake-kebab            | leading-snake-case             | trailing-underscore                                    |                             |
| gobeam               |                    | camel-snake-kebab            | leading-snake-case             | trailing-underscore                                    |                             |
| golang-cz            |                    | camel-snake-kebab            | leading-snake-case             | trailing-underscore                                    |                             |
| iancoleman           |                    | --camel-snake-kebab--        | -leading-snake-case-           | --trailing--underscore--                               |                             |
| resenje              |                    | camel-snake-kebab            | leading-snake-case             | trailing-underscore                                    |                             |
| resenje-camel        |                    | Camel-Snake-Kebab            | Leading-Snake-Case             | Trailing-Underscore                                    |                             |
| resenje-snake        |                    | __Camel_Snake_Kebab__        | _Leading_Snake_Case_           | __Trailing_Underscore__                                |                             |
| screaming-iancoleman |                    | --CAMEL-SNAKE-KEBAB--        | -LEADING-SNAKE-CASE-           | --TRAILING--UNDERSCORE--                               |                             |
| screaming-resenje    |                    | CAMEL-SNAKE-KEBAB            | LEADING-SNAKE-CASE             | TRAILING-UNDERSCORE                                    |                             |
| searking             |                    | _-_camel-_snake-_kebab-_-_   | _leading-_snake-_case-_        | _-_trailing-_-_underscore-_-_                          |                             |
| tomeh                |                    | camel-snake-kebab            | leading-snake-case             | trailing-underscore                                    |                             |
|                      | Hyphen Cases       | --camel-snake-kebab          | -leading-kebab-case-           | --trailing--hyphen--                                   |                             |
| ETTLE                |                    | CAMEL-SNAKE-KEBAB            | LEADING-KEBAB-CASE-            | TRAILING-HYPHEN-                                       |                             |
| ettle                |                    | camel-snake-kebab            | leading-kebab-case-            | trailing-hyphen-                                       |                             |
| go-ettle             |                    | camel-snake-kebab            | leading-kebab-case             | trailing-hyphen                                        |                             |
| gobeam               |                    | camel-snake-kebab            | leading-kebab-case             | trailing-hyphen                                        |                             |
| golang-cz            |                    | camel-snake-kebab            | leading-kebab-case             | trailing-hyphen                                        |                             |
| iancoleman           |                    | --camel-snake-kebab          | -leading-kebab-case-           | --trailing--hyphen--                                   |                             |
| resenje              |                    | --camel-snake-kebab          | -leading-kebab-case-           | --trailing-hyphen--                                    |                             |
| resenje-camel        |                    | --Camel-Snake-Kebab          | -Leading-Kebab-Case-           | --Trailing-Hyphen--                                    |                             |
| resenje-snake        |                    | Camel_Snake_Kebab            | Leading_Kebab_Case             | Trailing_Hyphen                                        |                             |
| screaming-iancoleman |                    | --CAMEL-SNAKE-KEBAB          | -LEADING-KEBAB-CASE-           | --TRAILING--HYPHEN--                                   |                             |
| screaming-resenje    |                    | --CAMEL-SNAKE-KEBAB          | -LEADING-KEBAB-CASE-           | --TRAILING-HYPHEN--                                    |                             |
| searking             |                    | ---camel--snake--kebab       | -leading--kebab--case--        | ---trailing----hyphen----                              |                             |
| tomeh                |                    | camel-snake-kebab            | leading-kebab-case             | trailing-hyphen                                        |                             |
|                      | Mixed Delimiters   | hello_world-and-kebab        | Mixed_Snake-Kebab--Case        | UPPER_snake-KEBAB_Case                                 |                             |
| ETTLE                |                    | HELLO-WORLD-AND-KEBAB        | MIXED-SNAKE-KEBAB-CASE         | UPPER-SNAKE-KEBAB-CASE                                 |                             |
| ettle                |                    | hello-world-and-kebab        | mixed-snake-kebab-case         | upper-snake-kebab-case                                 |                             |
| go-ettle             |                    | hello-world-and-kebab        | mixed-snake-kebab-case         | upper-snake-kebab-case                                 |                             |
| gobeam               |                    | hello-world-and-kebab        | Mixed-Snake-Kebab-Case         | UPPER-snake-KEBAB-Case                                 |                             |
| golang-cz            |                    | hello-world-and-kebab        | mixed-snake-kebab-case         | upper-snake-kebab-case                                 |                             |
| iancoleman           |                    | hello-world-and-kebab        | mixed-snake-kebab--case        | upper-snake-kebab-case                                 |                             |
| resenje              |                    | hello-world-and-kebab        | mixed-snake-kebab-case         | upper-snake-kebab-case                                 |                             |
| resenje-camel        |                    | Hello-World-And-Kebab        | Mixed-Snake-Kebab-Case         | Upper-Snake-Kebab-Case                                 |                             |
| resenje-snake        |                    | Hello_World_And_Kebab        | Mixed_Snake_Kebab_Case         | Upper_Snake_Kebab_Case                                 |                             |
| screaming-iancoleman |                    | HELLO-WORLD-AND-KEBAB        | MIXED-SNAKE-KEBAB--CASE        | UPPER-SNAKE-KEBAB-CASE                                 |                             |
| screaming-resenje    |                    | HELLO-WORLD-AND-KEBAB        | MIXED-SNAKE-KEBAB-CASE         | UPPER-SNAKE-KEBAB-CASE                                 |                             |
| searking             |                    | hello-_world--and--kebab     | mixed-_-snake---kebab-----case | u-p-p-e-r-_snake---k-e-b-a-b-_-case                    |                             |
| tomeh                |                    | hello-world-and-kebab        | mixed-snake-kebab-case         | u-p-p-e-r-snake-k-e-b-a-b-case                         |                             |
|                      | Acronym Handling   | ID_NUMBER                    | HTTP_Response_Code             | XML_HTTP_REQUEST                                       | API_Version_2               |
| ETTLE                |                    | ID-NUMBER                    | HTTP-RESPONSE-CODE             | XML-HTTP-REQUEST                                       | API-VERSION-2               |
| ettle                |                    | id-number                    | http-response-code             | xml-http-request                                       | api-version-2               |
| go-ettle             |                    | ID-number                    | HTTP-response-code             | XML-HTTP-request                                       | API-version-2               |
| gobeam               |                    | ID-NUMBER                    | HTTP-Response-Code             | XML-HTTP-REQUEST                                       | API-Version-2               |
| golang-cz            |                    | id-number                    | http-response-code             | xml-http-request                                       | api-version-2               |
| iancoleman           |                    | id-number                    | http-response-code             | xml-http-request                                       | api-version-2               |
| resenje              |                    | id-number                    | http-response-code             | xml-http-request                                       | api-version-2               |
| resenje-camel        |                    | Id-Number                    | Http-Response-Code             | Xml-Http-Request                                       | Api-Version-2               |
| resenje-snake        |                    | Id_Number                    | Http_Response_Code             | Xml_Http_Request                                       | Api_Version_2               |
| screaming-iancoleman |                    | ID-NUMBER                    | HTTP-RESPONSE-CODE             | XML-HTTP-REQUEST                                       | API-VERSION-2               |
| screaming-resenje    |                    | ID-NUMBER                    | HTTP-RESPONSE-CODE             | XML-HTTP-REQUEST                                       | API-VERSION-2               |
| searking             |                    | i-d-_-n-u-m-b-e-r            | h-t-t-p-_-response-_-code      | x-m-l-_-h-t-t-p-_-r-e-q-u-e-s-t                        | a-p-i-_-version-_-2         |
| tomeh                |                    | i-d-n-u-m-b-e-r              | h-t-t-p-response-code          | x-m-l-h-t-t-p-r-e-q-u-e-s-t                            | a-p-i-version-2             |
|                      | Spaces and Words   |    leading spaces            | trailing spaces                |    both sides                                          | This is a test sentence     |
| ETTLE                |                    | LEADING-SPACES               | TRAILING-SPACES                | BOTH-SIDES                                             | THIS-IS-A-TEST-SENTENCE     |
| ettle                |                    | leading-spaces               | trailing-spaces                | both-sides                                             | this-is-a-test-sentence     |
| go-ettle             |                    | leading-spaces               | trailing-spaces                | both-sides                                             | this-is-a-test-sentence     |
| gobeam               |                    | leading-spaces               | trailing-spaces                | both-sides                                             | This-is-a-test-sentence     |
| golang-cz            |                    | leading-spaces               | trailing-spaces                | both-sides                                             | this-is-a-test-sentence     |
| iancoleman           |                    | leading-spaces               | trailing-spaces                | both-sides                                             | this-is-a-test-sentence     |
| resenje              |                    | leading-spaces               | trailing-spaces                | both-sides                                             | this-is-a-test-sentence     |
| resenje-camel        |                    | Leading-Spaces               | Trailing-Spaces                | Both-Sides                                             | This-Is-A-Test-Sentence     |
| resenje-snake        |                    | Leading_Spaces               | Trailing_Spaces                | Both_Sides                                             | This_Is_A_Test_Sentence     |
| screaming-iancoleman |                    | LEADING-SPACES               | TRAILING-SPACES                | BOTH-SIDES                                             | THIS-IS-A-TEST-SENTENCE     |
| screaming-resenje    |                    | LEADING-SPACES               | TRAILING-SPACES                | BOTH-SIDES                                             | THIS-IS-A-TEST-SENTENCE     |
| searking             |                    |  - - leading- spaces         | trailing- spaces- - -          |  - - both- sides- - -                                  | this- is- a- test- sentence |
| tomeh                |                    | leading-spaces               | trailing-spaces                | both-sides                                             | this-is-a-test-sentence     |
|                      | Symbols and Random | 123_ABC-xyz--789             | !!special$$characters**        | ___!!!___weird___CASE___!!!___                         |                             |
| ETTLE                |                    | 123-ABC-XYZ-789              | !!SPECIAL$$CHARACTERS**        | !!!-WEIRD-CASE-!!!-                                    |                             |
| ettle                |                    | 123-abc-xyz-789              | !!special$$characters**        | !!!-weird-case-!!!-                                    |                             |
| go-ettle             |                    | 123-abc-xyz-789              | !!special$$characters**        | !!!-weird-case-!!!                                     |                             |
| gobeam               |                    | 123-ABC-xyz-789              | !!special$$characters**        | !!!-weird-CASE-!!!                                     |                             |
| golang-cz            |                    | 123-abc-xyz-789              | special-characters             | weird-case                                             |                             |
| iancoleman           |                    | 123-abc-xyz--789             | !!special$$characters**        | ---!!!---weird---case---!!!---                         |                             |
| resenje              |                    | 123-abc-xyz-789              | !-!special-$-$characters-**    | !!!-weird-case-!!!                                     |                             |
| resenje-camel        |                    | 123-Abc-Xyz-789              | !-!special-$-$characters-**    | !!!-Weird-Case-!!!                                     |                             |
| resenje-snake        |                    | 123_Abc_Xyz_789              | !_!special_$_$characters_**    | ___!!!_Weird_Case_!!!___                               |                             |
| screaming-iancoleman |                    | 123-ABC-XYZ--789             | !!SPECIAL$$CHARACTERS**        | ---!!!---WEIRD---CASE---!!!---                         |                             |
| screaming-resenje    |                    | 123-ABC-XYZ-789              | !-!SPECIAL-$-$CHARACTERS-**    | !!!-WEIRD-CASE-!!!                                     |                             |
| searking             |                    | 1-2-3-_-a-b-c--xyz-----7-8-9 | !-!special-$-$characters-*-*   | _-_-_-!-!-!-_-_-_weird-_-_-_-c-a-s-e-_-_-_-!-!-!-_-_-_ |                             |
| tomeh                |                    | 123-a-b-c-xyz-789            | !!special-characters**         | !!!-weird-c-a-s-e-!!!                                  |                             |
+----------------------+--------------------+------------------------------+--------------------------------+--------------------------------------------------------+-----------------------------+
exit status 1
```


## `rntolower`

```bash
$ rntolower
Error: No files provided.
Usage: rntolower [options] <file1> [<file2> ...]

Options:
  -algorithm string
    	Choose the lowercase algorithm to use, supported: resenje,gobeam,searking,go. (default "resenje")
  -dry-run
    	Display the intended changes without renaming.
  -interactive
    	Ask for confirmation before renaming each file.

Conversion Examples:
+-----------+--------------------+-----------------------+-----------------------------+--------------------------------+-------------------------+
| ALGORITHM |                    |                       |                             |                                |                         |
+-----------+--------------------+-----------------------+-----------------------------+--------------------------------+-------------------------+
|           | Basic Cases        | Hello World           | helloWorld                  | HelloWorld                     | Hello_WORLD             |
| go        |                    | hello world           | helloworld                  | helloworld                     | hello_world             |
| gobeam    |                    | hello world           | helloworld                  | helloworld                     | hello_world             |
| resenje   |                    | hello world           | hello world                 | hello world                    | hello world             |
| searking  |                    | hello world           | helloworld                  | helloworld                     | hello_world             |
|           | Underscore Cases   | __camel_snake_kebab__ | _leading_snake_case_        | __trailing__underscore__       |                         |
| go        |                    | __camel_snake_kebab__ | _leading_snake_case_        | __trailing__underscore__       |                         |
| gobeam    |                    | __camel_snake_kebab__ | _leading_snake_case_        | __trailing__underscore__       |                         |
| resenje   |                    | camel snake kebab     | leading snake case          | trailing underscore            |                         |
| searking  |                    | __camel_snake_kebab__ | _leading_snake_case_        | __trailing__underscore__       |                         |
|           | Hyphen Cases       | --camel-snake-kebab   | -leading-kebab-case-        | --trailing--hyphen--           |                         |
| go        |                    | --camel-snake-kebab   | -leading-kebab-case-        | --trailing--hyphen--           |                         |
| gobeam    |                    | --camel-snake-kebab   | -leading-kebab-case-        | --trailing--hyphen--           |                         |
| resenje   |                    | camel snake kebab     | leading kebab case          | trailing hyphen                |                         |
| searking  |                    | --camel-snake-kebab   | -leading-kebab-case-        | --trailing--hyphen--           |                         |
|           | Mixed Delimiters   | hello_world-and-kebab | Mixed_Snake-Kebab--Case     | UPPER_snake-KEBAB_Case         |                         |
| go        |                    | hello_world-and-kebab | mixed_snake-kebab--case     | upper_snake-kebab_case         |                         |
| gobeam    |                    | hello_world-and-kebab | mixed_snake-kebab--case     | upper_snake-kebab_case         |                         |
| resenje   |                    | hello world and kebab | mixed snake kebab case      | upper snake kebab case         |                         |
| searking  |                    | hello_world-and-kebab | mixed_snake-kebab--case     | upper_snake-kebab_case         |                         |
|           | Acronym Handling   | ID_NUMBER             | HTTP_Response_Code          | XML_HTTP_REQUEST               | API_Version_2           |
| go        |                    | id_number             | http_response_code          | xml_http_request               | api_version_2           |
| gobeam    |                    | id_number             | http_response_code          | xml_http_request               | api_version_2           |
| resenje   |                    | id number             | http response code          | xml http request               | api version 2           |
| searking  |                    | id_number             | http_response_code          | xml_http_request               | api_version_2           |
|           | Spaces and Words   |    leading spaces     | trailing spaces             |    both sides                  | This is a test sentence |
| go        |                    |    leading spaces     | trailing spaces             |    both sides                  | this is a test sentence |
| gobeam    |                    |    leading spaces     | trailing spaces             |    both sides                  | this is a test sentence |
| resenje   |                    | leading spaces        | trailing spaces             | both sides                     | this is a test sentence |
| searking  |                    |    leading spaces     | trailing spaces             |    both sides                  | this is a test sentence |
|           | Symbols and Random | 123_ABC-xyz--789      | !!special$$characters**     | ___!!!___weird___CASE___!!!___ |                         |
| go        |                    | 123_abc-xyz--789      | !!special$$characters**     | ___!!!___weird___case___!!!___ |                         |
| gobeam    |                    | 123_abc-xyz--789      | !!special$$characters**     | ___!!!___weird___case___!!!___ |                         |
| resenje   |                    | 123 abc xyz 789       | ! !special $ $characters ** | !!! weird case !!!             |                         |
| searking  |                    | 123_abc-xyz--789      | !!special$$characters**     | ___!!!___weird___case___!!!___ |                         |
+-----------+--------------------+-----------------------+-----------------------------+--------------------------------+-------------------------+
exit status 1
```


## `rntolowerleading`

```bash
$ rntolowerleading
Error: No files provided.
Usage: rntolowerleading [options] <file1> [<file2> ...]

Options:
  -algorithm string
    	Choose the lowerleading algorithm to use, supported: gobeam,searking. (default "searking")
  -dry-run
    	Display the intended changes without renaming.
  -interactive
    	Ask for confirmation before renaming each file.

Conversion Examples:
+-----------+--------------------+-----------------------+-------------------------+--------------------------------+-------------------------+
| ALGORITHM |                    |                       |                         |                                |                         |
+-----------+--------------------+-----------------------+-------------------------+--------------------------------+-------------------------+
|           | Basic Cases        | Hello World           | helloWorld              | HelloWorld                     | Hello_WORLD             |
| gobeam    |                    | hello World           | helloWorld              | helloWorld                     | hello_WORLD             |
| searking  |                    | hello World           | helloWorld              | helloWorld                     | hello_WORLD             |
|           | Underscore Cases   | __camel_snake_kebab__ | _leading_snake_case_    | __trailing__underscore__       |                         |
| gobeam    |                    | __camel_snake_kebab__ | _leading_snake_case_    | __trailing__underscore__       |                         |
| searking  |                    | __camel_snake_kebab__ | _leading_snake_case_    | __trailing__underscore__       |                         |
|           | Hyphen Cases       | --camel-snake-kebab   | -leading-kebab-case-    | --trailing--hyphen--           |                         |
| gobeam    |                    | --camel-snake-kebab   | -leading-kebab-case-    | --trailing--hyphen--           |                         |
| searking  |                    | --camel-snake-kebab   | -leading-kebab-case-    | --trailing--hyphen--           |                         |
|           | Mixed Delimiters   | hello_world-and-kebab | Mixed_Snake-Kebab--Case | UPPER_snake-KEBAB_Case         |                         |
| gobeam    |                    | hello_world-and-kebab | mixed_Snake-Kebab--Case | uPPER_snake-KEBAB_Case         |                         |
| searking  |                    | hello_world-and-kebab | mixed_Snake-Kebab--Case | uPPER_snake-KEBAB_Case         |                         |
|           | Acronym Handling   | ID_NUMBER             | HTTP_Response_Code      | XML_HTTP_REQUEST               | API_Version_2           |
| gobeam    |                    | iD_NUMBER             | hTTP_Response_Code      | xML_HTTP_REQUEST               | aPI_Version_2           |
| searking  |                    | iD_NUMBER             | hTTP_Response_Code      | xML_HTTP_REQUEST               | aPI_Version_2           |
|           | Spaces and Words   |    leading spaces     | trailing spaces         |    both sides                  | This is a test sentence |
| gobeam    |                    |    leading spaces     | trailing spaces         |    both sides                  | this is a test sentence |
| searking  |                    |    leading spaces     | trailing spaces         |    both sides                  | this is a test sentence |
|           | Symbols and Random | 123_ABC-xyz--789      | !!special$$characters** | ___!!!___weird___CASE___!!!___ |                         |
| gobeam    |                    | 123_ABC-xyz--789      | !!special$$characters** | ___!!!___weird___CASE___!!!___ |                         |
| searking  |                    | 123_ABC-xyz--789      | !!special$$characters** | ___!!!___weird___CASE___!!!___ |                         |
+-----------+--------------------+-----------------------+-------------------------+--------------------------------+-------------------------+
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
    	Choose the pascal algorithm to use, supported: searKing,ettle,go-ettle,resenje,tomeh,golang-cz,gobeam. (default "ettle")
  -dry-run
    	Display the intended changes without renaming.
  -interactive
    	Ask for confirmation before renaming each file.
  -word-seperators value
    	Word separators. (gobeam only) Default: "_-. "

Conversion Examples:
+-----------+--------------------+-----------------------+-------------------------+--------------------------------+-------------------------+
| ALGORITHM |                    |                       |                         |                                |                         |
+-----------+--------------------+-----------------------+-------------------------+--------------------------------+-------------------------+
|           | Basic Cases        | Hello World           | helloWorld              | HelloWorld                     | Hello_WORLD             |
| ettle     |                    | HelloWorld            | HelloWorld              | HelloWorld                     | HelloWorld              |
| go-ettle  |                    | HelloWorld            | HelloWorld              | HelloWorld                     | HelloWorld              |
| gobeam    |                    | HelloWorld            | HelloWorld              | HelloWorld                     | HelloWORLD              |
| golang-cz |                    | HelloWorld            | HelloWorld              | HelloWorld                     | HelloWorld              |
| resenje   |                    | HelloWorld            | HelloWorld              | HelloWorld                     | HelloWorld              |
| searKing  |                    | Hello World           | HelloWorld              | HelloWorld                     | Hello_WORLD             |
| tomeh     |                    | HelloWorld            | HelloWorld              | HelloWorld                     | HelloWORLD              |
|           | Underscore Cases   | __camel_snake_kebab__ | _leading_snake_case_    | __trailing__underscore__       |                         |
| ettle     |                    | CamelSnakeKebab       | LeadingSnakeCase        | TrailingUnderscore             |                         |
| go-ettle  |                    | CamelSnakeKebab       | LeadingSnakeCase        | TrailingUnderscore             |                         |
| gobeam    |                    | CamelSnakeKebab       | LeadingSnakeCase        | TrailingUnderscore             |                         |
| golang-cz |                    | CamelSnakeKebab       | LeadingSnakeCase        | TrailingUnderscore             |                         |
| resenje   |                    | CamelSnakeKebab       | LeadingSnakeCase        | TrailingUnderscore             |                         |
| searKing  |                    | __camel_snake_kebab__ | _leading_snake_case_    | __trailing__underscore__       |                         |
| tomeh     |                    | CamelSnakeKebab       | LeadingSnakeCase        | TrailingUnderscore             |                         |
|           | Hyphen Cases       | --camel-snake-kebab   | -leading-kebab-case-    | --trailing--hyphen--           |                         |
| ettle     |                    | CamelSnakeKebab       | LeadingKebabCase        | TrailingHyphen                 |                         |
| go-ettle  |                    | CamelSnakeKebab       | LeadingKebabCase        | TrailingHyphen                 |                         |
| gobeam    |                    | CamelSnakeKebab       | LeadingKebabCase        | TrailingHyphen                 |                         |
| golang-cz |                    | CamelSnakeKebab       | LeadingKebabCase        | TrailingHyphen                 |                         |
| resenje   |                    | CamelSnakeKebab       | LeadingKebabCase        | TrailingHyphen                 |                         |
| searKing  |                    | --camel-snake-kebab   | -leading-kebab-case-    | --trailing--hyphen--           |                         |
| tomeh     |                    | CamelSnakeKebab       | LeadingKebabCase        | TrailingHyphen                 |                         |
|           | Mixed Delimiters   | hello_world-and-kebab | Mixed_Snake-Kebab--Case | UPPER_snake-KEBAB_Case         |                         |
| ettle     |                    | HelloWorldAndKebab    | MixedSnakeKebabCase     | UpperSnakeKebabCase            |                         |
| go-ettle  |                    | HelloWorldAndKebab    | MixedSnakeKebabCase     | UpperSnakeKebabCase            |                         |
| gobeam    |                    | HelloWorldAndKebab    | MixedSnakeKebabCase     | UPPERSnakeKEBABCase            |                         |
| golang-cz |                    | HelloWorldAndKebab    | MixedSnakeKebabCase     | UpperSnakeKebabCase            |                         |
| resenje   |                    | HelloWorldAndKebab    | MixedSnakeKebabCase     | UpperSnakeKebabCase            |                         |
| searKing  |                    | Hello_world-and-kebab | Mixed_Snake-Kebab--Case | UPPER_snake-KEBAB_Case         |                         |
| tomeh     |                    | HelloWorldAndKebab    | MixedSnakeKebabCase     | UPPERSnakeKEBABCase            |                         |
|           | Acronym Handling   | ID_NUMBER             | HTTP_Response_Code      | XML_HTTP_REQUEST               | API_Version_2           |
| ettle     |                    | IdNumber              | HttpResponseCode        | XmlHttpRequest                 | ApiVersion2             |
| go-ettle  |                    | IDNumber              | HTTPResponseCode        | XMLHTTPRequest                 | APIVersion2             |
| gobeam    |                    | IDNUMBER              | HTTPResponseCode        | XMLHTTPREQUEST                 | APIVersion2             |
| golang-cz |                    | IdNumber              | HttpResponseCode        | XmlHttpRequest                 | ApiVersion2             |
| resenje   |                    | IdNumber              | HttpResponseCode        | XmlHttpRequest                 | ApiVersion2             |
| searKing  |                    | ID_NUMBER             | HTTP_Response_Code      | XML_HTTP_REQUEST               | API_Version_2           |
| tomeh     |                    | IDNUMBER              | HTTPResponseCode        | XMLHTTPREQUEST                 | APIVersion2             |
|           | Spaces and Words   |    leading spaces     | trailing spaces         |    both sides                  | This is a test sentence |
| ettle     |                    | LeadingSpaces         | TrailingSpaces          | BothSides                      | ThisIsATestSentence     |
| go-ettle  |                    | LeadingSpaces         | TrailingSpaces          | BothSides                      | ThisIsATestSentence     |
| gobeam    |                    | LeadingSpaces         | TrailingSpaces          | BothSides                      | ThisIsATestSentence     |
| golang-cz |                    | LeadingSpaces         | TrailingSpaces          | BothSides                      | ThisIsATestSentence     |
| resenje   |                    | LeadingSpaces         | TrailingSpaces          | BothSides                      | ThisIsATestSentence     |
| searKing  |                    |    leading spaces     | Trailing spaces         |    both sides                  | This is a test sentence |
| tomeh     |                    | LeadingSpaces         | TrailingSpaces          | BothSides                      | ThisIsATestSentence     |
|           | Symbols and Random | 123_ABC-xyz--789      | !!special$$characters** | ___!!!___weird___CASE___!!!___ |                         |
| ettle     |                    | 123AbcXyz789          | !!special$$characters** | !!!WeirdCase!!!                |                         |
| go-ettle  |                    | 123AbcXyz789          | !!special$$characters** | !!!WeirdCase!!!                |                         |
| gobeam    |                    | 123ABCXyz789          | !!special$$characters** | !!!WeirdCASE!!!                |                         |
| golang-cz |                    | 123AbcXyz789          | SpecialCharacters       | WeirdCase                      |                         |
| resenje   |                    | 123AbcXyz789          | !!special$$characters** | !!!WeirdCase!!!                |                         |
| searKing  |                    | 123_ABC-xyz--789      | !!special$$characters** | ___!!!___weird___CASE___!!!___ |                         |
| tomeh     |                    | 123ABCXyz789          | !!SpecialCharacters**   | !!!WeirdCASE!!!                |                         |
+-----------+--------------------+-----------------------+-------------------------+--------------------------------+-------------------------+
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
    	Choose the snake algorithm to use, supported: resenje,screaming-resenje,tomeh,screaming-iancoleman,lowercase-searking,ettle,go-ettle,ETTLE,golang-cz,gobeam,iancoleman,searking. (default "iancoleman")
  -dry-run
    	Display the intended changes without renaming.
  -interactive
    	Ask for confirmation before renaming each file.
  -word-seperators value
    	Word separators. (gobeam only) Default: "_-. "

Conversion Examples:
+----------------------+--------------------+----------------------------+------------------------------+----------------------------------+-----------------------------+
| ALGORITHM            |                    |                            |                              |                                  |                             |
+----------------------+--------------------+----------------------------+------------------------------+----------------------------------+-----------------------------+
|                      | Basic Cases        | Hello World                | helloWorld                   | HelloWorld                       | Hello_WORLD                 |
| ETTLE                |                    | HELLO_WORLD                | HELLO_WORLD                  | HELLO_WORLD                      | HELLO_WORLD                 |
| ettle                |                    | hello_world                | hello_world                  | hello_world                      | hello_world                 |
| go-ettle             |                    | hello_world                | hello_world                  | hello_world                      | hello_world                 |
| gobeam               |                    | Hello_World                | hello_World                  | Hello_World                      | Hello_WORLD                 |
| golang-cz            |                    | hello_world                | hello_world                  | hello_world                      | hello_world                 |
| iancoleman           |                    | hello_world                | hello_world                  | hello_world                      | hello_world                 |
| lowercase-searking   |                    | hello_ _world              | hello_world                  | hello_world                      | hello_w_o_r_l_d             |
| resenje              |                    | hello_world                | hello_world                  | hello_world                      | hello_world                 |
| screaming-iancoleman |                    | HELLO_WORLD                | HELLO_WORLD                  | HELLO_WORLD                      | HELLO_WORLD                 |
| screaming-resenje    |                    | HELLO_WORLD                | HELLO_WORLD                  | HELLO_WORLD                      | HELLO_WORLD                 |
| searking             |                    | hello_ _world              | hello_world                  | hello_world                      | hello_w_o_r_l_d             |
| tomeh                |                    | hello_world                | hello_world                  | hello_world                      | hello_w_o_r_l_d             |
|                      | Underscore Cases   | __camel_snake_kebab__      | _leading_snake_case_         | __trailing__underscore__         |                             |
| ETTLE                |                    | CAMEL_SNAKE_KEBAB_         | LEADING_SNAKE_CASE_          | TRAILING_UNDERSCORE_             |                             |
| ettle                |                    | camel_snake_kebab_         | leading_snake_case_          | trailing_underscore_             |                             |
| go-ettle             |                    | camel_snake_kebab_         | leading_snake_case_          | trailing_underscore_             |                             |
| gobeam               |                    | camel_snake_kebab          | leading_snake_case           | trailing_underscore              |                             |
| golang-cz            |                    | camel_snake_kebab          | leading_snake_case           | trailing_underscore              |                             |
| iancoleman           |                    | __camel_snake_kebab__      | _leading_snake_case_         | __trailing__underscore__         |                             |
| lowercase-searking   |                    | x_camel_snake_kebab        | x_leading_snake_case__       | x_trailing_underscore            |                             |
| resenje              |                    | __camel_snake_kebab__      | _leading_snake_case_         | __trailing_underscore__          |                             |
| screaming-iancoleman |                    | __CAMEL_SNAKE_KEBAB__      | _LEADING_SNAKE_CASE_         | __TRAILING__UNDERSCORE__         |                             |
| screaming-resenje    |                    | __CAMEL_SNAKE_KEBAB__      | _LEADING_SNAKE_CASE_         | __TRAILING_UNDERSCORE__          |                             |
| searking             |                    | x_camel_snake_kebab        | x_leading_snake_case__       | x_trailing_underscore            |                             |
| tomeh                |                    | camel_snake_kebab          | leading_snake_case           | trailing_underscore              |                             |
|                      | Hyphen Cases       | --camel-snake-kebab        | -leading-kebab-case-         | --trailing--hyphen--             |                             |
| ETTLE                |                    | CAMEL_SNAKE_KEBAB          | LEADING_KEBAB_CASE_          | TRAILING_HYPHEN_                 |                             |
| ettle                |                    | camel_snake_kebab          | leading_kebab_case_          | trailing_hyphen_                 |                             |
| go-ettle             |                    | camel_snake_kebab          | leading_kebab_case_          | trailing_hyphen_                 |                             |
| gobeam               |                    | camel_snake_kebab          | leading_kebab_case           | trailing_hyphen                  |                             |
| golang-cz            |                    | camel_snake_kebab          | leading_kebab_case           | trailing_hyphen                  |                             |
| iancoleman           |                    | __camel_snake_kebab        | _leading_kebab_case_         | __trailing__hyphen__             |                             |
| lowercase-searking   |                    | -_-camel_-snake_-kebab     | -leading_-kebab_-case_-      | -_-trailing_-_-hyphen_-_-        |                             |
| resenje              |                    | camel_snake_kebab          | leading_kebab_case           | trailing_hyphen                  |                             |
| screaming-iancoleman |                    | __CAMEL_SNAKE_KEBAB        | _LEADING_KEBAB_CASE_         | __TRAILING__HYPHEN__             |                             |
| screaming-resenje    |                    | CAMEL_SNAKE_KEBAB          | LEADING_KEBAB_CASE           | TRAILING_HYPHEN                  |                             |
| searking             |                    | -_-camel_-snake_-kebab     | -leading_-kebab_-case_-      | -_-trailing_-_-hyphen_-_-        |                             |
| tomeh                |                    | camel_snake_kebab          | leading_kebab_case           | trailing_hyphen                  |                             |
|                      | Mixed Delimiters   | hello_world-and-kebab      | Mixed_Snake-Kebab--Case      | UPPER_snake-KEBAB_Case           |                             |
| ETTLE                |                    | HELLO_WORLD_AND_KEBAB      | MIXED_SNAKE_KEBAB_CASE       | UPPER_SNAKE_KEBAB_CASE           |                             |
| ettle                |                    | hello_world_and_kebab      | mixed_snake_kebab_case       | upper_snake_kebab_case           |                             |
| go-ettle             |                    | hello_world_and_kebab      | mixed_snake_kebab_case       | upper_snake_kebab_case           |                             |
| gobeam               |                    | hello_world_and_kebab      | Mixed_Snake_Kebab_Case       | UPPER_snake_KEBAB_Case           |                             |
| golang-cz            |                    | hello_world_and_kebab      | mixed_snake_kebab_case       | upper_snake_kebab_case           |                             |
| iancoleman           |                    | hello_world_and_kebab      | mixed_snake_kebab__case      | upper_snake_kebab_case           |                             |
| lowercase-searking   |                    | hello_world_-and_-kebab    | mixed_snake_-_kebab_-_-_case | u_p_p_e_r_snake_-_k_e_b_a_b_case |                             |
| resenje              |                    | hello_world_and_kebab      | mixed_snake_kebab_case       | upper_snake_kebab_case           |                             |
| screaming-iancoleman |                    | HELLO_WORLD_AND_KEBAB      | MIXED_SNAKE_KEBAB__CASE      | UPPER_SNAKE_KEBAB_CASE           |                             |
| screaming-resenje    |                    | HELLO_WORLD_AND_KEBAB      | MIXED_SNAKE_KEBAB_CASE       | UPPER_SNAKE_KEBAB_CASE           |                             |
| searking             |                    | hello_world_-and_-kebab    | mixed_snake_-_kebab_-_-_case | u_p_p_e_r_snake_-_k_e_b_a_b_case |                             |
| tomeh                |                    | hello_world_and_kebab      | mixed_snake_kebab_case       | u_p_p_e_r_snake_k_e_b_a_b_case   |                             |
|                      | Acronym Handling   | ID_NUMBER                  | HTTP_Response_Code           | XML_HTTP_REQUEST                 | API_Version_2               |
| ETTLE                |                    | ID_NUMBER                  | HTTP_RESPONSE_CODE           | XML_HTTP_REQUEST                 | API_VERSION_2               |
| ettle                |                    | id_number                  | http_response_code           | xml_http_request                 | api_version_2               |
| go-ettle             |                    | id_number                  | http_response_code           | xml_http_request                 | api_version_2               |
| gobeam               |                    | ID_NUMBER                  | HTTP_Response_Code           | XML_HTTP_REQUEST                 | API_Version_2               |
| golang-cz            |                    | id_number                  | http_response_code           | xml_http_request                 | api_version_2               |
| iancoleman           |                    | id_number                  | http_response_code           | xml_http_request                 | api_version_2               |
| lowercase-searking   |                    | i_d_n_u_m_b_e_r            | h_t_t_p_response_code        | x_m_l_h_t_t_p_r_e_q_u_e_s_t      | a_p_i_version_2             |
| resenje              |                    | id_number                  | http_response_code           | xml_http_request                 | api_version_2               |
| screaming-iancoleman |                    | ID_NUMBER                  | HTTP_RESPONSE_CODE           | XML_HTTP_REQUEST                 | API_VERSION_2               |
| screaming-resenje    |                    | ID_NUMBER                  | HTTP_RESPONSE_CODE           | XML_HTTP_REQUEST                 | API_VERSION_2               |
| searking             |                    | i_d_n_u_m_b_e_r            | h_t_t_p_response_code        | x_m_l_h_t_t_p_r_e_q_u_e_s_t      | a_p_i_version_2             |
| tomeh                |                    | i_d_n_u_m_b_e_r            | h_t_t_p_response_code        | x_m_l_h_t_t_p_r_e_q_u_e_s_t      | a_p_i_version_2             |
|                      | Spaces and Words   |    leading spaces          | trailing spaces              |    both sides                    | This is a test sentence     |
| ETTLE                |                    | LEADING_SPACES             | TRAILING_SPACES              | BOTH_SIDES                       | THIS_IS_A_TEST_SENTENCE     |
| ettle                |                    | leading_spaces             | trailing_spaces              | both_sides                       | this_is_a_test_sentence     |
| go-ettle             |                    | leading_spaces             | trailing_spaces              | both_sides                       | this_is_a_test_sentence     |
| gobeam               |                    | leading_spaces             | trailing_spaces              | both_sides                       | This_is_a_test_sentence     |
| golang-cz            |                    | leading_spaces             | trailing_spaces              | both_sides                       | this_is_a_test_sentence     |
| iancoleman           |                    | leading_spaces             | trailing_spaces              | both_sides                       | this_is_a_test_sentence     |
| lowercase-searking   |                    |  _ _ leading_ spaces       | trailing_ spaces_ _ _        |  _ _ both_ sides_ _ _            | this_ is_ a_ test_ sentence |
| resenje              |                    | leading_spaces             | trailing_spaces              | both_sides                       | this_is_a_test_sentence     |
| screaming-iancoleman |                    | LEADING_SPACES             | TRAILING_SPACES              | BOTH_SIDES                       | THIS_IS_A_TEST_SENTENCE     |
| screaming-resenje    |                    | LEADING_SPACES             | TRAILING_SPACES              | BOTH_SIDES                       | THIS_IS_A_TEST_SENTENCE     |
| searking             |                    |  _ _ leading_ spaces       | trailing_ spaces_ _ _        |  _ _ both_ sides_ _ _            | this_ is_ a_ test_ sentence |
| tomeh                |                    | leading_spaces             | trailing_spaces              | both_sides                       | this_is_a_test_sentence     |
|                      | Symbols and Random | 123_ABC-xyz--789           | !!special$$characters**      | ___!!!___weird___CASE___!!!___   |                             |
| ETTLE                |                    | 123_ABC_XYZ_789            | !!SPECIAL$$CHARACTERS**      | !!!_WEIRD_CASE_!!!_              |                             |
| ettle                |                    | 123_abc_xyz_789            | !!special$$characters**      | !!!_weird_case_!!!_              |                             |
| go-ettle             |                    | 123_abc_xyz_789            | !!special$$characters**      | !!!_weird_case_!!!_              |                             |
| gobeam               |                    | 123_ABC_xyz_789            | !!special$$characters**      | !!!_weird_CASE_!!!               |                             |
| golang-cz            |                    | 123_abc_xyz_789            | special_characters           | weird_case                       |                             |
| iancoleman           |                    | 123_abc_xyz__789           | !!special$$characters**      | ___!!!___weird___case___!!!___   |                             |
| lowercase-searking   |                    | 1_2_3_a_b_c_-xyz_-_-_7_8_9 | !_!special_$_$characters_*_* | x_!_!_!_weird_c_a_s_e_!_!_!      |                             |
| resenje              |                    | 123_abc_xyz_789            | !_!special_$_$characters_**  | ___!!!_weird_case_!!!___         |                             |
| screaming-iancoleman |                    | 123_ABC_XYZ__789           | !!SPECIAL$$CHARACTERS**      | ___!!!___WEIRD___CASE___!!!___   |                             |
| screaming-resenje    |                    | 123_ABC_XYZ_789            | !_!SPECIAL_$_$CHARACTERS_**  | ___!!!_WEIRD_CASE_!!!___         |                             |
| searking             |                    | 1_2_3_a_b_c_-xyz_-_-_7_8_9 | !_!special_$_$characters_*_* | x_!_!_!_weird_c_a_s_e_!_!_!      |                             |
| tomeh                |                    | 123_a_b_c_xyz_789          | !!special_characters**       | !!!_weird_c_a_s_e_!!!            |                             |
+----------------------+--------------------+----------------------------+------------------------------+----------------------------------+-----------------------------+
exit status 1
```


## `rntotitle`

```bash
$ rntotitle
Error: No files provided.
Usage: rntotitle [options] <file1> [<file2> ...]

Options:
  -algorithm string
    	Choose the titlecase algorithm to use, supported: revett,searking,go,resenje,gobeam. (default "revett")
  -dry-run
    	Display the intended changes without renaming.
  -interactive
    	Ask for confirmation before renaming each file.

Conversion Examples:
+-----------+--------------------+-----------------------+-----------------------------+--------------------------------+-------------------------+
| ALGORITHM |                    |                       |                             |                                |                         |
+-----------+--------------------+-----------------------+-----------------------------+--------------------------------+-------------------------+
|           | Basic Cases        | Hello World           | helloWorld                  | HelloWorld                     | Hello_WORLD             |
| go        |                    | HELLO WORLD           | HELLOWORLD                  | HELLOWORLD                     | HELLO_WORLD             |
| gobeam    |                    | Hello World           | Helloworld                  | Helloworld                     | Hello_world             |
| resenje   |                    | Hello World           | Hello World                 | Hello World                    | Hello World             |
| revett    |                    | Hello World           | Helloworld                  | Helloworld                     | Hello_world             |
| searking  |                    | HELLO WORLD           | HELLOWORLD                  | HELLOWORLD                     | HELLO_WORLD             |
|           | Underscore Cases   | __camel_snake_kebab__ | _leading_snake_case_        | __trailing__underscore__       |                         |
| go        |                    | __CAMEL_SNAKE_KEBAB__ | _LEADING_SNAKE_CASE_        | __TRAILING__UNDERSCORE__       |                         |
| gobeam    |                    | __camel_snake_kebab__ | _leading_snake_case_        | __trailing__underscore__       |                         |
| resenje   |                    | Camel Snake Kebab     | Leading Snake Case          | Trailing Underscore            |                         |
| revett    |                    | __camel_snake_kebab__ | _leading_snake_case_        | __trailing__underscore__       |                         |
| searking  |                    | __CAMEL_SNAKE_KEBAB__ | _LEADING_SNAKE_CASE_        | __TRAILING__UNDERSCORE__       |                         |
|           | Hyphen Cases       | --camel-snake-kebab   | -leading-kebab-case-        | --trailing--hyphen--           |                         |
| go        |                    | --CAMEL-SNAKE-KEBAB   | -LEADING-KEBAB-CASE-        | --TRAILING--HYPHEN--           |                         |
| gobeam    |                    | --camel-snake-kebab   | -leading-kebab-case-        | --trailing--hyphen--           |                         |
| resenje   |                    | Camel Snake Kebab     | Leading Kebab Case          | Trailing Hyphen                |                         |
| revett    |                    | --camel-snake-kebab   | -leading-kebab-case-        | --trailing--hyphen--           |                         |
| searking  |                    | --CAMEL-SNAKE-KEBAB   | -LEADING-KEBAB-CASE-        | --TRAILING--HYPHEN--           |                         |
|           | Mixed Delimiters   | hello_world-and-kebab | Mixed_Snake-Kebab--Case     | UPPER_snake-KEBAB_Case         |                         |
| go        |                    | HELLO_WORLD-AND-KEBAB | MIXED_SNAKE-KEBAB--CASE     | UPPER_SNAKE-KEBAB_CASE         |                         |
| gobeam    |                    | Hello_world-and-kebab | Mixed_snake-kebab--case     | Upper_snake-kebab_case         |                         |
| resenje   |                    | Hello World And Kebab | Mixed Snake Kebab Case      | Upper Snake Kebab Case         |                         |
| revett    |                    | Hello_world-and-kebab | Mixed_snake-kebab--case     | Upper_snake-kebab_case         |                         |
| searking  |                    | HELLO_WORLD-AND-KEBAB | MIXED_SNAKE-KEBAB--CASE     | UPPER_SNAKE-KEBAB_CASE         |                         |
|           | Acronym Handling   | ID_NUMBER             | HTTP_Response_Code          | XML_HTTP_REQUEST               | API_Version_2           |
| go        |                    | ID_NUMBER             | HTTP_RESPONSE_CODE          | XML_HTTP_REQUEST               | API_VERSION_2           |
| gobeam    |                    | Id_number             | Http_response_code          | Xml_http_request               | Api_version_2           |
| resenje   |                    | Id Number             | Http Response Code          | Xml Http Request               | Api Version 2           |
| revett    |                    | ID_NUMBER             | Http_response_code          | XML_HTTP_REQUEST               | Api_version_2           |
| searking  |                    | ID_NUMBER             | HTTP_RESPONSE_CODE          | XML_HTTP_REQUEST               | API_VERSION_2           |
|           | Spaces and Words   |    leading spaces     | trailing spaces             |    both sides                  | This is a test sentence |
| go        |                    |    LEADING SPACES     | TRAILING SPACES             |    BOTH SIDES                  | THIS IS A TEST SENTENCE |
| gobeam    |                    | !!!Error!!!           | !!!Error!!!                 | !!!Error!!!                    | This Is A Test Sentence |
| resenje   |                    | Leading Spaces        | Trailing Spaces             | Both Sides                     | This Is A Test Sentence |
| revett    |                    | Leading Spaces        | Trailing Spaces             | Both Sides                     | This Is a Test Sentence |
| searking  |                    |    LEADING SPACES     | TRAILING SPACES             |    BOTH SIDES                  | THIS IS A TEST SENTENCE |
|           | Symbols and Random | 123_ABC-xyz--789      | !!special$$characters**     | ___!!!___weird___CASE___!!!___ |                         |
| go        |                    | 123_ABC-XYZ--789      | !!SPECIAL$$CHARACTERS**     | ___!!!___WEIRD___CASE___!!!___ |                         |
| gobeam    |                    | 123_abc-xyz--789      | !!special$$characters**     | ___!!!___weird___case___!!!___ |                         |
| resenje   |                    | 123 Abc Xyz 789       | ! !special $ $characters ** | !!! Weird Case !!!             |                         |
| revett    |                    | 123_abc-xyz--789      | !!special$$characters**     | ___!!!___weird___case___!!!___ |                         |
| searking  |                    | 123_ABC-XYZ--789      | !!SPECIAL$$CHARACTERS**     | ___!!!___WEIRD___CASE___!!!___ |                         |
+-----------+--------------------+-----------------------+-----------------------------+--------------------------------+-------------------------+
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
+-----------+--------------------+-----------------------+-----------------------------+--------------------------------+-------------------------+
| ALGORITHM |                    |                       |                             |                                |                         |
+-----------+--------------------+-----------------------+-----------------------------+--------------------------------+-------------------------+
|           | Basic Cases        | Hello World           | helloWorld                  | HelloWorld                     | Hello_WORLD             |
| go        |                    | HELLO WORLD           | HELLOWORLD                  | HELLOWORLD                     | HELLO_WORLD             |
| gobeam    |                    | HELLO WORLD           | HELLOWORLD                  | HELLOWORLD                     | HELLO_WORLD             |
| resenje   |                    | HELLO WORLD           | HELLO WORLD                 | HELLO WORLD                    | HELLO WORLD             |
| searking  |                    | Hello World           | HelloWorld                  | HelloWorld                     | Hello_WORLD             |
|           | Underscore Cases   | __camel_snake_kebab__ | _leading_snake_case_        | __trailing__underscore__       |                         |
| go        |                    | __CAMEL_SNAKE_KEBAB__ | _LEADING_SNAKE_CASE_        | __TRAILING__UNDERSCORE__       |                         |
| gobeam    |                    | __CAMEL_SNAKE_KEBAB__ | _LEADING_SNAKE_CASE_        | __TRAILING__UNDERSCORE__       |                         |
| resenje   |                    | CAMEL SNAKE KEBAB     | LEADING SNAKE CASE          | TRAILING UNDERSCORE            |                         |
| searking  |                    | __camel_snake_kebab__ | _leading_snake_case_        | __trailing__underscore__       |                         |
|           | Hyphen Cases       | --camel-snake-kebab   | -leading-kebab-case-        | --trailing--hyphen--           |                         |
| go        |                    | --CAMEL-SNAKE-KEBAB   | -LEADING-KEBAB-CASE-        | --TRAILING--HYPHEN--           |                         |
| gobeam    |                    | --CAMEL-SNAKE-KEBAB   | -LEADING-KEBAB-CASE-        | --TRAILING--HYPHEN--           |                         |
| resenje   |                    | CAMEL SNAKE KEBAB     | LEADING KEBAB CASE          | TRAILING HYPHEN                |                         |
| searking  |                    | --camel-snake-kebab   | -leading-kebab-case-        | --trailing--hyphen--           |                         |
|           | Mixed Delimiters   | hello_world-and-kebab | Mixed_Snake-Kebab--Case     | UPPER_snake-KEBAB_Case         |                         |
| go        |                    | HELLO_WORLD-AND-KEBAB | MIXED_SNAKE-KEBAB--CASE     | UPPER_SNAKE-KEBAB_CASE         |                         |
| gobeam    |                    | HELLO_WORLD-AND-KEBAB | MIXED_SNAKE-KEBAB--CASE     | UPPER_SNAKE-KEBAB_CASE         |                         |
| resenje   |                    | HELLO WORLD AND KEBAB | MIXED SNAKE KEBAB CASE      | UPPER SNAKE KEBAB CASE         |                         |
| searking  |                    | Hello_world-and-kebab | Mixed_Snake-Kebab--Case     | UPPER_snake-KEBAB_Case         |                         |
|           | Acronym Handling   | ID_NUMBER             | HTTP_Response_Code          | XML_HTTP_REQUEST               | API_Version_2           |
| go        |                    | ID_NUMBER             | HTTP_RESPONSE_CODE          | XML_HTTP_REQUEST               | API_VERSION_2           |
| gobeam    |                    | ID_NUMBER             | HTTP_RESPONSE_CODE          | XML_HTTP_REQUEST               | API_VERSION_2           |
| resenje   |                    | ID NUMBER             | HTTP RESPONSE CODE          | XML HTTP REQUEST               | API VERSION 2           |
| searking  |                    | ID_NUMBER             | HTTP_Response_Code          | XML_HTTP_REQUEST               | API_Version_2           |
|           | Spaces and Words   |    leading spaces     | trailing spaces             |    both sides                  | This is a test sentence |
| go        |                    |    LEADING SPACES     | TRAILING SPACES             |    BOTH SIDES                  | THIS IS A TEST SENTENCE |
| gobeam    |                    |    LEADING SPACES     | TRAILING SPACES             |    BOTH SIDES                  | THIS IS A TEST SENTENCE |
| resenje   |                    | LEADING SPACES        | TRAILING SPACES             | BOTH SIDES                     | THIS IS A TEST SENTENCE |
| searking  |                    |    leading spaces     | Trailing spaces             |    both sides                  | This is a test sentence |
|           | Symbols and Random | 123_ABC-xyz--789      | !!special$$characters**     | ___!!!___weird___CASE___!!!___ |                         |
| go        |                    | 123_ABC-XYZ--789      | !!SPECIAL$$CHARACTERS**     | ___!!!___WEIRD___CASE___!!!___ |                         |
| gobeam    |                    | 123_ABC-XYZ--789      | !!SPECIAL$$CHARACTERS**     | ___!!!___WEIRD___CASE___!!!___ |                         |
| resenje   |                    | 123 ABC XYZ 789       | ! !SPECIAL $ $CHARACTERS ** | !!! WEIRD CASE !!!             |                         |
| searking  |                    | 123_ABC-xyz--789      | !!special$$characters**     | ___!!!___weird___CASE___!!!___ |                         |
+-----------+--------------------+-----------------------+-----------------------------+--------------------------------+-------------------------+
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
+-----------+--------------------+-----------------------+-------------------------+--------------------------------+-------------------------+
| ALGORITHM |                    |                       |                         |                                |                         |
+-----------+--------------------+-----------------------+-------------------------+--------------------------------+-------------------------+
|           | Basic Cases        | Hello World           | helloWorld              | HelloWorld                     | Hello_WORLD             |
| gobeam    |                    | Hello World           | HelloWorld              | HelloWorld                     | Hello_WORLD             |
| searking  |                    | Hello World           | HelloWorld              | HelloWorld                     | Hello_WORLD             |
|           | Underscore Cases   | __camel_snake_kebab__ | _leading_snake_case_    | __trailing__underscore__       |                         |
| gobeam    |                    | __camel_snake_kebab__ | _leading_snake_case_    | __trailing__underscore__       |                         |
| searking  |                    | __camel_snake_kebab__ | _leading_snake_case_    | __trailing__underscore__       |                         |
|           | Hyphen Cases       | --camel-snake-kebab   | -leading-kebab-case-    | --trailing--hyphen--           |                         |
| gobeam    |                    | --camel-snake-kebab   | -leading-kebab-case-    | --trailing--hyphen--           |                         |
| searking  |                    | --camel-snake-kebab   | -leading-kebab-case-    | --trailing--hyphen--           |                         |
|           | Mixed Delimiters   | hello_world-and-kebab | Mixed_Snake-Kebab--Case | UPPER_snake-KEBAB_Case         |                         |
| gobeam    |                    | Hello_world-and-kebab | Mixed_Snake-Kebab--Case | UPPER_snake-KEBAB_Case         |                         |
| searking  |                    | Hello_world-and-kebab | Mixed_Snake-Kebab--Case | UPPER_snake-KEBAB_Case         |                         |
|           | Acronym Handling   | ID_NUMBER             | HTTP_Response_Code      | XML_HTTP_REQUEST               | API_Version_2           |
| gobeam    |                    | ID_NUMBER             | HTTP_Response_Code      | XML_HTTP_REQUEST               | API_Version_2           |
| searking  |                    | ID_NUMBER             | HTTP_Response_Code      | XML_HTTP_REQUEST               | API_Version_2           |
|           | Spaces and Words   |    leading spaces     | trailing spaces         |    both sides                  | This is a test sentence |
| gobeam    |                    |    leading spaces     | Trailing spaces         |    both sides                  | This is a test sentence |
| searking  |                    |    leading spaces     | Trailing spaces         |    both sides                  | This is a test sentence |
|           | Symbols and Random | 123_ABC-xyz--789      | !!special$$characters** | ___!!!___weird___CASE___!!!___ |                         |
| gobeam    |                    | 123_ABC-xyz--789      | !!special$$characters** | ___!!!___weird___CASE___!!!___ |                         |
| searking  |                    | 123_ABC-xyz--789      | !!special$$characters** | ___!!!___weird___CASE___!!!___ |                         |
+-----------+--------------------+-----------------------+-------------------------+--------------------------------+-------------------------+
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
+-----------+--------------------+-----------------------+-------------------------+--------------------------------+-------------------------+
| ALGORITHM |                    |                       |                         |                                |                         |
+-----------+--------------------+-----------------------+-------------------------+--------------------------------+-------------------------+
|           | Basic Cases        | Hello World           | helloWorld              | HelloWorld                     | Hello_WORLD             |
| go        |                    | Hello World           | helloWorld              | HelloWorld                     | Hello_WORLD             |
|           | Underscore Cases   | __camel_snake_kebab__ | _leading_snake_case_    | __trailing__underscore__       |                         |
| go        |                    | __camel_snake_kebab__ | _leading_snake_case_    | __trailing__underscore__       |                         |
|           | Hyphen Cases       | --camel-snake-kebab   | -leading-kebab-case-    | --trailing--hyphen--           |                         |
| go        |                    | --camel-snake-kebab   | -leading-kebab-case-    | --trailing--hyphen--           |                         |
|           | Mixed Delimiters   | hello_world-and-kebab | Mixed_Snake-Kebab--Case | UPPER_snake-KEBAB_Case         |                         |
| go        |                    | hello_world-and-kebab | Mixed_Snake-Kebab--Case | UPPER_snake-KEBAB_Case         |                         |
|           | Acronym Handling   | ID_NUMBER             | HTTP_Response_Code      | XML_HTTP_REQUEST               | API_Version_2           |
| go        |                    | ID_NUMBER             | HTTP_Response_Code      | XML_HTTP_REQUEST               | API_Version_2           |
|           | Spaces and Words   |    leading spaces     | trailing spaces         |    both sides                  | This is a test sentence |
| go        |                    | leading spaces        | trailing spaces         | both sides                     | This is a test sentence |
|           | Symbols and Random | 123_ABC-xyz--789      | !!special$$characters** | ___!!!___weird___CASE___!!!___ |                         |
| go        |                    | 123_ABC-xyz--789      | !!special$$characters** | ___!!!___weird___CASE___!!!___ |                         |
+-----------+--------------------+-----------------------+-------------------------+--------------------------------+-------------------------+
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
