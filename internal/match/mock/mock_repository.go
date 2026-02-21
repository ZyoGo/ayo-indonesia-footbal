// Code generated manually following mockgen patterns.

package mock

import (
	context "context"
	reflect "reflect"

	domain "github.com/ZyoGo/ayo-indonesia-footbal/internal/match/domain"
	gomock "go.uber.org/mock/gomock"
)

// ---------------------------------------------------------------------------
// MockMatchRepository
// ---------------------------------------------------------------------------

type MockMatchRepository struct {
	ctrl     *gomock.Controller
	recorder *MockMatchRepositoryMockRecorder
}

type MockMatchRepositoryMockRecorder struct {
	mock *MockMatchRepository
}

func NewMockMatchRepository(ctrl *gomock.Controller) *MockMatchRepository {
	mock := &MockMatchRepository{ctrl: ctrl}
	mock.recorder = &MockMatchRepositoryMockRecorder{mock}
	return mock
}

func (m *MockMatchRepository) EXPECT() *MockMatchRepositoryMockRecorder {
	return m.recorder
}

func (m *MockMatchRepository) Create(ctx context.Context, match *domain.Match) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, match)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockMatchRepositoryMockRecorder) Create(ctx, match any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockMatchRepository)(nil).Create), ctx, match)
}

func (m *MockMatchRepository) FindByID(ctx context.Context, id string) (*domain.Match, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, id)
	ret0, _ := ret[0].(*domain.Match)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockMatchRepositoryMockRecorder) FindByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockMatchRepository)(nil).FindByID), ctx, id)
}

func (m *MockMatchRepository) FindAll(ctx context.Context) ([]domain.Match, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx)
	ret0, _ := ret[0].([]domain.Match)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockMatchRepositoryMockRecorder) FindAll(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockMatchRepository)(nil).FindAll), ctx)
}

func (m *MockMatchRepository) Update(ctx context.Context, match *domain.Match) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, match)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockMatchRepositoryMockRecorder) Update(ctx, match any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockMatchRepository)(nil).Update), ctx, match)
}

func (m *MockMatchRepository) Delete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockMatchRepositoryMockRecorder) Delete(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockMatchRepository)(nil).Delete), ctx, id)
}

// ---------------------------------------------------------------------------
// MockMatchResultRepository
// ---------------------------------------------------------------------------

type MockMatchResultRepository struct {
	ctrl     *gomock.Controller
	recorder *MockMatchResultRepositoryMockRecorder
}

type MockMatchResultRepositoryMockRecorder struct {
	mock *MockMatchResultRepository
}

func NewMockMatchResultRepository(ctrl *gomock.Controller) *MockMatchResultRepository {
	mock := &MockMatchResultRepository{ctrl: ctrl}
	mock.recorder = &MockMatchResultRepositoryMockRecorder{mock}
	return mock
}

func (m *MockMatchResultRepository) EXPECT() *MockMatchResultRepositoryMockRecorder {
	return m.recorder
}

func (m *MockMatchResultRepository) Create(ctx context.Context, result *domain.MatchResult) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, result)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockMatchResultRepositoryMockRecorder) Create(ctx, result any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockMatchResultRepository)(nil).Create), ctx, result)
}

func (m *MockMatchResultRepository) FindByMatchID(ctx context.Context, matchID string) (*domain.MatchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByMatchID", ctx, matchID)
	ret0, _ := ret[0].(*domain.MatchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockMatchResultRepositoryMockRecorder) FindByMatchID(ctx, matchID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByMatchID", reflect.TypeOf((*MockMatchResultRepository)(nil).FindByMatchID), ctx, matchID)
}

func (m *MockMatchResultRepository) ExistsByMatchID(ctx context.Context, matchID string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExistsByMatchID", ctx, matchID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockMatchResultRepositoryMockRecorder) ExistsByMatchID(ctx, matchID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExistsByMatchID", reflect.TypeOf((*MockMatchResultRepository)(nil).ExistsByMatchID), ctx, matchID)
}

// ---------------------------------------------------------------------------
// MockReportRepository
// ---------------------------------------------------------------------------

type MockReportRepository struct {
	ctrl     *gomock.Controller
	recorder *MockReportRepositoryMockRecorder
}

type MockReportRepositoryMockRecorder struct {
	mock *MockReportRepository
}

func NewMockReportRepository(ctrl *gomock.Controller) *MockReportRepository {
	mock := &MockReportRepository{ctrl: ctrl}
	mock.recorder = &MockReportRepositoryMockRecorder{mock}
	return mock
}

func (m *MockReportRepository) EXPECT() *MockReportRepositoryMockRecorder {
	return m.recorder
}

func (m *MockReportRepository) GetMatchReport(ctx context.Context, matchID string) (*domain.MatchReportView, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMatchReport", ctx, matchID)
	ret0, _ := ret[0].(*domain.MatchReportView)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockReportRepositoryMockRecorder) GetMatchReport(ctx, matchID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMatchReport", reflect.TypeOf((*MockReportRepository)(nil).GetMatchReport), ctx, matchID)
}

func (m *MockReportRepository) GetAllMatchReports(ctx context.Context) ([]domain.MatchReportView, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllMatchReports", ctx)
	ret0, _ := ret[0].([]domain.MatchReportView)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockReportRepositoryMockRecorder) GetAllMatchReports(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllMatchReports", reflect.TypeOf((*MockReportRepository)(nil).GetAllMatchReports), ctx)
}
