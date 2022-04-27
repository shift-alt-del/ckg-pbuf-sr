from chalice import Chalice

app = Chalice(app_name='pbuf-sr')

@app.lambda_function(name='demo')
def demo(event, context):
    """
    A demo function to interact with Confluent Lambda sink connector.
    The `event` here are deserialized.
    """

    # to print out the event parameter.
    print(f'Raw event received: {event}')
    # to print out the context parameter.
    print(f'Raw event received: {context}')

    # how large the batch size is.
    batch_size = len(event)
    print(f'Batch size: {batch_size}')

    for message in event:
        # print the single messages.
        print(message)

    return {'this is': 'only a demo'}