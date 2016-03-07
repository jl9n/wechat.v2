package record

import (
	"github.com/chanxuehong/wechat/mp/core"
)

// RecordIterator
//
//  iter, err := NewRecordIterator(clt, request)
//  if err != nil {
//      // TODO: 增加你的代码
//  }
//
//  for iter.HasNext() {
//      records, err := iter.NextPage()
//      if err != nil {
//          // TODO: 增加你的代码
//      }
//      // TODO: 增加你的代码
//  }
type RecordIterator struct {
	clt *core.Client // 关联的微信 Client

	nextGetRecordRequest *GetRecordRequest // 上一次查询的 request

	lastGetRecordResult []Record // 上一次查询的 result
	nextPageCalled      bool     // NextPage() 是否调用过
}

func (iter *RecordIterator) HasNext() bool {
	if !iter.nextPageCalled { // 第一次调用需要特殊对待
		return len(iter.lastGetRecordResult) > 0
	}

	// 如果上一次读取的数据等于 PageSize, 则"可能"还有数据; 否则肯定是没有数据了.
	return len(iter.lastGetRecordResult) == iter.nextGetRecordRequest.PageSize
}

func (iter *RecordIterator) NextPage() (records []Record, err error) {
	if !iter.nextPageCalled { // 第一次调用需要特殊对待
		iter.nextPageCalled = true

		records = iter.lastGetRecordResult
		return
	}

	records, err = GetRecord(iter.clt, iter.nextGetRecordRequest)
	if err != nil {
		return
	}

	iter.nextGetRecordRequest.PageIndex++
	iter.lastGetRecordResult = records
	return
}

// 获取聊天记录遍历器.
func NewRecordIterator(clt *core.Client, request *GetRecordRequest) (iter *RecordIterator, err error) {
	// 逻辑上相当于第一次调用 RecordIterator.NextPage, 因为第一次调用 RecordIterator.HasNext 需要数据支撑, 所以提前获取了数据

	records, err := GetRecord(clt, request)
	if err != nil {
		return
	}

	request.PageIndex++

	iter = &RecordIterator{
		clt:                  clt,
		nextGetRecordRequest: request,
		lastGetRecordResult:  records,
		nextPageCalled:       false,
	}
	return
}