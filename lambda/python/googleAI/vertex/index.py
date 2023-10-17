import json
import os
import vertexai
import base64
from vertexai.language_models import ChatModel, InputOutputTextPair
from google.oauth2 import service_account


def handler(event, context):
    try:
        # WHEN RUNNING LOCALLY, MAKE SURE TO ADD THE ENV FILE PARAMETER ON LIKE SO:
        # sam local start-api  --env-vars .env.json
        encoded_credentials_json = os.environ.get('VERTEX_JSON')

        credentials_json = base64.b64decode(
            encoded_credentials_json).decode('utf-8')

        credentials = service_account.Credentials.from_service_account_info(
            json.loads(credentials_json),
            scopes=["https://www.googleapis.com/auth/cloud-platform"],
        )

        vertexai.init(credentials=credentials,
                      project="airy-task-399918", location="us-central1")
        chat_model = ChatModel.from_pretrained("chat-bison@001")

        parameters = {
            # "candidate_count": 1,
            "max_output_tokens": 256,
            "temperature": 0.2,
            "top_p": 0.8,
            "top_k": 40
        }

        if 'body' in event and isinstance(event.get('body'), str):
            json_body = json.loads(event['body'])
        elif 'body' in event and isinstance(event.get('body'), object):
            json_body = event['body']
        else:
            print('Event: ', event)

            return {
                'statusCode': 400,
                'body': json.dumps({'error': 'Missing JSON body'})
            }

        if 'prompt' in json_body:
            prompt_value = json_body['prompt']
        else:
            print('Event: ', event)

            return {
                'statusCode': 400,
                'body': json.dumps({'error': 'Key "prompt" not found in JSON body'})
            }

        chat = chat_model.start_chat(
            context="""
            You are a personal assistant helping the user book a reservation at a restaurant. Ask the user questions to find out what restaurant is the best fit for their needs. Your goal is to guide the user to book a reservation at the restaurant of their choosing. The user will provide criteria as listed below:

            User Criteria: location, type of food, time, date

            User can add other criteria as needed. 

            Always prioritize the criteria based on the order in which the user brings them up in conversation. 

            When you are sending back a list of recommended restaurants, always preface with "Here are your top picks from Google!".

            Send back the list in a bulleted list format grouped by restaurant choices.

            After responding with a message that includes the words "Here are your top picks from Google!", always mention "Click below to make a reservation at Goodfellas!".  There is another button at the bottom of the chat that will handle the reservation.

            The list should show 3 options of restaurants by default; always make sure one of the three options is Goodfellas Pizzaria.  Make sure to display hours of operation for the list of restaurants in the response.

            Do not ask the user if they want to see a list that meets the criteria, just show the list of restaurants.

            When the user is making a reservation, the following criteria are required:

            Reservation Criteria: Name, time, date
            
            If the user asks about topics outside of reservations at restaurants, politely say that you're not able to assist with that and that you are an AI assistant to make reservations at restaurants.

            """,
        )

        response = chat.send_message(prompt_value, **parameters)

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
