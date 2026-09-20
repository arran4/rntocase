package skill

import (
	"strings"
	"testing"
)

func TestParseAndValidateManifest(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "valid manifest",
			content: "---\nname: my-skill\ndescription: my description\n---\n# My Skill\n",
			wantErr: false,
		},
		{
			name:    "missing name",
			content: "---\ndescription: my description\n---\n",
			wantErr: true,
		},
		{
			name:    "missing description",
			content: "---\nname: my-skill\n---\n",
			wantErr: true,
		},
		{
			name:    "invalid name uppercase",
			content: "---\nname: My-Skill\ndescription: my description\n---\n",
			wantErr: true,
		},
		{
			name:    "invalid name consecutive hyphens",
			content: "---\nname: my--skill\ndescription: my description\n---\n",
			wantErr: true,
		},
		{
			name:    "no frontmatter",
			content: "# My Skill\n",
			wantErr: true,
		},
		{
			name:    "frontmatter not at start",
			content: "# Title\n---\nname: my-skill\ndescription: desc\n---\n",
			wantErr: true,
		},
		{
			name:    "whitespace description",
			content: "---\nname: my-skill\ndescription: \"   \"\n---\n",
			wantErr: true,
		},
		{
			name:    "too long description",
			content: "---\nname: my-skill\ndescription: " + strings.Repeat("a", 1025) + "\n---\n",
			wantErr: true,
		},
		{
			name:    "max length description",
			content: "---\nname: my-skill\ndescription: " + strings.Repeat("a", 1024) + "\n---\n",
			wantErr: false,
		},
		{
			name:    "valid delimiter within YAML scalar",
			content: "---\nname: my-skill\ndescription: embedded --- test\n---\n",
			wantErr: false,
		},
		{
			name:    "valid delimiter within Markdown body",
			content: "---\nname: my-skill\ndescription: valid\n---\n# Some body\n---\nembedded block",
			wantErr: false,
		},
		{
			name:    "valid closing delimiter at EOF",
			content: "---\nname: my-skill\ndescription: valid\n---",
			wantErr: false, // EOF termination is allowed
		},
		{
			name:    "CRLF endings",
			content: "---\r\nname: my-skill\r\ndescription: valid\r\n---\r\n# body",
			wantErr: false,
		},
		{
			name:    "malformed start delimiter",
			content: " --- \nname: my-skill\ndescription: valid\n---\n",
			wantErr: true,
		},
		{
			name:    "malformed end delimiter not alone",
			content: "---\nname: my-skill\ndescription: valid\n---not-a-delimiter\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseAndValidateManifest([]byte(tt.content))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAndValidateManifest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
