import os
import sys

# Demonstrate ENV and -e override capabilities
name = os.getenv("USER_NAME", "Stranger")
print(f"Hello, {name}! Docksmith is running successfully.")

# Demonstrate WORKDIR by checking the current path
print(f"Current Working Directory: {os.getcwd()}")

# Demonstrate visible output for the demo
print("Files in this directory:", os.listdir('.'))