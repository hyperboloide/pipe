# Piped
piped is an HTTP service that allows the transformation and transfer of binary streams (like files). It use the pipe library
to encode the streams (like gzip or aes) and transfer them to remote endpoints (ex: S3, GCS, or files).

to install:
```sh
go install github.com/hyperboloide/pipe/piped
```

See the usage bellow on how to use it (or enter `piped --help` in a terminal):
```
usage: piped [<flags>] [<config>]

Flags:
      --help       Show context-sensitive help (also try --help-long and --help-man).
  -p, --port=7890  Port number for of the HTTP service.
  -s, --silent     Do not log requests.
      --version    Show application version.

Args:
  [<config>]  Path to the configuration file.
```
The configuration is a simple json file (examples can be found in the [examples](https://github.com/hyperboloide/pipe/tree/master/piped/examples)
directory)
You can provide the path of the configuation file as an argument otherwise the program will search for the following locations:
1. `./piped.json`
2. `/etc/piped/config.json`
3. `$HOME/.piped/config.json`
