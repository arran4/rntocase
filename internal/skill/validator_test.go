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
