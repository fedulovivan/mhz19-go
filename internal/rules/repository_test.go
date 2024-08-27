package rules

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

type RepoSuite struct {
	suite.Suite
}

func (s *RepoSuite) Test10() {
	db, mock, err := sqlmock.New()
	s.NotNil(mock)
	s.Nil(err)
	NewRepository(db)
}

func (s *RepoSuite) Test11() {
	db, mock, err := sqlmock.New()
	s.NotNil(mock)
	s.Nil(err)
	mock.ExpectBegin().WillReturnError(errors.New("test"))
	r := NewRepository(db)
	_, _, _, _, _, err = r.Get()
	s.ErrorContains(err, "test")
}

func (s *RepoSuite) Test20() {
	db, mock, err := sqlmock.New()
	s.NotNil(mock)
	s.Nil(err)
	mock.ExpectBegin()
	mock.MatchExpectationsInOrder(false)
	mock.ExpectQuery(`rules`).WillReturnRows(sqlmock.NewRows([]string{"1", "test mapping", "1", "0"}))
	mock.ExpectQuery(`rule_conditions`).WillReturnRows(sqlmock.NewRows([]string{}))
	mock.ExpectQuery(`rule_actions`).WillReturnRows(sqlmock.NewRows([]string{}))
	mock.ExpectQuery(`rule_condition_or_action_arguments`).WillReturnRows(sqlmock.NewRows([]string{}))
	mock.ExpectQuery(`rule_action_argument_mappings`).WillReturnRows(sqlmock.NewRows([]string{}))
	mock.ExpectCommit()
	r := NewRepository(db)
	_, _, _, _, _, err = r.Get()
	s.Nil(err)
}

func (s *RepoSuite) Test30() {
	db, mock, err := sqlmock.New()
	s.NotNil(mock)
	s.Nil(err)
	mock.ExpectBegin()
	mock.ExpectExec("rules").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	r := NewRepository(db)
	err = r.Create(
		DbRule{},
		[]DbRuleCondition{},
		[]DbRuleAction{},
		[]DbRuleConditionOrActionArgument{},
		[]DbRuleActionArgumentMapping{},
	)
	s.Nil(err)
}

func (s *RepoSuite) Test40() {
	db, mock, err := sqlmock.New()
	s.NotNil(mock)
	s.Nil(err)
	mock.ExpectBegin()
	mock.ExpectExec("rules").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("rule_conditions").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("rule_actions").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("rule_condition_or_action_arguments").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("rule_action_argument_mappings").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	r := NewRepository(db)
	err = r.Create(
		DbRule{},
		[]DbRuleCondition{{}},
		[]DbRuleAction{{}},
		[]DbRuleConditionOrActionArgument{{}},
		[]DbRuleActionArgumentMapping{{}},
	)
	s.Nil(err)
}

func TestRepo(t *testing.T) {
	suite.Run(t, new(RepoSuite))
}
