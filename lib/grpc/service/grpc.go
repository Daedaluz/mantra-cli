package admin

import (
	"github.com/daedaluz/mantra-cli/lib/grpc/common"
	"github.com/daedaluz/mantra-cli/lib/location"
	"google.golang.org/grpc"
)

func RegisterPlatformAdminService(server *grpc.Server, impl PlatformAdminServiceServer) {
	server.RegisterService(&PlatformAdminService_ServiceDesc, impl)
}

func RegisterDomainAdminService(server *grpc.Server, impl DomainAdminServiceServer) {
	server.RegisterService(&DomainAdminService_ServiceDesc, impl)
}

// GetLocationCoordinates returns the location as coordinates, converting from geohash if necessary
func (r *CreateUserRequest) GetLocationCoordinates() (*common.LongLatitude, error) {
	return location.GetLocationCoordinates(r.Location)
}
