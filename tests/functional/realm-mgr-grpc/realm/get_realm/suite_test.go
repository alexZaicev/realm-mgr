package get_realm

import (
	"context"
	"fmt"
	"github.com/alexZaicev/realm-mgr/internal/adapters/realmmgrgrpc/models"
	"github.com/alexZaicev/realm-mgr/internal/domain/entities"
	realm_mgr_v1 "github.com/alexZaicev/realm-mgr/proto/go/realm_mgr/v1"
	"github.com/alexZaicev/realm-mgr/tests/functional/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

func TestRealmManagerGetRealmGRPCSuite(t *testing.T) {
	testSuite := NewGetRealmTestSuite(t)
	suite.Run(t, testSuite)
}

type GetRealmTestSuite struct {
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

func NewGetRealmTestSuite(t *testing.T) *GetRealmTestSuite {
	cfg, err := utils.LoadConfig()
	require.NoError(t, err, "error loading configuration file")

	db, err := utils.NewDB(cfg)
	require.NoError(t, err, "error creating database connection")

	client, err := utils.NewRealmManagerGRPCClient(cfg)
	require.NoError(t, err, "error creating gRPC client")

	return &GetRealmTestSuite{
		db:     db,
		client: client,
	}
}

func (s *GetRealmTestSuite) SetupSuite() {
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

func (s *GetRealmTestSuite) TearDownSuite() {
	err := s.db.Wipe()
	require.NoError(s.T(), err, "error wiping database")
}

func (s *GetRealmTestSuite) Test_GetRealm_Success() {
	testCases := []struct {
		name             string
		req              *realm_mgr_v1.GetRealmRequest
		expectedResponse *realm_mgr_v1.Realm
		skip             bool
	}{
		{
			name: "get single active realm (without status field provided)",
			req: &realm_mgr_v1.GetRealmRequest{
				Id: s.activeRealmID.String(),
			},
			expectedResponse: s.realmToGRPC(s.activeRealm),
			skip:             s.activeRealm == nil,
		},
		{
			name: "get single active realm (without status field provided)",
			req: &realm_mgr_v1.GetRealmRequest{
				Id:     s.activeRealmID.String(),
				Status: entities.StatusActive,
			},
			expectedResponse: s.realmToGRPC(s.activeRealm),
			skip:             s.activeRealm == nil,
		},
		{
			name: "get single draft realm",
			req: &realm_mgr_v1.GetRealmRequest{
				Id:     s.draftRealmID.String(),
				Status: entities.StatusDraft,
			},
			expectedResponse: s.realmToGRPC(s.draftRealm),
			skip:             s.draftRealm == nil,
		},
		{
			name: "get single disabled realm",
			req: &realm_mgr_v1.GetRealmRequest{
				Id:     s.disabledRealmID.String(),
				Status: entities.StatusDisabled,
			},
			expectedResponse: s.realmToGRPC(s.disabledRealm),
			skip:             s.disabledRealm == nil,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			// arrange
			ctx, err := utils.MakeGRPCRequestContext(context.Background())
			require.NoError(t, err)

			// act
			res, err := s.client.GetRealm(ctx, tc.req)

			// assert
			assert.NoError(t, err)

			require.NotNil(t, res)
			assert.Equal(t, tc.expectedResponse, res.GetRealm())
		})
	}
}

func (s *GetRealmTestSuite) Test_GetRealm_InvalidArgument() {
	testCases := []struct {
		name           string
		req            *realm_mgr_v1.GetRealmRequest
		expectedErrMsg string
	}{
		{
			name:           "no ID provided",
			req:            &realm_mgr_v1.GetRealmRequest{},
			expectedErrMsg: "invalid GetRealmRequest.Id: value must be a valid UUID | caused by: invalid uuid format",
		},
		{
			name: "malformed ID provided",
			req: &realm_mgr_v1.GetRealmRequest{
				Id: "not-valid-uuid",
			},
			expectedErrMsg: "invalid GetRealmRequest.Id: value must be a valid UUID | caused by: invalid uuid format",
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			// arrange
			ctx, err := utils.MakeGRPCRequestContext(context.Background())
			require.NoError(t, err)

			// act
			res, err := s.client.GetRealm(ctx, tc.req)

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

func (s *GetRealmTestSuite) Test_GetRealm_NotFound() {
	testCases := []struct {
		name string
		req  *realm_mgr_v1.GetRealmRequest
	}{
		{
			name: "non-existing realm",
			req: &realm_mgr_v1.GetRealmRequest{
				Id: uuid.New().String(),
			},
		},
		{
			name: "non-existing draft realm",
			req: &realm_mgr_v1.GetRealmRequest{
				Id:     uuid.New().String(),
				Status: realm_mgr_v1.EnumStatus_ENUM_STATUS_DRAFT,
			},
		},
		{
			name: "deleted realm not visible",
			req: &realm_mgr_v1.GetRealmRequest{
				Id: s.disabledRealm.ID.String(),
			},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			// arrange
			ctx, err := utils.MakeGRPCRequestContext(context.Background())
			require.NoError(t, err)

			// act
			res, err := s.client.GetRealm(ctx, tc.req)

			// assert
			assert.Nil(t, res)

			require.Error(t, err)

			gRPCError, ok := status.FromError(err)
			require.True(t, ok)

			assert.Equal(t, codes.NotFound, gRPCError.Code())
			assert.Equal(t, fmt.Sprintf("realm with ID not found: %s", tc.req.Id), gRPCError.Message())
		})
	}
}

func (s *GetRealmTestSuite) populateTestData() error {
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

func (s *GetRealmTestSuite) realmToGRPC(realm *entities.Realm) *realm_mgr_v1.Realm {
	if realm == nil {
		return nil
	}

	grpcStatus, ok := models.StatusEnumValues[realm.Status]
	require.True(s.T(), ok, "unexpected status type: %d", realm.Status)

	return &realm_mgr_v1.Realm{
		Id:          realm.ID.String(),
		Name:        realm.Name,
		Description: realm.Description,
		Status:      grpcStatus,
		CreatedAt:   timestamppb.New(realm.CreatedAt),
		UpdatedAt:   timestamppb.New(realm.UpdatedAt),
	}
}
