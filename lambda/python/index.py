def handler(event, context):
    try:
        return {
            'statusCode': 200,
            'body': 'Hello, World!'
        }
    except Exception as e:
        return {
            'statusCode': 500,
            'body': str(e)
        }
