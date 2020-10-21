import os
import sys
import json
import tempfile
from datetime import datetime

if __name__ == '__main__':
    result = None
    event = sys.argv[1]  # create, read, update or delete, maybe exists too

    id = os.environ.get("filename")  # Get the id if present else None
    script = os.environ.get("script")

    if event == "exists":
        # ignore stdin
        # Is file there?
        if id is None:
            result = False
        else:
            result = os.path.isfile(id)
        print('true' if result else 'false')
        exit(0)

    elif event == "delete":
        # Delete the file
        os.remove(id)
        exit(0)

    # Read the JSON from standard input
    input = sys.stdin.read()
    input_dict = json.loads(input)

    if event == "create":
        # Create a unique file /tmp/hello_world.pyXXXX and write the data to it
        ff = tempfile.NamedTemporaryFile(mode='w+', prefix=script, delete=False)
        input_dict["created"] = datetime.now().strftime("%d/%m/%Y %H:%M:%S")
        ff.write(json.dumps(input_dict))
        ff.close()
        input_dict.update({"filename": ff.name})  # Give the ID back to Terraform - it's the filename
        result = input_dict

    elif event == "read":
        # Open the file given by the id and return the data
        fr = open(id, mode='r+')
        data = fr.read()
        fr.close()
        if len(data) > 0:
            result = json.loads(data)
        else:
            result = {}

    elif event == "update":
        # write the data out to the file given by the Id
        fu = open(id, mode='w+')
        fu.write(json.dumps(input_dict))
        fu.close()
        result = input_dict

    print(json.dumps(result))
