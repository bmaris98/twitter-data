## Requirements
Set TWITTER_BEARER to bearer token from Twitter Dev Api. The VAR is mentioned in env.list so that it is forwarded to the docker container

## Run locally

`go mod tidy`
`go build`
`./controller` or open exe on win

## Docker
Build image
`docker build -t twitter-controller .`

Run container
`docker run -l twitter-controller --env-file env.list twitter-controller`