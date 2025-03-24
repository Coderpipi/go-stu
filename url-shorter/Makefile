install_mysql:
	docker run -p 3306:3306 --name mysql \
	-v /Users/pipi/data/docker/mysql/mysql-files:/var/lib/mysql-files \
	-v /Users/pipi/data/docker/mysql/conf:/etc/mysql \
	-v /Users/pipi/data/docker/mysql/logs:/var/log/mysql \
	-v /Users/pipi/data/docker/mysql/data:/var/lib/mysql \
	-e MYSQL_ROOT_PASSWORD=root \
	-d mysql:8.0.23

install_redis:
	docker run --name redis \
    -p 6379:6379 \
    -v /Users/pipi/data/docker/redis/redis.conf:/etc/redis/redis.conf \
    -v /Users/pipi/data/docker/redis:/data \
    -d redis redis-server /etc/redis/redis.conf --appendonly yes

install_postgresql:
	docker run -id --name=postgresql \
	-v /Users/pipi/data/docker/postgres/data:/var/lib/postgresql/data \
 	-p 5432:5432 -e POSTGRES_PASSWORD=123456 \
 	-e LANG=C.UTF-8 postgres
