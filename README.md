# parx

`parx` is a simple tool to run multiple commands in parallel while having the output structured like [Docker Compose](https://docs.docker.com/compose/) does that. This is useful when developing on a distributed system without using Docker and Docker-Compose. Inspired by [k-bx/par](https://github.com/k-bx/par).

To install this tool simply download the binary from the current release and put it somewhere in your `$PATH`. With `parx -h` you should now be able to see the help page.

The exit code of the command only refers to its own logic, not to the exit code of the processes. If any process fails, but the command generally succeeded, the exit code will still be 0.

## Quick Usage

The simplest usage is to append each command with `""` escaped to the arguments. Now two processes get started up and exit if they are done.

```sh
> parx "echo Hello World; sleep 5; echo Bye" "sleep 10; echo Ciao!"
process_0 | Hello World
process_0 | Bye
process_0 exited
process_1 | Ciao!
process_1 exited
```

The prefix for each process is colored, just like Docker Compose does it. This is only enabled, if a valid tty is attached. See [this package](https://github.com/fatih/color) for more. Also, it is generated, as there is no nice way to set the name via command line.

## Extended usage

For extended usage, it is recommended to use a `parx.yml` file to describe the processes and their commands.

Example from above as config:

```yaml
processes:
  - name: hello
    command: echo Hello $NAME; sleep 5; echo Bye
    env:
      NAME: Thomas
  - name: ciao
    command: sleep 10; echo Ciao!
```

You can now store the file whereever you like and reference it via `parx -f path/to/your-file.yml`. The prefixes are now set to the name of the process. For a full configuration, have a look at the `parx.example.yml` in this repo.