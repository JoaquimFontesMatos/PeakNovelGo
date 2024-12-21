```plantuml
left to right direction

package "Frontend" as frontend
package "Backend" as backend {
    component "API Gateway (controllers)" as apiGateway
    component "Services (services)" as services
    component "Database (database)" as database
    component "Middleware (middleware)" as middleware
    component "Models (models)" as models
    component "Config (config)" as config

    interface "DTOs" as dtos
    interface "Settings" as settings
}

interface "Request" as request
interface "Response" as response
interface "External API" as externalAPI

frontend -u- request : Sends HTTP requests
apiGateway -d-( request : Receives HTTP requests

apiGateway -u- response: Sends HTTP responses
frontend -d-( response : Receives HTTP responses

services -u- dtos : Provides DTOs
apiGateway -( dtos : Uses DTOs

services -d- settings : Uses settings
config - settings : Provides settings

services -d-> database : Interacts with database
apiGateway -l-> middleware : Applies middleware
services -u-> externalAPI : Communicates with external API
models -> services : Provides models for services

```
