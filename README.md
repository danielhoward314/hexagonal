## Hexagonal Design
I followed the great 3-part tutorial series from the Tensor Programming YouTube channel. [Part one link here](https://youtu.be/rQnTtQZGpg8). The series shows how to follow the hexagonal design pattern in Go. The tutorial covers connecting the core logic to Redis and Mongodb. I decided to try as an exercise following the same pattern to connect to Postgres and ElasticSearch.


### Prerequisites:

- Install redis: 
1. link that will download latest stable redis version is [here](http://download.redis.io/redis-stable.tar.gz)
2. in terminal: 
```
tar xvzf redis-stable.tar.gz
cd redis-stable
make
sudo make install
```
- Install mongodb:
1. in terminal: 
```
brew tap mongodb/brew
brew install mongodb-community@4.4 // or different version
```
- Install postgres
1. in terminal: ` brew install postgresql`
2. best practice would be to set up dedicated users (with a shared role(s) that users could be granted), but for the sake of simple setup create a superuser with name postgres with the following: `createuser -s postgres`
- Install Elasticsearch:
1. in terminal:
```
brew tap elastic/tap
brew install elastic/tap/elasticsearch-full
```

### Configuration

This repo expects a `.env` file in root directory with the following values:

```
PORT=8080
REDIS_URL=redis://localhost:6379
MONGO_URL=mongodb://localhost/shortener
MONGO_DB_NAME=shortener
MONGO_TIMEOUT=30
POSTGRES_DB_NAME=shortener
POSTGRES_HOST=localhost
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_PORT=5432
```

### To run each:

- for Redis, from the `redis-stable` directory, run: `redis-server`
- for Mongo, run: `brew services start mongodb-community@4.4` (same command with `stop` stops it)
- for Postgres, run: `brew services start postgresql` (same command with `stop` stops it)
    and if you want to run sql on the postgresql server, run `psql postgres`
- for ES, run: `brew services start elastic/tap/elasticsearch-full`

Run the go server:

- `go run main.go -r=<repoType>` where `<repoType>` is `redis|mongo|postgres|es` (defaults to `redis`)

Test it with curl and a browser:

`curl localhost:8080 -XPOST -d '{"url":"https://google.com"}'`

The response should include a `code`. In browser, go to `localhost/{code}` and it should redirect to google.