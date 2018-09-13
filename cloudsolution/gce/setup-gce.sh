
gcloud container clusters create passa-cluster
kubectl create deployment nginx --image=nginx
gcloud container node-pools create t2-micro --cluster passa-cluster --machine-type f1-micro
gcloud container node-pools create t2-large --cluster passa-cluster --machine-type g1-small
gcloud container node-pools delete default-pool --cluster=passa-cluster -q
