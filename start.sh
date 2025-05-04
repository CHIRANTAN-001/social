#!/bin/bash

python3 -m venv venv
. ./venv/bin/activate
pip install boto3
python3 secret_manager.py
deactivate
rm -rf venv
