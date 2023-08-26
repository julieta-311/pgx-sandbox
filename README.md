# pgx-sandbox

Playing around with pgx v5 and squirrel, nothing too interesting yet as for now it's
like a script inserting a thing into the db.

Requirements:

* Install docker if not yet installed. For example in Arch distros do

```
sudo pacman -Sy docker
systemctl enable docker
sudo usermod -aG docker $USER
sudo reboot
```

How to run:

* Start Postgres instance for testing:

```
./start_db.sh
```

Note that running `docker stop sandbox` or `./stop_and_rm_db.sh` will delete de instance
along with any test data that was added.

* Run it:

```
POSTGRES_URL=postgres://postgres@localhost:5432/pgxsandbox go run .
```
