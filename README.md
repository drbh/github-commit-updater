# github-commit-updater


Check if you should update based on Github commits.

Just save the last commit on Github, before you make your update. 

If the parent to the most recent commit is not the parent you saved, the repo should update.


```bash
go run main.go drbh/github-commit-updater compare
```