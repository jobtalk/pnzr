var vars = require("./vars/vars.js");

config.Version = 1.0;

service.body = {
    "Cluster": "backend",
    "DeploymentConfiguration": {
        "MaximumPercent": 200,
        "MinimumHealthyPercent": 100
    },
    "DesiredCount": 1,
    "LoadBalancers": [
        {
            "ContainerName": "nginx",
            "ContainerPort": 80,
            "TargetGroupArn": vars.targetGroupArn
        }
    ],
    "Role": "ecsServiceRole",
    "ServiceName": vars.serviceName,
    "TaskDefinition": vars.taskDefinitionName
};

vars.taskDefinition.body = {
    "ContainerDefinitions": [
        {
            "Cpu": 0,
            "Essential": true,
            "ExtraHosts": null,
            "Image": "ieee0824/nginx-template",
            "MemoryReservation": 128,
            "Name": "nginx",
            "Links": [
                "api"
            ],
            "PortMappings": [
                {
                    "HostPort": 0,
                    "ContainerPort": 80,
                    "Protocol": "tcp"
                }
            ]
        },
        {
            "Cpu": 0,
            "Essential": true,
            "ExtraHosts": null,
            "Image": "ieee0824/dummy-app",
            "MemoryReservation": 512,
            "Name": "api"
        }
    ],
    "Family": vars.family,
    "NetworkMode": "bridge"
}


