Here is a suggested GitHub project description for your Go application:

---

## GoMarketWatch

GoMarketWatch is a real-time market data tracking application built with Go, utilizing the power of concurrent websockets and scheduled tasks. It uses the Gin framework for handling HTTP requests, and integrates `tradingview`'s API to fetch and update financial symbol data in real time.

### Features
- **Real-Time Data Updates:** Utilizes websockets to receive updates about financial symbols from TradingView.
- **Concurrent Handling:** Manages multiple symbol data concurrently with robust synchronization using mutexes to ensure data consistency.
- **Automated Cleanup:** Scheduled tasks identify and remove inactive symbols that haven't been requested in over 30 seconds to maintain efficiency.
- **Performance Monitoring:** Integrated with `pprof` for profiling and monitoring the application's performance on the go.

### Getting Started
To get started with GoMarketWatch:
1. Ensure you have Go installed on your system.
2. Clone this repository and navigate into the project directory.
3. Load the necessary environment variables by creating a `.env` file based on the provided `.env.example`.
4. Run the application using:
   ```bash
   go run main.go
   ```

### Endpoints
- **GET `/latest-price`**: Fetches the latest price of a specified symbol. This endpoint expects a query parameter `symbol` which represents the symbol of interest.

### Environment Variables
- `PORT`: Specifies the port on which the server will run.
- `HOST`: Specifies the host address for the server.

### Dependencies
- **Gin-Gonic**: A high-performance web framework.
- **Go-Cron**: Used for scheduling tasks.
- **GoDotenv**: For loading environment variables from a `.env` file.
- **TradingView Websocket Client**: For real-time financial data streaming.

### Contributing
Contributions are welcome! Please feel free to submit pull requests or open issues to suggest improvements or add new features.

---

This description provides a clear overview of the project's purpose, features, setup instructions, and usage, making it easier for potential users or contributors to get started with the application.