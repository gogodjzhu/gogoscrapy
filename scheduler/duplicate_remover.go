package scheduler

import (
	"github.com/gogodjzhu/gogoscrapy/entity"
	"github.com/gogodjzhu/gogoscrapy/utils"
	"github.com/gogodjzhu/gogoscrapy/utils/redisUtil"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type DuplicateRemover interface {
	IsDuplicate(request entity.IRequest) (bool, error)
	ResetDuplicate() error
	GetTotalCount() (int, error)
}

type MemDuplicateRemover struct {
	remover *utils.AsyncSet
}

func NewMemDuplicateRemover() *MemDuplicateRemover {
	return &MemDuplicateRemover{
		remover: utils.NewAsyncSet(),
	}
}

func (mdr *MemDuplicateRemover) IsDuplicate(request entity.IRequest) (bool, error) {
	if noNeedToRemoveDuplicate(request) {
		return false, nil
	}
	return !mdr.remover.Add(request.GetUrl()), nil
}

func (mdr *MemDuplicateRemover) ResetDuplicate() error {
	mdr.remover.Clear()
	return nil
}

func (mdr *MemDuplicateRemover) GetTotalCount() (int, error) {
	return mdr.remover.Size(), nil
}

type RedisDuplicateRemover struct {
	pfkey string
	rs    *redisUtil.RedisClient
}

func NewRedisDuplicatedRemover(rs *redisUtil.RedisClient, pfkey string) (*RedisDuplicateRemover, error) {
	return &RedisDuplicateRemover{
		rs:    rs,
		pfkey: pfkey,
	}, nil
}

func (rdr *RedisDuplicateRemover) IsDuplicate(request entity.IRequest) (bool, error) {
	if noNeedToRemoveDuplicate(request) {
		return false, nil
	}
	conn, err := rdr.rs.GetConn()
	if err != nil {
		return false, err
	}
	defer conn.Close()
	res, err := conn.Do("PFADD", rdr.pfkey, request.GetUrl())
	if err != nil {
		return false, nil
	}
	return res == 1, nil
}

func (rdr *RedisDuplicateRemover) ResetDuplicate() error {
	conn, err := rdr.rs.GetConn()
	if err != nil {
		log.Warnf("failed to get redis conn, err:%+v", err)
		return errors.New("failed to get redis conn")
	}
	defer conn.Close()
	_, err = conn.Do("DEL", rdr.pfkey)
	return err
}

func (rdr *RedisDuplicateRemover) GetTotalCount() (int, error) {
	conn, err := rdr.rs.GetConn()
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	res, err := conn.Do("PFCOUNT", rdr.pfkey)
	if err != nil {
		return 0, err
	}
	return res.(int), nil
}

func noNeedToRemoveDuplicate(request entity.IRequest) bool {
	return http.MethodPost == request.GetMethod()
}
