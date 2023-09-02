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

The first implementation was just to making a `GET` request to the URL using the `http.Client` and parse the response
body as HTML using the `goquery` library, and use the same library to find all links.

Then I tried to recursively call all links in the response that matched the same domain, keeping a record of visited urls
to avoid calling them again.

That implementation worked fine, but it only requested one URL at a time. So, to try to have multiple request at the same time
I tried generating a goroutine on each call, but to summarize it, my implementation was a buggy and not easy to read.

Thinking on ways to solve it, I got inspired on message queues, having a queue of events and consumers of those events.
To implement that, I used a buffered channel as a queue and a number of goroutine (the number is the parallelism flag) as
consumers/workers. This solution was much cleaner (at least compared with my previous implementation), but with one flaw;
queues and consumers usually keeps running until the application stops or there is an error. Therefore, all links get
collected, but the application never stops, because there is no condition of done.
Trying to solve this issue, I took a chance using `sync.WaitGroup`, but again, the issue was to find when everything was done.
Next, I added an atomic counter. I increase the counter each time I enqueue a new URL, and decrease the counter each time
an URL gets handled. Now I know if the counter reaches 0 (zero), all links got processed. But still, I can't check that on
the worker, because not all workers might get the chance to check the condition at the right moment.  
At the end, I gave up on using a wait group and a fancy solution. So, I ended up using a `done` channel on all workers
and a for loop that sleeps for 500 ms and then checks if the counter value is 0 and length of the URLs channel is also 0
(meaning all links where processed), I close the `done` channel, making all workers to stop. That is because the workers
run an infinite for loop with a select statement that waits for values in the `urls` and `done` channels.  
Later I realized that I could remove the `done` channel, I just needed to close the `urls` channel to have the same effect.

Last thing, I implemented a `Storage` interface to handle the storage of the results (the url and the links on it).
The only current implementation is not really a storage, because it prints out the results to the `stdout`. The idea here
was to have an easy way to implement different alternatives to persist the result into a file or database, probably an
overkill for this challenge.
