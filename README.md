# Environment Variables Manager


Define environment variables in one place and make it easy to manage variables.

[![Build Status](https://drone.io/github.com/suin/envar/status.png)](https://drone.io/github.com/suin/envar/latest)

## Problems in web application environments


Web app environments are not single; development, staging and production. However, environment variables are sometimes not same among these environments. For example, there are database username and password, SaaS API tokens and so on. So, sometimes developers make mistake like forgetting to define variables of environments.
Envar solves these problems by defining environment variables in single YAML file.


## Features


* Manages environment variables in a single YAML file.
* Validates environment variables definition and warns errors if there are missing variables.
* A simple CLI tool which runs any major platforms: Linux, Mac and Windows.
* Desinged for working with Docker.

## YAML definition example

This is a practical example of database config.

```yaml
environments: [dev, stag, prod]
variables:
  DB_HOST:
    - localhost
    - staging.example.ap-northeast-1.rds.amazonaws.com
    - production.example.ap-northeast-1.rds.amazonaws.com
  DB_PORT: 3306
  DB_USER: [root, rdsadmin, {stag}]
  DB_PASS: [root, FzN9HUrTox, {stag}]
  DB_NAME: [myapp, myapp_stag, {dev}]
```

Using this config file, Envar prints variables.

```console
$ envar print dev
# environment: dev
export DB_HOST="localhost"
export DB_NAME="myapp"
export DB_PASS="root"
export DB_PORT=3306
export DB_USER="root"
```


```console
$ envar print prod
# environment: prod
export DB_HOST="production.example.ap-northeast-1.rds.amazonaws.com"
export DB_NAME="myapp"
export DB_PASS="FzN9HUrTox"
export DB_PORT=3306
export DB_USER="rdsadmin"
```

## Installation

Download envar binary from below.

* Windows 64bit: https://drone.io/github.com/suin/envar/files/artifacts/windows-amd64/envar.exe
* Windows 32bit: https://drone.io/github.com/suin/envar/files/artifacts/windows-386/envar.exe
* Linux 64bit: https://drone.io/github.com/suin/envar/files/artifacts/linux-amd64/envar
* Linux 32bit: https://drone.io/github.com/suin/envar/files/artifacts/linux-386/envar
* OSX 64bit: https://drone.io/github.com/suin/envar/files/artifacts/darwin-amd64/envar
* OSX 32bit: https://drone.io/github.com/suin/envar/files/artifacts/darwin-386/envar

And make it executable.


```
chmod +x envar
```

## Usage

At first, you need to create `envar.yml` file.

```yaml
environments: [dev, stag, prod]
variables:
  VAR1: foo
  VAR2: bar
  VAR3: [A, B, C]
```

Then, you can print variables with envrionment name

```console
$ envar print [environment_name]
```

### Usage examples

Print development environment's variables:

```
envar print dev
```

Use different YAML file:

```
envar print prod --file another-envar.yml
```

Specify output format ([Docker compatible env-file format](https://docs.docker.com/reference/commandline/run)):

```console
$ envar print stag --output envfile
# environment: stag
VAR1=foo
VAR2=bar
VAR3=B
```

Importing variables into shell session:

```console
$ eval "$(envar print dev)"
$ env | grep VAR
41:VAR1=foo
42:VAR2=bar
43:VAR3=A
```


## envar.yaml

enver.yml is an environment definition file.

The YAML data is consists of two sections: `environments` and `variables`.

### `environments` section

In `environments` field, environment names are defined.

```yaml
environments: [dev, stag, prod]
```

If you want to add some more environments, it is possible:

```yaml
environments: [dev, stag, testing, preview, prod]
```


### `variables` section

In `variables` section, variables are defined.

In the case that the same value are used in all environments, the variable simply can be defined as primitive value:

```yaml
variables:
  DB_HOST: 127.0.0.1
```

In the case that all environments use different values, the variable is defined an array. The order of elements corresponds to `environments` fields order.

```yaml
variables:
  DB_HOST: [127.0.0.1, staging.db.local, production.db.local]
```


#### Environment symbols

Some variables may be same among some environments. In this case, **environment symbol** can be used instead of defining same value twice.

Environment symbols are notated by environment name which is wraped with brace like `{dev}`, `{stag}` and `{prod}`. Environment symbols mean that a value is same as the other environment value.

Following definition is redundant way.

```yaml
variables:
  DB_PASS: [root, FzN9HUrTox, FzN9HUrTox]
```

This definition can be rewrited with environment symbol because staging value and production value are same.

```yaml
variables:
  DB_PASS: [root, FzN9HUrTox, {stag}]
```

#### Reuse values

Sometimes across some variables, same value may be defined. In a such case, by using YAML anchors and references values can be consolidated in one place and can be reused in several variables.


```yaml
# This is a example to reuse values by using YAML anchor and reference.

environments: [dev, stag, prod]
# define values to reuse
values:
  - &db_user root
  - &db_pass mypass
  - &db_host [127.0.0.1, staging.db.local, production.db.local]
  #    |       `--- values
  #    `--- anchor names
variables:
  BACKEND_DB_USER: *db_user
  BACKEND_DB_PASS: *db_pass
  BACKEND_DB_HOST: *db_host
  FRONTEND_DB_USER: *db_user
  FRONTEND_DB_PASS: *db_pass
  FRONTEND_DB_HOST: *db_host
  #                   `--- references
```

Output will be like following:

```console
$ envar print prod -f examples/reuse_values.yml
# environment: prod
export BACKEND_DB_HOST="production.db.local"
export BACKEND_DB_PASS="mypass"
export BACKEND_DB_USER="root"
export FRONTEND_DB_HOST="production.db.local"
export FRONTEND_DB_PASS="mypass"
export FRONTEND_DB_USER="root"
```

#### Variable aliases

With anchors and references, it is possible to make aliases for several variables. This tips is useful in the case of defining different name variable with same value.

```yaml
environments: [dev, stag, prod]
variables:
  DB_PASS: &db_pass [mypass, FzN9HUrTox, weTuCgRy7n]
  DB_PASSWORD: *db_pass
  DB_PASSWD: *db_pass
```

Output will be like following:

```
$ envar print prod -f examples/aliases.yml
# environment: prod
export DB_PASS="weTuCgRy7n"
export DB_PASSWD="weTuCgRy7n"
export DB_PASSWORD="weTuCgRy7n"
```

## Working with Docker

Envar is designed to cooperate with Docker.

Envar provides output formats for `docker run` command. The first option is `docker` format. This produces arguments which docker command recognizes as environment variables:

```console
$ envar print dev -o docker
-e VAR1=foo -e VAR2=bar -e VAR3=A
```

This output can be injected into `docker run` command by quoting backticks:

```console
$ docker run `envar print dev -o docker` ubuntu env
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
HOSTNAME=38cecdaf33a2
VAR1=foo
VAR2=bar
VAR3=A
HOME=/root
```

The second option is `envfile` format. This is also recognized by Docker and docker-compose.

```shell-session
$ envar print dev -o envfile > envar.env
$ docker run --env-file envar.env ubuntu env
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
HOSTNAME=7a7c1df9a6c1
VAR1=foo
VAR2=bar
VAR3=A
HOME=/root
```

For docker-compose, set envfile path in `env_file` field of `docker-compose.yml`:

```yaml
myapp:
  image: ubuntu
  command: env
  env_file: ./envar.env
```



```shell-session
$ envar print dev -o envfile > envar.env
$ docker-compose run myapp
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
HOSTNAME=01318c21469a
TERM=xterm
VAR1=foo
VAR2=bar
VAR3=A
HOME=/root
```



## Help English improvements

I'm not native English speaker, so English improvements are welcome!

## License

The MIT License (MIT)

Copyright (c) 2015 suin

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
