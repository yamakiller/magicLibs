package dbs

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//MongoDB doc
//@Summary Mongo DB Object
//@Struct MongoDB
//@Member *mongo.Client mongo client
//@Member *mongo.Database mongo database object
//@Member int  connection/operation/close time out unit/second
type MongoDB struct {
	_c       *mongo.Client
	_opt     *options.ClientOptions
	_db      *mongo.Database
	_timeOut time.Duration
}

//Initial doc
//@Summary initialization mongo db
//@Method Initial
//@Param string   uri
//@Param string   user
//@Param string   password
//@Param string   db name
//@Param int      max pool
//@Param int      min pool
//@Param int      connection/operation/close time out unit/millsecond
//@Param int      connection heart reate unit/second
//@Param int      connection time out unit/second
//@Param int      connection idle time unit/second
//@Return (error) initialization fail return error
func (slf *MongoDB) Initial(url string,
	user string,
	password string,
	dbname string,
	max int,
	min int,
	timeout int) error {

	mongoUrl := "mongodb://" + user + ":" + password + "@" + url

	slf._timeOut = time.Duration(timeout) * time.Second
	slf._opt = options.Client().ApplyURI(mongoUrl)
	ctx, cancel := context.WithTimeout(context.Background(), slf._timeOut)
	defer cancel()

	slf._opt.SetMaxPoolSize(uint64(max))
	slf._opt.SetMinPoolSize(uint64(min))

	c, err := mongo.Connect(ctx, slf._opt)
	if err != nil {
		return err
	}

	if err := c.Ping(context.TODO(), nil); err != nil {
		c.Disconnect(context.TODO())
		return err
	}

	slf._db = c.Database(dbname)
	if slf._db == nil {
		c.Disconnect(context.TODO())
		return fmt.Errorf("mongoDB Database %s does not exist", dbname)
	}
	slf._c = c

	return nil
}

//Close doc
//@Summary close mongo db
//@Method Close
func (slf *MongoDB) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), slf._timeOut)
	defer cancel()
	slf._c.Disconnect(ctx)
	slf._c = nil
	slf._db = nil
}

//InsertOne doc
//@Method InsertOne @Summary Insert a piece of data
//@Param (string) set/table name
//@Param (interface{}) data
//@Return (interface{}) insert result
//@Return (error) insert fail
func (slf *MongoDB) InsertOne(name string, document interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), slf._timeOut)
	defer cancel()

	r, rerr := slf._db.Collection(name).InsertOne(ctx, document)
	if rerr != nil {
		return nil, rerr
	}

	return r.InsertedID, nil
}

//InsertMany doc
//@Summary Insert multiple pieces of data
//@Method InsertMany
//@Param (string) set/table name
//@Param ([]interface{}) more data
//@Return (interface{}) insert result
//@Return (error) insert fail
func (slf *MongoDB) InsertMany(name string, document []interface{}) ([]interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), slf._timeOut)
	defer cancel()

	r, rerr := slf._db.Collection(name).InsertMany(ctx, document)
	if rerr != nil {
		return nil, rerr
	}
	return r.InsertedIDs, nil
}

//FindOne doc
//@Summary Query a piece of data
//@Method FindOne
//@Param (string) set/table name
//@Param (interface{}) filter
//@Param (interface{}) out result
//@Return (error) Return error
func (slf *MongoDB) FindOne(name string, filter interface{}, out interface{}) error {

	ctx, cancel := context.WithTimeout(context.Background(), slf._timeOut)
	defer cancel()

	r := slf._db.Collection(name).FindOne(ctx, filter)
	if derr := r.Decode(out); derr != nil {
		return derr
	}

	return nil
}

//Find doc
//@Summary Query multiple data
//@Method Find
//@Param (string) set/table name
//@Param (interface{}) filter
//@Param (interface{})
//@Return ([]interface{}) Return result
//@Return (error) Return error
func (slf *MongoDB) Find(name string, filter interface{}, decode interface{}) ([]interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), slf._timeOut)
	defer cancel()

	r, rerr := slf._db.Collection(name).Find(ctx, filter)
	if rerr != nil {
		return nil, rerr
	}

	defer r.Close(ctx)
	ary := make([]interface{}, 0, 4)
	for r.Next(ctx) {
		if derr := r.Decode(&decode); derr != nil {
			return nil, derr
		}

		ary = append(ary, decode)
	}

	return ary, nil
}

//UpdateOne doc
//@Summary update a piece of data
//@Method UpdateOne
//@Param (string) set/table name
//@Param (interface{}) filter
//@Param (interface{}) update informat
//@Return (int64) match of number
//@Return (int64) modify of number
//@Return (int64) update of number
//@Return (interface{}) update id
//@Return (error)
func (slf *MongoDB) UpdateOne(name string, filter interface{}, update interface{}) (int64, int64, int64, interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), slf._timeOut)
	defer cancel()

	r, rerr := slf._db.Collection(name).UpdateOne(ctx, filter, update)
	if rerr != nil {
		return 0, 0, 0, nil, rerr
	}

	return r.MatchedCount, r.ModifiedCount, r.UpsertedCount, r.UpsertedID, nil
}

func (slf *MongoDB) UpdateOrInsert(name string, filter interface{}, update interface{}) (int64, int64, int64, interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), slf._timeOut)
	defer cancel()

	updateOpts := options.Update().SetUpsert(true)
	r, rerr := slf._db.Collection(name).UpdateOne(ctx, filter, update, updateOpts)
	if rerr != nil {
		return 0, 0, 0, nil, rerr
	}

	return r.MatchedCount, r.ModifiedCount, r.UpsertedCount, r.UpsertedID, nil
}

//UpdateMany doc
//@Summary update multiple data
//@Method UpdateMany
//@Param (string) set/table name
//@Param (interface{}) filter
//@Param (interface{}) update informat
//@Return (int64) match of number
//@Return (int64) modify of number
//@Return (int64) update of number
//@Return (interface{}) update id
//@Return (error)
func (slf *MongoDB) UpdateMany(name string, filter interface{}, update interface{}) (int64, int64, int64, interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), slf._timeOut)
	defer cancel()

	r, rerr := slf._db.Collection(name).UpdateMany(ctx, filter, update)
	if rerr != nil {
		return 0, 0, 0, nil, rerr
	}

	return r.MatchedCount, r.ModifiedCount, r.UpsertedCount, r.UpsertedID, nil
}

//ReplaceOne doc
//@Summary replace a piece of data
//@Method ReplaceOne
//@Param (string) set/table name
//@Param (interface{}) filter
//@Param (interface{}) update informat
//@Return (int64) match of number
//@Return (int64) modify of number
//@Return (int64) update of number
//@Return (interface{}) update id
//@Return (error)
func (slf *MongoDB) ReplaceOne(name string, filter interface{}, replacement interface{}) (int64, int64, int64, interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), slf._timeOut)
	defer cancel()

	r, rerr := slf._db.Collection(name).ReplaceOne(ctx, filter, replacement)

	if rerr != nil {
		return 0, 0, 0, nil, rerr
	}

	return r.MatchedCount, r.ModifiedCount, r.UpsertedCount, r.UpsertedID, nil
}

//DeleteOne doc
//@Summary delete a piece of data
//@Method DeleteOne
//@Param (string) set/table name
//@Param (interface{}) filter
//@Return (int64) delete of number
//@Return (error)
func (slf *MongoDB) DeleteOne(name string, filter interface{}) (int64, error) {

	ctx, cancel := context.WithTimeout(context.Background(), slf._timeOut)
	defer cancel()

	r, rerr := slf._db.Collection(name).DeleteOne(ctx, filter)
	if rerr != nil {
		return 0, rerr
	}

	return r.DeletedCount, nil
}

//DeleteMany doc
//@Summary Delete multiple pieces of data
//@Method DeleteMany
//@Param (string) set/table name
//@Param (interface{}) filter
//@Return (int64) delete of number
//@Return (error)
func (slf *MongoDB) DeleteMany(name string, filter interface{}) (int64, error) {

	ctx, cancel := context.WithTimeout(context.Background(), slf._timeOut)
	defer cancel()

	r, rerr := slf._db.Collection(name).DeleteMany(ctx, filter)
	if rerr != nil {
		return 0, rerr
	}

	return r.DeletedCount, nil
}

//FindOneAndDelete doc
//@Summary find a piece of data and delete
//@Method FindOneAndDelete
//@Param (string) set/table name
//@Param (interface{}) filter
//@Param (out interface{}) One piece of data found
//@Return (error)
func (slf *MongoDB) FindOneAndDelete(name string, filter interface{}, out interface{}) error {

	ctx, cancel := context.WithTimeout(context.Background(), slf._timeOut)
	defer cancel()

	r := slf._db.Collection(name).FindOneAndDelete(ctx, filter)

	if derr := r.Decode(out); derr != nil {
		return derr
	}

	return nil
}

//FindOneAndUpdate doc
//@Summary find a piece of data and update
//@Method FindOneAndUpdate
//@Param (string) set/table name
//@Param (interface{}) filter
//@Param (interface{}) data to be updated
//@Param (out interface{}) One piece of data found
//@Return (error)
func (slf *MongoDB) FindOneAndUpdate(name string, filter interface{}, update interface{}, out interface{}) error {

	ctx, cancel := context.WithTimeout(context.Background(), slf._timeOut)
	defer cancel()

	r := slf._db.Collection(name).FindOneAndUpdate(ctx, filter, update)
	if derr := r.Decode(out); derr != nil {
		return derr
	}

	return nil
}

//FindOneAndReplace doc
//@Summary find a piece of data and replace
//@Method FindOneAndReplace
//@Param (string) set/table name
//@Param (interface{}) filter
//@Param (interface{}) data to be replace
//@Param (out interface{}) One piece of data found
//@Return (error)
func (slf *MongoDB) FindOneAndReplace(name string, filter interface{}, replacement interface{}, out interface{}) error {

	ctx, cancel := context.WithTimeout(context.Background(), slf._timeOut)
	defer cancel()

	r := slf._db.Collection(name).FindOneAndReplace(ctx, filter, replacement)
	if derr := r.Decode(out); derr != nil {
		return derr
	}

	return nil
}

//Distinct doc
//@Summary Find in the specified field
//@Method Distinct
//@Param (string) set/table name
//@Param (string) field name
//@Param (interface{}) filter
//@Return ([]interface{}) Return result
//@Return (error)
func (slf *MongoDB) Distinct(name string, fieldName string, filter interface{}) ([]interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), slf._timeOut)
	defer cancel()

	r, rerr := slf._db.Collection(name).Distinct(ctx, fieldName, filter)
	if rerr != nil {
		return nil, rerr
	}

	return r, nil
}

func (slf *MongoDB) SpawnIndex(name string, idx mongo.IndexModel) error {

	ctx, cancel := context.WithTimeout(context.Background(), slf._timeOut)
	defer cancel()

	if _, err := slf._db.Collection(name).Indexes().CreateOne(ctx, idx); err != nil {
		return err
	}

	return nil
}

//Drop @Summary
//@Summary Delete set/table
//@Method Drop
//@Param  (string) set/table name
//@Return (error)
func (slf *MongoDB) Drop(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), slf._timeOut)
	defer cancel()

	return slf._db.Collection(name).Drop(ctx)
}

//CountDocuments @Summary
//@Summary Return the total number of documents
//@Method CountDocuments
//@Param (string) set/table name
//@Param (interface{}) filter
//@Return (int64) a number
//@Return (error)
func (slf *MongoDB) CountDocuments(name string, filter interface{}) (int64, error) {

	ctx, cancel := context.WithTimeout(context.Background(), slf._timeOut)
	defer cancel()

	return slf._db.Collection(name).CountDocuments(ctx, filter)
}
