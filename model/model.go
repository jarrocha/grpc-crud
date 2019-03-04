package model

import (
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"gitlab.com/codelittinc/golang-interview-project-jaime/grpc-crud/hirepb"
)

type HireDataItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name"`
	Type     hirepb.HireType    `bson:"type"`
	Duration int32              `bson:"duration"`
	Role     string             `bson:"role"`
	Tags     []string           `bson:"tags"`
}

func DataToHirepb(data *HireDataItem) *hirepb.HireData {
	return &hirepb.HireData{
		Id:       data.ID.Hex(),
		Name:     data.Name,
		Type:     data.Type,
		Duration: data.Duration,
		Role:     data.Role,
		Tags:     data.Tags,
	}
}
