# Go Modules Branch Fetcher
*Go Modules Branch Fetcher* is a command-line tool that allows you to easily track all the modules used in your Go project, display their versions, identify the branches they are associated with, and handle errors gracefully. The tool fetches modules, retrieves branch information based on commit hashes, and highlights special branches (tags, pull requests, etc.).

Features
- List Go Modules: Use the go list -m -json all command to list all Go modules in your project in JSON format.
- Branch Information: Retrieves and displays the branch name associated with each module.
- Special Branch Detection: Detects special branches such as refs/tags/, refs/pull/, and other non-standard branches and displays them in different colors.
- Error Handling: If the branch information for a module cannot be retrieved, error messages are displayed.
- Module Prefix Filtering: Filters and lists modules that match a given prefix using the --prefix flag.
- Table Visualization: Displays module names, versions, and branch details in a clean and organized table format.
