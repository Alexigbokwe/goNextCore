package core

// Called when a module is initialized.
type OnModuleInit interface {
	OnModuleInit() error
}

// Called when a module is destroyed.
type OnModuleDestroy interface {
	OnModuleDestroy() error
}
