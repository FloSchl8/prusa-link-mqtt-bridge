# PrusaLink MQTT Bridge

[![CI](https://github.com/FloSchl8/prusa-link-mqtt-bridge/actions/workflows/ci.yml/badge.svg)](https://github.com/FloSchl8/prusa-link-mqtt-bridge/actions/workflows/ci.yml)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/FloSchl8/prusa-link-mqtt-bridge?style)](https://github.com/florianschlund/prusa-link-mqtt-bridge/releases)

This application acts as a bridge between a PrusaLink-enabled 3D printer and an MQTT broker. It periodically fetches the printer's status via the PrusaLink API and publishes it to a specified MQTT topic.

This allows for easy integration into various home automation systems like Home Assistant, Node-RED, or any other system that can subscribe to MQTT topics.

## Features

-   Fetches printer status from the PrusaLink API.
-   Publishes status data to an MQTT broker.
-   Configurable exclusively via environment variables.
-   Structured, configurable logging (text or JSON).
-   Easy to run with Docker and Docker Compose.
-   Lightweight and efficient, built with Go.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

-   A Prusa 3D printer with PrusaLink enabled.
-   Go (version 1.24 or later) for local development.
-   Docker and Docker Compose (or Podman and Podman Compose) for running the containerized application.

## Configuration

The application is configured exclusively via environment variables. The library `go-envconfig` is used to parse them.

| Variable             | Required | Default      | Description                               |
| -------------------- | -------- | ------------ | ----------------------------------------- |
| `PRUSALINK_HOST`     | **Yes**  | -            | Hostname or IP address of the printer     |
| `PRUSALINK_APIKEY`   | **Yes**  | -            | API key from PrusaLink settings           |
| `PRUSALINK_INTERVAL` | No       | `10`         | Fetch interval in seconds                 |
| `MQTT_BROKER`        | **Yes**  | -            | Hostname or IP address of the MQTT broker |
| `MQTT_PORT`          | No       | `1883`       | Port of the MQTT broker                   |
| `MQTT_USERNAME`      | No       | -            | MQTT username (optional)                  |
| `MQTT_PASSWORD`      | No       | -            | MQTT password (optional)                  |
| `MQTT_TOPIC`         | No       | `prusa-link` | MQTT topic to publish to                  |
| `LOG_LEVEL`          | No       | `INFO`       | Log level (`DEBUG`, `INFO`, `WARN`, `ERROR`) |
| `LOG_FORMAT`         | No       | `text`       | Log format (`text` or `json`)             |

## Important Notes

- **Printer Availability:** The PrusaLink-enabled printer must be online and accessible when the bridge starts. The application fetches the printer's serial number on startup to use it in the MQTT topic, ensuring a unique identifier for each printer. If the printer is not available, the application will exit with an error.

## Usage

### Running with Docker Compose (Recommended)

This is the easiest way to run the application, as it includes a pre-configured MQTT broker (Eclipse Mosquitto).

1.  **Create a `.env` file** in the project root with your PrusaLink credentials and any desired overrides:

    ```env
    # Required
    PRUSALINK_HOST=your-prusa-link-host
    PRUSALINK_APIKEY=your-prusa-link-api-key

    # Optional - for debugging
    LOG_LEVEL=DEBUG
    ```

2.  **Start the services:**

    ```bash
    # With Docker
    docker compose up --build

    # With Podman
    podman compose up --build
    ```

The bridge will start, connect to the included Mosquitto broker, and begin publishing data.

### Running Locally

1.  **Install dependencies:**

    ```bash
    go mod tidy
    ```

2.  **Set the required environment variables** in your shell:

    ```bash
    export PRUSALINK_HOST="your-prusa-link-host"
    export PRUSALINK_APIKEY="your-prusa-link-api-key"
    export MQTT_BROKER="your-mqtt-broker" # e.g., localhost
    ```

3.  **Run the application:**

    ```bash
    go run cmd/prusa-link-mqtt-bridge/main.go
    ```

## Development

### Project Structure

-   `cmd/prusa-link-mqtt-bridge`: Main application entry point.
-   `pkg/config`: Configuration loading logic (`go-envconfig`).
-   `pkg/prusalink`: Logic for interacting with the PrusaLink API.
-   `pkg/mqtt`: Logic for publishing data to MQTT.
-   `Dockerfile`: Defines the multi-stage Docker build.
-   `docker-compose.yml`: Defines the application and MQTT broker services.

### Running Tests

To run the unit tests and see the coverage report, use the following command:

```bash
go test -v -coverprofile=coverage.out ./...
```

The detailed code coverage report can be found in the CI workflow logs under the "Show Coverage Summary" step.

## MQTT Data

The application publishes a JSON payload to the configured MQTT topic. The full MQTT topic path will include the printer's serial number to ensure uniqueness, formatted as `MQTT_TOPIC/SERIAL_NUMBER/status`. For example, if `MQTT_TOPIC` is `prusa-link` and the printer's serial number is `XYZ123`, the status will be published to `prusa-link/XYZ123/status`.

Additionally, the application publishes an availability status to `MQTT_AVAILABILITY_TOPIC/SERIAL_NUMBER/availability`. It sends `online` when the bridge starts and `offline` when it disconnects.

Here is an example of the status data payload:

```json
{
  "state_text": "PRINTING",
  "temp_nozzle": 215.2,
  "target_nozzle": 215.0,
  "temp_bed": 60.1,
  "target_bed": 60.0,
  "axis_z": 3.1,
  "flow": 95,
  "speed": 100,
  "fan_hotend": 4939,
  "fan_print": 5632,
  "progress": 19.0,
  "print_time_left": 11520,
  "print_time": 2913
}
```

## License

This project is licensed under the GNU Affero General Public License v3 (AGPLv3) - see the [LICENSE.md](LICENSE.md) file for details.