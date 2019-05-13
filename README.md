# RSS Digest

Create a daily RSS Digest and email it to yourself.

Blog post:

## Usage

Setup a `.env` file with your email settings:

```
MAIL_FROM=sender@gmail.com
MAIL_TO=recipient@gmail.com
MAIL_PASSWORD=<sender password goes here>
MAIL_HOST=smtp.gmail.com
MAIL_PORT=587
```

Then run it with:
```
eval $(cat .env) go run *.go
```
