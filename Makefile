.PHONY: bash
bash:
	docker-compose run --rm --service-ports app bash

.PHONY: fmt
fmt:
	docker-compose run --rm app go fmt ./...

.PHONY: db
db:
	docker-compose up -d db

.PHONY: pgweb
pgweb:
	docker-compose up -d pgweb


.PHONY: reset-tables
reset-tables: ## DBをtruncate
	docker-compose exec db sh -c "psql -hlocalhost -Uhoge-user -dhoge-db -c \"TRUNCATE users, microposts RESTART IDENTITY;\""

.PHONY: sample-data-for-db
sample-data-for-db: reset-tables ## DBにサンプルデータを入れる
	docker-compose exec db sh -c "psql -hlocalhost -Uhoge-user -dhoge-db -c \"INSERT INTO users (name, email) VALUES ('名無しの権兵衛01', 'test01@example.com') RETURNING id,name,email,created_at,updated_at;\""
	docker-compose exec db sh -c "psql -hlocalhost -Uhoge-user -dhoge-db -c \"INSERT INTO users (name, email) VALUES ('名無しの権兵衛02', 'test02@example.com') RETURNING id,name,email,created_at,updated_at;\""
	docker-compose exec db sh -c "psql -hlocalhost -Uhoge-user -dhoge-db -c \"INSERT INTO users (name, email) VALUES ('名無しの権兵衛03', 'test03@example.com') RETURNING id,name,email,created_at,updated_at;\""
	docker-compose exec db sh -c "psql -hlocalhost -Uhoge-user -dhoge-db -c \"SELECT * FROM users;\""
	sleep 5
	docker-compose exec db sh -c "psql -hlocalhost -Uhoge-user -dhoge-db -c \"UPDATE users SET name = '福沢諭吉' WHERE id = 1 RETURNING id,name,email,created_at,updated_at;\""
	docker-compose exec db sh -c "psql -hlocalhost -Uhoge-user -dhoge-db -c \"SELECT * FROM users;\""

.PHONY: down
down: ## docker containerを落とす
	docker-compose down

.PHONY: clean
clean: ## docker containerを落とす(volumeも)
	docker-compose down --volumes

################################################################################
# Utility-Command help
################################################################################
.DEFAULT_GOAL := help

################################################################################
# マクロ
################################################################################
# Makefileの中身を抽出してhelpとして1行で出す
# $(1): Makefile名
define help
  grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(1) \
  | grep --invert-match "## non-help" \
  | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
endef
################################################################################
# タスク
################################################################################
.PHONY: help
help: ## Make タスク一覧
	@echo '######################################################################'
	@echo '# Makeタスク一覧'
	@echo '# $$ make XXX'
	@echo '# or'
	@echo '# $$ make XXX --dry-run'
	@echo '######################################################################'
	@echo $(MAKEFILE_LIST) \
	| tr ' ' '\n' \
	| xargs -I {included-makefile} $(call help,{included-makefile})
