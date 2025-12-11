package core

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	*fiber.App
}

func NewApp() *App {
	return &App{
		App: fiber.New(fiber.Config{
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			},
		}),
	}
}

func (a *App) Listen(addr string) error {
	a.App.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			fmt.Printf("=== PANIC RECOVERED ===\n")
			fmt.Printf("Panic: %v\n", e)
			fmt.Printf("Path: %s\n", c.Path())
			fmt.Printf("Method: %s\n", c.Method())
			fmt.Printf("Headers: %+v\n", c.GetReqHeaders())

			// Print stack trace
			fmt.Printf("Stack trace: %s\n", debug.Stack())
		},
	}))

	// Default GET route
	a.App.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to GoNext framework")
	})

	return a.App.Listen(addr)
}

// Called when a module is initialized.
func (app *App) InitModules(modules []Module, container *Container) error {
	var initErrors []error

	for _, module := range modules {
		moduleName := fmt.Sprintf("%T", module)
		log.Printf("Initializing module: %s", moduleName)

		// Initialize if the module supports it
		if hook, ok := module.(OnModuleInit); ok {
			if err := hook.OnModuleInit(); err != nil {
				log.Printf("[ERROR] Failed to initialize module %s: %v", moduleName, err)
				initErrors = append(initErrors, err)
				continue // Skip registration if init fails
			}
			log.Printf("Module %s initialized successfully\n", moduleName)
		}

		// Register and mount only if initialization succeeded
		module.Register(container)
		module.MountRoutes(app)
		log.Printf("Module %s registered and mounted successfully", moduleName)
	}

	if len(initErrors) > 0 {
		return fmt.Errorf("failed to initialize %d modules", len(initErrors))
	}

	// Autowire all pending components after all modules are registered
	container.AutowireAll()
	return nil
}

func (app *App) ShutdownModules(modules []Module) {
	for _, module := range modules {
		if hook, ok := module.(OnModuleDestroy); ok {
			if err := hook.OnModuleDestroy(); err != nil {
				log.Println("Error during shutdown:", err)
			}
		}
	}
}

func (app *App) ConnectToDataBase(connectionString string, databaseName string) (*pgxpool.Pool, context.Context, error) {
	var err error
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		log.Printf("Unable to connect to timeline ms database: %v", err)
		return nil, nil, fmt.Errorf("unable to connect to timeline ms database: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		log.Printf("Unable to ping timeline ms database: %v", err)
		return nil, nil, fmt.Errorf("unable to ping timeline ms database: %w", err)
	}

	fmt.Println("Connected to timeline ms database successfully!")
	return pool, ctx, nil
}

func (app *App) DisconnectFromDatabase(dbPool *pgxpool.Pool) {
	if dbPool != nil {
		dbPool.Close()
		fmt.Println("Timeline database connection closed successfully.")
	} else {
		fmt.Println("No Timeline database connection to close.")
	}
}
