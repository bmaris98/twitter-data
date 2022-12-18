## Infra docker files go here
i.e. Mongo or Kafka clusters

create docker network
docker network create twitter-net



start mongo:
docker run -d -p 27017:27017 -v C:\school\twitter-data\infra\volumes\mongo:/data/db --name big-data-mongo --env MONGO_INITDB_ROOT_USERNAME=mongoadmin --env MONGO_INITDB_ROOT_PASSWORD=admin --env MONGO_INITDB_DATABASE=TWITTERDATA mongo:5.0.14_DATABASE=TWITTERDATA mongo:5.0.14

