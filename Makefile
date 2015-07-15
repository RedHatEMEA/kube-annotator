GOPATH=$(PWD)

build:
	go get golang.org/x/tools/go/types
	go build kube-annotator

origin:
	git clone https://github.com/openshift/origin.git
	pushd origin; make; popd

out: build origin
	for i in origin/pkg/*/api/v1; do GOPATH=origin/Godeps/_workspace:origin/_output/local/go ./kube-annotator github.com/openshift/$$i; done >out/openshift-v1.txt
	GOPATH=origin/Godeps/_workspace:origin/_output/local/go ./kube-annotator github.com/GoogleCloudPlatform/kubernetes/pkg/api/v1 >out/kubernetes-v1.txt

.PHONY: out
