package api

import (
	"testing"

	"github.com/gsamokovarov/assert"
)

func TestSSHKeyFingerprint(t *testing.T) {
	t.Run("matches ssh-keygen", func(t *testing.T) {
		key := SSHKey{Key: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIFdn/jDueFVdtBLhbJBi0ouGHjg8R6LYf7dR5EF+e3wT test"}
		assert.Equal(t, "SHA256:2y/xG2se96jZsBj2UrV4+Uc4vXWR6AiWRR82/l0/lcQ", key.Fingerprint())
	})

	t.Run("invalid key", func(t *testing.T) {
		assert.Equal(t, "", SSHKey{Key: "ssh-rsa AAAA..."}.Fingerprint())
		assert.Equal(t, "", SSHKey{Key: ""}.Fingerprint())
	})
}
