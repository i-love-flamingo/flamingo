DOCKERREPO?=docker-om3-akl.aoe.com
TAG?=latest

.PHONY: run doc docs dep godoc docker docker-run docker-push

run:
	cd akl && go run akl.go serve

dep:
	glide i

doc: docs

docs:
	docker run --rm -v $(shell pwd):/work -p 8000:8000 -ti thebod/mkdocs bash -c 'cd /work/docs; mkdocs serve --dev-addr=0.0.0.0:8000'

godoc:
	(sleep 10 ; open http://localhost:6060/pkg/flamingo/) &
	godoc -http=:6060 -v

docker: Dockerfile
	docker build -t $(DOCKERREPO)/flamingo:$(TAG) .

docker-run: docker
	docker run -ti -p 3210:3210 -v $(shell pwd)/akl/frontend:/go/src/flamingo/akl/frontend $(DOCKERREPO)/flamingo

docker-push: docker
	docker push $(DOCKERREPO)/flamingo
