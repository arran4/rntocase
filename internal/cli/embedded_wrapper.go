package cli

import "github.com/arran4/rntocase/internal/skill"

// extractEmbeddedSkillFn is a package-level variable that can be overridden in tests
var extractEmbeddedSkillFn = skill.ExtractEmbeddedSkill
