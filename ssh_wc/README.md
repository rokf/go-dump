## wc via SSH

Assuming you are inside the `ssh_wc` folder:

- Run with `go run cmd/wcserve/wcserve.go`.
- Use with something like `cat cmd/wcserve/wcserve.go | ssh user@localhost -p 2222 --`. Pass additional `wc` params (like `-l`) after the `--` part.
- To remove the generated host key run `ssh-keygen -R "[localhost]:2222"`.
