# rntocase Agent Skill

Welcome to the `rntocase` skill for AI coding agents! `rntocase` is a CLI application that helps humans and automated agents batch-rename files by changing their casing (e.g., camelCase, snake_case, PascalCase, kebab-case, etc.).

## Core Concepts

*   **Subcommands:** `rntocase` is a single binary (`rntocase`) with a subcommand architecture. Commands like `rntocase camel`, `rntocase snake`, and `rntocase acronym` perform the actual renaming operations.
*   **Batch Operations:** The tool is designed to take multiple files as arguments: `rntocase snake file1.txt file2.txt file3.txt` or `rntocase snake *`.
*   **Extensions Preserved:** The renaming logic applies to the filename *before* the extension. Extensions are preserved automatically.
*   **Safety First:** By default, `rntocase` renames files immediately.

## Operational Guidance for Agents

As an automated agent, you should follow these rules when using `rntocase`:

1.  **Always use `--dry-run` first:** Before executing any destructive rename operation, run the command with `--dry-run`. This outputs the intended changes without modifying the filesystem. Examine the output to verify the result is what you expect.
    *   *Example:* `rntocase kebab --dry-run ./*`
2.  **Avoid Interactive Mode:** Do not use the `-interactive` flag. Interactive prompts will block your execution in headless or scripted environments, requiring human intervention. If you must be sure, use `--dry-run` to check, then run the command normally.
3.  **No In-Place Modification Tracking:** The tool modifies the filesystem directly. There is no internal "undo" command.
4.  **Handling Spaces:** Ensure paths and filenames containing spaces or special characters are properly quoted when invoking the tool via shell (e.g., `rntocase camel "My File.txt"`).

## Common Traps and Misuses

*   **Forgetting arguments:** A subcommand requires file arguments to operate on. Running `rntocase snake` without files will result in an error: `Error: No files provided.`
*   **Using `-interactive`:** As mentioned, this blocks execution.
*   **Misunderstanding delimitations:** When using `rntocase delimited` or `rntocase dot`, you may need to provide additional flags like `-delimiter "."`. Use `rntocase <subcommand> -h` to see available flags for a specific operation.

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
