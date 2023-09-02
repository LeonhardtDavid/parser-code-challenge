# parser-code-challenge

## Challenge Description

We’d like you to write a simple web crawler in Golang.
Given a starting URL, the crawler should visit each URL it finds on the same domain. It should
print each URL visited, and a list of links found on that page. The crawler should be limited to
one subdomain – so when you start with https://parserdigital.com/, do not follow external
links, for example to facebook.com or community.parserdigital.com.
We would like to see your own implementation of a web crawler. Please do not use frameworks
like scrappy or go-colly which handle all the crawling behind the scenes or someone else’s
code. You are welcome to use libraries to handle things like HTML parsing.
Ideally, write it as you would a production piece of code. This exercise is not meant to show us
whether you can write code – we are more interested in how you design software. This means
that we care less about a fancy UI or sitemap format, and more about how your program is
structured: the trade-offs you’ve made, what behaviour the program exhibits, and your use of
concurrency, test coverage, and so on.
Once you have submitted your task, we will then schedule a session with an engineer, during
which we all will discuss your implementation.
When you’re ready, please submit your solution as a ZIP file.

## Challenge Solution

### Running the crawler

_[Go](https://go.dev/) is required to run the application.
[make](https://man7.org/linux/man-pages/man1/make.1.html) is useful, but not required, you can run the commands that are in the Makefile.
Vendor directory was added to void download the dependencies before running the tests or build._

```shell
make build
./bin/crawler -url 'https://parserdigital.com/'
```

Available flags:

| Flag        | Default | Description                        |
|-------------|:-------:|------------------------------------|
| url         |         | URL on which the crawler will run. |
| parallelism |    5    | Max number of concurrent requests. |

Examples:

```shell
./bin/crawler -url 'https://google.com.ar'
```

```shell
./bin/crawler -url 'https://google.com.ar' -parallelism 50
```

### Running the tests

For just the tests:

```shell
make test
```

For tests with coverage:

```shell
make test-coverage
```

The coverage report will be generated in the `coverage` directory.

### The Journey
