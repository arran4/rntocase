package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestRunSkill_RequiresSubcommand(t *testing.T) {
	err := runSkill([]string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "requires a subcommand")
}

func TestRunSkill_UnknownSubcommand(t *testing.T) {
	err := runSkill([]string{"unknown_cmd"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown skill subcommand: unknown_cmd")
}

func TestRunSkillInstall_RequiresSource(t *testing.T) {
	err := runSkillInstall([]string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "usage: skill install <source>")
}

func TestRunSkillUpdate_RequiresName(t *testing.T) {
	err := runSkillUpdate([]string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "usage: skill update <name> or skill update --all")
}

func TestRunSkillRemove_RequiresName(t *testing.T) {
	err := runSkillRemove([]string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "usage: skill remove <name>")
}

func TestRunSkillInspect_RequiresName(t *testing.T) {
	err := runSkillInspect([]string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "usage: skill inspect <name>")
}
