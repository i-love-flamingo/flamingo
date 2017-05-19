DOCKERREPO?=docker-om3.aoe.com

.PHONY: run doc dep godoc docker docker-run docker-push

run:
	cd akl && go run akl.go serve

dep:
	glide i

doc:
	docker run --rm -v $(pwd):/work -p 8000:8000 python bash -c 'pip install mkdocs; cd /work; mkdocs serve --dev-addr=0.0.0.0:8000'

godoc:
	(sleep 10 ; open http://localhost:6060/pkg/flamingo/) &
	godoc -http=:6060 -v

docker: Dockerfile
	docker build -t $(DOCKERREPO)/flamingo/akl .

docker-run: docker
	docker run -ti -p 3210:3210 -v $(pwd)/akl/frontend:/go/src/flamingo/akl/frontend $(DOCKERREPO)/flamingo/akl

docker-push: docker
	docker push $(DOCKERREPO)/flamingo/akl
