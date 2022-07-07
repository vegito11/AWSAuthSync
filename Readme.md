
## How to Run Code

### 1) To Update/Generate the clientset/informer code run codegenerator.
```bash
hacks/codegenerator.sh
```
- It will create the clientset,informers,listers inside [client](pkg/client/) folder
- Also the CRD Manifest is created inside manifests folder

### 2) Create CRD by applying the manifests

```bash
kubectl apply -f manifests/vegito11.io_awsauthmaps/yaml
```

### 3) Either Run the code/deploy the application in K8S cluster

```bash
go run main.go 
## on another terminal create new CRD and obsever the logs
kubectl apply -f manifests/devs_map.yaml
kubectl apply -f manifests/aws-auth-cm.yaml
kubectl get AWSAuthMap
kg cm aws-auth-test -n kube-system -o yaml
```


## Referense

- [code generartor](https://github.com/kubernetes/code-generator/blob/master/examples/crd/apis/example/v1register.go)
- [Sample Controller](https://github.com/kubernetes/sample-controller)