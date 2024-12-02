# CM-BACKEND

### Requirenments
- PostgreSQL
- Go
- gin
- Postgres driver for Go

### How to run

1. Create the Postgres database
2. Import all tables and types from the /sql directory
3. Start the /cmd/pccore.go file

### Endpoints

- /users/register
    - Method: POST
    - Purpose: Register a new user
    - Params:
        - Name 
            - Type: String
        - Email
            - Type: String
            - Notes: Accepts only valid email addresses
        - Password
            - Type: String
    - Parameter passing: body
    - Parameter format: JSON
    - Required Role: None
    - Returns: `models.User` or Error
    - Notes: sets user role as `Default`

- /users/login
    - Method: GET
    - Purpose: Login user
    - Params:
        - Email
            - Type: String
            - Notes: Accepts only valid email addresses
        - Password
            - Type: String
    - Parameter passing: Query String
    - Parameter format: Query String
    - Required role: None
    - Returns: `models.User` or Error

- /laptops/add
    - Method: POST
    - Purpose: Add a new laptop
    - Params: 
        - Laptop data:
            - Name
                - Type: String
            - Cpu
                - Type: String
            - Ram
                - Type: int16
            - Gpu
                - Type: String
            - Price
                - Type: String
                - Notes: the price should be a numeric string (e.g., "199.99")
            - Discount
                - Type: int16
                - Notes: the value should be between 0 and 100 (inclusive)
        - Parameter passing: body
        - Parameter format: JSON

        - Authentification data:
            - Email
                - Type: String
                - Notes: Accepts only valid email addresses
            - Password
                - Type: String
        - Parameter passing: Query String
        - Format: Query String
    - Required Role: Admin
    - Returns: `models.Laptop` or Error
