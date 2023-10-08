## GenericHTTPPopulator

A Kubernetes volume populator to populate the volumes with content of an HTTP
endpoint. [This document](https://kubernetes.io/blog/2022/05/16/volume-populators-beta/)
can be referred to read more about the Volume Populators.

## Install

### Installing custom populator's (`GenericHTTPPopulator`) CRD

Install the CRD for our new volume populator (`GenericHTTPPopulator`). A CR of
this type (`GenericHTTPPopulator`) can be used in the PVC's `dataSourceRef`
field to specify this type of populator.

```bash!
kubectl apply -f https://raw.githubusercontent.com/viveksinghggits/generichttppopulator/master/manifests/genrichttppopulator-crd.yaml
```

### Register custom populator by creating `VolumePopulator` resource

Make sure the proper CRD is available on the cluster that provides the
`VolumePopulator` CR.

Create an instance of `VolumePopulator` CR to register our populator
(`GenericHTTPPopulator`).

```bash!
kubectl apply -f https://raw.githubusercontent.com/viveksinghggits/generichttppopulator/master/manifests/volumepopulator.yaml
```

### Install `GenericHTTPPopulator` controller

Below command can be used to create a namespace `controller` and install the
populator controller in that namespace. This command will also setup the RBAC
required to run the contoller.

```bash!
kubectl apply -f https://github.com/viveksinghggits/generichttppopulator/raw/master/manifests/controller-deploy.yaml
```

Make sure the controller is running by checking the status of the pod in
`controller` namespace.

### Install data source validator

Data source validator can also optionally be installed to know if the data
source specified in the PVC is valid populator or not.

[This link](https://kubernetes.io/blog/2022/05/16/volume-populators-beta/#trying-it-out)
can be followed to install the Data source validator.

## Usage

To create a PVC in default namespace and populate it with the content of
[this](https://raw.githubusercontent.com/viveksinghggits/akcess/master/README.md)
README file, create a resource of type `GenericHTTPPopulator` using below
command

```bash!
kubectl apply -f https://raw.githubusercontent.com/viveksinghggits/generichttppopulator/master/manifests/akcessreadmepopulator.yaml
```

Once `GenericHTTPPopulator` named `akcess-readme-pop` is created, it can be
referred in a PVC resource like below

```yaml!
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: myclaim
spec:
  accessModes:
    - ReadWriteOnce
  volumeMode: Filesystem
  resources:
    requests:
      storage: 8Gi
  dataSourceRef:
    apiGroup: k8s.viveksingh.dev
    kind: GenericHTTPPopulator
    name: akcess-readme-pop
```

Create the PVC resource using the command below

```bash!
kubectl apply -f https://github.com/viveksinghggits/generichttppopulator/raw/master/manifests/pvc.yaml
```

PVC would be initially in `Pending` state and the controller that we have
deployed above in `controller` namespace would create a populator pod that
would actually populate the data into the PVC.

Check the status of that pod to make sure that population has succeeded and
then check the status of the PVC to make sure that it's in `Bound` state.

## Check the populated data

After the PVC `myclaim` is in `Bound` state, create a test pod that would mount
this volume at location `/mnt/check`. So that, the data can be verified by
`exec`ing into the pod.

```bash!
kubectl apply -f https://github.com/viveksinghggits/generichttppopulator/raw/master/manifests/checker.yaml
```

Once the pod is running `exec` into it, and then `cd` to `/mnt/check` and make
sure there is a file called `README.md` with correct content.

## YouTube Video

First version of this project was written as a tutorial and the recordings are
available at [this link](https://www.youtube.com/playlist?list=PLh4KH3LtJvRSyRi5uDb43n8Nc5vKoQNss).