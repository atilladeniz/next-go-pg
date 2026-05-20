package email

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
)

var updateGolden = flag.Bool("update", false, "update golden files for email template tests")

const (
	sampleAppURL      = "https://app.example.com"
	sampleSettingsURL = "https://app.example.com/settings"
)

func TestRenderMagicLink(t *testing.T) {
	got, err := render("magic_link.html", magicLinkData{
		URL:    "https://app.example.com/magic-link/verify?token=tok",
		AppURL: sampleAppURL,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertGolden(t, "magic_link.golden", got)
}

func TestRenderVerification(t *testing.T) {
	got, err := render("verification.html", verificationData{
		URL:    "https://app.example.com/verify-email?token=tok",
		AppURL: sampleAppURL,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertGolden(t, "verification.golden", got)
}

func TestRender2FAOTP(t *testing.T) {
	got, err := render("2fa_otp.html", twoFactorOTPData{
		UserName: "Atilla",
		OTP:      "123456",
		AppURL:   sampleAppURL,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertGolden(t, "2fa_otp.golden", got)
}

func TestRender2FAEnabled(t *testing.T) {
	got, err := render("2fa_enabled.html", twoFactorEnabledData{
		UserName:    "Atilla",
		MethodName:  "Authenticator App",
		AppURL:      sampleAppURL,
		SettingsURL: sampleSettingsURL,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertGolden(t, "2fa_enabled.golden", got)
}

func TestRenderPasskeyAdded(t *testing.T) {
	got, err := render("passkey_added.html", passkeyAddedData{
		UserName:    "Atilla",
		PasskeyName: "MacBook Touch ID",
		Device:      "Chrome on macOS",
		Time:        "20.05.2026 16:30",
		AppURL:      sampleAppURL,
		SettingsURL: sampleSettingsURL,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertGolden(t, "passkey_added.golden", got)
}

func TestRenderLoginNotification(t *testing.T) {
	got, err := render("login_notification.html", loginNotificationData{
		UserName:    "Atilla",
		Device:      "Chrome on macOS",
		IPAddress:   "192.168.1.42",
		Time:        "20.05.2026 16:30",
		AppURL:      sampleAppURL,
		SettingsURL: sampleSettingsURL,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertGolden(t, "login_notification.golden", got)
}

func assertGolden(t *testing.T, name, got string) {
	t.Helper()
	path := filepath.Join("testdata", name)
	if *updateGolden {
		if err := os.MkdirAll("testdata", 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("missing golden file %s — run `go test -run %s -update`: %v", path, t.Name(), err)
	}
	if string(want) != got {
		t.Errorf("output mismatch for %s\n--- want\n%s\n--- got\n%s", name, want, got)
	}
}
