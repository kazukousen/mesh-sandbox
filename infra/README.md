
## Kind

```console
$ kind create cluster --name meshbox --config infra/kind/cluster.yaml
$ kind get kubeconfig --name meshbox > infra/kind/kubeconfig
```

## Contour

contour をインストール。  

```console
$ kubectl apply -f https://projectcontour.io/quickstart/contour.yaml
$ kubectl patch ds -n projectcontour envoy -p="$(cat infra/contour/patch-contour-envoy.yaml)"
```
Control Plane ノードのみに countour のL7プロキシ本体である envoy が立ち上がるようにパッチを当てる。  

```console
$ kubectl get po -n projectcontour -o wide | grep envoy
```

ノード毎にデプロイされる DaemonSet であるにも関わらず control-plane ノードにしかデプロイされていないことを確認できる。  
DaemonSet として動くこの envoy は port: 80 をリッスンしている。DaemonSet は hostPort を使用できる。  
そのため、 kind のクラスタ作成時にポートマッピングしておいたことで、 localhost:80 はこの envoy に疎通できる。  

contour のコントローラは、 Ingress リソースに反応し、 envoy に届くトラフィックを指定した Service に届ける。  
のちに Istio の `ingressgateway` Service にトラフィックを送る Ingress を記述していくこととする。  

```
...oO(
このやり方だと外部公開トラフィックを SPOF な DaemonSet で受け取ることになるし、
しかも Control Plane のネットワークを外部公開トラフィックに使ってしまうと Kubernetes クラスタの制御に影響を及ばすので、
まずこんな構成はプロダクションレベルとは全然違うなって感じ。
ローカルでの検証だけに用途を限定した構成ってことで。
```

## Istio

```console
$ istioctl manifest apply --set profile=demo
```

`profile=demo` は全部入り。  

```console
$ kubectl patch svc istio-ingressgateway -n istio-system --type json -p "$(cat infra/istio/patch-istio-ingressgateway.json)"
```

Service `istio-ingressgateway` は type が `LoadBalancer` になっているが、  
特にその必要はないので、 `NodePort` に変更するパッチを当てる。この操作自体は別に必須でわけではない。  

## Dashboard

```console
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.0-rc5/aio/deploy/recommended.yaml
$ kubectl create sa k8s-admin -n kube-system
$ kubectl create clusterrolebinding k8s-admin --clusterrole cluster-admin --serviceaccount=kube-system:k8s-admin
$ kubectl describe secret $(kubectl get secret -n kube-system | grep k8s-admin | awk '{print $1}') -n kube-system | grep token: | awk '{print $2}'
$ kubectl proxy
```

