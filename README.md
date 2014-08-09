# mocker

[![Build Status](https://travis-ci.org/ksoichiro/mocker.svg?branch=master)](https://travis-ci.org/ksoichiro/mocker)

mocker is a mock up framework for mobile apps.

## Install

```sh
$ go get github.com/ksoichiro/mocker
```

## Usage

Create mock definition file `Mockerfile`.  
Its contents is JSON format.
Then execute this:

```sh
$ mocker gen android
```

```sh
$ mocker gen ios
```

## License

Copyright (c) 2014 Soichiro Kashima  
Licensed under MIT license.  
See the bundled [LICENSE](https://github.com/ksoichiro/mocker/blob/master/LICENSE) file for details.
