GOPATH=${HOME}/go

%:
	@true

.PHONY: fmt
fmt:
	./scripts/fmt.sh $(filter-out $@,$(MAKECMDGOALS))

.PHONY: run
run:
	./scripts/fmt.sh $(filter-out $@,$(MAKECMDGOALS))
	./scripts/run.sh $(filter-out $@,$(MAKECMDGOALS))

.PHONY: proto
proto:
	./scripts/proto.sh $(filter-out $@,$(MAKECMDGOALS))

.PHONY: deploy
deploy:
	./scripts/deploy.sh $(filter-out $@,$(MAKECMDGOALS))
.PHONY: doc

.PHONY: build
build:
	./scripts/build.sh $(filter-out $@,$(MAKECMDGOALS))
