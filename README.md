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

### Boggle Package

- I created an abstraction for the boggle-solving algorithm because it doesn't need to be tightly coupled with the web service.
- This boggle package can be in a separate GitHub repository for anyone to use, and the service can just consume the package.
- The algorithm I used to solve the boggle board uses a Depth First Traversal from each cell, and checks the validity of a word via a hardcoded map of words.
- I chose to use this algorithm only for its ease of implementation, not for its performance.
- I chose to hardcode the map of valid words for the sake of time. Realistically the application should have access to some data store for quick lookups.

## Run Locally

```bash
# start the server
go run cmd/main.go

# send a request to the server
curl -d @data/request.json -H "Content-Type: application/json" http://localhost:8080/solve
```

>Note: The response time is long due to the inefficient algorithm that is used.

## Areas For Improvement

### Algorithm

- The first area I would like to improve is the performance of the algorithm. The article I referenced also includes a link to an optimized solution that uses a Trie data structure for the dictionary.

### Data Store

- Instead of hardcoding the map of valid words, the service should have an on-disk or in-memory data store that confirms the validitiy of words.
- If the application has access to a database of English words, it can run a query for each `word` that needs to be checked during the runtime of the algorithm. If the word is found in the database, then it can be cached as valid. If the word is not found, it can be cached as invalid. Before each call to the database, the application should check the cache for the existence of the potential `word`.
- I am curious to see the implementation of the Trie dictionary is performant enough to not need a cache, or if both can be used together to improve performance even further.

### Logging

- There are currently only a few print statements that make up the logging. A better logging tool should be leveraged to give us better insight on what the application is doing.

### Configuration

- The application should be configurable (e.g. host and port values, logging level).

## Sources

- [How to play Boggle](https://www.youtube.com/watch?v=BJAdXnGAb7k)
- [Find all possible words in a board of characters](https://www.geeksforgeeks.org/boggle-find-possible-words-board-characters/?ref=lbp)
