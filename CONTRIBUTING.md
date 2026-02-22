# Team Collaboration & Git Workflow Guide

Welcome to the project! Since we are working together on this project, we’ll use this workflow to keep our `main` branch stable and our code clean.

---

## 1. Initial Project Setup
If you haven't already, clone the repo and set up your remotes.

### Clone the Repository

```bash
git clone <GITEA_URL>
cd <GITEA_URL>
```

### Configure Dual Remotes (Optional)
We will use Gitea for collaboration and you can use GitHub for your personal portfolio.

> [!CAUTION]
> GitHub repository should be private.
> Also, make sure that Gitea is default repo.

```bash
# Gitea is 'origin' by default. Add GitHub as a second remote if you want to:
git remote add github <GITHUB_URL>

# To push to both (do this after merging to main):
git push origin main
git push github main
```

```bash
# To check what remotes you have, run the command:
git remote -v
```

```bash
# You can rename the remote with the command:
git remote rename <OLDNAME> <NEWNAME>
# e.g. git remote rename origin gitea
```

## 2. The Feature Branch Workflow

> [!CAUTION]
> Rule #1: Never commit directly to main. Always work on a branch.

### Step 1: Sync your local environment
Before starting new work, make sure your local main matches the repo.

```bash
git checkout main
git pull origin main
```

### Step 2: Create a descriptive branch

Please follow a simplified feature branch naming, and use kebab case.
Kebab case is a naming convention where all letters are lowercase and spaces are replaced by hyphens (`-`), such as `user-profile-settings`.

**Format:** `user-name(or account name)/type(optional)/short-description`
Examples: `user-name/feat/dijkstra-algorithm-logic`, `account-name/fix/heap-overflow`, `user-name/api-ref`

```bash
# Pattern: user-name/feat/feature-name or user-name/fix/bug-name
git checkout -b <name of the branch>
#e.g. git checkout -b user-name/feat/dijkstra-algorithm-logic
```
## 3. Development & Committing

Before committing, ensure your code quality:
1. Run `go fmt ./...` before committing. (Ensures standard style)
2. Run `go mod tidy` (Cleans up `go.mod` and `go.sum`)
3. Add `git add .` (adds all changed files) or `git add </path/to_file>` (adds single changed file, you can add multiple files)
4. Commit `git commit -m "feat: add adjacency list structure for map"` (Please write clear, concise messages). 

### Commit Message Convention
Check [Conventional Commits](https://www.conventionalcommits.org/).
Check [Semantic Commit Messages](https://gist.github.com/joshbuchea/6f47e86d2510bce28f8e7f42ae84c716)

Format: `<type>(<scope>): <subject>`
`<scope>` is optional

- `feat:` New feature (not a new feature for build script)
- `fix:` Bug fix (not a fix to a build script)
- `docs:` Documentation changes.
- `style:` Code style changes (formatting, missing semi colons, etc;)
- `refactor:` Refactoring production code (eg. renaming a variable)
- `test:` Adding missing tests, refactoring tests; no production code change
- `chore:` Routine tasks (updating dependencies, build tools, no production code change)
- `build:` Changes affecting the build system or external dependencies.
- `ci:` Changes to CI configuration files or scripts.
- `perf:` A code change that improves performance
- `revert:` Reverts to a previous commit. It should begin with `revert:`, followed by the header of the reverted commit. In the body, it should say: `This reverts commit <hash>`, where the hash is the SHA of the commit being reverted.

**Examples:**
Good: `"feat(algo): implement priority queue for Dijkstra"`
Good: `"feat: add adjacency list structure for map"`

Bad: `"updates"`
Bad: `"Add Some Fixes."`

### Step 4: Push changes to the branch
```bash
# Push changes to the branch, NOT TO THE MAIN!
# The -u flag is the standard shorthand for --set-upstream.
git push -u origin <name of the branch>
# e.g. git push -u origin account-name/fix/heap-overflow
```

> [!TIP]
> You can configure your global Git settings to automatically create the upstream branch on the remote server whenever you push a branch for the first time.
> Use command `git config --global push.autoSetupRemote true`

## 4. Merging Code (The Pull Request)

Once your feature is pushed to the branch:

1. Go to the Gitea web interface and navigate to project Pull Requests tab.
2. Click a `New Pull Request` button, and select the branch `merge into` and your branch `pull from`.
3. Review: The other partner should look at the code, leave comments, and then Approve.
4. Merge: Once approved, hit the Merge button.
5. Clean up: Delete your local branch:

```bash
# Do this only when the PR is approved, and changes are in main!
git checkout main
git pull origin main
git branch -d account-name/fix/heap-overflow # delete local branch if you don't need it anymore and fix or feature is done and merged
```

## 5. Common Git Commands Reference

| Action | Command |
| ----------- | ----------- |
| Check status | `git status` |
| See history | `git log --oneline --graph` |
| Undo last commit (keep code) | `git reset --soft HEAD~1` |
| Switch branches | `git checkout <branch-name>` |
| Discard local changes | `git checkout -- <file-name>` |

## 6. Other Tips

1. VS Code Extensions: Install the official Go extension and `GitLens`.
2. Mod Files: If you add a new dependency, run `go mod tidy` before committing.
3. Conflicts: If we both edit the same line, Git will ask us to choose. Use the VS Code "Merge Editor" to pick the best logic.