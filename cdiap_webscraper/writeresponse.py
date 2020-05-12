import requests

r = requests.get('https://api.covid19api.com/summary')

with open('/users/remsy/CDIAP/Research/covid_cases.json','w') as fd:
    fd.write(r.text)