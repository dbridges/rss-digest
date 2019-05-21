# RSS Digest

Create a daily RSS Digest and email it to yourself.

[Blog post](https://medium.com/parallel-thinking/building-a-personal-rss-email-digest-service-in-go-7c8b71ac5b89)

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
