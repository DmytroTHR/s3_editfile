To run the app, launch
```
./bin/app
```

`-cred` : path to aws credentials & config,   "`./aws/`" - default, if not provided tries to use default EC2 IAM roles

`-b` : aws s3 bucket name,   "`somebucket-123321`" - default

`-f` : object to be changed name from the bucket, "`index.txt`" - default

`-a` : string to append to the file,   "`fmt.Sprintf("%s - appended string", time.Now().Format(time.Stamp))`" - default