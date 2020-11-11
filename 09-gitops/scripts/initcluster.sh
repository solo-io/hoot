CURDIR=$(cd $(dirname $0); pwd)
kind create cluster

flux install --version=latest --namespace=flux-system

sleep 10

case "$(uname -s)" in
    Linux*)     DOCKERHOST=172.17.0.1;;
    Darwin*)    DOCKERHOST=host.docker.internal;;
    *)          echo unsupported platform; exit 1
esac
cat $CURDIR/gitrepo.yaml| sed s/HOSTIP/$DOCKERHOST/ |kubectl apply -f -

sleep 5

flux get sources git
