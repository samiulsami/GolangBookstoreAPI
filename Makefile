#Push changes to remote repository

commitMessage ?= "from makefile"
gitPush:
	git add .
	git commit -m $(commitMessage)
	git push
#Compile and run the binary
binaryName ?= GoBookstoreAPI
$(binaryName):
	go build -o ./$(binaryName)
build: $(binaryName)

port ?= 3000
run: build
	./$(binaryName) serve --port=$(port)

clean::
	-rm GoBookstoreAPI

#Run with helm
helmDeploymentName ?= bookstorehelm
helmInstall:
	helm install $(helmDeploymentName) deploy/helm
	-kubectl get nodes -o wide

clean::
	-helm uninstall $(helmDeploymentName)

#Run with docker
imageName ?= sami7786/gobookstoreapi
tag ?= latest
dockerImage:
ifeq ($(shell docker images $(imageName):$(tag) -a -q),)
	docker build -t $(imageName):$(tag) .
else
	@echo "image already exists"
endif

dockerRun: dockerImage
	docker run -d -p $(port):$(port) --name=$(binaryName) $(imageName):$(tag)
	@docker logs $(binaryName)

DOCKER_USER ?= sami7786
dockerLogin:
	@docker login -u $(DOCKER_USER) -p $(DOCKER_PASS)
dockerPush: dockerImage dockerLogin
	docker push $(imageName):$(tag)

clean::
	-docker kill $(binaryName)
	-docker rm $(binaryName)

cleanAll:: clean
	docker rmi -f $$(docker images $(imageName):$(tag) -a -q)