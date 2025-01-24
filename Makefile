compose-local:
	@docker compose --profile dev up --build

compose-local-d:
	@docker compose --profile dev up --build -d

compose-up:
	@docker compose --profile full up --build

compose-up-d:
	@docker compose --profile full up -d

generate:
	@make -C ./backend generate
	@make -C ./frontend generate
