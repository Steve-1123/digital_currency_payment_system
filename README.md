# digital_currency_payment_system

## Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/dcep-simulator.git
   cd dcep-simulator
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Set up MySQL and Redis:
   - Create a MySQL database named `dcep`.
   - Update `main.go` with your MySQL credentials.
   - Ensure Redis is running on `localhost:6379`.
4. Run the node:
   ```bash
   go run main.go
   ```

## Usage
- Start multiple nodes with different IDs and ports:
  ```bash
  go run main.go --id node1 --addr :8080
  go run main.go --id node2 --addr :8081
  ```
- Create and broadcast transactions via API (TBD in future versions).