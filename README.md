# Ploy

Ploy is a simple terminal-based deployment tool. It lets you define services in a configuration
file together with the desired versions of these services. It can verify if the desired versions
are deployed and, if not, it can deploy the right versions.

It can currently verify and deploy services of the following types:

- [AWS Lambda](https://aws.amazon.com/lambda/) - with the code packaged as a Docker image and
  published to ECR
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
  - id: my-container # in case of ECS, this is the service name
    type: ecs
    version: v666
```

Typically, you'll have one configuration file for each environment (e.g. dev, prod, staging).

After defining the configuration files, you can use the `ploy` command to verify or deploy the 
services.

```bash
ploy verify development.yml
ploy deploy development.yml
```

Ploy will only deploy new versions of services if the desired version is different from the 
currently deployed version.

## Contributing

Fork the repo, make your changes, and submit a pull request.

## TODO

- Update CLI help texts (currently, the Cobra defaults are used)
- Implement support for deploying ECS services (currently verification of ECS services works)
- Better error handling
- Add support for deploying new ECS task definitions for one-off tasks that are not part of a 
  service
- Add support for other deployment engines. See `github.com/DandyDev/ploy/engine` for examples of 
  how engines are implemented
- Create command that serves a simple dashboard the visualizes the services that are deployed 
  and their versions. Should do periodic checks in the background to verify versions

## Authors

- [@DandyDev](https://www.github.com/DandyDev)
