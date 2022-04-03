# roboCAPTCHA - Call Screening Service
_Subjecting unknown callers to human verification since 2022._

Our Call Screening service is responsible for weeding out robocalls by posing a human verification test to all callers. The human verification test will consist of either a number to be entered, or a word to be spoken by the caller.

Successful callers are routed to their respective roboCAPTCHA users. The service will communicate with the MongoDB database to obtain the required information about the user, it's whitelist/blacklist status, and send real-time notifications via the AWS Simple Queue Service.

All unsuccessful calls are recorded and stored in the database, allowing the user to view all their blocked calls. This allows the users to never miss an important call which may be inadvertently blocked by roboCAPTCHA.

### Environment Variables
| Name                     | Description                                                     |
| ------------------------ | --------------------------------------------------------------- |
| DB_CONN_STRING           | Connection string to establish an connection to MongoDB         |
| DB_NAME                  | roboCAPTCHA MongoDB Database name                               |
| AWS_ID                   | AWS Access Key ID                                               |
| AWS_SECRET               | AWS Secret                                                      |
| AWS_SQS_NOTIFICATION_URL | AWS SQS Notification URL                                        |
| ANON_FORWARD_TO          | For debugging purposes - anonymous call forwarding              |
| DEFAULT_NUMBER_TO        | For debugging purposes - default route                          |
| DEFAULT_NUMBER_FROM      | For debugging purposes - default route                          |


### Cloning the Repository
```bash
git clone https://github.com/phomies/robocaptcha-service.git
```

### Local Deployment
Requires installation of Go and project dependencies
```
go mod download
go run ./internal
```

### Local Deployment with Docker Compose
``` bash
docker-compose up -d --build
```
