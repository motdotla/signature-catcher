# signature-catcher

<img src="https://raw.githubusercontent.com/motdotla/signature-catcher/master/signature-catcher.jpg" alt="signature-catcher" align="right" width="200" />

Catches the webhook with the converted document arriving from [signature-api](https://github.com/motdotla/signature-api) and updates the database. Works in tandem with [signature-api](https://github.com/motdotla/signature-api).

## Installation

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)

### Development

```
git clone https://github.com/motdotla/signature-catcher.git
cd signature-catcher
go get 
cp .env.example .env
go run app.go
```

Edit the contents of `.env`.
