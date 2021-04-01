# terraform-tfstate-retriever
Retrieve terraform.tfstate file and write it to a ConfigMaps 

## Build image

```shell
$ docker build -t zzxwill/terraform-tfstate-retriever:v0.1 .

$ docker push zzxwill/terraform-tfstate-retriever:v0.1

```

## Run it in-cluster

```shell
$ kubectl create clusterrolebinding default-view --clusterrole=view --serviceaccount=default:default

```