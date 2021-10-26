kind create cluster

# tested with istio 1.11.4

$ISTIO_HOME/bin/istioctl install --set profile=minimal -y
kubectl label namespace default istio-injection=enabled
sleep 1
kubectl apply -f $ISTIO_HOME/samples/bookinfo/platform/kube/bookinfo.yaml
# for simplicity, delete other reviews:
kubectl delete deployment reviews-v2
kubectl delete deployment reviews-v3

kubectl apply -f srvconfig.yaml
kubectl apply -f $ISTIO_HOME/samples/ratelimit/rate-limit-service.yaml

export PATH=$ISTIO_HOME/bin:$PATH

