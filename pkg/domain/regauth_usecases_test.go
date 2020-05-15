// +build unit_tests all_tests

package domain_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"

	"github.com/ankur-anand/prod-todo/pkg/domain"
)

type dummyRepo struct {
	returnFunc  func() domain.UserModel
	returnStore func(domain.UserModel) (uuid.UUID, error)
}

func (d dummyRepo) Find(ctx context.Context, id uuid.UUID) (domain.UserModel, error) {
	panic("implement me")
}

func (d dummyRepo) FindByEmail(ctx context.Context,
	email string) (domain.UserModel, error) {
	return d.returnFunc(), nil
}

func (d dummyRepo) FindAll(ctx context.Context) (domain.UserIterator, error) {
	panic("implement me")
}

func (d dummyRepo) Update(ctx context.Context, user domain.UserModel) error {
	panic("implement me")
}

func (d dummyRepo) Store(ctx context.Context, user domain.UserModel) (uuid.UUID, error) {
	return d.returnStore(user)
}

func TestService_IsValidEmail(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		name  string
		email string
		want  bool
	}{
		{
			name:  "invalid email @missing",
			email: "ankur.com",
			want:  false,
		},
		{
			name:  "invalid email domain",
			email: "ankur@.com",
			want:  false,
		},
		{
			name:  "valid email",
			email: "ankur@example.com",
			want:  true,
		},
		{
			name:  "invalid user name is more than 64 characters",
			email: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@gm.com",
			want:  false,
		},
	}
	as := domain.NewRegAndAuthService(dummyRepo{})
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			if as.IsValidEmail(tc.email) != tc.want {
				t.Errorf("email validation failed for %s, want %v, got %v", tc.email, tc.want, as.IsValidEmail(tc.email))
			}
		})

	}
}

func TestService_IsValidPassword(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "invalid len password small",
			password: "ankur",
			want:     false,
		},
		{
			name:     "invalid length password too large",
			password: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			want:     false,
		},
		{
			name:     "valid password",
			password: "ankur@example.com",
			want:     true,
		},
	}
	as := domain.NewRegAndAuthService(dummyRepo{})
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			if as.IsValidPassword(tc.password) != tc.want {
				t.Errorf("email validation failed for %s, want %v, got %v", tc.password, tc.want, as.IsValidPassword(tc.password))
			}
		})

	}
}

func TestService_IsDuplicateRegistration(t *testing.T) {
	t.Parallel()
	dummyR := dummyRepo{}
	dummyR.returnFunc = func() domain.UserModel {
		return domain.UserModel{
			ID:       uuid.New(),
			Email:    "ankuranand@example.com",
			Password: "garbage",
			Username: "ankuranand",
		}
	}
	as := domain.NewRegAndAuthService(dummyR)
	ok, _ := as.IsDuplicateRegistration(context.Background(), "anKuranand@example.com")
	if !ok {
		t.Errorf("duplicate Registration validation failed for %s", "ankuranand@example.com")
	}

	ok, _ = as.IsDuplicateRegistration(context.Background(), "anKur@example.com")
	if ok {
		t.Errorf("duplicate Registration validation failed for %s", "ankur@example.com")
	}
}

func TestService_IsCredentialValid(t *testing.T) {
	t.Parallel()
	password := "ankuranand"
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(password),
		bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	dummyR := dummyRepo{}
	dummyR.returnFunc = func() domain.UserModel {
		return domain.UserModel{
			ID:       uuid.New(),
			Email:    "ankuranand@example.com",
			Password: string(encryptedPass),
			Username: "ankuranand",
		}
	}
	as := domain.NewRegAndAuthService(dummyR)
	ok, user, _ := as.IsCredentialValid(context.Background(), "ankuranand@example.com", password)
	if !ok && user.Email != "ankuranand@example.com" {
		t.Errorf("credentail validation failed")
	}

	ok, user, _ = as.IsCredentialValid(context.Background(), "ankuranand@example.com", "garbage")
	if ok && user != domain.NilUserModel {
		t.Errorf("credentail validation failed")
	}
}

func TestService_StoreUser(t *testing.T) {
	t.Parallel()
	var err error
	password := "ankuranand"

	dummyR := dummyRepo{}
	userReceived := make(chan domain.UserModel)

	dummyR.returnStore = func(model domain.UserModel) (uuid.UUID, error) {
		go func() {
			userReceived <- model
		}()
		return model.ID, nil
	}

	as := domain.NewRegAndAuthService(dummyR)
	usr := domain.UserModel{
		Email:     "AnkurananD@example.com", // email should be normalized
		Password:  password,
		FirstName: "Ankur",
		LastName:  "Anand",
	}
	_, err = as.StoreUser(context.Background(), usr)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case user := <-userReceived:
		close(userReceived)
		ok := user.Email == "ankuranand@example.com" && user.FirstName == "Ankur" && password != user.Password
		if !ok {
			t.Errorf("StoreUser failed to received the expected user model")
		}
	case <-time.After(time.Second * 1):
		t.Errorf("storeUser timedout")
	}
}
