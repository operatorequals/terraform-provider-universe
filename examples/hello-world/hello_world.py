import os
import sys
import json
import tempfile

if __name__ == '__main__':
   event = sys.argv[1]
   id = os.environ.get("id")
   script = os.environ.get("script")
   input = sys.stdin.read()
   input_dict = json.loads(input)
   result = None
   if event == "create":
        ff = tempfile.NamedTemporaryFile(mode = 'w+',  prefix=script, delete=False)
        ff.write(json.dumps(input_dict))
        ff.close()
        input_dict.update({ "id" : ff.name})
        result = input_dict
   elif event == "read":
        fr=open(id, mode='r+')
        data = fr.read()
        fr.close()
        if len(data) > 0:
            result = json.loads(data)
        else:
            result = {}
   elif event == "update":
       fu=open(id,mode='w+')
       fu.write(json.dumps(input_dict))
       fu.close()
       result = input_dict
   elif event == "delete":
       os.remove(id)
       result =  {}
   print(json.dumps(result))

