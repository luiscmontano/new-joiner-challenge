import json
import urllib.parse
import boto3
import PyPDF2
import io
import datetime
import re

s3 = boto3.client('s3')
client = boto3.client('sqs')
    
def getFileContent(event):
    bucket = event['Records'][0]['s3']['bucket']['name']
    key = urllib.parse.unquote_plus(event['Records'][0]['s3']['object']['key'], encoding='utf-8')
    
    response = s3.get_object(Bucket = bucket, Key = key)
    file = response["Body"].read()
    pdfReader = PyPDF2.PdfFileReader(io.BytesIO(file))
    
    count = pdfReader.numPages
    fileContent = ""
    for i in range(count):
        page = pdfReader.getPage(i)
        fileContent = fileContent + page.extractText().replace('\n', ' ')
    return fileContent

def getStringWithoutSpaces(str, strFrom, strTo):
    return getString(str, strFrom, strTo).replace(' ', '')

def getString(str, strFrom, strTo):
    matches = re.findall(strFrom + '(.*?)' + strTo, str)
    for match in matches:
        return match.replace(':', '').strip()

def generateJoiner(fileContent):
    
    identificationNumber = getStringWithoutSpaces(fileContent, "Identification Number", "Name")
    name = getStringWithoutSpaces(fileContent, "Name", "Last Name")
    lastName = getStringWithoutSpaces(fileContent, "Last Name", "Stack")
    stack = getString(fileContent, "Stack", "Role")
    role = getStringWithoutSpaces(fileContent, "Role", "English Level")
    englishLevel = getStringWithoutSpaces(fileContent, "English Level", "Domain Experience")
    domainExperience = getStringWithoutSpaces(fileContent, "Domain Experience", "")
    
    message = {
        "identificationNumber": identificationNumber,
        "name": name,
        "lastName": lastName,
        "stack": stack,
        "role": role,
        "englishLevel": englishLevel,
        "domainExperience": domainExperience,
    }
    
    return message
    
def sentMessageToQueue(message):

    client.send_message(
        QueueUrl = "https://sqs.us-east-2.amazonaws.com/944213679037/documentSQS",
        MessageBody = json.dumps(json.dumps(message))
    )

def lambda_handler(event, context):
    
    print("----- Process Document [START]: " + datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S') + " ----- ")
    
    fileContent = getFileContent(event);
    
    message = generateJoiner(fileContent)
    
    sentMessageToQueue(message)
    
    print("----- Process Document [END]:   " + datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S') + " ----- ")

    return "Ok"
    