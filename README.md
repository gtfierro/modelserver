# Model Server

BOSSWAVE interface on top of Clipper

`clipperserver.py`: sample script putting a model into Clipper

## Clipper Usage

Model calls for thermal model (Python):

```python
# start at midnight
normal_schedule = [
    # midnight - 8:00am
    (50, 90),(50, 90), (50, 90),(50, 90), (50, 90),(50, 90), (50, 90),(50, 90), (50, 90),(50, 90), (50, 90),(50, 90), (50, 90),(50, 90), (50, 90),(50, 90),
    # 8:00am - 4:00pm
    (70, 74),(70, 74),(70, 74),(70, 74),(70, 74),(70, 74),(70, 74),(70, 74),(70, 74),(70, 74),(70, 74),(70, 74),(70, 74),(70, 74),(70, 74),(70, 74),
    # 4:00pm - 6:00pm
    (70, 74),(70, 74),(70, 74),(70, 74),
    # 6:00pm - 12:00am
    (50, 90),(50, 90),(50, 90),(50, 90),(50, 90),(50, 90),(50, 90),(50, 90),(50, 90),(50, 90),(50, 90),(50, 90)
]
params = {
    'zone': 'http://buildsys.org/ontologies/ciee#CentralZone',
    'date': '2018-02-06 00:00:00 UTC',
    'schedule': normal_schedule
}

import requests
import json
resp = requests.post('http://localhost:1337/<model name>/predict', data=json.dumps({'input': params}))
print resp.json()['output']
```

## Model Descriptions

For each model, we keep the following metadata:
- model name
- model version (combination of name, version is unique and immutable)
- predictive range + resolution (list of days, date ranges that the model can predict)
    - possibly support other definitions of what are the 'good' ranges for this model
    - possibly a function that gets called? `model.can_predict(start date, end date, resolution) => {true, false}`
- training data streams
    - the URIs or UUIDs or Brick model names of the streams used to train this model
    - can be used to determine if user can access this model (based on the URIs and BOSSWAVE DoTs)
    - optional
- training data range
    - the range of data this model was trained on
    - optional

Models are defined as a Python closure running through Clipper. This gives us a very flexible execution model.
These models could be a simple evaluation such as using a weight vector, a more involved iteration (such as MPC
under some schedule), or even a call against an external API (like weather underground -- we'd likely use a caching
layer here too to avoid hitting the API too many times).

It should be super simple to deploy a model from a piece of Python code; likely just a simple wrapper for Clipper (make
the call idempotent, or automatically increment the version number).

## Prediction API

Query contains:
- model type: (consumption, thermal, occupancy, weather)
- model context: (zone, room, etc)
- prediction range: start, end timestamp e.g. `2018-01-30 00:00:00 PST - 2018-01-31 00:00:00 PST`
- prediction granularity: `30min`, `1h`

Idea is that you could use these parameters to have the model server choose the model.

Optional parameters:
- data to use to train the model
- entity to execute the model as a given entity
- model name (invoke a model directly)
