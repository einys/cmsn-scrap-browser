
```markdown
# Twitter Scraper Flask Application

This application is a Flask-based web service that uses Selenium WebDriver with Firefox to scrape data from Twitter pages. It extracts text content, images, usernames, and meta tags from a specified Twitter URL.

## Requirements

To run this application, you need the following:

- Python 3.x
- Flask
- Selenium
- WebDriver Manager for Selenium
- Mozilla Firefox browser
- GeckoDriver

## Setup

### 1. Clone the repository

```bash
git clone <repository-url>
cd <repository-directory>
```

### 2. Install Python dependencies

You can install the required dependencies using `pip`:

```bash
pip install Flask selenium webdriver-manager
```

### 3. Configure Firefox and GeckoDriver

Ensure that Firefox is installed on your system. The application sets the Firefox binary location based on the operating system:

- **macOS**: `/Applications/Firefox.app/Contents/MacOS/firefox`
- **Linux**: `/usr/bin/firefox` or the path specified by the `FIREFOX_BINARY_PATH` environment variable.
- **Windows**: `C:\Program Files\Mozilla Firefox\firefox.exe`

### 4. Run the Flask application

Start the Flask server by running:

```bash
python app.py
```

By default, the server will start on `http://0.0.0.0:18081`.

## Usage

To scrape data from a Twitter page, send a POST request to `/scrape-twitter` with a JSON payload containing the Twitter URL:

### Example Request

```bash
curl -X POST http://localhost:18081/scrape-twitter -H "Content-Type: application/json" -d '{"url": "https://twitter.com/username/status/1234567890"}'
```

### Example Response

```json
{
  "text": "This is an example tweet text.",
  "image": "https://pbs.twimg.com/media/example.jpg",
  "username": "exampleuser",
  "meta_tag": "An example tweet title"
}
```

## Error Handling

The application includes robust error handling for different types of exceptions:

- **Missing URL**: Returns a 400 error if the URL is not provided in the request.
- **Timeouts**: Returns a 500 error if loading elements (text, image, username) times out.
- **Internal Server Errors**: Captures and logs exceptions, returning detailed error messages and tracebacks in the response.

### Example Error Response

```json
{
  "error": "Internal server error",
  "message": "Timeout loading text content",
  "traceback": "<stack trace>"
}
```

## Logging

The application uses the `logging` module to log important events and errors. Logs are output to the console.

## Development

### Running Locally

To run the application locally for development:

1. Set up a virtual environment (optional but recommended):

    ```bash
    python -m venv venv
    source venv/bin/activate  # On Windows use `venv\Scripts\activate`
    ```

2. Install the dependencies:

    ```bash
    pip install -r requirements.txt
    ```

3. Run the Flask app:

    ```bash
    python app.py
    ```

### Notes

- This application uses Selenium WebDriver in headless mode, which means no browser window will open during scraping.
- Be aware of Twitter's terms of service and policies when scraping data from their platform.

## License

This project is licensed under the MIT License. See the `LICENSE` file for more details.
```

### Additional Notes:

- Replace `<repository-url>` and `<repository-directory>` with your actual repository URL and directory names.