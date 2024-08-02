# ngql

ngql is the NebulaGraph console for NebulaGraph 5.0. With ngql, you can create a graph schema, import the demonstration `movie` dataset, and retrieve data.

## Features

- Supports interactive and non-interactive mode.
- Supports searching the history statements.
- Supports autocompletion.
- Supports multiple OS and architecture (We recommend Linux/AMD64).

## How to Install

### From Source Code

> NOTE: Go version provided with apt on ubuntu is usually "outdated".

Run the following command to examine if Go is installed on your machine.

```bash
$ go version
```

The version should be newer than 1.19.

Use Git to clone the source code of Nebula Graph Console to your host.

Before build, you should download the nebula-golang source code.

```bash
# suppose nebula-golang-5.0.0.tar.gz & nebula-ngql-5.0.0.tar.gz
tar zxvf nebula-golang-5.0.0.tar.gz -o /tmp/golang
tar zxvf nebula-ngql-5.0.0.tar.gz -o /app/ngql


cd /app/ngql
go mod edit -replace github.com/vesoft-inc/nebula-ng-tools/golang=/tmp/golang
make
```

You can find a binary named `ngql`.

### From Binary

- Download the binaries on the [Releases page](https://github.com/vesoft-inc/nebula-ng-tools/ngql/releases)

- Add execute permissions to the binary 

```bash
chmod +x ./ngql
```

## Usage

### Connect to Nebula Graph

To connect to your Nebula Graph services, use the following command.

```bash
$ ./ngql --host <host> --port <port> --user <username> --password <password>
    [-t 120] [-e "nGQL_statement" | -f filename.nGQL]
```

Options

```bash
-e, --eval string       The GQL directly
-f, --file string       The GQL script file name
-h, --help              help for ngql
-H, --host string       Nebula Graph host (default "127.0.0.1")
-p, --password string   The Nebula Graph login password
-P, --port int          The Nebula Graph port (default 9669)
-t, --timeout int       The Nebula Graph client connection timeout in seconds, 0 means never timeout
-u, --user string       The Nebula Graph login user name (default "root")
--width-max int     The max width of the column of the execution plan (default 100)
```

E.g.,

```bash
./ngql -H 192.168.8.6 -P 17163 -u root -p NebulaGraph01
Welcome to NebulaGraph 5.0, the distributed graph database offering native GQL support!
:help for help.

(root@nebula) [(none)]>
```

Check options for `./ngql -h`:

- try `./ngql` in interactive mode directly.

- And try `./ngql -e 'show graphs'` for the direct script mode.

- And try `./ngql -f demo.nGQL` for the script file mode.

## Docker

Create a container:

```bash
$ docker run --rm -ti --network nebula-docker-compose_nebula-net --entrypoint=/bin/sh vesoft/ngql:nightly
```

To connect to your Nebula Graph services, run the follow command in the container:

```bash
docker> ngql -u <user> -p <password> -H <graphd> -P 9669
```

## Print result vertically

Under regular operation, the console displays the results in a tabular format.
For instance, when executing the command `show create graph ldbc`, the output is shown as follows:

```gql
nebula> show create graph ldbc
+------------+-------------------------------------------------------+
| graph_name | create_graph_statement                                |
+------------+-------------------------------------------------------+
| "ldbc"     | "CREATE GRAPH IF NOT EXISTS `ldbc` TYPED `ldbc_type`" |
+------------+-------------------------------------------------------+
```

However, when dealing with a large number of columns or when the contents of a column are too long,
it may be more convenient to display the data vertically. This can be achieved by appending the `\G` command at the end of your query.

For example, the output of above query is displayed vertically as shown below:

```gql
nebula> show create graph ldbc \G
*************************** 1. row ***************************
            graph_name: "ldbc"
create_graph_statement: "CREATE GRAPH IF NOT EXISTS `ldbc` TYPED `ldbc_type`"
```

## Console side commands

> **NOTE**:
> The following commands are case insensitive.
> You can show all commands by `:help`

```
(root@nebula) [(none)]> :help
+---------+-------+------------------------------+------------------------------------------------------------+
| Command | Alias | Usage                        | Description                                                |
+---------+-------+------------------------------+------------------------------------------------------------+
| help    | :h    | :help                        | Show this help.                                            |
| sleep   |       | :sleep 5                     | Sleep N seconds.                                           |
| play    |       | :play movie                  | Playing the dateset                                        |
| tee     |       | :tee [-o] <filename>         | Append all results to an output file (overwrite using -o). |
| notee   |       | :notee                       | Stop writing to the output file.                           |
| pager   |       | :pager <commnad> <row_limit> | Set pager for result, default: ":pager less 200"           |
| nopager |       | :nopager                     | No pager                                                   |
| exit    | :e    | :exit                        | Exit.                                                      |
| quit    | :q    | :quit                        | Quit.                                                      |
+---------+-------+------------------------------+------------------------------------------------------------+
```

### sleep

Sleep N seconds

### play

Load the demonstration `movie` dataset

```bash
(root@nebula) [(none)]> :play movie
Playing dataset: movie...
Play dataset: movie done.
(root@nebula) [(none)]>
```

### tee & notee

Append all result to an output file

### pager & nopage

by default, use `less` command for pager if the execute result rows are more than 200.

e.g.
```bash
(root@nebula) [(none)]> :pager less 2
Pager set to less with row limit 2

(root@nebula) [(none)]> use ldbc match(v) return v limit 3
Got 3 rows (time spent 16.884ms/17.659322ms)

Wed, 15 May 2024 11:46:46 CST

```

and in pager
```bash
+--------------------------------------------------------------------------------------------------------------------------------------------------------------+
| v                                                                                                                                                            |
+--------------------------------------------------------------------------------------------------------------------------------------------------------------+
| (288232351836667905@Place:City&Continent&Country{id:6,kind:city,name:Shenzhen,url:https://shenzhen.com})                                                     |
| (288415253018968065@Place:City&Continent&Country{id:5,kind:city,name:Chengdu,url:https://chengdu.com})                                                       |
| (288826320043900931@Comment:Comment&Message{browserUsed:Chrome,content:comment1,creationDate:1991-01-01T10:00:40.213000,extent:8,id:1,locationIP:192.168.1}) |
+--------------------------------------------------------------------------------------------------------------------------------------------------------------+
(END)
```

### exit & quit

```nGQL
(root@nebula) [(none)]> :q
Bye root!

Wed, 15 May 2024 11:48:09 CST
```

## Keyboard Shortcuts

| Key Binding                                 | Description                                                |
| ------------------------------------------- | ---------------------------------------------------------- |
| <kbd>Ctrl-A</kbd>, <kbd>Home</kbd>          | Move cursor to beginning of line                           |
| <kbd>Ctrl-E</kbd>, <kbd>End</kbd>           | Move cursor to end of line                                 |
| <kbd>Ctrl-B</kbd>, <kbd>Left</kbd>          | Move cursor one character left                             |
| <kbd>Ctrl-F</kbd>, <kbd>Right</kbd>         | Move cursor one character right                            |
| <kbd>Ctrl-Left</kbd>, <kbd>Alt-B</kbd>      | Move cursor to previous word                               |
| <kbd>Ctrl-Right</kbd>, <kbd>Alt-F</kbd>     | Move cursor to next word                                   |
| <kbd>Ctrl-D</kbd>, <kbd>Del</kbd>           | (if line is *not* empty) Delete character under cursor     |
| <kbd>Ctrl-D</kbd>                           | (if line *is* empty) End of File --- quit from the console |
| <kbd>Ctrl-C</kbd>                           | Reset input (create new empty prompt)                      |
| <kbd>Ctrl-L</kbd>                           | Clear screen (line is unmodified)                          |
| <kbd>Ctrl-T</kbd>                           | Transpose previous character with current character        |
| <kbd>Ctrl-H</kbd>, <kbd>BackSpace</kbd>     | Delete character before cursor                             |
| <kbd>Ctrl-W</kbd>, <kbd>Alt-BackSpace</kbd> | Delete word leading up to cursor                           |
| <kbd>Alt-D</kbd>                            | Delete word following cursor                               |
| <kbd>Ctrl-K</kbd>                           | Delete from cursor to end of line                          |
| <kbd>Ctrl-U</kbd>                           | Delete from start of line to cursor                        |
| <kbd>Ctrl-P</kbd>, <kbd>Up</kbd>            | Previous match from history                                |
| <kbd>Ctrl-N</kbd>, <kbd>Down</kbd>          | Next match from history                                    |
| <kbd>Ctrl-R</kbd>                           | Reverse Search history (Ctrl-S forward, Ctrl-G cancel)     |
| <kbd>Ctrl-Y</kbd>                           | Paste from Yank buffer (Alt-Y to paste next yank instead)  |
| <kbd>Tab</kbd>                              | Next completion                                            |
| <kbd>Shift-Tab</kbd>                        | (after Tab) Previous completion                            |
