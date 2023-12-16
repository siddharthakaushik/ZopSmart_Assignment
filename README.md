# ZopSmart_Assignment
#Real time car tracking system
Code Functionality

Certainly! Let's go through the functionality of the provided code:

1. **Database Initialization:**
   - The `initDB` function initializes a SQLite database using the `gorm` package.
   - It connects to the SQLite database stored in the file named "test.db."
   - It auto-migrates the `Car` model, ensuring that the corresponding table is created in the database.

2. **Car Model (car.go):**
   - The `Car` struct defines the model for cars.
   - It embeds the `gorm.DB` type, indicating that it is a GORM model.
   - It has fields representing the brand, model, and status of a car.

3. **API Endpoints (main.go):**
   - **Create Car Endpoint (`/cars` - POST):**
     - Parses JSON data from the request body to create a new car.
     - Saves the car data to the database.

   - **Get Car List Endpoint (`/cars` - GET):**
     - Retrieves a list of cars from the database.
     - Returns the list of cars as a JSON response.

   - **Update Car Endpoint (`/cars/:id` - PUT):**
     - Updates the status of a specific car based on the provided ID.
     - Saves the updated car data to the database.

   - **Delete Car Endpoint (`/cars/:id` - DELETE):**
     - Deletes a specific car based on the provided ID from the database.

4. **WebSocket Endpoint (`/ws`):**
   - Handles WebSocket connections using the `github.com/gorilla/websocket` package.
   - Clients can connect to this endpoint to receive real-time updates.
   - Each WebSocket connection is added to a map (`clientSockets`) for broadcasting updates.

5. **WebSocket Broadcast (`handleWebSocketBroadcast`):**
   - Periodically retrieves the list of cars from the database.
   - Broadcasts the list of cars to all connected WebSocket clients.
   - The broadcast occurs every 5 seconds (configurable).

**Instructions for Running and Testing:**
- Run the server using `go run main.go`.
- Use API endpoints (e.g., cURL commands) to create, retrieve, update, and delete cars.
- Connect to the WebSocket endpoint (`ws://localhost:8080/ws`) for real-time updates on car data.

Expected Results:

The server starts and initializes the SQLite database.
RESTful APIs allow you to perform CRUD operations on car data.
WebSocket connections receive periodic updates on the list of cars.
Real-time communication between the server and WebSocket clients provides an up-to-date view of the car data.
  
This project provides a basic example of building a web server with RESTful APIs and WebSocket functionality in Go.
