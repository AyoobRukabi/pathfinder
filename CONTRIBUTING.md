# Team Collaboration & Git Workflow Guide

Welcome to the project! Since we are working together, we’ll use this workflow to keep our `main` branch stable, our code clean, and our Git history untangled.

---

## 1. Initial Project Setup
If you haven't already, clone the repo and set up your remotes.

### Clone the Repository

```bash
git clone <GITEA_URL>
cd pathfinder # Make sure to cd into the actual folder name, not the URL!
```

### Configure Dual Remotes (Optional but Recommended)
We will use Gitea as our main source of truth, but you can sync to your private GitHub repository for your portfolio. 

> [!CAUTION]
> Your GitHub repository must be private so other students cannot copy our code.

To avoid forgetting to push to GitHub and causing branch divergence, we can configure Git to push to *both* servers simultaneously every time you type `git push`:

```bash
# 1. Add GitHub as a second remote named 'github'
git remote add github <GITHUB_URL>

# 2. Tell Git to push to both Gitea and GitHub whenever you push to 'origin'
git remote set-url --add --push origin <GITEA_URL>
git remote set-url --add --push origin <GITHUB_URL>
```
*(Now, a standard `git push origin main` will update both automatically.)*

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

Please follow a simplified feature branch naming, and use kebab case (all lowercase, spaces replaced by hyphens `-`).

**Format:** `user-name/type/short-description`
Examples: `korben/feat/dijkstra-algorithm`, `korben/fix/heap-overflow`, `korben/docs/api-ref`

```bash
git checkout -b <name-of-the-branch>
# e.g., git checkout -b korben/feat/dijkstra-algorithm
```

## 3. Development & Committing

Before committing, ensure your code quality:
1. Run `go fmt ./...` (Ensures standard Go style)
2. Run `go mod tidy` (Cleans up `go.mod` and `go.sum`)
3. **Always run `git status` first** to see exactly what changed so you don't accidentally commit secret files or compiled binaries.
4. Stage files: `git add <file-name>` (or `git add .` only if you are 100% sure all changed files should be committed).
5. Commit with a clear message: `git commit -m "feat: add adjacency list structure"` 

### Commit Message Convention
Check [Conventional Commits](https://www.conventionalcommits.org/).

Format: `<type>(<scope>): <subject>`

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `style:` Code style changes (formatting, missing semicolons, etc.)
- `refactor:` Refactoring production code (e.g., renaming a variable)
- `test:` Adding or updating tests
- `chore:` Routine tasks (updating dependencies)

**Examples:**
Good: `"feat(algo): implement priority queue for Dijkstra"`
Bad: `"updates"`, `"Add Some Fixes."`

### Step 4: Push changes to the branch
```bash
# Push changes to the remote branch, NOT TO MAIN!
git push -u origin <name-of-the-branch>
```
> [!TIP]
> You can configure Git to automatically create the upstream branch so you never have to type `-u origin ...` again:
> `git config --global push.autoSetupRemote true`

## 4. Merging Code (The Pull Request)

Once your feature is ready:

1. Go to the Gitea web interface and navigate to the Pull Requests tab.
2. Click **New Pull Request**. Select `main` as the target and your branch as the source.
3. **Write a good description:** Explicitly state *What* you changed and *Why* you changed it. Do not force the reviewer to guess your logic.
4. **Review:** Your partner will review the code, leave comments, and approve it.
5. **Merge:** Once approved, hit "Squash and Merge".
6. **Clean up:** Delete your local and remote branches to keep the repo clean:

```bash
git checkout main
git pull origin main
git branch -d <name-of-the-branch> # Deletes local branch
git push origin --delete <name-of-the-branch> # Deletes remote branch
```

## 5. Common Git Commands Reference

| Action | Command |
| ----------- | ----------- |
| Check status | `git status` |
| See history | `git log --oneline --graph` |
| Fix last commit message | `git commit --amend` |
| Undo last commit | `git reset --soft HEAD~1` **(ONLY use this if you haven't pushed yet!)** |
| Switch branches | `git checkout <branch-name>` |
| Discard local changes | `git checkout -- <file-name>` |

## 6. Other Tips

1. **VS Code Extensions:** Install the official Go extension and `GitLens`.
2. **Conflicts:** If we both edit the same line, Git will ask us to choose. Use the VS Code "Merge Editor" to pick the correct logic.