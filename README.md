# devoops
CLI tool to facilitate EKS local development written in GOLANG

# Install from git repo

[Download and install GO](https://go.dev/doc/install)

### Add GO executables to your PATH

```bash
# Discover the install path with
go list -f '{{.Target}}'

EXPORT PATH=$PATH:/path/to/your/install/dir
```

### Build tool
Inside project root directory
````bash
go build -o /path/to/your/install/dir/ku
````

### Use the tool
```bash
# login with your AWS credentials
ku login -p PROFILE

# get the current context
ku current-context
ku cc

# switch to another context defined in your kube config
ku switch-context
ku sc
```
