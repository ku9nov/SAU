[![Stand With Ukraine](https://raw.githubusercontent.com/vshymanskyy/StandWithUkraine/main/banner2-direct.svg)](https://stand-with-ukraine.pp.ua)

<div align="center">
  
  [![Documentation](https://img.shields.io/badge/Documentation-available-brightgreen)](https://ku9nov.github.io/faynoSync-site/docs/intro)
  ![Docker Pulls](https://img.shields.io/docker/pulls/ku9nov/faynosync)
  ![GitHub Release](https://img.shields.io/github/v/release/ku9nov/faynoSync)
  ![Docker Compose Test](https://github.com/ku9nov/faynoSync/actions/workflows/tests.yml/badge.svg)

</div>

# FaynoSync

![225881501-b8aab72a-31e7-45ec-9340-4cca2a7893e9 (1)](https://github.com/ku9nov/faynoSync/assets/69673517/59ee4531-5d6c-4bc3-8aab-96854e2a4844)

"FaynoSync" is derived from the Ukrainian word "файно" (fayno), which is transliterated as "fayno" in English. In the Ukrainian language, "файно" (fayno) is an informal term used to describe something as excellent, fine, or great, indicating a positive, satisfactory, or enjoyable state or experience.

This application is a simple auto updater service written in Golang. It allows you to upload your application to S3 and set the version number. The client application can then check the version number against the auto updater service API. If the auto updater service has a newer version, it will return a link to the updated service, and the client application will show an alert.

## Installation

To use this application, you will need to have Golang installed on your machine. You can install it from the official [website](https://golang.org/doc/install).

Once you have installed Golang, clone this repository to your local machine:

```
git clone https://github.com/ku9nov/faynoSync.git
```

## Configuration
To configure the `faynoSync`, you will need to set the following environment variables:
```
STORAGE_DRIVER (`minio` or `aws`)
S3_ACCESS_KEY (Your AWS or Minio access key ID.)
S3_SECRET_KEY (Your AWS or Minio secret access key.)
S3_REGION (The AWS region in which your S3 bucket is located. For Minio this value should be empty.)
S3_BUCKET_NAME (The name of your S3 bucket.)
S3_ENDPOINT= (s3 endpoint, check documentation of your cloud provider)
DASHBOARD_URL= (dashboard url to allow CORS configuration)
PORT (The port on which the auto updater service will listen. Default: 9000)
MONGODB_URL=mongodb://root:MheCk6sSKB1m4xKNw5I@127.0.0.1/cb_faynosync_db?authSource=admin (see docker-compose file)
SYSTEM_KEY= (generated by 'openssl rand -base64 16')
API_KEY= (generated by 'openssl rand -base64 16') Used for SignUp
```

You can set these environment variables in a `.env` file in the root directory of the application. You can use the `.env.local` file, which contains all filled variables.

### Docker configuration
To build and run the API with all dependencies, you can use the following command:
```
docker-compose up --build
```
You can now run tests using this command (please wait until the `s3-service` successfully creates the bucket):
```
docker exec -it faynoSync_backend "/usr/bin/faynoSync_tests"
```
If you only want to run dependency services (for local development without Docker), use this command:
```
docker-compose -f docker-compose.yaml -f docker-compose.development.yaml up
```
## Usage
To use the auto updater service, follow these steps:
1. Build the application:
```
go build -o faynoSync faynoSync.go
```

2. Start the auto updater service with migrations:
```
./faynoSync --migration
```
Note: To rollback your migrations run:
```
./faynoSync --migration --rollback
```

3. Upload your application to S3 and set the version number in Admin Api.

4. In your client application, make a POST request to the auto updater service API, passing the current version number as a query parameter:
```
http://localhost:9000/checkVersion?app_name=myapp&version=0.0.1
```

The auto updater service will return a JSON response with the following structure:

```
{
    "update_available": false,
    "update_url_deb": "https://<bucket_name>.s3.amazonaws.com/secondapp/myapp-0.0.1.deb",
    "update_url_rpm": "https://<bucket_name>.s3.amazonaws.com/secondapp/myapp-0.0.1.rpm"
}
```

If an update is available, the update_available field will be true, and the update_url field will contain a link to the updated application.

5. In your client application, show an alert to the user indicating that an update is available and provide a link to download the updated application.

## Testing
Run e2e tests:
```
go test
```
Build test binary file:
```
go test -c -o faynoSync_tests
```
**Test Descriptions**

To successfully run the tests and have them pass, you need to populate the .env file.

The tests verify the implemented API using a test database and an existing S3 bucket.

**List of Tests**

    - TestHealthCheck
    - TestLogin
    - TestFailedLogin (expected result from API "401")
    - TestListApps
    - TestListAppsWithInvalidToken (expected result from API "401")
    - TestAppCreate
    - TestSecondaryAppCreate (expected result from API "failed")
    - TestUploadApp
    - TestUploadDuplicateApp (expected result from API "failed")
    - TestDeleteApp
    - TestChannelCreateNightly
    - TestChannelCreateStable
    - TestUploadAppWithoutChannel (expected result from API "failed")
    - TestMultipleUploadWithChannels
    - TestSearchApp
    - TestCheckVersionLatestVersion
    - TestMultipleDelete
    - TestDeleteNightlyChannel
    - TestDeleteStableChannel
    - TestPlatformCreate
    - TestUploadAppWithoutPlatform
    - TestArchCreate
    - TestUploadAppWithoutArch
    - TestDeletePlatform
    - TestDeleteArch
    - TestListArchs
    - TestListPlatforms
    - TestListChannels
    - TestListArchsWhenExist
    - TestListPlatformsWhenExist
    - TestListChannelsWhenExist
    - TestSignUp
    - TestFailedSignUp (expected result from API "401")
    - TestUpdateSpecificApp
    - TestListAppsWhenExist
    - TestDeleteAppMeta
    - TestUpdateChannel
    - TestUpdateApp
    - TestUpdatePlatform
    - TestUpdateArch
    
## Create new migrations
Install migrate tool [here](https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md).
```
cd mongod/migrations
migrate create -ext json name_of_migration
```
Then run the migrations again.
## License
This application is licensed under the Apache license. See the LICENSE file for more details
