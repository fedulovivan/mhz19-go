package db

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/sync/errgroup"

	"github.com/stretchr/testify/suite"
)

type DbSuite struct {
	suite.Suite
}

func (s *DbSuite) SetupSuite() {
}

func (s *DbSuite) TeardownSuite() {
}

func (s *DbSuite) Test10() {
	mdatabase, mock, err := sqlmock.New()
	s.NotNil(mock)
	s.Nil(err)
	mock.ExpectBegin()
	mock.ExpectCommit()
	start := time.Now()
	err = RunTx(mdatabase, func(ctx CtxEnhanced) error {
		g, _ := errgroup.WithContext(ctx)
		g.Go(func() (e error) {
			time.Sleep(time.Millisecond * 50)
			return
		})
		g.Go(func() (e error) {
			time.Sleep(time.Millisecond * 60)
			return
		})
		g.Go(func() (e error) {
			time.Sleep(time.Millisecond * 70)
			return
		})
		return g.Wait()
	})
	elapsed := time.Since(start)
	s.Greater(elapsed, time.Millisecond*70)
	s.Less(elapsed, time.Millisecond*100)
	s.Nil(err)
}

func (s *DbSuite) Test20() {
	mdatabase, mock, err := sqlmock.New()
	s.NotNil(mock)
	s.Nil(err)
	mock.ExpectBegin()
	mock.ExpectRollback()
	err = RunTx(mdatabase, func(ctx CtxEnhanced) error {
		return errors.New("some error")
	})
	s.ErrorContains(err, "some error")
}

func (s *DbSuite) Test30() {
	mdatabase, mock, err := sqlmock.New()
	s.NotNil(mock)
	s.Nil(err)
	mock.ExpectBegin().WillReturnError(errors.New("begin failed"))
	mock.ExpectRollback()
	err = RunTx(mdatabase, func(ctx CtxEnhanced) error {
		return nil
	})
	s.ErrorContains(err, "begin failed")
}

func TestDb(t *testing.T) {
	suite.Run(t, new(DbSuite))
}
