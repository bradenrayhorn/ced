# Quickstart

This guide will walkthrough getting started with ced using Docker Compose.

## Step 1 - Navigate to project directory

Create a project directory to store ced files. Navigate to this directory
within your terminal.

## Step 2 - Download files

Within your project directory, download the Docker Compose file and example env file.

```
curl -O https://raw.githubusercontent.com/bradenrayhorn/ced/main/docker/docker-compose.yml
curl -O https://raw.githubusercontent.com/bradenrayhorn/ced/main/docker/.env.example
```

## Step 3 - Configure environment

Rename the `.env.example` to `.env` and adjust configuration as needed.

```
mv .env.example .env
```

The directory specified at `APP_DATA_LOCATION` must be created.
This directory will contain the SQLite database.

```
mkdir data/
```

## Step 4 - Start ced

Start ced by starting the docker compose project.

For example, this could be done with the following command.

```
docker-compose up -d
```

If you are using the example env file, ced will be available at http://localhost:8080.

## Step 5 - Import or create groups

See [usage in README](README.md#usage) for additional info.

Example creating a group with CLI:

```
docker compose exec ced-server /app/ced-cli group create --name="Pamela Palm" --max-attendees=1 --search-hints="Pamela Palm, Pam Palm"
```

Example importing a group csv with CLI:

```
echo "Pamela Palm,1," > mydata.csv

docker compose cp mydata.csv ced-server:/mydata.csv
docker compose exec ced-server sh -c "cat /mydata.csv | /app/ced-cli group import"
```
