Backend.
<br>
- (Will) Exposes an HTTP API which enables/disables queries (See tweet search API,q param)
- Fetches data periodically and outputs csv files to `/tmp/{query}/{timestamp}.csv`
- Sends the csv files to hadoop
- Syncs HDFS
- Is notified when Hadoop Job exits
- Manages statuses & responses w/ Mongo


<hr>

## Requirements
Set TWITTER_BEARER to bearer token from Twitter Dev Api. The VAR is mentioned in env.list so that it is forwarded to the docker container

## Run locally

`go mod tidy` <br>
`go build` <br>
`./controller` or open exe on win

## Docker
Build image
`docker build -t twitter-controller .`

Run container
`docker run -l twitter-controller --env-file env.list twitter-controller`
