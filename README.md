# Go-JupyterHub

🚧 **Disclaimer: This project is in the alpha stage and is under active development. Features may change, and stability is not guaranteed. Use at your own risk.** 🚧

Go-JupyterHub is a lightweight, Go-based alternative to JupyterHub, designed to manage and serve Jupyter notebooks efficiently. Leveraging Go's concurrency features, this project aims to provide a performant and scalable environment for interactive computing.

## Features

- **Lightweight and Fast**: Built with Go, ensuring minimal resource usage and quick response times.
- **Concurrent Notebook Management**: Utilizes Go's goroutines to handle multiple notebook sessions simultaneously.
- **Simplified Deployment**: Easier setup and deployment compared to traditional JupyterHub installations.

## Installation

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/viktorszagorskis/go-jupyterhub.git
   cd go-jupyterhub
   ```

2. **Build the Application**:

   ```bash
   go build -o go-jupyterhub main.go
   ```

3. **Run the Application**:

   ```bash
   ./go-jupyterhub
   ```

   By default, the server will start on `http://localhost:8080`.

## Usage

- **Access the Web Interface**: Navigate to `http://localhost:8080` to access the main dashboard.
- **Create a New Notebook**: Click on "New Notebook" to start a new Jupyter notebook session.
- **Manage Existing Notebooks**: View and manage your active notebook sessions from the dashboard.

## Configuration

Configuration options can be set via environment variables or a configuration file. Refer to the `config.example.json` file for available settings.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request for any enhancements or bug fixes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

