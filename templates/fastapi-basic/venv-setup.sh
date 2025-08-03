#!/bin/bash

# FastAPI Virtual Environment Setup Script
# This script sets up a Python virtual environment for the FastAPI project

echo "Setting up Python virtual environment..."

# Create virtual environment
python3 -m venv venv

# Activate virtual environment
source venv/bin/activate

# Upgrade pip
pip install --upgrade pip

# Install dependencies
pip install -r requirements.txt

echo "Virtual environment setup complete!"
echo "To activate the virtual environment, run: source venv/bin/activate" 