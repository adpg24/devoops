# devoops
CLI tool with multiple useful commands to work with Kubernetes(EKS) and AWS.

## Install

### Linux
```bash
# Download the .tar.gz file
curl -LO https://github.com/adpg24/devoops/releases/download/0.0.1/devoops_0.0.1_linux_amd64.tar.gz
# Extract the executable and copy it to a destination of your choice
tar xvf -C /opt/bin devoops_0.0.1_linux_amd64.tar.gz
```

#### Dependency for Linux
This tool depends on this [clipboard package](https://pkg.go.dev/golang.design/x/clipboard)

- macOS: require Cgo, no dependency
- Linux: require X11 dev package. For instance, install libx11-dev or xorg-dev or libX11-devel to access X window system.
- Windows: no Cgo, no dependency
- iOS/Android: collaborate with gomobile

## Usage

#### Kubernetes tools

```bash
# Show the current context as defined in your kube config
devoops current-context
devoops cc

# Switch to another context defined in your kube config
devoops switch-context
devoops sc
```

#### AWS tools

##### login
Authenticate to AWS with MFA and generate short-term credentials.

Create an AWS profile(~/.aws/credentials) with suffix `-mfa`. E.g. my-profile-mfa
```ini
[my-profile-mfa]
aws_access_key_id=XXXXX
aws_secret_access_key=XXXX
aws_mfa_device=arn:aws:iam::123456789012:mfa/user
```
Login with the profile you created, leave out the `-mfa` suffix.\
This command will request short-term credentials and create a new profile `[my-profile]`.
```bash
devoops login -p my-profile
```

##### tag

Add a new tag for an existing image in an ECR repository.\
You must have selected the profile you want to use prior to using this command, e.g. `export AWS_PROFILE=my-profile`

```bash
devoops tag -r my-repository tag newTag
```

##### select-profile

Select a profile from the profiles defined in `~/.aws/credentials`. The export command (linux) will be copied to your clipboard.
```bash
devoops select-profile
devoops sp
```

## Development

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
go build -o /path/to/your/install/dir/devoops
````

### Cobra CLI
This project uses [Cobra CLI](https://github.com/spf13/cobra) for the interface.
