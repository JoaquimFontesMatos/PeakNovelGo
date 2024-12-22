# PeakNovelGo

## Description

PeakNovelGo is a backend and frontend application built using Go and the Gin web framework. It is designed to provide a platform for users to create and manage their own novels, with features such as chapter creation, reading, and sharing. The application is built using the MVC (Model-View-Controller) architectural pattern and uses the GORM ORM for database interactions.

## Features

- User registration and authentication
- User profile management
- Novel creation and management
- Chapter creation and management
- Novel reading and sharing
- Text-to-speech (TTS) functionality
- Social features (likes, comments, and followers)
- User settings and preferences
- Gamification and rewards
- Supported platforms (web, desktop, and mobile)
- User feedback and support

## Requirements

The application requires a PostgreSQL database to store user data.

## Installation

To install the application, follow these steps:

1. Clone the repository:

    ```bash
    git clone https://github.com/joaquim/PeakNovelGo.git
    ```

2. Change to the project directory:

    ```bash
    cd PeakNovelGo
    ```

3. Create a `.env` file in the root directory and add the following environment variables:

    ```bash
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=postgres
    DB_PASSWORD=your_password
    DB_NAME=PeakNovel

    SMTP_USERNAME=your_email@example.com
    SMTP_PASSWORD=your_password
    SMTP_HOST=smtp.example.com
    SMTP_PORT=587
    ```

4. Run the application:

    ```bash
    go run cmd/server/main.go
    ```

5. Access the application at `http://localhost:8080`.

6. Follow the instructions to register an account and start using the application.

## Contributing

Contributions are welcome! If you find a bug or have a suggestion, please open an issue or submit a pull request.

## Technologies

- Go
- Gin
- GORM
- PostgreSQL (for production)
- SQLite (for testing)
- Vue.js
