import os
import sys
import json

if __name__ == '__main__':
    result = None
    event = sys.argv[1]  # create, read, update or delete, maybe exists too
    ident = os.environ.get("id")  # Get the id if present else None

    if event == "delete":
        exit(0)

    # Read the JSON from standard input
    entre = sys.stdin.read()
    input_dict = json.loads(entre)

    if event == "exists":
        print('true' if ident == "42" else 'false')
        exit(0)

    if event in ["create", "update"]:
        input_dict["@created"] = "26/10/2020 18:55:51"
        input_dict.update({"id": "42"})
    result = input_dict
    print(json.dumps(result))
