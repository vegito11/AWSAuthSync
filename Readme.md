## Info 
- Using this operator, we can effectively maintain the aws-auth configmap without having to manually update it.

- After creation of  AWSAuthMap (CR) object. According to the entries specified in our AWSAuthMap Object, our controller will update the aws-auth configmap.

## How to Run Code

### 1) To Update/Generate the clientset/informer code run codegenerator.
```bash
hacks/codegenerator.sh
```
- It will create the clientset,informers,listers inside [client](pkg/client/) folder
- Also the CRD Manifest is created inside manifests folder

### 2) Create CRD by applying the manifests

```bash
kubectl apply -f manifests/vegito11.io_awsauthmaps.yaml
```

### 3) To test without deploying controller 

```bash
go run main.go 
## on another terminal create new CRD and obsever the logs
kubectl apply -f manifests/aws-auth-cm.yaml
kubectl get AWSAuthMap
kg cm aws-auth-test -n kube-system -o yaml
kubectl apply -f manifests/devs_map.yaml
kg cm aws-auth-test -n kube-system -o yaml
```

### 4) To deploy controller on k8s cluster.    

1. create docker image and push it (update the image repo accordingly). 
   
    ```bash 
    docker build -t vegito/awsauthsyncer:0.0.1 . 
    ```
    
2. Now create CRD and aws-auth-test cm on cluster

    ```bash
    kubectl apply -f manifests/aws-auth-cm.yaml
    kubectl apply -f manifests/vegito11.io_awsauthmaps.yaml
    kubectl get AWSAuthMap
    ```

3. Now deploy the Controller with necessary role and binding

    ```bash
    kubectl apply -f manifests/templates
    ```
## Examples

1. To see the CRD it register in cluser

```bash
❯ kubectl api-resources | grep aws
awsauthmaps                       awth         vegito11.io/v1beta                     true         AWSAuthMap
```

2. Create new Object using CRD and verify that configmap is updated

```bash
# 1) before creating object 
> kubectl cm aws-auth-test -n kube-system -o yaml

> kubectl apply -f manifests/devs_map.yaml

# 2) after creating object 
> kubectl cm aws-auth-test -n kube-system -o yaml

```

3. Delete the CR and verify that entries has been removed from aws-auth-test

```bash
> kubectl delete awth devauth 
> kubectl cm aws-auth-test -n kube-system -o yaml
```

4. List the object 

```bash
> kubectl get AWSAuthMap

# Using shortname
> kubectl get awth 
NAME     ROLE                USER   AGE
qaauth   aws::/role/qarole          49m
```

5. To see hidden columns
```bash
❯ kubectl get awth -o wide
NAME      ROLE                  USER                AGE   ALLUSER
devauth   aws::/role/testrole   aws::/user/devops   6s    [{"groups":["devops"],"userarn":"aws::/user/devops","username":"devops"}]
qaauth    aws::/role/qarole                         52m   
```

## Limitation and Improvement

1. Whenever an existing CR object is updated. The aws-auth map will retain any entries that are deleted from our CR object.

2. Because it hasn't been thoroughly tested and was created only for learning purposes, we should implement it in production right now.

3. We can have admission webhook which can verify the Userarn/Rolearn , So that we can prevent the common typos from happening

4. We can have the Mutating webhook to set the username/group by parsing Userarn/Rolearn

## Referense

- [code generartor](https://github.com/kubernetes/code-generator/blob/master/examples/crd/apis/example/v1register.go)
- [Sample Controller](https://github.com/kubernetes/sample-controller)
- [Kube Builder Tags for Auto Generation](https://book.kubebuilder.io/reference/markers/crd.html)