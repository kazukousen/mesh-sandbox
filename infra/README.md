
## Kind

```console
$ kind create cluster --name meshbox --config infra/kind/cluster.yaml
$ kind get kubeconfig --name meshbox > infra/kind/kubeconfig
```

## MetalLB

[docs](https://metallb.universe.tf/installation/#installation-with-kubernetes-manifests)

```console
$ docker network inspect bridge | jq .[0].IPAM.Config[0].Subnet
```

kind は `bridge` という名前の docker ネットワークに紐づけられているので、  
サブネットのIP範囲を把握し、 `infra/metallb/condfigmap.yaml` に設定する。  

```console
$ kubectl apply -f https://raw.githubusercontent.com/google/metallb/v0.9.3/manifests/namespace.yaml
$ kubectl apply -f https://raw.githubusercontent.com/google/metallb/v0.9.3/manifests/metallb.yaml
# On first install only
$ kubectl create secret generic -n metallb-system memberlist --from-literal=secretkey="$(openssl rand -base64 128)"
$ kubectl apply -f infra/metallb/configmap.yaml
```

## Istio

```console
$ istioctl manifest apply --set profile=demo
```

## Dashboard

```console
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.0-rc5/aio/deploy/recommended.yaml
$ kubectl create sa k8s-admin -n kube-system
$ kubectl create clusterrolebinding k8s-admin --clusterrole cluster-admin --serviceaccount=kube-system:k8s-admin
$ kubectl describe secret $(kubectl get secret -n kube-system | grep k8s-admin | awk '{print $1}') -n kube-system | grep token: | awk '{print $2}'
$ kubectl proxy
```

