package createrealm

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	realm_mgr_v1 "github.com/alexZaicev/realm-mgr/proto/go/realm_mgr/v1"
	"github.com/alexZaicev/realm-mgr/tests/functional/utils"
)

func TestRealmManagerCreateRealmGRPCSuite(t *testing.T) {
	testSuite := NewCreateRealmTestSuite(t)
	suite.Run(t, testSuite)
}

type CreateRealmTestSuite struct {
	suite.Suite

	db     *utils.DB
	client realm_mgr_v1.RealmManagerServiceClient
}

func NewCreateRealmTestSuite(t *testing.T) *CreateRealmTestSuite {
	cfg, err := utils.LoadConfig()
	require.NoError(t, err, "error loading configuration file")

	db, err := utils.NewDB(cfg)
	require.NoError(t, err, "error creating database connection")

	client, err := utils.NewRealmManagerGRPCClient(cfg)
	require.NoError(t, err, "error creating gRPC client")

	return &CreateRealmTestSuite{
		db:     db,
		client: client,
	}
}

func (s *CreateRealmTestSuite) TearDownSuite() {
	err := s.db.Wipe()
	require.NoError(s.T(), err, "error wiping database")
}

func (s *CreateRealmTestSuite) Test_CreateRealm_Success() {
	testCases := []struct {
		name           string
		expectedRealm  *realm_mgr_v1.Realm
		expectedErrMsg string
	}{
		{
			name: "create realm without description",
			expectedRealm: &realm_mgr_v1.Realm{
				Name:   "CreateRealmTestSuite Test",
				Status: realm_mgr_v1.EnumStatus_ENUM_STATUS_DRAFT,
			},
		},
		{
			name: "create realm with description",
			expectedRealm: &realm_mgr_v1.Realm{
				Name:        "CreateRealmTestSuite Test",
				Description: "Test realm description",
				Status:      realm_mgr_v1.EnumStatus_ENUM_STATUS_DRAFT,
			},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			// arrange
			ctx, err := utils.MakeGRPCRequestContext(context.Background())
			require.NoError(s.T(), err)

			req := &realm_mgr_v1.CreateRealmRequest{
				Name:        tc.expectedRealm.Name,
				Description: tc.expectedRealm.Description,
			}

			// act
			res, err := s.client.CreateRealm(ctx, req)

			// assert
			assert.NoError(s.T(), err)

			require.NotNil(s.T(), res)
			require.NotNil(s.T(), res.GetRealm())

			assert.NotEqual(s.T(), uuid.Nil, res.GetRealm().Id)
			assert.Equal(s.T(), tc.expectedRealm.Name, res.GetRealm().Name)
			assert.Equal(s.T(), tc.expectedRealm.Description, res.GetRealm().Description)
			assert.Equal(s.T(), tc.expectedRealm.Status, res.GetRealm().Status)
		})
	}
}

func (s *CreateRealmTestSuite) Test_CreateRealm_InvalidArgument() {
	testCases := []struct {
		name           string
		realmName      string
		realmDesc      string
		expectedErrMsg string
	}{
		{
			name:           "no name or description",
			expectedErrMsg: "invalid CreateRealmRequest.Name: value length must be at least 1 runes",
		},
		{
			name:           "no name",
			realmDesc:      "Test realm description",
			expectedErrMsg: "invalid CreateRealmRequest.Name: value length must be at least 1 runes",
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			// arrange
			ctx, err := utils.MakeGRPCRequestContext(context.Background())
			require.NoError(s.T(), err)

			req := &realm_mgr_v1.CreateRealmRequest{
				Name:        tc.realmName,
				Description: tc.realmDesc,
			}

			// act
			res, err := s.client.CreateRealm(ctx, req)

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
