.PHONY: kafka-run
kafka-run:
	docker compose -f docker-compose.yml up -d

.PHONY: kafka-rerun
kafka-return: kafka-delete kafka-run

.POHNY: kafka-delete
kafka-delete:
	docker stop kafka-test zookeeper-test
	docker rm kafka-test zookeeper-test
	docker volume rm kafka_zookeeper_vol kafka_kafka_vol

.PHONY: etcd-start
etcd-start:
	docker compose -f docker-compose-etcd.yml -d

.PHONY: influxdb-start
influxdb-start:
	docker run -d -p 8083:8083 -p 8086:8086 --expose 8090 --expose 8099 --name infulxDbService influxdb

.PHONY: grafana-start
grafana-start:
	docker run -d -p 3000:3000 --name=grafana -v /Users/pipi/data/grafana:/var/lib/grafana grafana/grafana
