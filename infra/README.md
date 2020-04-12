
## Kind

```console
$ kind create cluster --name meshbox --config infra/kind/cluster.yaml
$ kind get kubeconfig --name meshbox > infra/kind/kubeconfig
```

## Contour

```console
$ kubectl apply -f https://projectcontour.io/quickstart/contour.yaml
$ kubectl patch ds -n projectcontour envoy -p="$(cat infra/contour/patch-contour-envoy.yaml)"
```

contour をインストール。  
Control Plane ノードのみに countour のL7プロキシ本体である envoy が立ち上がるようにパッチを当てる。  

```console
$ kubectl get po -n projectcontour -o wide | grep envoy
```

ノード毎にデプロイされる DaemonSet であるにも関わらず control-plane ノードにしかデプロイされていないことを確認できる。  

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

