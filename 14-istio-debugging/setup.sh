kind create cluster

$ISTIO_HOME/bin/istioctl install --set profile=minimal --set meshConfig.accessLogFile=/dev/stdout
kubectl label namespace default istio-injection=enabled
sleep 1
kubectl apply -f $ISTIO_HOME/samples/bookinfo/platform/kube/bookinfo.yaml

export PATH=$ISTIO_HOME/bin:$PATH


kubectl port-forward -n default deploy/productpage-v1 9080 &


# apply vs with the destination rule to show an error
kubectl apply -f virtual-service-all-v1.yaml

# you can fix this error with:
# kubectl apply -f destination-rule-all.yaml

# after the error is fixed you can see the route name in the logs
# kubectl logs -n default deploy/productpage-v1 -c istio-proxy
