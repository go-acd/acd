# Amazon Cloud Drive client for Go
[![Build Status](https://travis-ci.org/go-acd/acd.svg?branch=master)](https://travis-ci.org/go-acd/acd) [![GoDoc](https://godoc.org/gopkg.in/acd.v0?status.png)](https://godoc.org/gopkg.in/acd.v0)

Amazon Cloud Drive uses
[oAuth 2.0 for authentication](https://developer.amazon.com/public/apis/experience/cloud-drive/content/restful-api-getting-started).
The [token server](https://github.com/go-acd/token-server) takes care of
the oAuth authentication. For your convenience, an instance of the
server is deployed at:

https://go-acd.appspot.com

# Usage

In order to use this library, you must authenticate through the token server. A  (e.g. https://go-acd.appspot.com).

Get and build:

```
$ go get gopkg.in/acd.v0/cmd/acd
$ go build gopkg.in/acd.v0/cmd/acd
```

Run:

```
$ $GOPATH/bin/acd ls acd://
```

Example output:

```
Backups Archived Music Pictures My Send-to-Kindle Docs
```

You may run `acd help` for additional usage information.

## Config

You must create a config file like the following and either store it as ~/.config/acd.json or provide an alternate file-path with "--config-file":

```js
{
    "tokenFile": "/home/<user>/.config/acd-token.json",
    "cacheFile": "/home/<user>/.config/acd-cache.json",
    "timeout": 0,
    "Oauth2RefreshUrl": "https://go-acd.appspot.com/refresh"
}
```

`tokenFile` represents the file containing the oauth settings which must be present on disk and has permissions 0600. The file is used by the token package to produce a valid access token by calling the oauthServer with the refresh token. You may use the public server, which that has been provided for your convenience, to generate this file: https://go-acd.appspot.com . You should deploy your own server if the security of your files is a concern.

`cacheFile` represents the file used by the client to cache the NodeTree. This file is not assumed to be present and will be created on the first run. It is gob-encoded node.Node.

`timeout` configures the HTTP Client with a timeout after which the client will cancel the request and return. It is in `time.Duration` units (nanoseconds). A timeout of 0 means no timeout. See http://godoc.org/net/http#Client for more information.

`defaultOauth2RefreshUrl` is the URL for the token-server.

`timeout` and `defaultOauth2RefreshUrl` are optional. The default values are those shown. 

If need help troubleshooting (it refuses to load the config or nothing seems to be happening), you may enable debug logging by passing "--log-level 4":

```
[ERROR] 2017/02/12 00:57:36 file has wrong permissions: want 0600 got -rw-rw-r--
```

## Library

Consult the [Godoc](https://godoc.org/gopkg.in/acd.v0) for information
on how to use the library.

# Contributions

Contributions are welcome as pull requests.

# Commit Style Guideline

We follow a rough convention for commit messages borrowed from Deis who
borrowed theirs from CoreOS, who borrowed theirs from AngularJS. This is
an example of a commit:

    feat(token): remove dependency on file system.

    use an IO.Reader and IO.Writer to deal with the token.

To make it more formal, it looks something like this:
    {type}({scope}): {subject}
    <BLANK LINE>
    {body}
    <BLANK LINE>
    {footer}

The {scope} can be anything specifying place of the commit change.

The {subject} needs to use imperative, present tense: “change”, not “changed” nor
“changes”. The first letter should not be capitalized, and there is no dot (.) at the end.

Just like the {subject}, the message {body} needs to be in the present tense, and includes
the motivation for the change, as well as a contrast with the previous behavior. The first
letter in a paragraph must be capitalized.

All breaking changes need to be mentioned in the {footer} with the description of the
change, the justification behind the change and any migration notes required.

Any line of the commit message cannot be longer than 72 characters, with the subject line
limited to 50 characters. This allows the message to be easier to read on github as well
as in various git tools.

The allowed {types} are as follows:

    feat -> feature
    fix -> bug fix
    docs -> documentation
    style -> formatting
    ref -> refactoring code
    test -> adding missing tests
    chore -> maintenance

# Credits

Although this project was built from scratch, it was inspired by the
following:

- [sgeb/go-acd](https://github.com/sgeb/go-acd)
- [yadayada/acd_cli](https://github.com/yadayada/acd_cli)
- [caseymrm/drivesink](https://github.com/caseymrm/drivesink)

# License ![License](https://img.shields.io/badge/license-MIT-blue.svg?style=plastic)

The MIT License (MIT) - see LICENSE for more details
