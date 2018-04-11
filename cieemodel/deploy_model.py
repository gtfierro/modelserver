import json
from datetime import datetime
from clipper_admin import ClipperConnection, XBOSContainerManager

print "Connecting to XBOS container manager"
clipper_conn = ClipperConnection(XBOSContainerManager())
clipper_conn.connect()

print "Building models"
from thermal_model import get_model_per_zone, normal_schedule, dr_schedule, execute_schedule
models = get_model_per_zone("2018-01-30 00:00:00 PST")

# define our closure that will return the predictions
def execute_thermal_model(params):
    """
    Accepts list of JSON string as argument
    model parameters:
        zone: string
        date: string
        schedule: [(hsp, csp), ... x 24 ...]
    """
    ret = []
    for param in params:
        args = json.loads(param)
        zone = args['zone']
        date = str(args['date'])
        schedule = args['schedule']
        temps, actions = execute_schedule(date, schedule, models[zone], 65)
        ret.append(temps)
    return ret

from clipper_admin.deployers import python as python_deployer
python_deployer.deploy_python_closure(clipper_conn,
                                      name='thermal-model-ciee-test',
                                      version=3,
                                      input_type='strings',
                                      func=execute_thermal_model,
                                      base_image="gtfierro/xbospy")
