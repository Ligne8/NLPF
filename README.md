## Best practices
Everything should be written in english (even comments).
### Naming convention
We use camel case in the code and snake case in the databse.  
For example:
```
const (
	StateAvailable State = "available"
	StateInTransit State = "in_transit"
	StateArchive   State = "archive"
	StateOnMarket  State = "on_market"
	StateAtTrader  State = "at_trader"
)
```
### Starting a new ticket
You need to create a branch from `dev`:  
```
git switch -c <your_branch_name>
```
Your branch name needs to look like the following.  
Fix: `[ticket_number] fix`  
Feature: `[ticket_number] feat`
### Commiting your ticket
Add your modifications:  
```
git add <path_to_file>
```
Commit your files:  
```
git commit -m "[<ticket_number>] <fix/feat>: small description of the ticket"
```
### Pushing your code
In order to push your modifications there are `5 steps` to follow:  
`Step 1` Clone the repository or update your local repository with the latest changes  
```
git pull origin dev
```
`Step 2` Switch to the head branch of the pull request.  
```
git checkout <your_branch_name>
```
`Step 3` Merge branch into the head branch.  
```
git merge dev
```
`Step 4` Fix the conflicts and commit the result.  
`Step 5` Push the changes.
```
git push -u origin <your_branch_name>
```
