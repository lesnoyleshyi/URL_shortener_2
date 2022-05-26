IMG_NAME	=	url_shortener
IMG_VERSION	=	latest
IMG			=	$(IMG_NAME):$(IMG_VERSION)
LAST2IMG	=	$(shell docker images | awk 'FNR>1 && FNR<4 {print $$3}')

build:
			docker build -t lesnoyleshyi/$(IMG) .

push:
			docker push lesnoyleshyi/url_shortener:latest

run_cache:	build
			docker run --rm --detach --name $(IMG_NAME) -p "8080:8080" $(IMG) --storage cache

stop_cache:
			docker stop $(IMG_NAME)

rm_cache:
			docker rmi $(LAST2IMG)

run_pg:
			docker-compose up --build

rm_pg:
			docker rm api_shortener
			docker rmi $(LAST2IMG)

rmrf_pg:	rm_pg
			docker rm postgres_shortener
			docker volume rm shortener2_volume

post:
			curl -L -X POST 'localhost:8080' -H 'Content-Type: application/json' --data-raw \
  			'{"url": "https://github.com/lesnoyleshyi/URL_shortener_2/"}'

get:
			curl -L -X GET 'localhost:8080' -H 'Content-Type: application/json' --data-raw \
  			'{"url": "GET_IT_FROM_POST_ANSWER"}'


.PHONY:		build run_cache stop_cache rm_cache run_pg rm_pg rmrf_pg post get