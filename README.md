# API Service 

## Tech Stack
- Go 1.24
- MySql 9.4
- Fiber (HTTP framework)
- GORM
- Viper (configuration)
- sqlmock, gomock (testing/mocking)
- k6 (stress test)

## Project Layout
- src/controller/ - request/application orchestration
- src/model/mysql/ - MySQL data repository (GORM)
- src/interface/http/v1/ - HTTP handlers (Fiber), response models
- src/global/ - shared types and constants (e.g., error type)
- src/migration/ - additional DB migrations after dumped provided mock data
- src/mocks/ - mock logger and model for unit test
- scripts/ - stores sql script to use on DB initialization and k6 stress test files

## Configuration
Configuration is read via Viper. A common setup is a YAML file stored in src/config consists of 2 files:

- config.yaml (used in local run from terminal)
- config_docker_compose.yaml (used in docker services)

## Running the Service (Docker Compose)
### Prerequisite
Please paste mock sql files in to directory `scripts/mysql` so it can be inserted into DB when executing docker compose.
Mock sql file can be found in https://drive.google.com/drive/folders/1Htg0KFHUgU8jrdGwGEdIbSs99z7I_JPC filename `mock.zip`. 
Please download and extract it to `scripts/mysql` directory.

### Running the Service
Build the project
```sh
docker compose build --no-cache
```
Start project (The first time run should take about 10 minutes since it will have to dump mock sql data into the DB)
```sh
docker compose up -d --wait
```
For shutting down 
```sh
docker compose down
```
The service will use volumes to store DB data if you want to erase DB volumes execute this command
```sh
docker compose down --volumes
```

## Running the Service (Local)
### Prerequisite
Running locally can also be done. You can start sql service from `docker_compose.yaml` file (by commenting out other services) or connect to your own local DB of the same schema.
DB connection config can be edited in `src/config/config.yaml` file
### Running the Service
First, go to the src directory
```sh
cd src
```
Before running the service, it will be necessary to execute DB migration first by executing 
```sh
go run main.go migrate --config=config/config.yaml 
```
After the migration is completed, you can then start the service
```sh
go run main.go serve --config=config/config.yaml 
```
The server then should be started and ready to use

## API Specs
This project consists of 6 total APIs to serve a given interface.

### Login
By providing user_id and pin (mocked as 123456 for all users) the API will give response with token to use on other APIs
#### Request
```sh
curl --location 'localhost:3000/api/v1/login' \
--header 'Content-Type: application/json' \
--header 'Authorization: ••••••' \
--data '{
    "user_id": "fffeb5b4e1a111ef95a30242ac180002",
    "pin": "123456"
}'
```
#### Response
```sh
{
    "code": 0,
    "message": "success",
    "data": {
        "greeting": "Hello User_fffeb5b4e1a111ef95a30242ac180002",
        "token": "c86ec84199c1045cc119acff0006306c73840365d02ef4922fe1013d85418894"
    }
}
```

### Get User Accounts
This API will return all accounts owned by user (user will be validated from bearer token).
#### Request
```sh
curl --location 'localhost:3000/api/v1/get-user-accounts' \
--header 'Authorization: ••••••'
```
#### Response
```sh
{
    "code": 0,
    "message": "success",
    "data": {
        "accounts": [
            {
                "account_id": "fffeb6d4e1a111ef95a30242ac180002",
                "type": "saving-account",
                "currency": "THB",
                "account_number": "568-2-90992",
                "issuer": "TestLab",
                "amount": 17730,
                "color": "#24c875",
                "is_main_account": true,
                "progress": 15,
                "flags": [
                    {
                        "flag_type": "system",
                        "flag_value": "Disbursement"
                    },
                    {
                        "flag_type": "system",
                        "flag_value": "Flag4"
                    }
                ]
            },
            {
                "account_id": "fffeba4de1a111ef95a30242ac180002",
                "type": "saving-account",
                "currency": "THB",
                "account_number": "568-2-71318",
                "issuer": "TestLab",
                "amount": 96707.92,
                "color": "#24c875",
                "is_main_account": false,
                "progress": 69,
                "flags": [
                    {
                        "flag_type": "system",
                        "flag_value": "Disbursement"
                    },
                    {
                        "flag_type": "system",
                        "flag_value": "Flag3"
                    }
                ]
            }
}
```

### Get User Debit Cards
This API will return all debit cards owned by user (user will be validated from bearer token).
#### Request
```sh
curl --location 'localhost:3000/api/v1/get-user-debit-cards' \
--header 'Authorization: ••••••'
```
#### Response
```sh
{
    "code": 0,
    "message": "success",
    "data": {
        "debit_cards": [
            {
                "card_id": "fffeb5d1e1a111ef95a30242ac180002",
                "name": "My Debit Card",
                "status": "Active",
                "number": "4772 **** **** 1428",
                "issuer": "TestLab",
                "color": "#00a1e2",
                "border_color": "#ffffff"
            }
        ]
    }
}
```

### Get User Banners
This API will return all banners under user_id (user will be validated from bearer token).
#### Request
```sh
curl --location 'localhost:3000/api/v1/get-user-banners' \
--header 'Authorization: ••••••'
```
#### Response
```sh
{
    "code": 0,
    "message": "success",
    "data": {
        "banners": [
            {
                "banner_id": "fffeb5d1e1a111ef95a30242ac180002",
                "title": "Want some money?",
                "description": "You can start applying",
                "image": "https://dummyimage.com/54x54/999/fff"
            }
        ]
    }
}
```

### Get User Saved Accounts
This API will return all saved accounts of user (user will be validated from bearer token).
#### Request
```sh
curl --location 'localhost:3000/api/v1/get-user-saved-accounts' \
--header 'Authorization: ••••••'
```
#### Response
```sh
{
    "code": 0,
    "message": "success",
    "data": {
        "saved_accounts": [
            {
                "name": "Dummy Name",
                "number": "1234567890",
                "image": "https://dummyimage.com/54x54/999/fff"
            }
        ]
    }
}
```

### Get User Info
This API will return user information. The purpose of this API is to use for getting user's name to display 
on entering pin page. So info given from this API will be only name and non-sensitive values, hence it does not require bearer token
since user hasn't logged in yet.
#### Request
```sh
curl --location 'localhost:3000/api/v1/get-user-by-id' \
--header 'Content-Type: application/json' \
--data '{
    "user_id": "fffeb5b4e1a111ef95a30242ac180002"
}'
```
#### Response
```sh
{
    "code": 0,
    "message": "success",
    "data": {
        "user_info": {
            "name": "User_fffeb5b4e1a111ef95a30242ac180002",
            "dummy_col1": "dummy_value_1"
        }
    }
}
```

## Stress Test (k6)
located in `scripts/k6` with test results