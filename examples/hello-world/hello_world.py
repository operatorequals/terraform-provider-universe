import os
import sys
import json

if __name__ == '__main__':
   event = sys.argv[1]
   config = os.environ.get("multiverse")
   if len(sys.argv) > 2:
       config = sys.argv[2]
   input = sys.stdin.read()
   input_dict = json.loads(input)
   if event == "create":
        input_dict.update({ "id" : "1"})
        result = input_dict
   elif event == "read":
        result =  input_dict
   elif event == "update":
        result =  input_dict
   elif event == "delete":
        result =  {}
   else:
       sys.exit(1)
   print(json.dumps(result))

