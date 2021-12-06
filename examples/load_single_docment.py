import base64, json
import requests

user = "admin"
password = "Complexpass#123"
bas64encoded_creds = base64.b64encode(bytes(user + ":" + password, "utf-8")).decode("utf-8")


data = {
    "Athlete": "DEMTSCHENKO, Albert",
    "City": "Turin",
    "Country": "RUS",
    "Discipline": "Luge",
    "Event": "Singles",
    "Gender": "Men",
    "Medal": "Silver",
    "Season": "winter",
    "Sport": "Luge",
    "Year": 2006
  }

headers = {"Content-type": "application/json", "Authorization": "Basic " + bas64encoded_creds}
index = "games3"
zinc_host = "http://localhost:4080"
zinc_url = zinc_host + "/api/" + index + "/document"

res = requests.put(zinc_url, headers=headers, data=json.dumps(data))
