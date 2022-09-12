[![Go Report Card](https://goreportcard.com/badge/github.com/springload/ssm-parent)](https://goreportcard.com/report/github.com/springload/ssm-parent)

## SSM Parent

This is wrapper entrypoint for Docker to do one thing: fetch parameters from SSM Parameter store and expose them as environment variables to the underlying process.

Please note, that it still requires a proper `init` process, for example the one embedded into Docker can be used with `docker run --init`.

```
SSM-Parent is a docker entrypoint.

It gets specified parameters (possibly secret) from AWS SSM Parameter Store,
then exports them to the underlying process. Or creates a .env file to be consumed by an application.

It reads parameters in the following order: path->name->plain-path->plain-name.
So that every rightmost parameter overrides the previous one.

Usage:
  ssm-parent [command]

Available Commands:
  dotenv      Writes dotenv file
  help        Help about any command
  print       Prints the specified parameters.
  run         Runs the specified command

Flags:
  -c, --config string        Path to the config file (optional). Allows to set transformations
  -d, --debug                Turn on debug logging
  -e, --expand               Expand arguments and values using shell-style syntax
  -h, --help                 help for ssm-parent
  -n, --name strings         Name of the SSM parameter to retrieve. Expects JSON in the value. Can be specified multiple times.
  -p, --path strings         Path to a SSM parameter. Expects JSON in the value. Can be specified multiple times.
      --plain-name strings   Name of the SSM parameter to retrieve. Expects actual parameter in the value. Can be specified multiple times.
      --plain-path strings   Path to a SSM parameter. Expects actual parameter in the value. Can be specified multiple times.
  -r, --recursive            Walk through the provided SSM paths recursively.
  -s, --strict               Strict mode. Fail if found less parameters than number of names.
      --version              version for ssm-parent

Use "ssm-parent [command] --help" for more information about a command.
```

The SSM parameter names or paths can be specified with `-p` or `-n` flags. In this case all parameters must be in JSON format, i.e.:

```
    {
        "ENVIRONMENT": "production"
    }
```

If several parameters are specified, all JSON entities will be read and merged into one, overriding existing keys, i.e.

Parameter one:

```
    {
        "USERNAME": "myuser",
        "DATABASE": "production"
    }
```

Parameter two:

```
    {
        "DATABASE": "test"
    }
```

The result will be merged as this:

```
    {
        "USERNAME": "myuser",
        "DATABASE": "test"
    }
```

One can also specify `--plain-name` and `--plain-path` command line options to read _plain_ parameters that are not in JSON format.
`ssm-parent` takes the value as is, and constructs a key name from the `basename parameter`,
i.e. a SSM Parameter `/project/environment/myParameter` with value `supervalue` will be exported as `myParameter=supervalue`.

### How to use

Determine the paths you want to read and try it out with `ssm-parent print` to see the resulting JSON output.
Then use `ssm-parent run` or `ssm-parent dotenv`.

### Variables transformations

To transform variables, a config file is needed due to the complex nature of it. `ssm-parent` supports all config formats supported by https://github.com/spf13/viper, i.e. `.toml`, `.yaml`, `.json`.

All configuration entities can be specified in there rather than as flags.
The supported transformations are:

1. rename - renames env vars
2. delete - deletes env vars
3. template - templates env vars
4. trim_name_prefix - removes a prefix from variable names

Rename, template, trim_name_prefix transformations expect a dictionary rule. The delete transformation expects an array.
Template transformation uses [Go templates](https://golang.org/pkg/text/template/), and the environment variables map is passed as `.`.

There are the following extra functions available in templates: url_host, url_user, url_password, url_path, url_scheme and trim_prefix. The current list of the custom functions can be found here https://github.com/springload/ssm-parent/blob/master/ssm/transformations/template_funcs.go#L9

trim_name_prefix will match any variables starting with `starts_with` and will remove the `trim` string from the start of the corresponding variable names.

There is practically no limit on the number of transformations and they are applied in order from top to the bottom.

Below there is an example that recursively gets parameters from `/$PROJECT/common/` and `/$PROJECT/$ENVIRONMENT` and constructs variables out of
`DATABASE_URL` to be consumed by an PHP application. It also renames `AWS_BUCKET` to `AWS_S3_BUCKET`, removes `DATABASE_URL` and trims a leading underscore from any variable name that may start with `_PHP`.

```yaml
recursive: true
expand: true
debug: true
paths: ["/$PROJECT/common/", "/$PROJECT/$ENVIRONMENT"]

transformations:
    - action: template
      rule:
          SS_DATABASE_SERVER: "{{ url_host .DATABASE_URL }}"
          SS_DATABASE_USERNAME: "{{ url_user .DATABASE_URL }}"
          SS_DATABASE_PASSWORD: "{{ url_password .DATABASE_URL }}"
          SS_DATABASE_NAME: '{{ with $x := url_path .DATABASE_URL }}{{ trim_prefix $x "/" }}{{end}}'
    - action: rename
      rule:
          AWS_BUCKET: AWS_S3_BUCKET
    - action: delete
      rule:
          - DATABASE_URL
    - action: trim_name_prefix
      rule:
          trim: "_"
          starts_with: "_PHP"
```

### Example Dockerfile part

```
ENV PROJECT myproject
ENV ENVIRONMENT production

RUN wget -O /tmp/ssm-parent.tar.gz https://github.com/springload/ssm-parent/releases/download/v1.4.1/ssm-parent_1.4.1_linux_amd64.tar.gz && \
    tar xvf /tmp/ssm-parent.tar.gz && mv ssm-parent /usr/bin/ssm-parent && rm /tmp/ssm-parent.tar.gz

ENTRYPOINT ["/usr/bin/ssm-parent", "run", "-e", "-p", "/$PROJECT/$ENVIRONMENT/backend/", "-r",  "--"]
CMD ["caddy" , "--conf", "/etc/Caddyfile", "--log", "stdout"]
```

### Use as a Docker stage

```
# get the ssm-parent as a Docker stage
FROM springload/ssm-parent:1.4.1 as ssm-parent

# your main stage
FROM alpine
ENV PROJECT myproject
ENV ENVIRONMENT production

COPY --from=ssm-parent /usr/bin/ssm-parent /usr/bin/ssm-parent

ENTRYPOINT ["/usr/bin/ssm-parent", "run", "-e", "-p", "/$PROJECT/$ENVIRONMENT/backend/", "-r",  "--"]
CMD ["caddy" , "--conf", "/etc/Caddyfile", "--log", "stdout"]
```

### Config generation

If your application can't be configured via environment variables, then the following script, utilising `envsubst`, can be used to generate configs.

```
#!/bin/sh

echo "Bootstrapping Caddy"
envsubst < /etc/Caddyfile.env > /etc/Caddyfile

exec $@
```

### .env file generation

Sometimes you just want a .env file, and it is also possible.

Just specify all the same parameters, but use `dotenv` command instead with a filename to generate `.env` file.

```
./ssm-parent dotenv -r -p /project/environment dotenv.env
2018/10/01 16:37:59  info Wrote the .env file       filename=dotenv.env
```

### How to build

This project uses `go mod` as a dependency manager. Go v.1.13 was used.

```
    $git clone https://github.com/springload/ssm-parent.git
    $go build
    # (after some hacking)
    $git tag vXXX && git push && git push --tags
    $goreleaser # to create a new release
```
