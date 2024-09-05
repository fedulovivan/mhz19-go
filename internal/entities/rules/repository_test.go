package rules

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/fedulovivan/mhz19-go/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

type RepoSuite struct {
	suite.Suite
}

func (s *RepoSuite) Test10() {
	mdatabase, mock, err := sqlmock.New()
	s.NotNil(mock)
	s.Nil(err)
	NewRepository(mdatabase)
}

func (s *RepoSuite) Test11() {
	mdatabase, mock, err := sqlmock.New()
	s.NotNil(mock)
	s.Nil(err)
	mock.ExpectBegin().WillReturnError(errors.New("test"))
	r := NewRepository(mdatabase)
	_, _, _, _, _, err = r.Get(sql.NullInt32{})
	s.ErrorContains(err, "test")
}

func (s *RepoSuite) Test20() {
	mdatabase, mock, err := sqlmock.New()
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
	r := NewRepository(mdatabase)
	_, _, _, _, _, err = r.Get(sql.NullInt32{})
	s.Nil(err)
}

func (s *RepoSuite) Test21() {
	mdatabase, mock, err := sqlmock.New()
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
	r := NewRepository(mdatabase)
	_, _, _, _, _, err = r.Get(db.NewNullInt32(123))
	s.Nil(err)
}

func (s *RepoSuite) Test30() {
	mdatabase, mock, err := sqlmock.New()
	s.NotNil(mock)
	s.Nil(err)
	mock.ExpectBegin()
	mock.ExpectExec("rules").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	r := NewRepository(mdatabase)
	_, err = r.Create(
		DbRule{},
		[]DbRuleCondition{},
		[]DbRuleAction{},
		[]DbRuleConditionOrActionArgument{},
		[]DbRuleActionArgumentMapping{},
	)
	s.Nil(err)
}

func (s *RepoSuite) Test40() {
	mdatabase, mock, err := sqlmock.New()
	s.NotNil(mock)
	s.Nil(err)
	mock.ExpectBegin()
	mock.ExpectExec("rules").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("rule_conditions").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("rule_actions").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("rule_condition_or_action_arguments").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("rule_action_argument_mappings").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	r := NewRepository(mdatabase)
	_, err = r.Create(
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
