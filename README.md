## DOCKER
1. Build Docker image:

```
docker build -t simplebank .
```

2. Run images
   
```
docker run --name simplebank_app -p 8080:8080 simplebank
```

3. View network container
   
```
docker container inspect [OPTIONS] CONTAINER [CONTAINER...]
```

```
docker container inspect postgres12
```
4. Run docker image with new env

```
docker run --name simplebank_app  -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:root@172.17.0.2:5432/simple_bank?sslmode=disable" simplebank
```


5. View networks inspect

```
docker network inspect [OPTIONS] NETWORK [NETWORK...]
```

```
docker network inspect bank_network
```

6. Create new network

```
docker network create [OPTIONS] NETWORK
```

```
docker network create bank-network
```

7. Connect to container

```
docker network connect [OPTIONS] NETWORK CONTAINER
```

```
docker network connect bank-network  postgres12
```


8. Run container with new network

```
docker run --name simplebank_app --network bank-network -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:root@postgres12:5432/simple_bank?sslmode=disable" simplebank
```



## K9s

1. Connect to AWS EKS cluter

```
aws eks update-kubeconfig --name simple-bank --region ap-southeast-1
```
-   Get the services that are running on the cluster

```
kubectl get service
```

-   Get list of running pods

```
kubectl get pods
```