
killall gitserver

CURDIR=$(cd $(dirname $0);cd ..; pwd)
cd gitserver && go build .

cd $(mktemp -d)
echo Initializing git repo in:
echo $PWD

git init .

touch README.md
git add README.md
git commit -m "initial commit"
echo revision: $(git rev-parse HEAD)

cp -r $CURDIR/prod-cluster prod-cluster


echo starting git server
$CURDIR/gitserver/gitserver 2>&1 > /dev/null &

# wait for input
echo Waiting to add app to cluster
read

git add .
git commit -am "prod"
echo revision: $(git rev-parse HEAD)

echo commited app. waiting to commit additional changes
read

cp $CURDIR/scripts/anotherapp.yaml prod-cluster/app/
echo "- ./anotherapp.yaml" >> prod-cluster/app/kustomization.yaml

git add .
git commit -am "prod - 2nd app"
echo revision: $(git rev-parse HEAD)

echo commited app. wait to enter shell
read
$SHELL

kill %1