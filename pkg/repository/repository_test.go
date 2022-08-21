package repository

import (
	"context"
	"testing"
	"time"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/suite"
)

var testSliceSignUpGoodUser []*models.User = []*models.User{
	{
		User_id: uuid.New(),
		User_name: "TestName1",
		Email: "testemail@gmail.com",
		Password: "12345",
		Time: time.Now(),
	},
	{
		User_id: uuid.New(),
		User_name: "TestName2",
		Email: "1testemail@gmail.com",
		Password: "12345",
		Time: time.Now(),
	},
	{
		User_id: uuid.New(),
		User_name: "TestName3",
		Email: "2testemail@gmail.com",
		Password: "12345",
		Time: time.Now(),
	},
	{
		User_id: uuid.New(),
		User_name: "TestName4",
		Email: "3testemail@gmail.com",
		Password: "12345",
		Time: time.Now(),
	},
	{
		User_id: uuid.New(),
		User_name: "TestName5",
		Email: "4testemail@gmail.com",
		Password: "12345",
		Time: time.Now(),
	},
}

type PostgresSuit struct {
	suite.Suite
	repository Repository
	conn *pgx.Conn
}

func (p *PostgresSuit) SetupSuite() {
	db, err := pgxmock.NewConn()
	if err != nil {
		p.FailNow(err.Error())
	} 
	p.conn = db.Conn()
	p.repository = NewRepository(db.Conn())
}

func (p *PostgresSuit) TearDownSuite() {
	p.conn.Close(context.Background())
}

func (p *PostgresSuit) TestRepositorySignUpGoodUser() {
	for _, v := range testSliceSignUpGoodUser {
		res := p.repository.SignUpUser(v)
		p.Nil(res)
	}
}

func TestRepo(t *testing.T) {
	suite.Run(t, new(PostgresSuit))
}