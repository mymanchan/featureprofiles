/*
 Copyright 2022 Google LLC

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package system_base_test

import (
	"context"
	"testing"
	"time"

	"github.com/openconfig/ondatra"
	"github.com/openconfig/ondatra/binding/introspect"
	"github.com/openconfig/ondatra/knebind/creds"
	"google.golang.org/grpc"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
	spb "github.com/openconfig/gnoi/system"
	authzpb "github.com/openconfig/gnsi/authz"
	gribipb "github.com/openconfig/gribi/v1/proto/service"
	p4rtpb "github.com/p4lang/p4runtime/go/p4/v1"
)

func dialConn(t *testing.T, dut *ondatra.DUTDevice, service introspect.Service, wantPort uint32) *grpc.ClientConn {
	t.Helper()
	cd := introspect.Introspect(t, dut, service)
	if cd.DevicePort != int(wantPort) {
		t.Fatalf("DUT is not listening on correct port for %q: got %d, want %d", service, cd.DevicePort, wantPort)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := cd.DefaultDial(ctx, cd.DefaultTarget, cd.DefaultDialOpts...)
	if err != nil {
		t.Fatalf("grpc.Dial failed to: %q", cd.DefaultTarget)
	}
	return conn
}

type rpcCredentials struct {
	*creds.UserPass
}

func (r *rpcCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"username": "admin",
		"password": "admin",
	}, nil
}

func (r *rpcCredentials) RequireTransportSecurity() bool {
	return true
}

// TestGNMIClient validates that the DUT listens on standard gNMI Port.
func TestGNMIClient(t *testing.T) {
	dut := ondatra.DUT(t, "dut")
	conn := dialConn(t, dut, introspect.GNMI, 9339)
	c := gpb.NewGNMIClient(conn)
	if _, err := c.Get(context.Background(), &gpb.GetRequest{}); err != nil {
		t.Fatalf("gnmi.Get failed: %v", err)
	}
}

// TestGNOIClient validates that the DUT listens on standard gNMI Port.
func TestGNOIClient(t *testing.T) {
	dut := ondatra.DUT(t, "dut")
	conn := dialConn(t, dut, introspect.GNOI, 9339)
	c := spb.NewSystemClient(conn)
	if _, err := c.Ping(context.Background(), &spb.PingRequest{}); err != nil {
		t.Fatalf("gnoi.system.Time failed: %v", err)
	}
}

// TestGNSIClient validates that the DUT listens on standard gNMI Port.
func TestGNSIClient(t *testing.T) {
	dut := ondatra.DUT(t, "dut")
	conn := dialConn(t, dut, introspect.GNSI, 9339)
	c := authzpb.NewAuthzClient(conn)
	if _, err := c.Get(context.Background(), &authzpb.GetRequest{}); err != nil {
		t.Fatalf("gnsi.authz.Get failed: %v", err)
	}
}

// TestGRIBIClient validates that the DUT listens on standard gNMI Port.
func TestGRIBIClient(t *testing.T) {
	dut := ondatra.DUT(t, "dut")
	conn := dialConn(t, dut, introspect.GRIBI, 9340)
	c := gribipb.NewGRIBIClient(conn)
	if _, err := c.Get(context.Background(), &gribipb.GetRequest{}); err != nil {
		t.Fatalf("gribi.Get failed: %v", err)
	}
}

// TestP4RTClient validates that the DUT listens on standard gNMI Port.
func TestP4RTClient(t *testing.T) {
	dut := ondatra.DUT(t, "dut")
	conn := dialConn(t, dut, introspect.P4RT, 9559)
	c := p4rtpb.NewP4RuntimeClient(conn)
	if _, err := c.Capabilities(context.Background(), &p4rtpb.CapabilitiesRequest{}); err != nil {
		t.Fatalf("p4rt.Capabilites failed: %v", err)
	}
}
