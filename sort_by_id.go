package sortby

import (
	"errors"
	"reflect"

	"github.com/globalsign/mgo/bson"
)

// MongoID will sort pointer slice of struct based on field name. if struct doesn't have field name will sort by Id
func MongoID(ps ParamMongoID) error {
	resultv := reflect.ValueOf(ps.Value)
	if resultv.Kind() != reflect.Ptr || resultv.Elem().Kind() != reflect.Slice {
		return errors.New("sort value not pointer")
	}
	exec := byID(ps)
	return executeSort(&exec)
}

func (s *byID) SetIndicatorIndex() {
	indicatorIndex := make(map[bson.ObjectId]int)
	for i, v := range s.Indicator {
		indicatorIndex[v] = i
	}
	s.indicatorIndex = indicatorIndex
}

func (s *byID) CheckFieldName() error {
	if s.FieldName == "" {
		s.FieldName = "Id"
	}
	v := reflect.ValueOf(s.Value).Elem()
	getField := v.Index(0).FieldByName(s.FieldName)
	if !getField.IsValid() {
		return errors.New("field not found")
	}
	return nil
}

func (s *byID) Len() int {
	return reflect.ValueOf(s.Value).Elem().Len()
}

func (s *byID) Less(i, j int) bool {
	v := reflect.ValueOf(s.Value).Elem()
	indexI := s.indicatorIndex[v.Index(i).FieldByName(s.FieldName).Interface().(bson.ObjectId)]
	indexJ := s.indicatorIndex[v.Index(j).FieldByName(s.FieldName).Interface().(bson.ObjectId)]
	if s.Reverse {
		return indexI > indexJ
	}
	return indexI < indexJ
}

func (s *byID) Swap(i, j int) {
	v := reflect.ValueOf(s.Value).Elem()
	tempI := v.Index(i)
	tempJ := v.Index(j)
	temp := tempI.Interface()
	v.Index(i).Set(tempJ)
	v.Index(j).Set(reflect.ValueOf(temp))
}
