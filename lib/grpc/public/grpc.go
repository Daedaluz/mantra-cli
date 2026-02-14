package public

import (
	"github.com/daedaluz/mantra/internal/location"
	"github.com/daedaluz/mantra/lib/grpc/common"
	"github.com/go-webauthn/webauthn/protocol"
	"google.golang.org/grpc"
)

func RegisterPublicAuthService(server *grpc.Server, impl AuthServiceServer) {
	server.RegisterService(&AuthService_ServiceDesc, impl)
}

func RegisterPublicRegistrationService(server *grpc.Server, impl RegisterServiceServer) {
	server.RegisterService(&RegisterService_ServiceDesc, impl)
}

func (c *RegisterRequest) ToParsedCredential() (*protocol.ParsedCredentialCreationData, error) {
	x := protocol.CredentialCreationResponse{
		PublicKeyCredential: protocol.PublicKeyCredential{
			Credential: protocol.Credential{
				ID:   c.GetId(),
				Type: c.GetType(),
			},
			RawID:                   c.GetRawId(),
			AuthenticatorAttachment: c.GetAttachment(),
		},
		AttestationResponse: protocol.AuthenticatorAttestationResponse{
			AuthenticatorResponse: protocol.AuthenticatorResponse{
				ClientDataJSON: protocol.URLEncodedBase64(c.GetClientDataJson()),
			},
			AttestationObject: protocol.URLEncodedBase64(c.GetAttestationObject()),
		},
	}
	var response *protocol.ParsedCredentialCreationData
	var err error
	if response, err = x.Parse(); err != nil {
		return nil, err
	}
	response.Raw.AttestationResponse.PublicKeyAlgorithm = c.GetPublicKeyAlgorithm()
	return response, nil
}

func (s *SignRequest) ToParsedCredentialAssertionData() (*protocol.ParsedCredentialAssertionData, error) {
	x := protocol.CredentialAssertionResponse{
		PublicKeyCredential: protocol.PublicKeyCredential{
			Credential: protocol.Credential{
				ID:   s.GetId(),
				Type: "public-key",
			},
			RawID:                   s.GetRawId(),
			ClientExtensionResults:  nil,
			AuthenticatorAttachment: s.GetAuthenticatorAttachment(),
		},
		AssertionResponse: protocol.AuthenticatorAssertionResponse{
			AuthenticatorResponse: protocol.AuthenticatorResponse{
				ClientDataJSON: s.GetClientDataJson(),
			},
			AuthenticatorData: s.AuthenticatorData,
			Signature:         s.Signature,
			UserHandle:        s.UserHandle,
		},
	}
	var response = &protocol.ParsedCredentialAssertionData{}
	var err error
	if response, err = x.Parse(); err != nil {
		return nil, err
	}
	return response, nil
}

// GetLocationCoordinates returns the location as coordinates, converting from geohash if necessary
func (s *SignRequest) GetLocationCoordinates() (*common.LongLatitude, error) {
	return location.GetLocationCoordinates(s.Location)
}

// GetLocationCoordinates returns the location as coordinates, converting from geohash if necessary
func (r *RegisterRequest) GetLocationCoordinates() (*common.LongLatitude, error) {
	return location.GetLocationCoordinates(r.Location)
}
