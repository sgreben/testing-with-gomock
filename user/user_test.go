package user_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sgreben/testing-with-gomock/match"
	"github.com/sgreben/testing-with-gomock/mocks"
	"github.com/sgreben/testing-with-gomock/user"
)

func TestUse(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDoer := mocks.NewMockDoer(mockCtrl)
	testUser := &user.User{Doer: mockDoer}

	// Expect Do to be called once with 123 and "Hello GoMock" as parameters, and return nil from the mocked call.
	mockDoer.EXPECT().DoSomething(123, "Hello GoMock").Return(nil).Times(1)

	testUser.Use()
}

func TestUseReturnsErrorFromDo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dummyError := errors.New("dummy error")
	mockDoer := mocks.NewMockDoer(mockCtrl)
	testUser := &user.User{Doer: mockDoer}

	// Expect Do to be called once with 123 and "Hello GoMock" as parameters, and return dummyError from the mocked call.
	mockDoer.EXPECT().DoSomething(123, "Hello GoMock").Return(dummyError).Times(1)

	err := testUser.Use()

	if err != dummyError {
		t.Fail()
	}
}

func TestUseMatchersExample1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDoer := mocks.NewMockDoer(mockCtrl)
	testUser := &user.User{Doer: mockDoer}

	mockDoer.EXPECT().DoSomething(gomock.Any(), "Hello GoMock")

	testUser.Use()
}

func TestUseMatchersExample2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDoer := mocks.NewMockDoer(mockCtrl)
	testUser := &user.User{Doer: mockDoer}

	mockDoer.EXPECT().
		DoSomething(123, match.OfType("string")).
		Return(nil).
		Times(1)

	testUser.Use()
}

func TestUseOrderExample1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDoer := mocks.NewMockDoer(mockCtrl)

	callFirst := mockDoer.EXPECT().DoSomething(1, "first this")
	mockDoer.EXPECT().DoSomething(2, "then this").After(callFirst)
	mockDoer.EXPECT().DoSomething(2, "or this").After(callFirst)

	mockDoer.DoSomething(1, "first this")
	mockDoer.DoSomething(2, "or this")
	mockDoer.DoSomething(2, "then this")
}

func TestUseOrderExample2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDoer := mocks.NewMockDoer(mockCtrl)

	gomock.InOrder(
		mockDoer.EXPECT().DoSomething(1, "first this"),
		mockDoer.EXPECT().DoSomething(2, "then this"),
		mockDoer.EXPECT().DoSomething(3, "then this"),
		mockDoer.EXPECT().DoSomething(4, "finally this"),
	)

	mockDoer.DoSomething(1, "first this")
	mockDoer.DoSomething(2, "then this")
	mockDoer.DoSomething(3, "then this")
	mockDoer.DoSomething(4, "finally this")
}

func TestUseActionExamples(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDoer := mocks.NewMockDoer(mockCtrl)

	mockDoer.EXPECT().
		DoSomething(gomock.Any(), gomock.Any()).
		Return(nil).
		Do(func(x int, y string) {
			fmt.Println("Called with x =", x, "and y =", y)
		})

	mockDoer.EXPECT().
		DoSomething(gomock.Any(), gomock.Any()).
		Return(nil).
		Do(func(x int, y string) {
			if x > len(y) {
				t.Fail()
			}
		})

	mockDoer.DoSomething(123, "test")
	mockDoer.DoSomething(2, "test")

}
