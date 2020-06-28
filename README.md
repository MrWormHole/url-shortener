# url-shortener
Url shortener microservice with multi database(mongo and redis) and multi format(json and msgpack) support 🐳

##### To use this, run the command down below. Default selected database is redis for now. JSON and MSGPACK formats are supported as long as right content type is received.
##### For mongodb, you need to add mongodb to docker-compose file and provide environment values(URL_DB: "mongo", MONGO_URL, MONGO_DB, MONGO_TIMEOUT)
```
docker-compose -d up
```

##### POST /
##### request:
```
{
  "url" : "www.website.com/asdasdasdadasd"
}
```
##### response:
```
{
  "hash": "ASDXASD"
  "url" : "www.website.com/asdasdasdadasd"
  "created_at": 3241233443521
}
```
##### GET /:hash
```
  This will redirect you to saved url.
```
