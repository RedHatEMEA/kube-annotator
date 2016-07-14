GOPATH=$(PWD)

build:
	go build kube-annotator

origin:
	git clone https://github.com/openshift/origin.git
	pushd origin; make; popd

alpaca: build origin
	for i in origin/pkg/*/api/v1; do GOPATH=origin/Godeps/_workspace:origin/_output/local/go ./kube-annotator alpaca github.com/openshift/$$i; done
	GOPATH=origin/Godeps/_workspace:origin/_output/local/go ./kube-annotator alpaca k8s.io/kubernetes/pkg/api/v1
	GOPATH=origin/Godeps/_workspace:origin/_output/local/go ./kube-annotator alpaca k8s.io/kubernetes/pkg/apis/extensions/v1beta1

out: build origin
	for i in origin/pkg/*/api/v1; do GOPATH=origin/Godeps/_workspace:origin/_output/local/go ./kube-annotator doc github.com/openshift/$$i; done >out/openshift-v1.txt
	GOPATH=origin/Godeps/_workspace:origin/_output/local/go ./kube-annotator doc k8s.io/kubernetes/pkg/api/v1 >out/kubernetes-v1.txt
	GOPATH=origin/Godeps/_workspace:origin/_output/local/go ./kube-annotator doc k8s.io/kubernetes/pkg/apis/extensions/v1beta1 >out/kubernetes-ext-v1beta1.txt

pyout:
	./kube-annotator.py origin/api/swagger-spec/api-v1.json >out/kubernetes-v1-py.txt
	./kube-annotator.py origin/api/swagger-spec/oapi-v1.json >out/openshift-v1-py.txt


.PHONY: alpaca out pyout
