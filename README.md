# hbn
**hbn** - http benchmark tool written on golang

---

Fast, simple and powerfull!!!

---

## Installation

In UNIX based:
```shell
$ git clone https://github.com/out-of-mind/hbn && cd hbn/
```
```shell
$ make
```

In windows:
download zip file in git repo, than unzip it
than change your current dir to hbn
```shell
$ go build -v ./src/hbn.go
```

## Usage

```shell
$ hbn -url <host to attack> -w <count of workers> -d <duration of testing> -c <path/to/your/config.json> -m <http method to use (only get method support, sorry)> -uu <use useragents> -uh <use headers> -uc <use cookies>
```

## Help

```shell
$ hbn -h
```

or
```shell
$ hbn --help
```

or
```shell
$ hbn -help
```

## License

The MIT License
