package oraclerepo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tijanadmi/ddn_rdc/models"
	"github.com/tijanadmi/ddn_rdc/util"
)



func TestGetUser(t *testing.T) {
	testUser :=models.User{
		Username: "test",
		Password: "test",
	}
	user1, err := testRepo.GetUserByUsername(context.Background(), testUser.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user1)
    
	require.Equal(t, user1.Username, testUser.Username)

	hashedPassword, _ := util.HashPassword(testUser.Password)
	
	// if err != nil {
	// 	return 0, err
	// }
	//hashPassword,_:=util.HashPassword(testUser.Password)
	require.NotEqual(t, user1.Password, string(hashedPassword))
	//err=util.CheckPassword(user1.Password, string(hashedPassword))
	//require.NoError(t, err)
}