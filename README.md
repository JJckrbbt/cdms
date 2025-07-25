## Chargeback and Delinquency Management System (CDMS)

### A new approach to managing the collection of Chargebacks and Non-IPAC collections

#### Demo available @ [cdms.jjckrbbt.dev](https://cdms.jjckrbbt.dev)

CDMS is my attempt at delivering a significantly improved experience to my colleagues managing Chargebacks and Delinquencies. Currently this work is managed on Google Sheets with the help of an RPA.  While this was a big improvement over the former process, an individual manually performing it in Excel over several hours, it still takes over an hour to complete each update.  In addition, in sheets users are making their comments and notes in cells, managing permissions and protecting sheets is laborsome, and users are expected to download and work offline if they need to complete any action while the Google Sheet is down for maintanence. 

This is where CDMS comes in.  CDMS provides easy and fast updates, with no manual manipulation of the source csv files.  Context won't be lost from moving reconciled items off the sheet.  Reconciled items will be marked inactive, but still available for historical reporting.  Better attribution of user actions.  Comments easily attributed to to those who made them.  A detailed and complete audit history for every item.  

CDMS is built on Postgres, Go, and React/Typescript.  

This project is still under development, the current Roadmap includes:

- User Authentication & Authorization (RBAC) using OIDC
- Upload & Report Merge Reporting
- Comments & User Tags in Comments
- Alerts

Prerequisites:
- sqlc
- goose
- postgresql database

To run:
1. Clone the repo into a folder.
2. Start a Postgres database, run and connect.  Use goose to run migrations from `cdms/backend/sql/schema`
3. From the 'backend' directory `go run cmd/server/main.go`
4. From the 'frontend' directory `npm install` && `npm run dev`
5. Some adjustments may need to be made in `cmd/server/main.go` for CORS to allow your local host.



