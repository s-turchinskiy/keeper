package client

import (
	"github.com/golang/mock/gomock"
	"github.com/s-turchinskiy/keeper/internal/functional_tests"
	"testing"
)

func TestFunctionalGRPC(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	functional_tests.FunctionalTestGRPC(t, functional_tests.UserMockRepository(ctrl), functional_tests.SecretMockRepository(ctrl), false)

}
