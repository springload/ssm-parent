[![Go Report Card](https://goreportcard.com/badge/github.com/springload/ssm-parent)](https://goreportcard.com/report/github.com/springload/ssm-parent)

SSM Parent
----------

This is mostly a parent process for Docker with one addition: it can read from AWS SSM Parameter store.

Please note, that it still requires a proper `init` process, for example the one embedded into Docker can be used with `docker run --init`.

The way it works is that ssm-parent can be used as an entrypoint for Docker. Firstly, it retrieves all specified parameters, then injects them to the environment,
and finally runs the command using `execve` syscall.

All parameters must be in JSON format, i.e.:

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


That should be pretty self-explanatory.

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
  -e, --expand                   Expand arguments and values using shell-style syntax
  -h, --help                     help for ssm-parent
  -n, --name stringArray         Name of the SSM parameter to retrieve. Expects JSON in the value. Can be specified multiple times.
  -p, --path stringArray         Path to a SSM parameter. Expects JSON in the value. Can be specified multiple times.
      --plain-name stringArray   Name of the SSM parameter to retrieve. Expects actual parameter in the value. Can be specified multiple times.
      --plain-path stringArray   Path to a SSM parameter. Expects actual parameter in the value. Can be specified multiple times.
  -r, --recursive                Walk through the provided SSM paths recursively.
  -s, --strict                   Strict mode. Fail if found less parameters than number of names.
      --version                  version for ssm-parent

Use "ssm-parent [command] --help" for more information about a command.
```

The command `ssm-parent print` can be used to check the result.

### Example Dockerfile part

```
ENV PROJECT myproject
ENV ENVIRONMENT production

RUN wget -O /tmp/ssm-parent.tar.gz https://github.com/springload/ssm-parent/releases/download/v0.9/ssm-parent_0.9_linux_amd64.tar.gz && \
    tar xvf /tmp/ssm-parent.tar.gz && mv ssm-parent /usr/bin/ssm-parent && rm /tmp/ssm-parent.tar.gz

ENTRYPOINT ["/usr/bin/ssm-parent", "run", "-e", "-p", "/$PROJECT/$ENVIRONMENT/backend/", "-r",  "--"]
CMD ["caddy" , "--conf", "/etc/Caddyfile", "--log", "stdout"]
```

### Use as a Docker stage

```
# get the ssm-parent as a stage
FROM springload/ssm-parent as ssm-parent

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

Sometimes you just want a .env file, and now it is also possible.

Just specify all the same parameters, but use `dotenv` command instead with a filename to generate `.env` file.
```
./ssm-parent dotenv -r -p /project/environment dotenv.env
2018/10/01 16:37:59  info Wrote the .env file       filename=dotenv.env
```

### How to build

This project uses `go mod` as a dependency manager. Go v.1.12 was used.

```
    $git clone https://github.com/springload/ssm-parent.git
    $go build
    # (after some hacking)
    $git tag vXXX && git push && git push --tags
    $goreleaser # to create a new release
```
