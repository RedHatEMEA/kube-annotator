GOPATH=$(PWD)

build:
	go get golang.org/x/tools/go/types
	go build kube-annotator

origin:
	git clone https://github.com/openshift/origin.git
	pushd origin; make; popd

out: build origin
	for i in origin/pkg/*/api/v1beta1; do GOPATH=origin/Godeps/_workspace:origin/_output/local/go ./kube-annotator github.com/openshift/$$i; done >out/openshift-v1beta1.txt
	for i in origin/pkg/*/api/v1beta3; do GOPATH=origin/Godeps/_workspace:origin/_output/local/go ./kube-annotator github.com/openshift/$$i; done >out/openshift-v1beta3.txt
	GOPATH=origin/Godeps/_workspace:origin/_output/local/go ./kube-annotator github.com/GoogleCloudPlatform/kubernetes/pkg/api/v1beta1 >out/kubernetes-v1beta1.txt
	GOPATH=origin/Godeps/_workspace:origin/_output/local/go ./kube-annotator github.com/GoogleCloudPlatform/kubernetes/pkg/api/v1beta2 >out/kubernetes-v1beta2.txt
	GOPATH=origin/Godeps/_workspace:origin/_output/local/go ./kube-annotator github.com/GoogleCloudPlatform/kubernetes/pkg/api/v1beta3 >out/kubernetes-v1beta3.txt

.PHONY: out
