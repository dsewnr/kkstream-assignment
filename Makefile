test:
	buffalo test actions -m Test_HomeHandler -count=1 \
	&& buffalo test actions -m Test_Api_Upload -count=1

fake_users:
	export PGPASSWORD=postgres \
	&& psql -h localhost -d assignment_development -U postgres -p 5432 -a -w -f migrations/fake_users.sql
