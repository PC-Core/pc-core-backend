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
    - Required Role: None / Default
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
    - Required role: None / Default
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
                - Type: Float64
            - Stock
                - Type: Uint64
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
    - Returns: `models.Product` and `models.LaptopChars` or Error

- /laptops/chars/:id
    - Method: GET
    - Purpose: Get laptop characteristics
    - Params:
        - Id:
            - Type: Uint64
    - Parameter passing: Query
    - Parameter format: none
    - Required Role: None / Default
    - Returns: `models.LaptopChars` or Error

- /products/
    - Method: GET
    - Purpose: Get products from a certain index in a certain quantity
    - Params: 
        - Page
            - Type: Uint64
        - Count
            - Type: Uint64
    - Parameter passing: Query String
    - Parameter format: Query String
    - Required Role: None / Default
    - Returns: `[]models.Product` or Error

- /products/:id
    - Method: GET
    - Purpose: Get the specific product by id
    - Params:
        - Id
            - Type: Uint64
    - Parameter passing: Query
    - Parameter format: none
    - Required Role: None / Default
    - Returns: `models.Product` or Error
