package security

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type KeycloakUser struct {
	ID            string `json:"sub,omitempty"`
	GivenName     string `json:"given_name,omitempty"`
	Name          string `json:"name,omitempty"`
	FamilyName    string `json:"family_name,omitempty"`
	Username      string `json:"preferred_username,omitempty"`
	Email         string `json:"email,omitempty"`
	Audience      string `json:"audience"`
	EmailVerified bool   `json:"emailVerified,omitempty"`
	PhoneNumber   string `json:"phone_number,omitempty"`
	RealmAccess   struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
}

func GetUserClaims(ctx *fiber.Ctx) *KeycloakUser {
	kcUser := new(KeycloakUser)
	ctx.Locals("keycloakUser", kcUser)
	if kcUser == nil {
		return nil
	}

	return kcUser
}

func GetAccessToken(ctx *fiber.Ctx) string {
	rawAccessToken, ok := ctx.GetReqHeaders()["Authorization"]
	if !ok {
		return ""
	}

	return rawAccessToken[0]
}

func IsAuthorizedJWT(ctx *fiber.Ctx) error {
	var keycloakURL = fmt.Sprintf("%s/realms/%s", os.Getenv("KC.BASE_URL"), os.Getenv("KC.REALM"))
	var clientID = os.Getenv("KC.CLIENT_ID")

	log.Println("Keycloak BaseURL:", keycloakURL)
	log.Println("Keycloak ClientID:", clientID)

	rawAccessToken, ok := ctx.GetReqHeaders()["Authorization"]
	if !ok {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   time.Duration(6000) * time.Second,
		Transport: tr,
	}

	c := oidc.ClientContext(context.Background(), client)
	provider, err := oidc.NewProvider(c, keycloakURL)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(err.Error())
	}

	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}
	verifier := provider.Verifier(oidcConfig)

	accessToken := strings.Split(rawAccessToken[0], " ")[1]
	idToken, err := verifier.Verify(c, accessToken)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(err.Error())
	}

	user := new(KeycloakUser)
	if err := idToken.Claims(user); err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(err.Error())
	}

	ctx.Locals("keycloakUser", user)

	return ctx.Next()
}