package updaterealm

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

func TestRealmManagerUpdateRealmGRPCSuite(t *testing.T) {
	testSuite := NewUpdateRealmTestSuite(t)
	suite.Run(t, testSuite)
}

type UpdateRealmTestSuite struct {
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

func NewUpdateRealmTestSuite(t *testing.T) *UpdateRealmTestSuite {
	cfg, err := utils.LoadConfig()
	require.NoError(t, err, "error loading configuration file")

	db, err := utils.NewDB(cfg)
	require.NoError(t, err, "error creating database connection")

	client, err := utils.NewRealmManagerGRPCClient(cfg)
	require.NoError(t, err, "error creating gRPC client")

	return &UpdateRealmTestSuite{
		db:     db,
		client: client,
	}
}

func (s *UpdateRealmTestSuite) SetupSuite() {
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

func (s *UpdateRealmTestSuite) TearDownSuite() {
	err := s.db.Wipe()
	require.NoError(s.T(), err, "error wiping database")
}

func (s *UpdateRealmTestSuite) Test_UpdateRealm_Success() {
	testCases := []struct {
		name  string
		realm entities.Realm
		skip  bool
	}{
		{
			name: "update active realm",
			realm: entities.Realm{
				ID:          s.activeRealmID,
				Name:        "ActiveRealmUpdated",
				Description: "Updated active realm",
			},
			skip: s.activeRealm == nil,
		},
		{
			name: "update draft realm",
			realm: entities.Realm{
				ID:          s.draftRealmID,
				Name:        "DraftRealmUpdated",
				Description: "Updated draft realm",
			},
			skip: s.draftRealm == nil,
		},
		{
			name: "update draft realm with no description",
			realm: entities.Realm{
				ID:   s.activeRealmID,
				Name: "ActiveRealmUpdated",
			},
			skip: s.activeRealm == nil,
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
			res, err := s.client.UpdateRealm(ctx, &realm_mgr_v1.UpdateRealmRequest{
				Realm: &realm_mgr_v1.Realm{
					Id:          tc.realm.ID.String(),
					Name:        tc.realm.Name,
					Description: tc.realm.Description,
				},
			})

			// assert
			assert.NoError(t, err)

			require.NotNil(t, res)
			require.NotNil(t, res.GetRealm())

			assert.Equal(t, tc.realm.ID.String(), res.GetRealm().Id)
			assert.Equal(t, tc.realm.Name, res.GetRealm().Name)
			assert.Equal(t, tc.realm.Description, res.GetRealm().Description)
			assert.Equal(t, realm_mgr_v1.EnumStatus_ENUM_STATUS_DRAFT, res.GetRealm().Status)
		})
	}
}

func (s *UpdateRealmTestSuite) Test_UpdateRealm_InvalidArgument() {
	testCases := []struct {
		name           string
		req            *realm_mgr_v1.UpdateRealmRequest
		expectedErrMsg string
		skip           bool
	}{
		{
			name:           "realm to be updated is nil",
			req:            &realm_mgr_v1.UpdateRealmRequest{},
			expectedErrMsg: "invalid realm data supplied",
		},
		{
			name: "realm to be updated is without data",
			req: &realm_mgr_v1.UpdateRealmRequest{
				Realm: &realm_mgr_v1.Realm{},
			},
			expectedErrMsg: "invalid UpdateRealmRequest.Realm: embedded message failed validation | caused by: invalid Realm.Id: " +
				"value must be a valid UUID | caused by: invalid uuid format; invalid Realm.Name: value length must be at least 1 runes",
		},
		{
			name: "malformed ID provided",
			req: &realm_mgr_v1.UpdateRealmRequest{
				Realm: &realm_mgr_v1.Realm{
					Id: "not-valid-uuid",
				},
			},
			expectedErrMsg: "invalid UpdateRealmRequest.Realm: embedded message failed validation | caused by: invalid Realm.Id: " +
				"value must be a valid UUID | caused by: invalid uuid format; invalid Realm.Name: value length must be at least 1 runes",
		},
		{
			name: "realm with blank name",
			req: &realm_mgr_v1.UpdateRealmRequest{
				Realm: &realm_mgr_v1.Realm{
					Id:          s.draftRealmID.String(),
					Description: "This is test",
				},
			},
			expectedErrMsg: "invalid UpdateRealmRequest.Realm: embedded message failed validation | caused by: invalid Realm.Name: " +
				"value length must be at least 1 runes",
			skip: s.draftRealm == nil,
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
			res, err := s.client.UpdateRealm(ctx, tc.req)

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

func (s *UpdateRealmTestSuite) Test_UpdateRealm_NotFound() {
	testCases := []struct {
		name string
		req  *realm_mgr_v1.UpdateRealmRequest
		skip bool
	}{
		{
			name: "non-existing realm",
			req: &realm_mgr_v1.UpdateRealmRequest{
				Realm: &realm_mgr_v1.Realm{
					Id:   uuid.New().String(),
					Name: "Test realm",
				},
			},
		},
		{
			name: "disabled realm",
			req: &realm_mgr_v1.UpdateRealmRequest{
				Realm: &realm_mgr_v1.Realm{
					Id:   s.disabledRealmID.String(),
					Name: "Test realm",
				},
			},
			skip: s.disabledRealm == nil,
		},
		{
			name: "deleted realm",
			req: &realm_mgr_v1.UpdateRealmRequest{
				Realm: &realm_mgr_v1.Realm{
					Id:   s.deletedRealmID.String(),
					Name: "Test realm",
				},
			},
			skip: s.deletedRealm == nil,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			// arrange
			ctx, err := utils.MakeGRPCRequestContext(context.Background())
			require.NoError(t, err)

			// act
			res, err := s.client.UpdateRealm(ctx, tc.req)

			// assert
			assert.Nil(t, res)

			require.Error(t, err)

			gRPCError, ok := status.FromError(err)
			require.True(t, ok)

			assert.Equal(t, codes.NotFound, gRPCError.Code())
			assert.Equal(t, fmt.Sprintf("realm with ID not found: %s", tc.req.Realm.Id), gRPCError.Message())
		})
	}
}

func (s *UpdateRealmTestSuite) populateTestData() error {
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
