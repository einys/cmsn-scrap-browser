import csv
import requests
import json
import random
import time
import os

# File paths
csv_file_path = 'testapi.items.csv'  # Change this to your CSV file path
output_csv_path = "output_results.csv"  # The file where results will be saved

# API URL
api_url = "http://localhost:18081/scrape-twitter"

# Reading the CSV file
with open(csv_file_path, 'r', encoding='utf-8') as csvfile:
    reader = csv.DictReader(csvfile)
    ids = [row['source_item_id'] for row in reader]

# Infinite loop to send POST requests at random intervals
while True:
    # Check if the output CSV file already exists and has content
    file_exists = os.path.isfile(output_csv_path)
    
    # Open the output CSV file in append mode
    with open(output_csv_path, 'a', newline='', encoding='utf-8') as csvfile:
        fieldnames = ['URL', 'Status Code', 'Response', 'Error', 'Timestamp']
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)

        # Write the header only if the file does not exist or is empty
        if not file_exists or os.stat(output_csv_path).st_size == 0:
            writer.writeheader()
            
        # Shuffle the list of tweet IDs randomly
        random.shuffle(ids)

        # Pick one random tweet_id from the shuffled list
        tweet_id = ids[0]

        # Prepare the payload
        url = f"https://x.com/i/status/{tweet_id}"
        payload = json.dumps({"url": url})
        
        # print the id
        print(f"Requesting Tweet ID: {tweet_id}")

        try:
            # Send the POST request
            response = requests.post(api_url, headers={"Content-Type": "application/json"}, data=payload)

            # Log the response into the CSV file
            writer.writerow({
            'URL': url,
            'Status Code': response.status_code,
            'Response': response.text,
            'Error': '',  # No error
            'Timestamp': time.strftime("%Y-%m-%d %H:%M:%S %Z%z")  # Add timestamp with timezone
            })

            # Print the response (optional)
            print(f"URL: {url}")
            print(f"Status Code: {response.status_code}")
            print(f"Response: {response.text}")
            print("\n")

        except requests.exceptions.RequestException as e:
            # Handle any exceptions and log the error into the CSV file
            writer.writerow({
            'URL': url,
            'Status Code': '',  # No status code
            'Response': '',  # No response
            'Error': str(e),
            'Timestamp': time.strftime("%Y-%m-%d %H:%M:%S %Z%z")  # Add timestamp with timezone
            })

            # Print the error (optional)
            print(f"URL: {url}")
            print(f"Error: {str(e)}")
            print("\n")

        # Sleep for a random time between 10 to 50 seconds to avoid rate-limiting issues
        sec = random.uniform(10, 50)
        print(f"Sleeping for {sec:.2f} seconds...\n\n")
        time.sleep(sec)
