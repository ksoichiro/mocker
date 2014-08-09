# CONTRIBUTING

This project is still under development, this file is just for my memo.

## Build

```sh
$ go build ./...
```

## Generate codes for Android

```sh
$ mocker gen android
```

Change output directory:

```sh
$ mocker gen android -out foo
```

## Generate and build Android app

```sh
$ mocker gen android
$ cd out
$ ./gradlew assemble
```

## Generate, build and install Android app

```sh
$ ./build.sh
```
