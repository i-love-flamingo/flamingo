DOCKERREPO?=docker-om3.aoe.com

.PHONY: run doc dep dep-core dep-akl docker docker-run

run:
	cd akl && go run akl.go || true

dep: dep-core dep-akl

dep-core:
	cd core && glide i

dep-akl:
	cd akl && glide i

doc:
	godoc -http=:6060 -v
	open http://localhost:6060/pkg/flamingo/

docker: Dockerfile
	docker build -t $(DOCKERREPO)/flamingo/akl .

docker-run: docker
	docker run -ti -p 3210:3210 -v $(pwd)/akl/frontend:/go/src/flamingo/akl/frontend $(DOCKERREPO)/flamingo/akl
