package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/mapper"
	"github.com/itsLeonB/billsplittr/internal/util"
	"github.com/itsLeonB/cocoon-protos/gen/go/profile"
	"github.com/itsLeonB/ezutil"
	"github.com/rotisserie/eris"
)

type profileServiceGrpc struct {
	svcClient profile.ProfileServiceClient
}

func NewProfileService(svcClient profile.ProfileServiceClient) ProfileService {
	return &profileServiceGrpc{svcClient}
}

func (ps *profileServiceGrpc) GetByID(ctx context.Context, id uuid.UUID) (dto.ProfileResponse, error) {
	response, err := ps.svcClient.Get(ctx, &profile.ProfileRequest{ProfileId: id.String()})
	if err != nil {
		return dto.ProfileResponse{}, eris.Wrap(err, appconstant.ErrServiceClient)
	}

	return mapper.FromProfileProto(response)
}

func (ps *profileServiceGrpc) GetNames(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]string, error) {
	req := &profile.GetNamesRequest{
		ProfileIds: ezutil.MapSlice(ids, util.ToString),
	}

	response, err := ps.svcClient.GetNames(ctx, req)
	if err != nil {
		return nil, eris.Wrap(err, appconstant.ErrServiceClient)
	}

	namesByProfileID := make(map[uuid.UUID]string, len(response.GetNamesByProfileId()))
	for id, name := range response.GetNamesByProfileId() {
		parsedID, err := ezutil.Parse[uuid.UUID](id)
		if err != nil {
			return nil, err
		}
		namesByProfileID[parsedID] = name
	}

	return namesByProfileID, nil
}

func (ps *profileServiceGrpc) GetByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]dto.ProfileResponse, error) {
	profileIDs := ezutil.MapSlice(ids, util.ToString)
	response, err := ps.svcClient.GetByIDs(ctx, &profile.GetByIDsRequest{ProfileIds: profileIDs})
	if err != nil {
		return nil, eris.Wrap(err, appconstant.ErrServiceClient)
	}

	profilesByID := make(map[uuid.UUID]dto.ProfileResponse, len(response.GetProfiles()))
	for _, profile := range response.GetProfiles() {
		profileResponse, err := mapper.FromProfileProto(profile)
		if err != nil {
			return nil, err
		}

		profilesByID[profileResponse.ID] = profileResponse
	}

	return profilesByID, nil
}
