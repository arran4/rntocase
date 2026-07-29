package cli

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestRunSkill_RequiresSubcommand(t *testing.T) {
	err := RunSkill([]string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "requires a subcommand")
}

func TestRunSkill_UnknownSubcommand(t *testing.T) {
	err := RunSkill([]string{"unknown_cmd"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown skill subcommand: unknown_cmd")
}

func TestRunSkillInstall_RequiresSource(t *testing.T) {
	err := RunSkillInstall([]string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "usage: skill install <source>")
}

func TestRunSkillUpdate_RequiresName(t *testing.T) {
	err := RunSkillUpdate([]string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "usage: skill update <name> or skill update --all")
}

func TestRunSkillRemove_RequiresName(t *testing.T) {
	err := RunSkillRemove([]string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "usage: skill remove <name>")
}

func TestRunSkillInspect_RequiresName(t *testing.T) {
	err := RunSkillInspect([]string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "usage: skill inspect <name>")
}
