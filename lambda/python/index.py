import json
import os
import vertexai
import base64
from vertexai.language_models import ChatModel, InputOutputTextPair
from google.oauth2 import service_account


def handler(event, context):
    try:
        encoded_credentials_json = os.environ.get('VERTEX_JSON')
        credentials_json = base64.b64decode(encoded_credentials_json).decode('utf-8')

        credentials = service_account.Credentials.from_service_account_info(
            json.loads(credentials_json),
            scopes=["https://www.googleapis.com/auth/cloud-platform"],
        )

        vertexai.init(credentials=credentials, project="airy-task-399918", location="us-central1")
        chat_model = ChatModel.from_pretrained("chat-bison@001")

        parameters = {
            # "candidate_count": 1,
            "max_output_tokens": 256,
            "temperature": 0.2,
            "top_p": 0.8,
            "top_k": 40
        }
        chat = chat_model.start_chat(
            context="""After giving an itinerary of where to eat at, offer to make a reservation at those restaurants.  Make sure to ask if they want you to make a reservation at every restaurant that is confirmed the user wants to try.""",
        )

        response = chat.send_message("""Give me a list of restaurants to try out on my 2 day trip to indianapolis""", **parameters)
    
        response = chat.send_message("""Yes at the rathskeller and slippery noodle""", **parameters)
        
        return {
            'statusCode': 200,
            'body': response.text
        }
    except Exception as e:
        print('\nERROR:  ', e)
        return {
            'statusCode': 500,
            'body': str(e)
        }
