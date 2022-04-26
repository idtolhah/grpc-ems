# grpc-ems

# Development

## Run Mongo & Mysql
docker-compose --file local.yml up -d

## Run Local
sh start-local.sh 

# Production
## Via Docker Compose
docker-compose up -d

# Run On K8s
## Build Image and Run
docker compose -f docker-compose-k8s.yml up -d
sh start-k8s.sh

## Live building & reloading
cd packing-query
skaffold dev