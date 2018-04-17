SSM Parent
----------

This is a parent process for Docker with one addition: it can read from AWS SSM Parameter store.

The way it works is that ssm-parent can be used as an entrypoint for Docker. Firstly, it retrieves all specified parameters, then injects them to the environment,
and finally runs the command.

All parameters must be in JSON format, i.e.:

```
    {
        "ENVIRONMENT": "production"
    }
```

If a few parameters are specified, all JSON entities will be read and merged into one, overriding existing keys, i.e.

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

### How to use


That should be pretty self-explanatory.

```
$ssm-parent help                                                                                                         <aws:hosting>
SSM-Parent is a docker entrypoint.

It gets specified parameters (possibly secret) from AWS SSM Parameter Store,
then exports them to the underlying process.

Usage:
  ssm-parent [command]

Available Commands:
  help        Help about any command
  print       Prints the specified parameters.
  run         Runs the specified command

Flags:
  -h, --help               help for ssm-parent
  -n, --name stringArray   Name of the SSM parameter to retrieve. Can be specified multiple times.
  -p, --path stringArray   Path to a SSM parameter. Can be specified multiple times.
  -r, --recursive          Walk through the provided SSM paths recursively.
  -s, --strict             Strict mode. Fail if found less parameters than number of names.

Use "ssm-parent [command] --help" for more information about a command.
```

The command `ssm-parent print` can be used to check the result.

### Example Dockerfile part

```
ENV PROJECT myproject
ENV ENVIRONMENT production

RUN wget -O /tmp/ssm-parent.tar.gz https://github.com/springload/ssm-parent/releases/download/v0.4/ssm-parent_0.4_linux_amd64.tar.gz && \
    tar xvf /tmp/ssm-parent.tar.gz && mv ssm-parent /sbin/ssm-parent && rm /tmp/ssm-parent.tar.gz

ENTRYPOINT ["/sbin/ssm-parent", "run", "-e", "-p", "/$PROJECT/$ENVIRONMENT/backend/", "-r",  "--"]
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

### How to build

This project uses https://github.com/golang/dep as a dependency manager. Go v.1.10.1 was used.

```
    $git clone https://github.com/springload/ssm-parent.git
    $cd ssm-parent && dep ensure
    $go build
    # (after some hacking)
    $git tag vXXX && git push && git push --tags
    $goreleaser # to create a new release
```
