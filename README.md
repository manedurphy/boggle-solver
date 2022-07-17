# Boggle Solver

## Purpose

- The purpose of this application is to provide a web service to a user that solves a boggle board.

## Implmentation

### Service

- The service includes one endpoint that a user can submit a boggle board to.

```txt
POST http://<host>:<port>/solve
```

- The boggle board is a 2-dimensional array.
- The request body can be seen in the dropdown below.

<details>

<summary>request</summary>

```json
{
    "board": [
        ["D", "A", "T"],
        ["U", "A", "K"],
        ["P", "L", "A"],
        ["M", "Y", "O"],
        ["O", "G", "G"],
        ["L", "A", "N"]
    ]
}
```

</details>

- The response can be seen in the dropdown below.

<details>

<summary>response</summary>

```json
{
    "words_found": [
        "DATA",
        "PLAY",
        "GOLANG"
    ]
}
```

</details>

### Boggle Package (Revisited)

- I created an abstraction for the boggle-solving algorithm because it doesn't need to be tightly coupled with the web service.
- This boggle package can be in a separate GitHub repository for anyone to use, and the service can just consume the package.
- The algorithm I used to solve the boggle board uses a Depth First Traversal from each cell, and checks the validity of a word via a hardcoded map of words.
- I chose to use this algorithm only for its ease of implementation, not for its performance.
- I chose to hardcode the map of valid words for the sake of time. Realistically the application should have access to some data store for quick lookups.
- (Revisited) - The application now uses a Trie data structure as a word dictionary. The words are read from a `zip` file and the dictionary is initialized when the application starts. The algorithm that is used to do word lookups also limits the amount of recurive calls that are made, making it more performant than the previous implementation.

## Run Locally (Revisited)

```bash
# start the server
go run cmd/main.go

# send a request to the server
curl -d @data/request.json -H "Content-Type: application/json" http://localhost:8080/solve
```

>Note: The algorithm that is used in the "rough_draft" branch will have a slow response time. The new optimizations should return a result within 5-15ms.

## Run Using Docker

```bash
# Command for unoptimized image
docker run --name boggle --rm -d -p 8080:8080 manedurphy/boggle-solver:rough_draft

# Command for optimized image
docker run --name boggle --rm -d -p 8080:8080 manedurphy/boggle-solver:optimization
```

## Areas For Improvement (Revisited)

### Algorithm

- The first area I would like to improve is the performance of the algorithm. The article I referenced also includes a link to an optimized solution that uses a Trie data structure for the dictionary.

- (Revisited) - The algorithm in place now uses a Trie data structure for the dictionary of words. Prior to these changes, a hardcoded map was used as the dictionary, which gave the program a constant time lookup for words. But the issue was that each element triggered several recusive calls to the `findWords` method. With the new optimizations in place, recurive calls only occur when the Trie dictionary has a node with a child that matches an adjacent cell on the boggle board.

### Data Store

- Instead of hardcoding the map of valid words, the service should have an on-disk or in-memory data store that confirms the validitiy of words.
- If the application has access to a database of English words, it can run a query for each `word` that needs to be checked during the runtime of the algorithm. If the word is found in the database, then it can be cached as valid. If the word is not found, it can be cached as invalid. Before each call to the database, the application should check the cache for the existence of the potential `word`.
- I am curious to see the implementation of the Trie dictionary is performant enough to not need a cache, or if both can be used together to improve performance even further.
- (Revisited) - Based on the performance of the Trie dictionary, I know longer feel that a cache is needed to keep track of previous results. Instead, this application can be scaled horizontally and sit behind a load balancer to handle incoming traffic. The [Load Testing](#load-testing-optimization) section covers the performance of this architecture.
- (Revisited) - Another point to make is that the words that we want to use to build the Trie dictionary can be stored in a plain text `zip` file, and the dictionary can be initialized during the launch of the application. No need for an external data store.

### Logging

- There are currently only a few print statements that make up the logging. A better logging tool should be leveraged to give us better insight on what the application is doing.
- (Revisited) - The [go-hclog](https://github.com/hashicorp/go-hclog) logging tool has been included in these changes.

### Configuration

- The application should be configurable (e.g. host and port values, logging level).
- (Revisited) - A configuration file has been added with these optimizations.

<details>

<summary>config.hcl</summary>

```hcl
host = "0.0.0.0"
port = 8080
words_zip_file = "data/words.zip"
log_level = "info"
```

</details>

## Load Testing 

### Rough Draft

- The following load tests were performed with the original implementation of the boggle solver.
- From the `Apache Bench` command we can see that the number of requests sent is a low 50. This is because the application could not consistently handle more than this.

<details>

<summary>t3.medium</summary>

- 6 t2.medium EC2 nodes behind an Elastic Load Balancer (ELB)

```bash
# Apache Bench Command
ab -n 50 -c 25 -T "application/json" -p data/request.json http://boggle-114079829.us-east-1.elb.amazonaws.com/solve
```

- Results

```txt
Server Software:        
Server Hostname:        boggle-114079829.us-east-1.elb.amazonaws.com
Server Port:            80

Document Path:          /solve
Document Length:        41 bytes

Concurrency Level:      25
Time taken for tests:   55.404 seconds
Complete requests:      50
Failed requests:        0
Total transferred:      9150 bytes
Total body sent:        15850
HTML transferred:       2050 bytes
Requests per second:    0.90 [#/sec] (mean)
Time per request:       27701.911 [ms] (mean)
Time per request:       1108.076 [ms] (mean, across all concurrent requests)
Transfer rate:          0.16 [Kbytes/sec] received
                        0.28 kb/s sent
                        0.44 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:       74   82   3.7     83      92
Processing:  5378 19976 5256.2  19991   29415
Waiting:     5377 19976 5256.3  19991   29415
Total:       5452 20058 5257.8  20074   29498

Percentage of the requests served within a certain time (ms)
  50%  20074
  66%  21070
  75%  21882
  80%  25875
  90%  27805
  95%  28628
  98%  29498
  99%  29498
 100%  29498 (longest request)
```

</details>

### Optimization

- The following load tests were performed with the new code optimizations, including the Trie Dictionary.
- From the `Apache Bench` commands we can see that the number of requests sent is a high 10000. The application was capable of handling this kind of traffic with the optimizations in place.

<details>

<summary>t2.micro</summary>

- 3 t2.micro EC2 nodes behind an Elastic Load Balancer (ELB)

```bash
# Apache Bench Command
ab -n 10000 -c 1000 -T "application/json" -p data/request.json http://boggle-362284770.us-east-1.elb.amazonaws.com/solve
```

- Results

```txt
Server Software:        
Server Hostname:        boggle-362284770.us-east-1.elb.amazonaws.com
Server Port:            80

Document Path:          /solve
Document Length:        2190 bytes

Concurrency Level:      1000
Time taken for tests:   9.800 seconds
Complete requests:      10000
Failed requests:        0
Total transferred:      23120000 bytes
Total body sent:        3470000
HTML transferred:       21900000 bytes
Requests per second:    1020.41 [#/sec] (mean)
Time per request:       980.002 [ms] (mean)
Time per request:       0.980 [ms] (mean, across all concurrent requests)
Transfer rate:          2303.89 [Kbytes/sec] received
                        345.78 kb/s sent
                        2649.67 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:       70   85  11.2     82     186
Processing:    72  789 956.5    286    3761
Waiting:       72  789 956.5    286    3761
Total:        150  874 956.5    368    3847

Percentage of the requests served within a certain time (ms)
  50%    368
  66%    725
  75%   1652
  80%   2032
  90%   2517
  95%   2871
  98%   2917
  99%   3586
 100%   3847 (longest request)
```

</details>

<details>

<summary>t3.medium</summary>

- 3 t3.medium EC2 nodes behind an Elastic Load Balancer (ELB)

```bash
# Apache Bench Command
ab -n 10000 -c 1000 -T "application/json" -p data/request.json http://boggle-362284770.us-east-1.elb.amazonaws.com:81/solve
```

- Results

```txt
Server Software:        
Server Hostname:        boggle-362284770.us-east-1.elb.amazonaws.com
Server Port:            81

Document Path:          /solve
Document Length:        2190 bytes

Concurrency Level:      1000
Time taken for tests:   4.483 seconds
Complete requests:      10000
Failed requests:        0
Total transferred:      23120000 bytes
Total body sent:        3500000
HTML transferred:       21900000 bytes
Requests per second:    2230.81 [#/sec] (mean)
Time per request:       448.267 [ms] (mean)
Time per request:       0.448 [ms] (mean, across all concurrent requests)
Transfer rate:          5036.76 [Kbytes/sec] received
                        762.48 kb/s sent
                        5799.24 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:       70   82  18.3     81    1094
Processing:    74  322 169.9    281    1171
Waiting:       74  322 169.9    281    1171
Total:        151  404 172.3    363    1466

Percentage of the requests served within a certain time (ms)
  50%    363
  66%    453
  75%    507
  80%    534
  90%    621
  95%    711
  98%    861
  99%    922
 100%   1466 (longest request)
```

- 6 t3.medium EC2 nodes behind an Elastic Load Balancer (ELB)

```bash
# Apache Bench Command
ab -n 10000 -c 1000 -T "application/json" -p data/request.json http://boggle-114079829.us-east-1.elb.amazonaws.com/solve
```

```txt
Server Software:        
Server Hostname:        boggle-114079829.us-east-1.elb.amazonaws.com
Server Port:            80

Document Path:          /solve
Document Length:        2190 bytes

Concurrency Level:      1000
Time taken for tests:   3.369 seconds
Complete requests:      10000
Failed requests:        0
Total transferred:      23120000 bytes
Total body sent:        3470000
HTML transferred:       21900000 bytes
Requests per second:    2967.91 [#/sec] (mean)
Time per request:       336.938 [ms] (mean)
Time per request:       0.337 [ms] (mean, across all concurrent requests)
Transfer rate:          6700.98 [Kbytes/sec] received
                        1005.73 kb/s sent
                        7706.70 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:       75  105  77.8     99    1134
Processing:    75  168  83.3    147     924
Waiting:       75  167  83.2    147     924
Total:        162  273 113.8    247    1313

Percentage of the requests served within a certain time (ms)
  50%    247
  66%    257
  75%    266
  80%    274
  90%    324
  95%    440
  98%    639
  99%    780
 100%   1313 (longest request)
```

</details>

<details>

<summary>t2.micro+t3.medium</summary>

- 3 t2.micro and 3 t3.medium EC2 nodes behind an Elastic Load Balancer (ELB)

```bash
# Apache Bench Command
ab -n 10000 -c 1000 -T "application/json" -p data/request.json http://boggle-362284770.us-east-1.elb.amazonaws.com/solve
```

- Results

```txt
Server Software:        
Server Hostname:        boggle-362284770.us-east-1.elb.amazonaws.com
Server Port:            80

Document Path:          /solve
Document Length:        2190 bytes

Concurrency Level:      1000
Time taken for tests:   3.771 seconds
Complete requests:      10000
Failed requests:        0
Total transferred:      23120000 bytes
Total body sent:        3470000
HTML transferred:       21900000 bytes
Requests per second:    2651.50 [#/sec] (mean)
Time per request:       377.145 [ms] (mean)
Time per request:       0.377 [ms] (mean, across all concurrent requests)
Transfer rate:          5986.58 [Kbytes/sec] received
                        898.51 kb/s sent
                        6885.09 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:       72  160 120.6    102     536
Processing:    73  176 123.4    119     545
Waiting:       73  176 123.4    118     544
Total:        150  336 236.0    231    1016

Percentage of the requests served within a certain time (ms)
  50%    231
  66%    303
  75%    362
  80%    432
  90%    799
  95%    884
  98%    926
  99%    960
 100%   1016 (longest request)

```

</details>

## Sources

- [How to play Boggle](https://www.youtube.com/watch?v=BJAdXnGAb7k)
- [Find all possible words in a board of characters](https://www.geeksforgeeks.org/boggle-find-possible-words-board-characters/?ref=lbp)
- [English Words](https://github.com/dwyl/english-words)
