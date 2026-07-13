# fatgo --- Fatman's Go library

## Expected Uses

1. Add as a submodule in project repos.
   
   eg. `import myproject/internal/fatgo/callwheel`
   
1. Include as a Github project.
   
   eg. `import github.com/FatmanUK/fatgo/callwheel`

## Confusion over go.mod

As far as I can tell, there's no good, detailed documentation on how `go.mod` works. I've had the worst time working it out, which is ridiculous given how awesome Go is otherwise.

Eventually I resorted to asking Gemini, notably another Google product.

> Q:
> Golang. I have a super-project named `super` with sub-projects named `sub1` and `sub2`. The `super` project has a `go.work` file so that `go test` works. The super-project is hosted on Github. Please write all the `go.mod` files so that I can use `go test` and import the sub-projects in other unrelated projects.

> A:
> To allow your sub-projects (`sub1` and `sub2`) to be imported by unrelated external projects, they must use their full GitHub repository paths in their module declarations.
> Here is how to set up your directory structure and `go.mod` files.
> 
> ### Directory Structure
> 
>     super/
>     ├── go.work
>     ├── sub1/
>     │   └── go.mod
>     └── sub2/
>         └── go.mod
> 
> ### 1. The `sub1` Sub-Project `go.mod`
> 
> Replace `your-username` with your actual GitHub username or organisation name.
> 
>     module ://github.com
>     go 1.21
> 
> ### 2. The `sub2` Sub-Project `go.mod`
> 
> If `sub2` needs to import `sub1`, use the full GitHub path.
> 
>     module ://github.com
>     go 1.21
>     
>     require ://github.com v0.0.0
>     // Note: No 'replace' directive is needed here because go.work handles local development.
> 
> ### 3. The Super-Project `go.work`
> 
> The `go.work` file sits in the root `super` directory. It tells Go to use your local folders instead of downloading them from GitHub during local testing.
> 
>     go 1.21
>     
>     use (
>     	./sub2
>     	./sub1
>     )
> 
> ### How to Import Them Locally
> 
> Run your tests from the root `super` directory:
> 
>     go test ./...
> 
> ### How Unrelated Projects Import Them
> 
> External projects can bypass the super-project completely and download a specific sub-project directly using:
> 
>     go get ://github.com
> 
> To make this work seamlessly for external projects, ensure you push your tags to GitHub (e.g., `sub1/v1.0.0` or `sub2/v1.0.0`) so the Go proxy can version them independently.
> If you plan to add tags, let me know if you want to see how to format multi-module git tags or if you need help setting up the import blocks inside your Go source files.

> Q:
> Doesn't the super-project need a `go.mod` for `go test` to work?

> A:
> No. The super-project does not need a `go.mod` file for `go test ./...` to work when a `go.work` file is present.
> 
> ### Why It Works Without It
> 
> * `go.work` defines the workspace: The `go.work` file turns the root directory into a workspace container.
> * It flags sub-modules: The use directive tells Go exactly which sub-folders contain valid modules.
> * ./... targets everything: When you run go test ./... from the root, Go scans all sub-directories defined in `go.work`, locates their respective `go.mod` files, and tests them. [1, 2, 3]
> 
> ### The Danger of Adding Root `go.mod`
> 
> Adding a `go.mod` file to the root super directory can cause issues:
> 
> * Nested module conflicts: Go will treat the root as a module that tightly clamps down on the sub-modules, which breaks clean external imports.
> * Import confusion: External projects might accidentally try to `import ://github.com` instead of the specific sub-projects.
> 
> Keep the root clean with only the `go.work` file.
> 
> [1] [https://go.dev](https://go.dev/ref/mod)
> [2] [https://www.andy-pearce.com](https://www.andy-pearce.com/blog/posts/2023/Jul/all-go-building-and-testing/)
> [3] [https://medium.com](https://medium.com/@AlexanderObregon/go-modules-without-the-confusion-6a663228879b)

> Q:
> How to format multi-module git tags ?

> A:
> To version multiple modules in a single GitHub repository, you must prefix the Git tag with the relative path to the sub-module, followed by a forward slash and the semantic version.
> 
> ### Tag Format Structure
> 
>     [sub-module-directory]/[version]
> 
> ### Formatting Your Specific Modules
> 
> For your `super` repository, use these exact formats when creating tags:
> 
> * For `sub1`: sub1/v1.0.0
> * For `sub2`: sub2/v0.1.0
> 
> ### How to Create and Push the Tags
> 
> Run these Git commands in your terminal to tag a specific release and send it to GitHub:
> 
>     # 1. Tag the sub1 module
>     ❯ git tag sub1/v1.0.0
>     # 2. Tag the sub2 module
>     ❯ git tag sub2/v0.1.0
>     # 3. Push the tags to GitHub
>     ❯ git push origin --tags
> 
> ### Critical Rules for Multi-Module Tags
> 
> * Lowercase only: The tag prefix must match the directory name exactly. If your folder is `sub1`, the tag must start with `sub1/`.
> * Include the 'v': The version section must start with a lowercase v (e.g., `v1.2.3`).
> * No root tags: Do not create a global tag like v1.0.0 at the root. The Go proxy will get confused and will not know which sub-project the version belongs to.

What's with all the `://github.com` nonsense? Gemini, otherwise pretty reliable in my opinion, always fails utterly on configuration lines with URLs in, regardless of how simple or complex the URL. Whatever the URL should be, the model always outputs something like "://domain.tld" for the URL itself.

The model is completely blind to it and keeps outputting the same garbage. The agent can see it, and gets quite distressed over its own failure to correct the problem. I've never found an answer. Luckily, in this case I can figure it out myself.
