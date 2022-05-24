IMG_NAME	=	shortener
IMG_VERSION	=	latest
IMG			=	$(IMG_NAME):$(IMG_VERSION)
LAST2IMG	=	$(shell docker images | awk 'FNR>1 && FNR<4 {print $$3}')

build:
			docker build -t $(IMG) .

run_cache:	build
			docker run --rm --detach --name $(IMG_NAME) -p "8080:8080" $(IMG) --storage cache

run_pg:
			docker-compose up --build

stop_cache:
			docker stop $(IMG_NAME)
			docker rmi $(LAST2IMG)


.PHONY:		build run_cache run_pg stop