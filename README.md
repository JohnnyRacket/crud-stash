# more here later ;)

# running

build

`docker build --tag=crud-stash .`

then run

`docker run -d -p 8080:8080 crud-stash`

# running with docker compose

`docker compose down`

`docker compose build`

`docker compose up`

## getting into the seed cassandra node to rune nodetool

`docker exec -it cassandra-seed-node bash`
