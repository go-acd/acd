# Amazon Cloud Drive client for Go
[![Build Status](https://travis-ci.org/go-acd/acd.svg?branch=master)](https://travis-ci.org/go-acd/acd) [![GoDoc](https://godoc.org/gopkg.in/acd.v0?status.png)](https://godoc.org/gopkg.in/acd.v0)

Amazon Cloud Drive uses
[oAuth 2.0 for authentication](https://developer.amazon.com/public/apis/experience/cloud-drive/content/restful-api-getting-started).
The [token server](https://github.com/go-acd/token-server) takes care of
the oAuth authentication. For your convenience, an instance of the
server is deployed at:

https://go-acd.appspot.com

# Install

This project is go-gettable:

```
go get gopkg.in/acd.v0/...
```

# Usage

In order to use this library, you must authenticate through the [token server](https://go-acd.appspot.com).

## CLI

Run `acd help` for usage.

## Library

Consult the [Godoc](https://godoc.org/gopkg.in/acd.v0) for information
on how to use the library.

# Contributions

This repository does not accept pull requests. All contributions must go
through [Phabricator](http://phabricator.nasreddine.com). Please refer
to [Arcanist Quick Start](https://secure.phabricator.com/book/phabricator/article/arcanist_quick_start/)
to install arc and learn basic commands.

# Credits

Although this project was built from scratch, it was inspired by the
following:

- [sgeb/go-acd](https://github.com/sgeb/go-acd)
- [yadayada/acd_cli](https://github.com/yadayada/acd_cli)
- [caseymrm/drivesink](https://github.com/caseymrm/drivesink)

# License

Refer to the LICENSE file.
