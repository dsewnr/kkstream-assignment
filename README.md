# KKS Assignment

### Environments

Generate a GCP authentication credentials file with Google Cloud Storage permissions.

##### Shell
```
$ export GOOGLE_APPLICATION_CREDENTIALS={GCP Authentication Credentials File.json}

```

##### .env
```
SESSION_SECRET={Session secret}
SESSION_ID={Session ID key}
UPLOAD_DIR={Upload folder}
LINE_AUTH_URL=https://access.line.me/oauth2/v2.1/authorize
LINE_TOKEN_URL=https://api.line.me/oauth2/v2.1/token
LINE_PROFILE_URL=https://api.line.me/v2/profile
LINE_CHANNEL_ID={Line channel ID}
LINE_CHANNEL_SECRET={Line channel secret}
LINE_CALLBACK={Line login callback}
GOOGLE_CLOUD_STORAGE_BUCKET={GCS bucket}
```
For this repo run `.env` file is ready.


##### Initial
```
$ docker-compose up -d
$ buffalo db create -a
$ buffalo db migrate
$ make fake_users
```

### Test APIs
```
make test
```

### Run
```
$ buffalo dev
```

### Default Role
```
Username: admin
Password: admin

Username: user
Password: user
```
