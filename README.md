# Distributor Permissions API

This API enables the management and verification of distributor access permissions. Each distributor can have specific regions `INCLUDE` and `EXCLUDE`. A `Check Permission` endpoint allows you to validate if a distributor has permission to operate in a specified region.

## Table of Contents
- [Installation](#installation)
- [Endpoints](#endpoints)
  - [Add Distributor](#add-distributor)
  - [Get Distributor](#get-distributor)
  - [Check Permission](#check-permission)
- [Usage Examples](#usage-examples)
- [Error Codes](#error-codes)

---

## Installation

1. **Clone the Repository**:
   ```bash
   git clone <repo-url>
   cd <repo-folder>
   ```
2. **Install Dependencies**:
   ```bash
   go mod tidy
   ```
3. **Run the Server**:
   ```bash
   go run main.go
   ```
   The server will start on `http://localhost:8000`.

---

## Endpoints

### Add Distributor

**Description**: Adds a new distributor with specified `INCLUDE` and `EXCLUDE` permissions.

- **URL**: `/distributor`
- **Method**: `POST`
- **Headers**:
  - `Content-Type`: `application/json`
- **Request Body**:
  - `name` (string): Name of the distributor.
  - `parent` (string): Name of the parent distributor.
  - `include` (array of strings): Regions where the distributor has permission to operate. Accepted formats: `COUNTRY`, `PROVINCE-COUNTRY`, or `CITY-PROVINCE-COUNTRY`.
  - `exclude` (array of strings): Regions where the distributor is restricted from operating. Accepted formats: `COUNTRY`, `PROVINCE-COUNTRY`, or `CITY-PROVINCE-COUNTRY`.

#### Example Request
```bash
curl --location 'http://localhost:8000/distributor' \
--header 'Content-Type: application/json' \
--data '{
    "name": "DIST",
    "include": [
        "INDIA"
    ],
    "exclude": [
        "KARNATAKA-INDIA"
    ]
}'
```

#### Example Response
- **201 Created** on success.
- **400 Bad Request** if input data is invalid.
- **404 Not Found** if parent distributor not found.

---

### Get Distributor

**Description**: Retrieves details of a distributor, with all their permissions.

- **URL**: `/distributor/{name}`
- **Method**: `GET`
- **Parameters**:
  - `name` (string): Name of the distributor.

#### Example Request
```bash
curl --location 'http://localhost:8000/distributor/DIST'
```

#### Example Response
```json
{
  "data": {
    "name": "DIST",
    "locations": {
      "INDIA": {
        "code": "IN",
        "provinces": {
          "KARNATAKA": {
            "code": "KA",
            "cities": {
              "BENGALURU": {
                "code": "BENAU"
              }
            }
          }
        }
      }
    }
  }
}
```

---

### Check Permission

**Description**: Checks if a distributor has access to a specific region.

- **URL**: `/distributor/{name}/permission`
- **Method**: `GET`
- **Parameters**:
  - `name` (string): Name of the distributor.
  - `region` (query parameter, string): Region to check for access, formatted as `COUNTRY`, `PROVINCE-COUNTRY`, or `CITY-PROVINCE-COUNTRY`.


---

#### Example Request
```bash
curl --location 'http://localhost:8000/distributor/DIST/permission?region=TAMILNADU-INDIA'
```

#### Example Response
```json
{
  "data": "NO"
}
```

- **Possible Responses**:
  - `"YES"` if the distributor has permission.
  - `"NO"` if the distributor is restricted from the region.
  - **400 Bad Request** if the region format is invalid.

---

## Usage Examples

### Adding a Distributor
```bash
curl --location 'http://localhost:8000/distributor' \
--header 'Content-Type: application/json' \
--data '{
    "name": "DIST",
    "include": [
        "INDIA"
    ],
    "exclude": [
        "KARNATAKA-INDIA"
    ]
}'
```

### Retrieving a Distributor
```bash
curl --location 'http://localhost:8000/distributor/DIST'
```

### Checking Permission for a Region
```bash
curl --location 'http://localhost:8000/distributor/DIST/permission?region=BANGALORE-KARNATAKA-INDIA'
```

---

## Error Codes

- **400 Bad Request**: Invalid input or region format.
- **404 Not Found**: Distributor not found.

---


