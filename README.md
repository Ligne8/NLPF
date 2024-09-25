## Best practices
Everything should be written in english (even comments).
### Naming convention
We use `camelCase` in the code and `snake_case` in the databse.  

### Starting a new ticket
You need to create a branch from `dev`:  
```
git switch -c <your_branch_name>
git checkout -B <you_branch_name>
```

On each Linear issue, you can associate a Git branch corresponding to `feature/identifier-title`. To make this easier, you can copy the branch name via `cmd + shift + .` or use the branch icon directly on the issue.

### Commit convention
```
git commit -m "[<ticket_number>] <fix/feat>: small description of what you have done"
```

### Pushing your code

1. **Update your local repository with the latest changes**:
   ```bash
   git pull origin dev
   ```

2. **Switch to your feature branch**:
   ```bash
   git checkout <your_branch_name>
   ```

3. **Rebase onto `dev`** (or the target branch):
   ```bash
   git rebase dev
   ```

4. **Resolve conflicts** if any, then continue the rebase:
   ```bash
   git add <conflicted_file>
   git rebase --continue
   ```

5. **Push your changes**:
    - If this is your first push after the rebase, use the `--force-with-lease` option to avoid data loss:
   ```bash
   git push --force-with-lease origin <your_branch_name>
   ``` 
