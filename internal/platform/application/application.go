package application

import (
	"blog/configs"
	"blog/internal/domain/posts"
	"blog/internal/domain/users"
	"blog/internal/interactors"
	"blog/internal/platform/pkg/jwt"
	"blog/internal/platform/pkg/lang"
	"blog/internal/platform/pkg/password"
	"blog/internal/platform/pkg/validators"
	"blog/internal/platform/repositories"
	"blog/internal/platform/seeder"
	"context"
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	mongodbadapter "github.com/casbin/mongodb-adapter/v3"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type App struct {
	Config *configs.Config
	DB     struct {
		Mongo *mongo.Database
	}
	Repositories struct {
		UsersRepository users.Repository
		PostsRepository posts.Repository
	}
	Interactors struct {
		UserInteractor users.Interactor
		PostInteractor posts.Interactor
	}
	Router         *gin.Engine
	PasswordHasher users.PasswordHasher
	AuthService    jwt.Service
	Validator      *validators.Validator
	Enforcer       *casbin.Enforcer
	Seeder         *seeder.Seeder
}

func NewApp(cfg *configs.Config) (*App, error) {
	app := &App{
		Config: cfg,
	}
	if err := app.registerMongoDB(); err != nil {
		return nil, err
	}
	app.registerRepositories()
	app.registerPasswordHasher()
	if err := app.registerAuthService(); err != nil {
		return nil, err
	}
	if err := app.registerCasbinAdapter(); err != nil {
		return nil, err
	}
	if err := app.registerInteractors(); err != nil {
		return nil, err
	}
	if err := app.registerValidator(); err != nil {
		return nil, err
	}
	app.registerSeeder()
	app.registerRouter()
	app.RegisterRoutes()
	return app, nil
}

func (a *App) registerMongoDB() error {
	log.Println("Database is Connecting...")
	if a.Config.DB.Mongo.Username == "" ||
		a.Config.DB.Mongo.Password == "" ||
		a.Config.DB.Mongo.HostName == "" ||
		a.Config.DB.Mongo.Port == "" ||
		a.Config.DB.Mongo.Database == "" ||
		a.Config.DB.Mongo.AuthSource == "" {
		return errors.New("mongo config is required")
	}
	dbUri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=%s",
		a.Config.DB.Mongo.Username,
		a.Config.DB.Mongo.Password,
		a.Config.DB.Mongo.HostName,
		a.Config.DB.Mongo.Port,
		a.Config.DB.Mongo.Database,
		a.Config.DB.Mongo.AuthSource)
	clientOptions := options.Client().ApplyURI(dbUri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	a.DB.Mongo = client.Database(a.Config.DB.Mongo.Database)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}
	log.Println("Database Connected Successfully.")
	return err
}

func (a *App) registerCasbinAdapter() error {
	if a.Config.Casbin.DB.CollectionName == "" {
		return errors.New("casbin collection name is required")
	}
	config := &mongodbadapter.AdapterConfig{
		DatabaseName:   a.Config.DB.Mongo.Database,
		CollectionName: a.Config.Casbin.DB.CollectionName,
		Timeout:        30 * time.Second,
		IsFiltered:     false,
	}
	adapter, err := mongodbadapter.NewAdapterByDB(a.DB.Mongo.Client(), config)
	if err != nil {
		return err
	}
	enforcer, err := casbin.NewEnforcer("./configs/rbac_model.conf", adapter)
	if err != nil {
		return err
	}
	a.Enforcer = enforcer
	return nil
}

func (a *App) registerInteractors() error {
	a.Interactors.UserInteractor = interactors.NewUserInteractor(a.Repositories.UsersRepository, a.PasswordHasher, a.AuthService, a.Enforcer)
	a.Interactors.PostInteractor = interactors.NewPostInteractor(a.Repositories.PostsRepository, a.Enforcer)
	return nil
}

func (a *App) registerRepositories() {
	a.Repositories.UsersRepository = repositories.NewMongoUserRepository(a.DB.Mongo)
	a.Repositories.PostsRepository = repositories.NewMongoPostRepository(a.DB.Mongo)
}

func (a *App) registerValidator() error {
	a.Validator = validators.NewCustomValidator(a.Repositories.UsersRepository)
	return a.Validator.RegisterValidation()
}

func (a *App) registerPasswordHasher() {
	a.PasswordHasher = password.NewPasswordHasher()
}

func (a *App) registerAuthService() error {
	if len(a.Config.Auth.JWT.SecretKey) < 1 ||
		a.Config.Auth.JWT.ExpiredTime < 1 ||
		len(a.Config.Auth.JWT.Issuer) < 1 ||
		a.Config.Auth.JWT.ExpiredTime < 1 {
		return errors.New("the configuration values of the auth service are not done")
	}
	exp := time.Hour * time.Duration(a.Config.Auth.JWT.ExpiredTime)
	a.AuthService = jwt.NewJWTService(
		a.Config.Auth.JWT.SecretKey,
		a.Config.Auth.JWT.Issuer,
		exp,
	)
	return nil
}

func (a *App) registerSeeder() {
	a.Seeder = seeder.NewSeeder(a.DB.Mongo, a.Enforcer)
}

func (a *App) registerRouter() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery(), lang.GinLocaleDetectorMiddleware())
	a.Router = router
}

func (a *App) RunRouter() {
	log.Println("Router is Running ...")
	if a.Config.Router.Address == "" {
		log.Fatal("Router Address isn't Set")
	}
	if err := a.Router.Run(a.Config.Router.Address); err != nil {
		log.Fatal(err)
	}
	log.Printf("Server Started On Port %s", a.Config.Router.Address)
}

func (a *App) Close() error {
	log.Debug("System Successfully Closed.")
	return nil
}
