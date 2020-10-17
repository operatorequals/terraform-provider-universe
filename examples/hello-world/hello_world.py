import sys
import json


def handler(event, data):
   data_as_json = json.loads(data)
   if event == "create":
        result = { "id" : "1" }   
   elif event == "read":
        result =  {"id": "1" }   
   elif event == "update":
        result =  { "id": "1" }
   elif event == "delete":
        result =  {}
   return result

def read_data():
    data = ''
    for line in sys.stdin:
        data += line

    return data
   
if __name__ == '__main__':
    context = sys.stdin.read()
    print(json.dumps(handler(sys.argv[1], context)))
