package skill

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveSkillPathAllowsSafeLegacyDotDotPrefix(t *testing.T) {
	target := &Target{Path: t.TempDir()}

	resolved, err := ResolveSkillPath(target, "..legacy")
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(target.Path, "..legacy"), resolved)

	assert.Error(t, ValidateSkillName("..legacy"), "new Agent Skills installs should still reject the legacy name syntax")
}
