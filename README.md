# terraform-tfstate-retriever
Retrieve terraform.tfstate file and write it to a ConfigMaps

Several environments:

- CONFIGMAPS_NAMESPACE
  
ConfigMap namespace

- CONFIGMAPS_NAME
  
ConfigMap name

- TF_STATE_DIR

Terraform state file directory

- TF_STATE_NAME

Terraform state file name

## Build image


```shell
# It's so slow to build binary in the docker, so build it locally
$ GOOS=linux go build -o ./terraform-tfstate-retriever .


$ docker build -t zzxwill/terraform-tfstate-retriever:v0.1 .

$ docker push zzxwill/terraform-tfstate-retriever:v0.1

```

## Run it in-cluster

```shell
✗ kubectl apply -f rbac.yaml
clusterrole.rbac.authorization.k8s.io/tf-clusterrole created
clusterrolebinding.rbac.authorization.k8s.io/tf-binding created

✗ kubectl apply -f sample-deployment.yaml
deployment.apps/poc created
```


```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: poc
  namespace: default
spec:
  selector:
    matchLabels:
      name: poc
  template:
    metadata:
      labels:
        name: poc
    spec:
      containers:
        - env:
            - name: CONFIGMAPS_NAMESPACE
              value: default
            - name: CONFIGMAPS_NAME
              value: aliyun-oss-tf-state
            - name: TF_STATE_DIR
              value: /go/src/app
            - name: TF_STATE_NAME
              value: go.mod
          image: zzxwill/terraform-tfstate-retriever:v0.1
          imagePullPolicy: Always
          name: terraform-tfstate-retriever

```

```shell
✗ kubectl get cm aliyun-oss-tf-state
NAME                  DATA   AGE
aliyun-oss-tf-state   1      8m37s
```