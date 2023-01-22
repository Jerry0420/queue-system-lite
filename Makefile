up_db:
	docker-compose -f ./compose/docker-compose.db.yml up -d --build

exec_db:
	docker exec -it migration_tools sh

down_db:
	docker-compose -f ./compose/docker-compose.db.yml down

# ========================================================

up_test:
	docker-compose -f ./compose/docker-compose.test.yml up -d --build

exec_test:
	docker exec -it backend_test sh

down_test:
	docker-compose -f ./compose/docker-compose.test.yml down

# ========================================================

up_build_dev:
	docker-compose -f ./compose/docker-compose.dev.yml up -d --build

re_build_nginx:
	docker-compose -f ./compose/docker-compose.dev.yml up -d --force-recreate --build nginx

down_dev:
	docker-compose -f ./compose/docker-compose.dev.yml down

exec_backend:
	docker exec -it backend sh

# =========================================================

up_build:
	docker-compose -f ./compose/docker-compose.yml up -d --build

down:
	docker-compose -f ./compose/docker-compose.yml down

# ==========================================================

add_crontab:
	(crontab -l 2>/dev/null; echo "* * * * * curl --max-time 30 --connect-timeout 5 -X DELETE --url $(server)/api/v1/routine/stores") | crontab -

# ==========================================================

renew_certbot:
	docker-compose -f ./compose/docker-compose.cert.yml run --rm certbot renew
	docker restart nginx

# ==========================================================
