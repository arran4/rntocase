package rntocase

type ExampleGroup struct {
	Name   string   // Name of the example group
	Values []string // Input examples within the group
}

var (
	ExampleGroups = []ExampleGroup{
		{"Basic Cases", []string{"Hello World", "helloWorld", "HelloWorld", "Hello_WORLD"}},
		{"Underscore Cases", []string{"__camel_snake_kebab__", "_leading_snake_case_", "__trailing__underscore__"}},
		{"Hyphen Cases", []string{"--camel-snake-kebab", "-leading-kebab-case-", "--trailing--hyphen--"}},
		{"Mixed Delimiters", []string{"hello_world-and-kebab", "Mixed_Snake-Kebab--Case", "UPPER_snake-KEBAB_Case"}},
		{"Acronym Handling", []string{"ID_NUMBER", "HTTP_Response_Code", "XML_HTTP_REQUEST", "API_Version_2"}},
		{"Spaces and Words", []string{"   leading spaces", "trailing spaces   ", "   both sides   ", "This is a test sentence"}},
		{"Symbols and Random", []string{"123_ABC-xyz--789", "!!special$$characters**", "___!!!___weird___CASE___!!!___"}},
	}
)
