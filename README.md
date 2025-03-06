# Ascenda-Loyalty-Assignment

## Overview

The project is designed to provide hotel information based on specific criteria such as hotel IDs and destination IDs. This API supports filtering and merging of hotel data to provide comprehensive responses.

## Features

- Retrieve hotel information by hotel ID or destination ID.
- Merge hotel data from multiple sources.
- Validate input parameters.
- Error handling for consistent API responses.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Error Handling](#error-handling)
- [Testing](#testing)
- [Design Considerations](#design-considerations)
- [Contributing](#contributing)
- [License](#license)

## Installation

1. Clone the repository:
   ```bash
   https://github.com/mle003/Ascenda-loyalty-assignment.git

2. Install dependencies:
    ```bash
    go mod tidy -v
    go mod download
    go mod verify

## Usage

1. To start the server run:
    ```bash
    go run cmd/server/server.go

The server will start on http://localhost:8000.

2. To test with completely new data remove reset the internal/data/hotels.json file to blank file

## API Endpoints

1. Get Hotel by ID

- Endpoint: /hotel_id
- Method: GET
- Parameters:
    - `hotelIds` (optional): A comma-separated list of hotel IDs to filter by.  Example: `hotelIds=hotel1,hotel2,hotel3`
    - `destinationIds` (optional): A comma-separated list of destination IDs to
    - Response: 
        ```json
        [
           {
            "id": "SjyX",
            "destination_id": 5432,
            "hotel_name": "InterContinental",
            "location": {
                "address": "1 Nanson Rd, Singapore 238909",
                "country": "Singapore"
            },
            "description": [
                "InterContinental Singapore Robertson Quay is luxury's preferred address offering stylishly cosmopolitan riverside living for discerning travelers to Singapore. Prominently situated along the Singapore River, the 225-room inspiring luxury hotel is easily accessible to the Marina Bay Financial District, Central Business District, Orchard Road and Singapore Changi International Airport, all located a short drive away. The hotel features the latest in Club InterContinental design and service experience, and five dining options including Publico, an Italian landmark dining and entertainment destination by the waterfront."
            ],
            "amenities": {
                "general": [
                    "outdoor pool",
                    "business center",
                    "childcare",
                    "parking",
                    "bar",
                    "dry cleaning",
                    "wifi",
                    "breakfast",
                    "concierge"
                ],
                "room": [
                    "aircon",
                    "minibar",
                    "tv",
                    "bathtub",
                    "hair dryer"
                ]
            },
            "images": {
                "rooms": [
                    {
                        "link": "https://d2ey9sqrvkqdfs.cloudfront.net/Sjym/i93_m.jpg",
                        "description": "Double room"
                    },
                    {
                        "link": "https://d2ey9sqrvkqdfs.cloudfront.net/Sjym/i94_m.jpg",
                        "description": "Bathroom"
                    }
                ],
                "site": [
                    {
                        "link": "https://d2ey9sqrvkqdfs.cloudfront.net/Sjym/i1_m.jpg",
                        "description": "Restaurant"
                    },
                    {
                        "link": "https://d2ey9sqrvkqdfs.cloudfront.net/Sjym/i2_m.jpg",
                        "description": "Hotel Exterior"
                    },
                    {
                        "link": "https://d2ey9sqrvkqdfs.cloudfront.net/Sjym/i5_m.jpg",
                        "description": "Entrance"
                    },
                    {
                        "link": "https://d2ey9sqrvkqdfs.cloudfront.net/Sjym/i24_m.jpg",
                        "description": "Bar"
                    }
                ]
            },
            ...
        ]

2Update Hotel Data
- Endpoint: /update_data
- Method: POST
- Description: This endpoint queries a list of external endpoints to fetch the latest hotel data and updates the local data.
- Parameters: None required in the request body.
- Response: 
    ```json
    {
    "message": "Hotel data updated successfully",
    "sources": [
        "http://externalapi1.com",
        "http://externalapi2.com",
        ...
    ]
    }


## Error Handling

The API uses a centralized error-handling middleware to provide consistent error responses. Common error responses include:

- 400 Bad Request: Missing or invalid parameters.
- 500 Internal Server Error: Unexpected server errors.

## Testing 

Run the tests using:
    ```bash
    go test ./...

## Design Considerations
- Data storage: Since no database implementation required, hotels data is stored as json file as map[string]Hotel data format to maintain the unique of hotel id
- Description is stored as []string type instead of string in example response format for multiple descriptions 
- Scalability: The API is designed to handle a large number of requests efficiently.
- Extensibility: The merging logic is modular to allow easy addition of new data sources.
- Security: Ensure no sensitive data is exposed via the API. Details of errors will be logged for troubleshooting and only returns generic error message to client

