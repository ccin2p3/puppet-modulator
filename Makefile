#
UID=$(shell id -u)
GID=$(shell id -g)

# Naive but enough
PWD=$(shell pwd)/docs

MKDOCS_MATERIAL_DOCKER_IMAGE="squidfunk/mkdocs-material:7.3.0"

serve-doc:
	docker run --rm \
	  --user $(UID):$(GID) \
	  -v $(PWD):/src \
	  -w '/src' \
	  --network host \
	  $(MKDOCS_MATERIAL_DOCKER_IMAGE) "serve"
