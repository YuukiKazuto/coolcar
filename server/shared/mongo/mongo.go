package mgutil

import (
	"coolcar/shared/mongo/objid"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	IDFieldName        = "_id"
	UpdatedAtFieldName = "updatedat"
)

type IDField struct {
	ID primitive.ObjectID `bson:"_id"`
}

type UpdatedAtField struct {
	UpdatedAt int64 `bson:"updatedat"`
}

var NewObjID = primitive.NewObjectID

func NewObjIDWithValue(id fmt.Stringer) {
	NewObjID = func() primitive.ObjectID {
		return objid.MustFromID(id)
	}
}

var UpdatedAt = func() int64 {
	return time.Now().UnixNano()
}

func Set(v interface{}) bson.M {
	return bson.M{
		"$set": v,
	}
}
func SetOnInsert(v interface{}) bson.M {
	return bson.M{
		"$setOnInsert": v,
	}
}
func ZeroOrDoesNotExist(field string, zero interface{}) bson.M {
	return bson.M{
		"$or": []bson.M{
			{
				field: zero,
			},
			{
				field: bson.M{
					"$exists": false,
				},
			},
		},
	}
}
