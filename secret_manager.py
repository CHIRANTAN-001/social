import boto3
from botocore.exceptions import ClientError
import json


def get_secret():

    secret_name = "goec2"
    region_name = "ap-south-1"

    session = boto3.session.Session()
    client = session.client(
        service_name='secretsmanager',
        region_name=region_name
    )

    try:
        get_secret_value_response = client.get_secret_value(
            SecretId=secret_name
        )
    except ClientError as e:
        raise e

    secret = get_secret_value_response['SecretString']
    return secret


def write_to_env_file(secret_json):
    try:
        secret_dict = json.loads(secret_json)

        with open('/home/ubuntu/.env', 'w') as env_file:
            for key, value in secret_dict.items():
                if isinstance(value, str):
                    if '"' in value or "'" in value or " " in value or "\n" in value:
                        formatted_value = f'"{value}"'
                    else:
                        formatted_value = value
                else:
                    formatted_value = str(value)

                env_file.write(f"{key}={formatted_value}\n")

        print(f"Successfully wrote {len(secret_dict)} environment variables to .env file")
        return True
    except json.JSONDecodeError as e:
        print(f"Error parsing JSON: {e}")
        return False
    except IOError as e:
        print(f"Error writing to .env file: {e}")
        return False

write_to_env_file(get_secret())