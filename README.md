# Poll

`poll` is a niche tool for automating calls to asynchronous REST APIs on the Availity ARIES platform running locally on your machine. Unless you work for Availity, you probably will have no interest, unless you:

* Want to learn a little about making HTTP GET requests and parsing JSON in Go
* Want to make fun of my n00b Go code

## Installation

To install from source, you must install a Go environment [https://golang.org/doc/install](https://golang.org/doc/install), and then type:

```
$ go get github.com/hoop33/poll
$ cd $GOPATH/src/github.com/hoop33/poll
$ go build
$ go install
```

Or, you can download a pre-built binary:

* Mac OS X: [https://raw.githubusercontent.com/hoop33/poll/master/bin/macosx/poll](https://raw.githubusercontent.com/hoop33/poll/master/bin/macosx/poll)
* Windows: [https://raw.githubusercontent.com/hoop33/poll/master/bin/windows/poll.exe](https://raw.githubusercontent.com/hoop33/poll/master/bin/windows/poll.exe)
 
## Usage

To get help, run:

```
$ poll --help
```

To run `poll`, you specify a user, a base URL, a polling interval in seconds, a maximum number of times to poll the callback URL, and the resource and parameters to append to the initial URL. For example, if I wanted to use these values:

* user: `rwarner`
* base URL: `http://localhost:8280/v1/`
* interval: `2`
* maximum tries: `10`
* resource and parameters: `codes?list=ndc&q=tylenol`

I would type:

```
$ poll --user rwarner --baseurl http://localhost:8280/v1/ --interval 2 --maxtries 10 codes\?list=ndc\&q=tylenol
```

## Shortcuts, or How Do I Type Less?

Each option supports a single-letter version as well, so this command is identical to the one above:

```
$ poll -u rwarner -b http://localhost:8280/v1/ -i 2 -m 10 codes\?list=ndc\&q=tylenol
```
 
Definitely better.

But not good enough. You can also set default values through environment values, as follows:

* `user`: `$REMOTE_USER`
* `baseurl`: `$POLL_BASE_URL`
* `interval`: `$POLL_INTERVAL`
* `maxtries`: `$POLL_MAX_TRIES`

You obviously can set the environment values in `.zshrc`, `.bashrc`, or however you normally set your environment values, and then not worry about them again. Here's the same command using environment variables:

```
$ export REMOTE_USER=rwarner
$ export POLL_BASE_URL=http://localhost:8280/v1/
$ export POLL_INTERVAL=2
$ export POLL_MAX_TRIES=10
$ poll codes\?list=ndc\&q=tylenol
```

Command-line arguments override environment variables.

## Other Options

`poll` also supports the following options:

* `--verbose, -V`: Output more verbose information
* `--version, -v`: Print the version and exit
* `--help, -h`: Show help and exit

## Defaults

* Base URL: `http://localhost:8280/v1/`
* Interval: `1`
* Max Tries: `10`
* Verbose: `false`
* User: `(none)`

## License

`poll` is released under the MIT License.

## Contributing

Fork and make pull requests. I welcome them!

## FAQ

### What if the URL I'm hitting isn't asynchronous?

You can still use `poll` -- it starts polling only if the initial response returns a status code 202. Otherwise, it displays the status code and payload of the initial request/response.

