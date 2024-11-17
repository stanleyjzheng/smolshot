# SmolShot
a smol MoonShot

## about
- `api/` contains a simple go api
- `app` is the flutter app

## instructions
fill out `api/.env.default`

create the users table (TODO: move to migrations)
```sql
CREATE TABLE IF NOT EXISTS accounts (
    user_id TEXT PRIMARY KEY,
    private_key TEXT NOT NULL,
    public_key TEXT NOT NULL
);
```

for the flutter app,
```sh
cd app
flutter pub get
# open your simulator of choice
flutter run
```

for the go api,
```sh
go run main.go
```
