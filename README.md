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

1. **Push your latest change on your branche**:

   ```bash
   git push
   ```

2. **Switch to the branch you want merge**:

   ```bash
   git checkout <dev?>
   ```

3. **pull `dev`** (or the target branch):

   ```bash
   git pull dev
   ```

4. **Checkout again on your feature branch**

   ```
   git checkout <your_feature_branch>
   ```

5. **Merge dev into your feature branche** :

   ```
   git merge dev
   ```

6. **Resolve conflits and push and PR**:
