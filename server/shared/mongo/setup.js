use("coolcar");
db.account.createIndex({
    open_id: 1,
}, {
    unique: true,
});

db.trip.createIndex({
    "trip.accountid": 1,
    "trip.status": 1,
}, {
    unique: true,
    partialFilterExpression: {
        "trip.status": 1,
    }
});

db.profile.createIndex({
    accountid: 1,
}, {
    unique: true,
});

db.car.insertMany([
    {
        "_id": ObjectId("6321e91019021112b3c519e3"),
        "car": {
            "status": 1,
            "driver": {
                "id": "",
                "avatarurl": ""
            },
            "position": {
                "latitude": 22.87863720136816,
                "longitude": 118.48985338020975
            },
            "tripid": ""
        }
    },
    {
        "_id": ObjectId("6321e91019021112b3c519e4"),
        "car": {
            "status": 1,
            "driver": {
                "id": "",
                "avatarurl": ""
            },
            "position": {
                "latitude": 29.750839648108627,
                "longitude": 122.099106938638
            },
            "tripid": ""
        }
    },
    {
        "_id": ObjectId("6321e91019021112b3c519e5"),
        "car": {
            "status": 1,
            "driver": {
                "id": "",
                "avatarurl": ""
            },
            "position": {
                "latitude": 23.910601318436978,
                "longitude": 118.027081993758
            },
            "tripid": ""
        }
    },
    {
        "_id": ObjectId("6321e91019021112b3c519e6"),
        "car": {
            "status": 1,
            "driver": null,
            "position": {
                "latitude": 24,
                "longitude": 118
            },
            "tripid": ""
        }
    },
    {
        "_id": ObjectId("6321e91019021112b3c519e7"),
        "car": {
            "status": 1,
            "driver": null,
            "position": {
                "latitude": 24,
                "longitude": 118
            },
            "tripid": ""
        }
    },
    {
        "_id": ObjectId("6321e91019021112b3c519e8"),
        "car": {
            "status": 1,
            "driver": {
                "id": "",
                "avatarurl": ""
            },
            "position": {
                "latitude": 24.022740927174564,
                "longitude": 117.99698619456917
            },
            "tripid": ""
        }
    },
    {
        "_id": ObjectId("6321e91019021112b3c519e9"),
        "car": {
            "status": 1,
            "driver": {
                "id": "",
                "avatarurl": ""
            },
            "position": {
                "latitude": 24.031295858170257,
                "longitude": 118.00985773858335
            },
            "tripid": ""
        }
    },
    {
        "_id": ObjectId("6321e91019021112b3c519ea"),
        "car": {
            "status": 1,
            "driver": {
                "id": "",
                "avatarurl": ""
            },
            "position": {
                "latitude": 23.292037964575684,
                "longitude": 118.58477881672204
            },
            "tripid": ""
        }
    },
    {
        "_id": ObjectId("6321e91019021112b3c519eb"),
        "car": {
            "status": 1,
            "driver": null,
            "position": {
                "latitude": 24,
                "longitude": 118
            },
            "tripid": ""
        }
    },
    {
        "_id": ObjectId("6321e91019021112b3c519ec"),
        "car": {
            "status": 1,
            "driver": null,
            "position": {
                "latitude": 24,
                "longitude": 118
            },
            "tripid": ""
        }
    }
]);