package validate_test

import (
	"os"
	"testing"

	"github.com/opencontainers/runc/libcontainer/configs"
	"github.com/opencontainers/runc/libcontainer/configs/validate"
)

func TestValidate(t *testing.T) {
	config := &configs.Config{
		Rootfs: "/var",
	}

	validator := validate.New()
	err := validator.Validate(config)
	if err != nil {
		t.Errorf("Expected error to not occur: %+v", err)
	}
}

func TestValidateWithInvalidRootfs(t *testing.T) {
	dir := "rootfs"
	os.Symlink("/var", dir)
	defer os.Remove(dir)

	config := &configs.Config{
		Rootfs: dir,
	}

	validator := validate.New()
	err := validator.Validate(config)
	if err == nil {
		t.Error("Expected error to occur but it was nil")
	}
}

func TestValidateNetworkWithoutNETNamespace(t *testing.T) {
	network := &configs.Network{Type: "loopback"}
	config := &configs.Config{
		Rootfs:     "/var",
		Namespaces: []configs.Namespace{},
		Networks:   []*configs.Network{network},
	}

	validator := validate.New()
	err := validator.Validate(config)
	if err == nil {
		t.Error("Expected error to occur but it was nil")
	}
}

func TestValidateNetworkRoutesWithoutNETNamespace(t *testing.T) {
	route := &configs.Route{Gateway: "255.255.255.0"}
	config := &configs.Config{
		Rootfs:     "/var",
		Namespaces: []configs.Namespace{},
		Routes:     []*configs.Route{route},
	}

	validator := validate.New()
	err := validator.Validate(config)
	if err == nil {
		t.Error("Expected error to occur but it was nil")
	}
}

func TestValidateHostname(t *testing.T) {
	config := &configs.Config{
		Rootfs:   "/var",
		Hostname: "runc",
		Namespaces: configs.Namespaces(
			[]configs.Namespace{
				{Type: configs.NEWUTS},
			},
		),
	}

	validator := validate.New()
	err := validator.Validate(config)
	if err != nil {
		t.Errorf("Expected error to not occur: %+v", err)
	}
}

func TestValidateHostnameWithoutUTSNamespace(t *testing.T) {
	config := &configs.Config{
		Rootfs:   "/var",
		Hostname: "runc",
	}

	validator := validate.New()
	err := validator.Validate(config)
	if err == nil {
		t.Error("Expected error to occur but it was nil")
	}
}

func TestValidateSecurityWithMaskPaths(t *testing.T) {
	config := &configs.Config{
		Rootfs:    "/var",
		MaskPaths: []string{"/proc/kcores"},
		Namespaces: configs.Namespaces(
			[]configs.Namespace{
				{Type: configs.NEWNS},
			},
		),
	}

	validator := validate.New()
	err := validator.Validate(config)
	if err != nil {
		t.Errorf("Expected error to not occur: %+v", err)
	}
}

func TestValidateSecurityWithROPaths(t *testing.T) {
	config := &configs.Config{
		Rootfs:        "/var",
		ReadonlyPaths: []string{"/proc/sys"},
		Namespaces: configs.Namespaces(
			[]configs.Namespace{
				{Type: configs.NEWNS},
			},
		),
	}

	validator := validate.New()
	err := validator.Validate(config)
	if err != nil {
		t.Errorf("Expected error to not occur: %+v", err)
	}
}

func TestValidateSecurityWithoutNEWNS(t *testing.T) {
	config := &configs.Config{
		Rootfs:        "/var",
		MaskPaths:     []string{"/proc/kcores"},
		ReadonlyPaths: []string{"/proc/sys"},
	}

	validator := validate.New()
	err := validator.Validate(config)
	if err == nil {
		t.Error("Expected error to occur but it was nil")
	}
}

func TestValidateUsernamespace(t *testing.T) {
	config := &configs.Config{
		Rootfs: "/var",
		Namespaces: configs.Namespaces(
			[]configs.Namespace{
				{Type: configs.NEWUSER},
			},
		),
	}

	validator := validate.New()
	err := validator.Validate(config)
	if err != nil {
		t.Errorf("expected error to not occur %+v", err)
	}
}

func TestValidateUsernamespaceWithoutUserNS(t *testing.T) {
	uidMap := configs.IDMap{ContainerID: 123}
	config := &configs.Config{
		Rootfs:      "/var",
		UidMappings: []configs.IDMap{uidMap},
	}

	validator := validate.New()
	err := validator.Validate(config)
	if err == nil {
		t.Error("Expected error to occur but it was nil")
	}
}

func TestValidateSysctl(t *testing.T) {
	sysctl := map[string]string{
		"fs.mqueue.ctl": "ctl",
		"net.ctl":       "ctl",
		"kernel.ctl":    "ctl",
	}

	for k, v := range sysctl {
		config := &configs.Config{
			Rootfs: "/var",
			Sysctl: map[string]string{k: v},
		}

		validator := validate.New()
		err := validator.Validate(config)
		if err == nil {
			t.Error("Expected error to occur but it was nil")
		}
	}
}
