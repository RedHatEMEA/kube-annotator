kube-annotator
==============

```
$ git clone https://github.com/openshift/origin.git
$ pushd origin; make; popd

$ make

$ export GOPATH=origin/Godeps/_workspace:origin/_output/local/go
$ ./kube-annotator github.com/GoogleCloudPlatform/kubernetes/pkg/api/v1beta2
$ for i in origin/pkg/*/api/v1beta1; do ./kube-annotator github.com/openshift/$i; done
```
