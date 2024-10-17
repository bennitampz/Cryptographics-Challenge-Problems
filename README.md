# Cryptographics-Challenge-Problems

This repository contains a User Registration and Authentication System implemented in Go. It allows users to register with a username and password and authenticate themselves later. The application uses SQLite as the database and Argon2 for password hashing, providing a secure and straightforward web interface.

## Table of Contents
- [Overview](#overview)
- [Requirements](#requirements)
- [Setup](#setup)
- [Usage](#usage)
  - [Registering a User](#registering-a-user)
  - [Authenticating a User](#authenticating-a-user)
- [Verifying Data in SQLite](#verifying-data-in-sqlite)
- [Code Structure](#code-structure)
- [Endpoints](#endpoints)

## Overview

This application allows users to register and authenticate using a web interface. It ensures secure handling of user credentials and provides feedback to users during the registration and authentication processes.

## Requirements

- **Go**: Version 1.16 or higher
- **SQLite3**: For database management

## Setup

1. **Install Go**: Ensure Go is installed on your machine. Download it from [golang.org](https://golang.org/dl/).

2. **Install SQLite Driver**: Open your terminal and run:
   ```bash
   go get github.com/mattn/go-sqlite3

3. **Clone the Repository: Download the code files to your local machine.**
   ```bash
   git clone https://github.com/bennitampz/Cryptographics-Challenge-Problems
   
   cd Cryptographics-Challenge-Problems

## Usage

1. Open your terminal.
2. Navigate to the directory containing main.go.
3. Run the application using:
   ```bash
   go run main.go

4. You should see a message indicating the server has started at :8080

## Registering a User

1. Open your web browser and navigate to http://localhost:8080/register.
2. Fill in the registration form:

     Username: Enter a username (1-15 alphanumeric characters).

     Password: Enter a password (minimum 8 characters, must include letters, numbers, and special characters).

3. Click the "Daftar" button to submit the form.
   
4. You will receive feedback indicating whether the registration was successful or if there were errors.

## Authenticating a User

1. Navigate to http://localhost:8080/authenticate.
2. Fill in the authentication form:

    Username: Enter your registered username.

    Password: Enter your password.
   
3. Click the "Masuk" button to log in.
4. You will receive feedback indicating whether authentication was successful or if there were errors.

## Verifying Data in SQLite

To check if the user data has been successfully inserted into the SQLite database use Browser SQLite or :

1. Open SQLite: In your terminal, navigate to the directory containing users.db and start the SQLite shell:
   
   ```bash
 
    sqlite3 users.db

   SELECT * FROM users;

## Endpoints

a. GET /register: Serves the registration HTML page.

b. GET /authenticate: Serves the authentication HTML page.

c. POST /register_process: Processes user registration.

d. POST /authenticate_process: Processes user authentication.
