// package contains handlers
package controller

import (
	"context"
	"fmt"
	"log"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"gitlab.com/codelittinc/golang-interview-project-jaime/grpc-crud/hirepb"
	"gitlab.com/codelittinc/golang-interview-project-jaime/grpc-crud/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	SrvClient     *mongo.Client
	SrvCol        *mongo.Collection
	SrcCtx        context.Context
	SrvCancelFunc context.CancelFunc
)

type HireDataController struct {
}

// CreateHire is called to create a hire. It uses a unary operation
func (*HireDataController) CreateHire(ctx context.Context,
	req *hirepb.CreateHireRequest) (*hirepb.CreateHireResponse, error) {

	log.Println("Create hire request received.")

	hire := req.GetData()

	data := model.HireDataItem{
		Name:     hire.GetName(),
		Type:     hire.GetType(),
		Duration: hire.GetDuration(),
		Role:     hire.GetRole(),
		Tags:     hire.GetTags(),
	}

	res, err := SrvCol.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintln("Insert error: ", err))
	}

	objID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, fmt.Sprintln("ObjectID conversion error"))
	}
	data.ID = objID

	return &hirepb.CreateHireResponse{
		Data: model.DataToHirepb(&data),
	}, nil
}

// FindHire first argument is the string pattern which implements the filter, the second
// is the stream channel
func (*HireDataController) FindHire(req *hirepb.FindHireRequest,
	stream hirepb.HireDataService_FindHireServer) error {

	log.Println("Find request received.")

	var filter bson.M
	if req.GetFindPattern() == "duration" {
		filter = bson.M{req.GetFindPattern(): req.GetFindNumber()}
	} else {
		filter = bson.M{req.GetFindPattern(): req.GetFindText()}
	}

	cursor, find_err := SrvCol.Find(context.Background(), filter)
	if find_err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintln("Cannot find hires"))
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		data := &model.HireDataItem{}

		decode_err := cursor.Decode(data)
		if decode_err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintln("Cannot decode hire data"))
		}

		resp := &hirepb.FindHireResponse{
			Data: model.DataToHirepb(data),
		}

		stream.Send(resp)
	}

	if err := cursor.Err(); err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintln("Cursor error", err))
	}

	return nil
}

// ListHires is called to open a stream channel that will send the information to the client
// program
func (*HireDataController) ListHires(req *hirepb.ListHireRequest,
	stream hirepb.HireDataService_ListHiresServer) error {

	log.Println("List hires request received.")

	cursor, find_err := SrvCol.Find(context.Background(), nil)
	if find_err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintln("Cannot find hires"))
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		data := &model.HireDataItem{}

		decode_err := cursor.Decode(data)
		if decode_err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintln("Cannot decode hire data"))
		}

		resp := &hirepb.ListHireResponse{
			Data: model.DataToHirepb(data),
		}

		stream.Send(resp)
	}

	if err := cursor.Err(); err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintln("Cursor error", err))
	}

	return nil
}

// DeleteHire deletes a single document from the DB
func (*HireDataController) DeleteHire(ctx context.Context,
	req *hirepb.DeleteHireRequest) (*hirepb.DeleteHireResponse, error) {

	log.Println("Delete hire request received.")

	filter := bson.M{"name": req.GetHireName()}

	_, derr := SrvCol.DeleteOne(context.Background(), filter)
	if derr != nil {
		return nil, status.Errorf(codes.Internal,
			fmt.Sprintln("Cannot delete from hire: ", derr))
	}

	return &hirepb.DeleteHireResponse{HireName: req.GetHireName()}, nil

}

// UpdateHire takes a single document referenced with his ID to update all his data members.
func (*HireDataController) UpdateHire(ctx context.Context,
	req *hirepb.UpdateHireRequest) (*hirepb.UpdateHireResponse, error) {

	log.Println("Update hire request received.")

	hire := req.GetData()
	hire_id := hire.GetId()

	oid, err := primitive.ObjectIDFromHex(hire_id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintln("Cannot parse hire id"))
	}

	// finc by object id
	filter := bson.M{"_id": oid}
	data := &model.HireDataItem{}

	data.Name = hire.GetName()
	data.Role = hire.GetRole()
	data.Type = hire.GetType()
	data.Duration = hire.GetDuration()
	data.Tags = hire.GetTags()

	_, uerr := SrvCol.ReplaceOne(context.Background(), filter, data)
	if uerr != nil {
		return nil, status.Errorf(codes.Internal,
			fmt.Sprintln("Cannot update object: ", uerr))
	}

	return &hirepb.UpdateHireResponse{
		Data: model.DataToHirepb(data),
	}, nil

}

// FindOneHire returns true if the document was found, false otherwise. It also returns
// the information of the document found
func (*HireDataController) FindOneHire(ctx context.Context,
	req *hirepb.FindOneHireRequest) (*hirepb.FindOneHireResponse, error) {

	log.Println("FindOne hire request received.")

	filter := bson.M{"name": req.GetHireName()}

	output := SrvCol.FindOne(context.Background(), filter)
	data := &model.HireDataItem{}
	derr := output.Decode(data)

	resp := &hirepb.FindOneHireResponse{
		Found: true,
		Data:  model.DataToHirepb(data),
	}

	if derr != nil {
		resp.Found = false
	}

	return resp, nil
}
