./cloud_sql_proxy --private-ip cdms-463617:us-central1:cdms-db & 
sleep 5
goose -dir "backend/sql/schema" postgres "postgres://cdms_user1:Bas0123@/cdms?host=/cloudsql/cdms-463617:us-central1:cdms-db" up
