package cmd

import (
	"encoding/base64"
	"fmt"

	admin "github.com/daedaluz/mantra-cli/lib/grpc/service"
)

func printAuditEntry(e *admin.AuditEntry) {
	fmt.Printf("Challenge: %s\n", e.ChallengeId)
	fmt.Printf("Title:     %s\n", e.Title)
	fmt.Printf("Message:   %s\n", e.Message)
	fmt.Printf("Status:    %s\n", e.Status.String())
	fmt.Printf("Verified:  %v\n", e.Verified)
	if e.CreatedAt != nil {
		fmt.Printf("Created:   %s\n", e.CreatedAt.AsTime().Format("2006-01-02 15:04:05"))
	}
	if e.SignedAt != nil {
		fmt.Printf("Signed:    %s\n", e.SignedAt.AsTime().Format("2006-01-02 15:04:05"))
	}
	if e.Key != nil {
		fmt.Printf("Key:       %s (alg=%d, revoked=%v, name=%s)\n",
			base64.RawURLEncoding.EncodeToString(e.Key.KeyId),
			e.Key.Algorithm,
			e.Key.Revoked,
			e.Key.Name,
		)
	}
	if e.Client != nil {
		fmt.Printf("Client:    %s (%s, admin=%v)\n", e.Client.Name, e.Client.ClientId, e.Client.Admin)
	}
	if e.Device != nil {
		fmt.Printf("Device:    %s (%s)\n", e.Device.Name, e.Device.DeviceId)
	}
	if e.SignatureLocation != nil {
		loc := e.SignatureLocation
		if loc.Position != nil {
			fmt.Printf("Location:  lat=%.6f lon=%.6f\n", loc.Position.Latitude, loc.Position.Longitude)
		} else if loc.Geohash != "" {
			fmt.Printf("Location:  geohash=%s\n", loc.Geohash)
		} else if loc.Ip != "" {
			fmt.Printf("Location:  ip=%s\n", loc.Ip)
		}
	}
	if len(e.AuthenticatorData) > 0 {
		fmt.Printf("AuthData:  %s\n", base64.RawURLEncoding.EncodeToString(e.AuthenticatorData))
	}
	if len(e.ClientDataJson) > 0 {
		fmt.Printf("ClientData: %s\n", base64.RawURLEncoding.EncodeToString(e.ClientDataJson))
	}
	if len(e.Signature) > 0 {
		fmt.Printf("Signature: %s\n", base64.RawURLEncoding.EncodeToString(e.Signature))
	}
}
