# sil-backend-assessment

This rest api allow for the creation of uploading of product and the creation of orders. Order creation triggers an email to the admin in this case it is me with the order information. the application has been instrumented using open telemetry.

to run the application make sure to have the correct .env file for envronment variables and run `make start`

you should run the docker compose file to start a database instance before running the up, make sure to create the db.

On startup the application will make all the migrations