# Test crawler

Test realization of crawler to search all links

## Run 

```bash
docker build --tag test-crawler . && \
docker run test-crawler ./crawler -url http://some.url -rps 10
```
Where:
 - `url` is start host 
 - `rps` limitation of request per seconds.  Can be omit to no limitations

## Installation

```
go get -u github.com/an-death/go-crawler
```
## Test
```bash
docker run test-crawler go test -v -mod vendor -race
```