# rntocase Agent Skill

Welcome to the `rntocase` skill for AI coding agents! `rntocase` is a CLI application that helps humans and automated agents batch-rename files by changing their casing (e.g., camelCase, snake_case, PascalCase, kebab-case, etc.).

## Core Concepts

*   **Subcommands:** `rntocase` is a single binary (`rntocase`) with a subcommand architecture. Commands like `rntocase camel`, `rntocase snake`, and `rntocase acronym` perform the actual renaming operations.
*   **Batch Operations:** The tool is designed to take multiple files as arguments: `rntocase snake file1.txt file2.txt file3.txt` or `rntocase snake *`.
*   **Extensions Preserved:** The renaming logic applies to the filename *before* the extension. Extensions are preserved automatically. This includes standard extensions, dotfiles (e.g., `.env`), and common compound extensions (e.g., `.tar.gz`, `.tar.bz2`).
*   **Safety First:** By default, `rntocase` renames files immediately but uses a preflight planning phase to check for destination collisions (e.g., multiple source files mapping to the same name, or a destination that already exists). If any collision is detected, the entire batch will abort safely before making changes.
*   **Recursive Directory Renaming:** Directory recursion is explicit and opt-in via `-R` or `--recursive`. When enabled, `rntocase` traverses supplied directory roots to discover target files. By default, recursive mode selects files only; directories are not renamed. Discovery is gathered across all roots before the safe rename planner executes collision preflight checks. Existing explicit non-recursive positional-file behavior remains unchanged when `-R` / `--recursive` is omitted.
*   **Conservative Symlink Policy:** In recursive mode, recursively discovered symlinks are skipped and directory symlinks are not followed to prevent cycles and accidental mutations. Explicitly supplied symlinks are also skipped in recursive mode.
*   **Glob Filtering (`--include` / `--exclude`):** Filter discovered files using repeatable `--include` and `--exclude` flags. Filters use shell-style doublestar globs rather than regexes, evaluated against slash-normalized paths relative to each supplied directory root. Patterns support `**` (e.g., `--exclude '.git/**'` skips `.git` directories and all nested contents, and `--include '**/*.jpg'` matches `.jpg` files at any depth). Multiple `--include` flags behave as alternatives (union), and exclusions take precedence over inclusions.

## Operational Guidance for Agents

As an automated agent, you should follow these rules when using `rntocase`:

1.  **Always use `--dry-run` first:** Before executing any destructive rename operation, run the command with `--dry-run`. This outputs the intended changes without modifying the filesystem. Examine the output to verify the result is what you expect.
    *   *Example:* `rntocase kebab --dry-run ./*`
2.  **Inspect Recursive Operations with `--dry-run`:** When using `-R` or `--recursive`, always run with `--dry-run` first (or pair with `--json` for structured verification) before executing mutations.
    *   *Example (recursive dry run):*
        ```sh
        rntocase snake --dry-run --recursive ./assets
        ```
    *   *Example (recursive filtered dry run):*
        ```sh
        rntocase kebab --dry-run -R \
          --include '*.jpg' \
          --exclude '.git/**' \
          ./photos
        ```
3.  **Avoid Interactive Mode:** Do not use the `-interactive` flag. Interactive prompts will block your execution in headless or scripted environments, requiring human intervention. If you must be sure, use `--dry-run` to check, then run the command normally.
4.  **No In-Place Modification Tracking:** The tool modifies the filesystem directly. There is no internal "undo" command.
5.  **Handling Spaces:** Ensure paths and filenames containing spaces or special characters are properly quoted when invoking the tool via shell (e.g., `rntocase camel "My File.txt"`).

## Common Traps and Misuses

*   **Forgetting arguments:** A subcommand requires file arguments to operate on. Running `rntocase snake` without files will result in an error: `Error: No files provided.`
*   **Expecting directories to be renamed in recursive mode:** Recursive mode selects regular files only; directories are never renamed.
*   **Using regex instead of glob patterns:** `--include` and `--exclude` accept doublestar glob patterns (such as `*.jpg` or `.git/**`), not regular expressions.
*   **Assuming symlinks are followed recursively:** Directory symlinks are not followed and symlink files are skipped in recursive mode.
*   **Using `-interactive`:** As mentioned, this blocks execution.
*   **Misunderstanding delimitations:** When using `rntocase delimited`, you may need to provide an additional flag like `-delimiter "_"`. Use `rntocase <subcommand> -h` to see available flags for a specific operation.

## Available Subcommands (Examples)

*   `acronym`: Rename files by acronym
*   `camel`: Rename files to camel case
*   `constant`: Rename files to constant case
*   `darwin`: Rename files to darwin case
*   `delimited`: Rename files with a custom delimiter
*   `dot`: Rename files to dot case
*   `kebab`: Rename files to kebab case
*   `lower`: Rename files to lower case
*   `lowerleading`: Rename files with a lower leading character
*   `pascal`: Rename files to pascal case
*   `reverse`: Reverse characters or words in file names
*   `snake`: Rename files to snake case
*   `title`: Rename files to title case
*   `upper`: Rename files to upper case
*   `upperleading`: Rename files with an upper leading character
*   `trim`: Trim whitespace or specific characters from file names

Use `rntocase help` or `rntocase <subcommand> -h` for more details on flags and usage.
