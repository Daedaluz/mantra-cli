package client

import (
	"github.com/daedaluz/mantra/internal/location"
	"github.com/daedaluz/mantra/lib/grpc/common"
	"github.com/go-webauthn/webauthn/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RegisterClientService(server *grpc.Server, impl ClientServiceServer) {
	server.RegisterService(&ClientService_ServiceDesc, impl)
}

func (r *SignRequest) Validate() error {
	if r.GetTimeout() < 0 {
		return status.Errorf(codes.InvalidArgument, "timeout must be non-negative")
	}
	if len(r.GetData()) > 0 && len(r.GetMessage()) == 0 {
		return status.Errorf(codes.InvalidArgument, "message must not be empty if data is provided")
	}

	switch protocol.CredentialMediationRequirement(r.Mediation) {
	case protocol.MediationDefault:
	case protocol.MediationOptional:
	case protocol.MediationRequired:
	case protocol.MediationSilent:
	case protocol.MediationConditional:
	default:
		return status.Errorf(codes.InvalidArgument, "invalid mediation: %s", r.Mediation)
	}

	switch protocol.UserVerificationRequirement(r.UserVerification) {
	case protocol.VerificationPreferred:
	case protocol.VerificationRequired:
	case protocol.VerificationDiscouraged:
	default:
		return status.Errorf(codes.InvalidArgument, "invalid user verification: %s", r.UserVerification)
	}

	// Validate location if provided
	if r.Location != nil {
		if err := location.ValidateLocation(r.Location); err != nil {
			return err
		}

		// Validate geohash if provided
		if r.Location.Geohash != "" {
			if err := location.ValidateGeohash(r.Location.Geohash); err != nil {
				return err
			}
		}

		// Validate coordinates if provided
		if r.Location.Position != nil {
			if err := location.ValidateCoordinates(r.Location.Position.Latitude, r.Location.Position.Longitude); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetLocationCoordinates returns the location as coordinates, converting from geohash if necessary
func (r *SignRequest) GetLocationCoordinates() (*common.LongLatitude, error) {
	return location.GetLocationCoordinates(r.Location)
}
