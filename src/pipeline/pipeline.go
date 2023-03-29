package pipeline

import (
	"context"
	"github.com/gogodjzhu/gogoscrapy/src/entity"
	"github.com/gogodjzhu/gogoscrapy/src/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

var LOG = utils.NewLogger()

type IPipeline interface {
	Process(items entity.IResultItems) error
}

type ConsolePipeline struct {
}

func NewConsolePipeline() ConsolePipeline {
	return ConsolePipeline{}
}

func (ConsolePipeline) Process(items entity.IResultItems) error {
	LOG.Infof("items :%+v", items)
	return nil
}

type MongoPipeline struct {
	database *mongo.Database
}

func NewMongoPipeline(ctx context.Context, address string, dbName string) (*MongoPipeline, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+address))
	if err != nil {
		return nil, err
	}
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	database := client.Database(dbName)
	return &MongoPipeline{database: database}, nil
}

func (this *MongoPipeline) Process(items entity.IResultItems) error {
	collectionName := items.Get("tbl")
	if collectionName == nil {
		collectionName = "gogoscrapy_" + time.Now().Format("2006-01-02")
	}
	item := bson.M{}
	for k, v := range items.All() {
		if k != "tbl" {
			item[k] = v
		}
	}
	item["create_time"] = time.Now().Format("2006-01-02")
	collection := this.database.Collection(collectionName.(string))
	if _, err := collection.InsertOne(nil, item); err != nil {
		return err
	}
	return nil
}
