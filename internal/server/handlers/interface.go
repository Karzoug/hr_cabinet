// TODO: Перенести в интерфейсы internal/delivery/http/internal/handlers
type KeyRepository interface {
	Set(ctx context.Context, key string, value int, duration time.Duration) error
	Get(ctx context.Context, key string) (int, error)
	Delete(ctx context.Context, key string) error
}

type DBRepository interface {
	ExistUser(ctx context.Context, userID int) (bool, error)
	GetAuthnData(ctx context.Context, login string) (model.AuthnDAO, error)
	// TODO: Добавить в интерфейсы internal/delivery/http/internal/handlers
	ExistEmployee(ctx context.Context, workEmail string) (bool, int, error)
	ChangePass(ctx context.Context, userID int, hash string) error
}
