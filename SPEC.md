# XBOS / Clipper Integration

## Overview

This document describes the integration of Clipper into a model/prediction service that leverages XBOS metadata and data.

Clipper will be deployed in "the cloud" rather than locally (likely on a Kubernetes cluster) and will serve models (Python closures)
and predictions on those models. Models are currently identified by an "application name", but soon Clipper will support referring to models
explicitly by name and version.

The XBOS prediction API will be a layer above this. Models will be defined/created locally and will need to be sent to the cluster over BOSSWAVE so they can be used. These models will have the following metadata:
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
under some schedule), or even a call against an external API (like weather underground -- we'd likely use a caching layer here too to avoid hitting the API too many times).

The XBOS frontend will accept prediction queries with the following parameters:
- model type: (e.g. consumption, thermal, occupancy, weather)
- model context: (e.g. zone, room, floor)
- prediction range: start, end timestamp e.g. `2018-01-30 00:00:00 PST - 2018-01-31 00:00:00 PST`
- prediction granularity: `30min`, `1h`

It will find the most appropriate model or set of models using the model metadata, query those, and merge the results appropriately.

## Architecture

```
+-----------------------------------------------------------------+
|        Model Containers                                         |
|                                |         Clipper Processes      |
|      +--------------------+    |                                |
|      |Thermal Model v1    |    |         +---------------+      |
|      +--------------------+    |         | Clipper Mgmt  |      |
|      +--------------------+    |         +---------------+      |
|      |Thermal Model v2    |    |                                |
|      +--------------------+    |         +---------------+      |
|      +--------------------+    |         | Clipper Redis |      |
|      |Occupancy Model v1  |    |         +---------------+      |
|      +--------------------+    |                                |
|      +--------------------+    |         +---------------+      |
|      |Consumption Model v1|    |         | Clipper Query |      |
|      +--------------------+    |         +---------------+      |
|                                                                 |
|      ----------------  Clipper HTTP API  -----------------      |
|                                                                 |
|       +---------------+            +-------------+              |
|       |               |<---------->| Brick Model |              |
|       | Model Chooser |            |   (HodDB)   |              |
|       |               |            +-------------+              |
|       +---------------+                  ^                      |
|              ^                           |                      |
|              |                           |                      |
|              v                           v                      |
|       +------------------------------------------+              |
|       |             XBOS Prediction API          |              |
|       +------------------------------------------+              |
|                             ^                                   |
+-----------------------------|---------- XBOS Model Cluster -----+
                              |
                              |
                              |
                              v
             +---------------------------------+
             |                                 |
             |       BOSSWAVE Message Bus      |
             |                                 |
             +---------------------------------+
```

## Implementation Notes

There are two parts to the implementation:
- server/cluster side:
    - proxy BOSSWAVE requests to the Clipper API
    - deserialize containers
    - metadata description of models + query capabilities
- client side:
    - container manager implementation for BOSSWAVE
    - keep most of the Clipper client implementation


### Container Manager

- `start_clipper`: 
    - no-op; the cluster is always on
- `connect(base_uri)`:
    - emulates a session over BOSSWAVE
    - check service liveness messages as a heartbeat (in BOSSWAVE)
- `build and deploy model`:
    - create docker container locally using existing clipper code
    - `docker save` to serialize docker container to `.tar` file
    - deploy: transmit `.tar` file to cloud service using BOSSWAVE;
      the `register` call will be transmitted over BOSSWAVE as well
    - will need to augment the call to accept additional metadata about the model (see above)
- `register app`, `link app to model`, etc:
    - REST calls get authenticated and proxied over BOSSWAVE to the Clipper API

### Server Implementation

- proxy received BOSSWAVE calls to the Clipper API and return the response over BOSSWAVE
- serialized Docker containers are unpacked using `docker load` so they can be added to the Clipper cluster
- we will save all of the model Docker containers so we can query historical versions of the models

- need to proxy parts of the Docker API from the XBOS container manager (client side) through to the server implementation
    - list containers with label filter
    - run container with image, name, env vars
