help:
	@echo "Invalid args!"
	@echo "make <args>"
	@echo ""
	@echo "u             |    Build with upx "
	@echo "d             |    Standart build"
	@echo ""
	@echo "Examples:"
	@echo "make u"
	@echo "make d"


launch_docker:
	nim c -d:quiet --hints:off docker/docker.nim
	./docker/docker cbl
	@rm docker/docker
ubuntu:
	docker build -f tests/ubuntu.dockerfile -t ubuntu .
	docker run -it ubuntu

arch:
	docker build -f tests/arch.dockerfile -t arch .
	docker run -it arch
	
fedora:
	docker build -f tests/fedora.dockerfile -t fedora .
	docker run -it fedora
	# sudo yum install snapd

docker_clear:
	docker rmi -f $(docker images -aq)	

payload: 
	nim c -d:quiet --hints:off scripts/payload.nim 
	./scripts/payload 
	@rm scripts/payload

e:
	cd enc; go run .
clear: 
	rm *.bin
u:
	bash scripts/build.sh upx
d:
	bash scripts/build.sh standart

