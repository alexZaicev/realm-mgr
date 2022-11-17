package releaserealm

import (
	"context"
	"fmt"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/alexZaicev/realm-mgr/internal/domain/entities"
	realm_mgr_v1 "github.com/alexZaicev/realm-mgr/proto/go/realm_mgr/v1"
	"github.com/alexZaicev/realm-mgr/tests/functional/utils"
)

func TestRealmManagerReleaseRealmGRPCSuite(t *testing.T) {
	testSuite := NewReleaseRealmTestSuite(t)
	suite.Run(t, testSuite)
}

type ReleaseRealmTestSuite struct {
	suite.Suite

	db     *utils.DB
	client realm_mgr_v1.RealmManagerServiceClient

	activeRealm   *entities.Realm
	activeRealmID uuid.UUID

	draftRealm   *entities.Realm
	draftRealmID uuid.UUID

	disabledRealm   *entities.Realm
	disabledRealmID uuid.UUID

	deletedRealm   *entities.Realm
	deletedRealmID uuid.UUID
}

func NewReleaseRealmTestSuite(t *testing.T) *ReleaseRealmTestSuite {
	cfg, err := utils.LoadConfig()
	require.NoError(t, err, "error loading configuration file")

	db, err := utils.NewDB(cfg)
	require.NoError(t, err, "error creating database connection")

	client, err := utils.NewRealmManagerGRPCClient(cfg)
	require.NoError(t, err, "error creating gRPC client")

	return &ReleaseRealmTestSuite{
		db:     db,
		client: client,
	}
}

func (s *ReleaseRealmTestSuite) SetupSuite() {
	require.NoError(s.T(), s.populateTestData(), "error populating test data")

	// get realm IDs for later tests
	realms, err := s.db.GetRealms(utils.GetRealmsQuery())
	require.NoError(s.T(), err, "error listing realms from database")

	for _, realm := range realms {
		switch realm.Status {
		case entities.StatusActive:
			if s.activeRealm == nil {
				s.activeRealm = realm
				s.activeRealmID = realm.ID
			}
		case entities.StatusDraft:
			if s.draftRealm == nil {
				s.draftRealm = realm
				s.draftRealmID = realm.ID
			}
		case entities.StatusDisabled:
			if s.disabledRealm == nil {
				s.disabledRealm = realm
				s.disabledRealmID = realm.ID
			}
		case entities.StatusDeleted:
			if s.deletedRealm == nil {
				s.deletedRealm = realm
				s.deletedRealmID = realm.ID
			}
		default:
			require.Fail(s.T(), "unhandled entity status", "status: %d", realm.Status)
		}
	}
}

func (s *ReleaseRealmTestSuite) TearDownSuite() {
	err := s.db.Wipe()
	require.NoError(s.T(), err, "error wiping database")
}

func (s *ReleaseRealmTestSuite) Test_ReleaseRealm_Success() {
	testCases := []struct {
		name             string
		realmID          uuid.UUID
		expectedResponse *realm_mgr_v1.Realm
		skip             bool
	}{
		{
			name:             "release draft realm",
			realmID:          s.draftRealmID,
			expectedResponse: s.realmToGRPC(s.draftRealm),
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			if tc.skip {
				s.T().Skip("environment not setup for this test case")
			}

			// arrange
			ctx, err := utils.MakeGRPCRequestContext(context.Background())
			require.NoError(t, err)

			// act
			res, err := s.client.ReleaseRealm(ctx, &realm_mgr_v1.ReleaseRealmRequest{
				Id: tc.realmID.String(),
			})

			// assert
			assert.NoError(t, err)

			require.NotNil(t, res)
			require.NotNil(t, res.GetRealm())

			assert.Equal(t, tc.expectedResponse.Id, res.GetRealm().Id)
			assert.Equal(t, tc.expectedResponse.Name, res.GetRealm().Name)
			assert.Equal(t, tc.expectedResponse.Description, res.GetRealm().Description)
			assert.Equal(t, tc.expectedResponse.Status, res.GetRealm().Status)
		})
	}
}

func (s *ReleaseRealmTestSuite) Test_ReleaseRealm_InvalidArgument() {
	testCases := []struct {
		name           string
		req            *realm_mgr_v1.ReleaseRealmRequest
		expectedErrMsg string
	}{
		{
			name:           "no ID provided",
			req:            &realm_mgr_v1.ReleaseRealmRequest{},
			expectedErrMsg: "invalid ReleaseRealmRequest.Id: value must be a valid UUID | caused by: invalid uuid format",
		},
		{
			name: "malformed ID provided",
			req: &realm_mgr_v1.ReleaseRealmRequest{
				Id: "not-valid-uuid",
			},
			expectedErrMsg: "invalid ReleaseRealmRequest.Id: value must be a valid UUID | caused by: invalid uuid format",
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			// arrange
			ctx, err := utils.MakeGRPCRequestContext(context.Background())
			require.NoError(t, err)

			// act
			res, err := s.client.ReleaseRealm(ctx, tc.req)

			// assert
			assert.Nil(t, res)

			require.Error(t, err)

			gRPCError, ok := status.FromError(err)
			require.True(t, ok)

			assert.Equal(t, codes.InvalidArgument, gRPCError.Code())
			assert.Equal(t, tc.expectedErrMsg, gRPCError.Message())
		})
	}
}

func (s *ReleaseRealmTestSuite) Test_ReleaseRealm_NotFound() {
	testCases := []struct {
		name string
		req  *realm_mgr_v1.ReleaseRealmRequest
		skip bool
	}{
		{
			name: "non-existing realm",
			req: &realm_mgr_v1.ReleaseRealmRequest{
				Id: uuid.New().String(),
			},
		},
		{
			name: "active realm",
			req: &realm_mgr_v1.ReleaseRealmRequest{
				Id: s.activeRealmID.String(),
			},
			skip: s.activeRealm == nil,
		},
		{
			name: "disabled realm",
			req: &realm_mgr_v1.ReleaseRealmRequest{
				Id: s.disabledRealmID.String(),
			},
			skip: s.disabledRealm == nil,
		},
		{
			name: "deleted realm",
			req: &realm_mgr_v1.ReleaseRealmRequest{
				Id: s.deletedRealmID.String(),
			},
			skip: s.deletedRealm == nil,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			if tc.skip {
				s.T().Skip("environment not setup for this test case")
			}

			// arrange
			ctx, err := utils.MakeGRPCRequestContext(context.Background())
			require.NoError(t, err)

			// act
			res, err := s.client.ReleaseRealm(ctx, tc.req)

			// assert
			assert.Nil(t, res)

			require.Error(t, err)

			gRPCError, ok := status.FromError(err)
			require.True(t, ok)

			assert.Equal(t, codes.NotFound, gRPCError.Code())
			assert.Equal(t, fmt.Sprintf("no releasable realm with ID found: %s", tc.req.Id), gRPCError.Message())
		})
	}
}

func (s *ReleaseRealmTestSuite) populateTestData() error {
	realms := []entities.Realm{
		{
			ID:          uuid.New(),
			Name:        "Test Realm 1",
			Description: "Functional test realm #1",
			Status:      entities.StatusActive,
			CreatedAt:   time.Date(2022, 01, 01, 12, 30, 30, 0, time.UTC),
			UpdatedAt:   time.Date(2022, 04, 01, 13, 30, 30, 0, time.UTC),
		},
		{
			ID:          uuid.New(),
			Name:        "Test Realm 1",
			Description: "Functional test realm #1 for unexpected cases",
			Status:      entities.StatusDraft,
			CreatedAt:   time.Date(2022, 01, 01, 12, 30, 30, 0, time.UTC),
			UpdatedAt:   time.Date(2022, 01, 14, 12, 30, 30, 0, time.UTC),
		},
		{
			ID:          uuid.New(),
			Name:        "Test Realm 2",
			Description: "Functional test realm #2",
			Status:      entities.StatusDisabled,
			CreatedAt:   time.Date(2022, 01, 01, 12, 30, 30, 0, time.UTC),
			UpdatedAt:   time.Date(2022, 01, 14, 12, 30, 30, 0, time.UTC),
		},
		{
			ID:          uuid.New(),
			Name:        "Test Realm 3",
			Description: "Functional test realm #3",
			Status:      entities.StatusDeleted,
			CreatedAt:   time.Date(2022, 10, 01, 12, 30, 30, 0, time.UTC),
			UpdatedAt:   time.Date(2022, 01, 14, 12, 30, 30, 0, time.UTC),
			DeletedAt:   time.Date(2022, 02, 14, 12, 30, 30, 0, time.UTC),
		},
	}

	queries, err := utils.GenerateRealmInsertQueries(realms...)
	if err != nil {
		return err
	}

	_, err = s.db.ExecuteInsertQueries(context.Background(), queries...)
	if err != nil {
		return err
	}

	return nil
}

func (s *ReleaseRealmTestSuite) realmToGRPC(realm *entities.Realm) *realm_mgr_v1.Realm {
	if realm == nil {
		return nil
	}

	return &realm_mgr_v1.Realm{
		Id:          realm.ID.String(),
		Name:        realm.Name,
		Description: realm.Description,
		Status:      realm_mgr_v1.EnumStatus_ENUM_STATUS_ACTIVE,
	}
}
