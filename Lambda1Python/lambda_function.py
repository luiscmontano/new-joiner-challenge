import json
import urllib.parse
import boto3
import PyPDF2
import io

s3 = boto3.client('s3')

def lambda_handler(event, context):

    bucket = event['Records'][0]['s3']['bucket']['name']
    key = urllib.parse.unquote_plus(event['Records'][0]['s3']['object']['key'], encoding='utf-8')
    
    try:
        
        response = s3.get_object(Bucket = bucket, Key = key)
        file = response["Body"].read()
        pdfReader = PyPDF2.PdfFileReader(io.BytesIO(file))
        firstPage = pdfReader.getPage(0)
        firstPageText = firstPage.extractText()
        print(firstPageText)
        
        print("CONTENT TYPE: " + response['ContentType'])
        return "All is OK"
    except Exception as e:
        print(e)
        print('Error getting object {} from bucket {}. Make sure they exist and your bucket is in the same region as this function.'.format(key, bucket))
        raise e
