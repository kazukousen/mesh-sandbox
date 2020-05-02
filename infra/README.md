
## Kind

```console
$ kind create cluster --name meshbox --config infra/kind/cluster.yaml
```

Dcoker for Desktop (Mac) では、 docker コンテナとホストOSが疎通するためには `--port` オプションによるポートマッピングが必要。  
Dockerコンテナでノードを構成する Kind では、クラスタ作成時に、疎通したいノードに対しポートマッピングの設定をしておく。  

今回は Control Plane ノードの :80, :443 をホストOSとポートマッピングする。  

```console
$ kind get kubeconfig --name meshbox > infra/kind/kubeconfig
```

## Metrics-Server

```console
$ kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/download/v0.3.6/components.yaml
```

`kubectl top nodes|pods` コマンドがうまくいくはずが・・・  

```console
$ kubectl get po -n kube-system | grep metrics-server | awk '{print $1}'
$ kubectl logs -f metrics-server-XXXXXXX -n kube-system
```

`no such host` 的なエラーがでた。  
Issue を参照。  
https://github.com/kubernetes-sigs/metrics-server/issues/131

```console
$ kubectl edit deploy metrics-server -n kube-system
```

`metrics-server` へ渡すコマンドライン引数を追加。  

```yaml
        args:
        - --kubelet-insecure-tls
        - --kubelet-preferred-address-types=InternalIP
```

起動後すぐは `unable to fetch node metrics for node "XXX": no metrics known for node` と出るが、ちょっと待てばOK.  


## Contour

https://kind.sigs.k8s.io/docs/user/ingress/

Ingress Controller 実装である [contour](https://projectcontour.io/) を使っていく。  

```console
$ kubectl apply -f https://projectcontour.io/quickstart/contour.yaml
```
contour をインストール。  
このマニフェストは、 namespace `projectcountour` に、  
envoy を Workload は DaemonSet でデプロイし、コントローラである contour を Workload は Deployment で 2 replicas デプロイする。  

```console
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
まずこんな構成ではプロダクションレディとはいえないだろーなって感じ。
検証用途での最低限の構成ってことで。
```

## Istio

```console
$ istioctl manifest apply --set profile=default --set values.prometheus.enabled=false
```

default は `istio-ingressgateway` と Prometheus が有効になっている。 Prometheus は別で入れたいので無効にする。  


`istiod` は memory の requests が 2Gi に設定されてるので、  
`insuffient memory` になっちゃったことがあった。節約したいときは edit するといいかも。  

```console
$ kubectl patch svc istio-ingressgateway -n istio-system --type json -p "$(cat infra/istio/patch-istio-ingressgateway.json)"
```

Service `istio-ingressgateway` は type が `LoadBalancer` になっているが、  
特にその必要はないので、 `NodePort` に変更するパッチを当てる。この操作自体は別に必須でわけではない。  

## Dashboard

```console
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.0/aio/deploy/recommended.yaml
$ kubectl create sa k8s-admin -n kube-system
$ kubectl create clusterrolebinding k8s-admin --clusterrole cluster-admin --serviceaccount=kube-system:k8s-admin
$ kubectl describe secret $(kubectl get secret -n kube-system | grep k8s-admin | awk '{print $1}') -n kube-system | grep token: | awk '{print $2}'
$ kubectl proxy
```

