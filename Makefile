# поднимает локальную кафку и постгресы
start-env:
	docker-compose -p async_arch up -d

# останавливает локальную кафку и постгресы
stop-env:
	docker-compose -p async_arch down
