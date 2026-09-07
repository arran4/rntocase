package skill

import "os"

// osRename is a package-level variable that can be overridden in tests
var osRename = os.Rename
