package core

func RegisterModuleComponents(container *Container, components ...any) {
	// Register each component in the container
	for _, c := range components {
		container.Register(c)
	}

	// Try autowiring each component, but don't fail if dependencies aren't ready yet
	for _, c := range components {
		if err := container.Autowire(c); err != nil {
			// Store for later autowiring if dependencies aren't available yet
			container.AddForAutowiring(c)
		}
	}
}
