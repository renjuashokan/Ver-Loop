import requests
import json
import time

content = []
filename = r'C:\Users\Renju\Documents\mystory.txt'
with open(filename) as f:
    content = f.readlines()

f.close()
content = [x.strip() for x in content]

url = "http://localhost:9090/add"
headers = {'content-type': 'application/json'}
mydict = {}
cnt = 0
for item in content:
    ww = item.split()
    for i in ww:
        print(i)

        mydict["word"] = i
        response = requests.post(
            url=url, data=json.dumps(mydict), headers=headers)
        print("time elapse ", response.elapsed.microseconds)
        print("Status code: ", response.status_code)
        print("Printing Entire Post Request")
        print(response.json())
        time.sleep(0.1)
        cnt = cnt + 1

print("total number of sents", cnt)
