# Ploy

Ploy is a simple terminal-based deployment tool. It lets you define services in a configuration
file together with the desired versions of these services. It can verify if the desired versions
are deployed and, if not, it can deploy the right versions.

It can currently verify and deploy services of the following types:

- [AWS Lambda](https://aws.amazon.com/lambda/) - with the code packaged either as a Docker image and published to 
  ECR or as a zip file on S3
- [AWS ECS](https://aws.amazon.com/ecs/)

## Installation

Install the tool with `go`

```bash
go install github.com/DandyDev/ploy
```

## Usage

First, you have to create one or more Ploy configuration files. A Ploy configuration file
contains a list of services and their desired versions. It's written in YAML.

Example:

```yaml
deployments:
  - id: my-lambda # in case of Lambda, this is the function name
    type: lambda
    version: v2
    version-environment-key: VERSION # optional, updates the given environment variable with the version when deploying
    pre-deploy-command: ["./my-script.sh", "arg1", "arg2"] # optional, runs the specified command before deployment. The to be deployed version is available as the $VERSION environment variable 
    post-deploy-command: ["./my-script.sh", "arg1", "arg2"] # optional, runs the specified command after successful deployment. The deployed version is available as the $VERSION environment variable 
  - id: my-zipped-lambda # in case of Lambda, this is the function name
    type: lambda
    version: v3
    version-environment-key: VERSION # optional, updates the given environment variable with the version when deploying
    bucket: my-bucket # S3 bucket that contains the zipped deployable
    key: "{{.Id}}/{{.Version}}.zip" # S3 key of zipped deployable. This is a template that is resolved at deployment time. Supported variables are the Id of the deployment and the Version
  - id: my-container # in case of ECS, this is the service name
    type: ecs
    cluster: my-cluster
    version: v666
    version-environment-key: VERSION # optional, updates the given environment variable in the container with the version when deploying
    wait-for-service-stability: true # optional, defaults to false
    wait-for-minutes: 5 # optional, how long to wait for service stability, defaults to 30
    force-new-deployment: true # optional, defaults to false
```

Typically, you'll have one configuration file for each environment (e.g. dev, prod, staging).

After defining the configuration files, you can use the `ploy` command to verify or deploy the
services.

### Verify

```bash
ploy verify development.yml
```

This will verify if all the services specified in the deployment file are running at the versions 
specified in the deployment file

### Deploy

```bash
ploy deploy development.yml
```

This will check for each service in the deployment file which version is currently deployed. If 
the deployed version differs from the version specified in the deployment file, it will deploy 
the desired version of that service. If the desired version is already running, it will do nothing.

### Update version

```bash
ploy update development.yml my-service v123
```

This will update the deployment file so that the specified version is set at the specified 
version. **This will not do an actual deployment**. You can run `ploy deploy` afterwards. 

This command was created to make it easy to update deployment files through CI/CD.

## Engines

There are currently two supported deployment engines:

- [AWS Lambda](https://aws.amazon.com/lambda/) (type: `lambda`) - with the code packaged as a Docker 
  image. Version is the image tag.
- [AWS ECS](https://aws.amazon.com/ecs/) (type: `ecs`) - with the code packaged as a Docker image. 
  Version is the image tag.

## Contributing

Fork the repo, make your changes, and submit a pull request.

## TODO

- Better error handling
- Add support for deploying new ECS task definitions for one-off tasks that are not part of a
  service
- Add support for other deployment engines. See `github.com/DandyDev/ploy/engine` for examples of
  how engines are implemented
- Create command that serves a simple dashboard the visualizes the services that are deployed
  and their versions. Should do periodic checks in the background to verify versions
- Create a Homebrew tap + formula for Ploy

## Authors

- [@DandyDev](https://www.github.com/DandyDev)
