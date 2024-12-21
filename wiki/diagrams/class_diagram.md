# Class Diagram

```mermaid
classDiagram

UserController -- UserService: Uses

UserService -- UserRepository: Interacts with
UserService -- EmailManager: Uses
UserService -- ExpirationValidator: Uses
UserService -- TokenGenerator: Uses

UserRepository -- UserValidator: Uses
UserRepository -- User: Interacts with

EmailManager -- EmailSender: Sends email using

class UserController {
    + HandleCreateUser(ctx *gin.Context)
    + HandleGetUser(ctx *gin.Context)
    + HandleGetUserByEmail(ctx *gin.Context)
    + HandleGetUserByUsername(ctx *gin.Context)
    + HandleUpdateUser(ctx *gin.Context)
    + HandleDeleteUser(ctx *gin.Context)
    + HandleVerifyEmail(ctx *gin.Context)
}

class UserService {
    + GetUser(id uint) (*models.User, error)
    + RegisterUser(user *models.User) (*models.User, error)
    + GetUserByEmail(email string) (*models.User, error)
    + GetUserByUsername(username string) (*models.User, error)
    + UpdateUser(user *models.User) (*models.User, error)
    + DeleteUser(id uint) error
    + VerifyEmail(token string) error
}

class UserRepository {
    + CreateUser(user *models.User) (*models.User, error)
    + GetUserByID(id uint) (*models.User, error)
    + GetUserByVerificationToken(token string) (*models.User, error)
    + UpdateUser(user *models.User) (*models.User, error)
    + DeleteUser(user *models.User) (*models.User, error)
    + GetUserByEmail(email string) (*models.User, error)
    + GetUserByUsername(username string) (*models.User, error)
}

class User {
    + ID    string
    + CreatedAt time.Time
    + UpdatedAt time.Time
    + DeletedAt gorm.DeletedAt
    + Username  string
    + Email string
    + Password string
    + EmailVerified bool
    + VerificationToken string
    + ProfilePicture string
    + Bio string
    + Roles string
    + LastLogin time.Time
    + DateOfBirth time.Time
    + PreferredLanguage string
    + ReadingPreferences string
    + IsDeleted bool
}

class EmailManager {
    + SendVerificationEmail(user models.User, sender EmailSender) error
}

class EmailSender {
    <<interface>>
    + SendMail(addr string, auth smtp.Auth, from string, to []string, msg []byte) error
}

class     TokenGenerator {
    + GenerateVerificationToken() string
}

class UserValidator {
    + ValidateUser(user *models.User) error
}

class ExpirationValidator {
    + IsVerificationTokenExpired(createdAt time.Time, emailVerified bool) bool
}


```
