package pipeline

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestNewMongoPipeline(t *testing.T) {
	ctx := context.Background()
	mongoPipeline, err := NewMongoPipeline(ctx, "127.0.0.1:27017", "douban")
	if err != nil {
		t.Error(err)
	}
	if res, err := mongoPipeline.database.Collection("test").
		InsertOne(ctx, bson.M{"id": 1, "name": "name", "chinese": "中文"}); err != nil {
		t.Error(err)
	} else {
		fmt.Println(res.InsertedID)
	}
}
