# Build and start the database container
run_db:
	docker-compose up -d db

# Run database migrations
migrate:
	docker-compose exec app go run cmd/migration/main.go up

# Build and start the application along with all services
run:
	# Start app container in detached mode (in background)
	docker-compose up -d app
	# Wait for the DB to be ready
	sleep 5
	# Run migrations
	make migrate
	# Now bring up all services
	docker-compose up --build
