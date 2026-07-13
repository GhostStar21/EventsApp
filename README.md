# EventsApp
This is a web application where events can be registered, and eventually registrations can be done, such that users can have all events on one website, instead of events being spread out through different social media pages (e.g. IG accounts, Discord, etc.)

## Prerequisites
- React v18.3.1
- Golang 1.26
- PostgresSQL Database of your choice (ensure that it is running)
- npm v11.16.0
- Install dependencies via npm: ``` npm run dev ```
- Your env file which should contain the `SECRET_KEY` and `DATABASE_URL`.

## How to run
To run the web application, you have to do the following:
1. Go to the `/cmd` directory in the `backend` folder and run:
```bash
go run main.go
``` 
If no port is provided, the backend will listen on port 8080.
2. In a new terminal, go to the `frontend` folder and run
```bash
npm run dev
``` 
This will run at the port 5173.
If no domain or IP address is provided, then the local URL is shown by vite i.e., http://localhost:5173 

## Features
- Sign in or register an account
- View all events returned by the backend

