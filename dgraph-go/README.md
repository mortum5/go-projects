# Dgraph-go

Simple url-slug store based on Dgraph db.

## Api

```sh
GET /urls/{slug} # Get url instance

POST /urls/      # Create new url instance
{
    "url": "some url",
    "slug": "some slug"
}
```

## Run

```sh
$> make dgraph    # create docker instance of dgraph 
$> go run main.go # run server
```