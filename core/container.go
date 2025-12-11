package core

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// ServiceFactory is a function that creates new instances
type ServiceFactory func() any

// ServiceScope defines the lifetime of services
type ServiceScope int

const (
	Singleton ServiceScope = iota
	Transient
	Scoped // Request scoped
)

// ServiceRegistration holds service metadata
type ServiceRegistration struct {
	Instance any
	Factory  ServiceFactory
	Scope    ServiceScope
}

func NewContainer() *Container {
	return &Container{
		typeServices:    make(map[reflect.Type]*ServiceRegistration),
		tokenServices:   make(map[string]*ServiceRegistration),
		scopedInstances: make(map[string]map[reflect.Type]any),
		scopedTokens:    make(map[string]map[string]any),
		pendingAutowire: make([]any, 0),
		processing:      make(map[reflect.Type]bool),
	}
}

// Add processing map to Container struct
type Container struct {
	typeServices    map[reflect.Type]*ServiceRegistration
	tokenServices   map[string]*ServiceRegistration
	scopedInstances map[string]map[reflect.Type]any // scopeKey -> type -> instance
	scopedTokens    map[string]map[string]any       // scopeKey -> token -> instance
	lock            sync.RWMutex
	pendingAutowire []any
	processing      map[reflect.Type]bool
}

// Register by type as singleton (default behavior)
func (c *Container) Register(service any) {
	c.RegisterScoped(service, Singleton)
}

// RegisterScoped registers a service with a specific scope
func (c *Container) RegisterScoped(service any, scope ServiceScope) {
	c.lock.Lock()
	defer c.lock.Unlock()
	serviceType := reflect.TypeOf(service)
	c.typeServices[serviceType] = &ServiceRegistration{
		Instance: service,
		Scope:    scope,
	}
}

// RegisterFactory registers a factory function for creating instances
func (c *Container) RegisterFactory(serviceType reflect.Type, factory ServiceFactory, scope ServiceScope) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.typeServices[serviceType] = &ServiceRegistration{
		Factory: factory,
		Scope:   scope,
	}
}

// RegisterTransient registers a service as transient (new instance each time)
func (c *Container) RegisterTransient(prototype any) {
	c.RegisterTransientFactory(reflect.TypeOf(prototype), func() any {
		// Create a new instance by copying the prototype
		prototypeVal := reflect.ValueOf(prototype)
		if prototypeVal.Kind() == reflect.Ptr {
			// If prototype is a pointer, create new instance of the underlying type
			elemType := prototypeVal.Elem().Type()
			newInstance := reflect.New(elemType)
			return newInstance.Interface()
		}
		// If prototype is a value, create new instance of its type
		newInstance := reflect.New(prototypeVal.Type())
		return newInstance.Interface()
	})
}

// RegisterTransientFactory registers a factory for transient services
func (c *Container) RegisterTransientFactory(serviceType reflect.Type, factory ServiceFactory) {
	c.RegisterFactory(serviceType, factory, Transient)
}

// RegisterScopedFactory registers a factory for request-scoped services
func (c *Container) RegisterScopedFactory(serviceType reflect.Type, factory ServiceFactory) {
	c.RegisterFactory(serviceType, factory, Scoped)
}

// Bind by string token as singleton (default behavior)
func (c *Container) Bind(token string, service any) {
	c.BindScoped(token, service, Singleton)
}

// BindScoped binds a service with a specific scope
func (c *Container) BindScoped(token string, service any, scope ServiceScope) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.tokenServices[token] = &ServiceRegistration{
		Instance: service,
		Scope:    scope,
	}
}

// BindFactory binds a factory function for creating instances
func (c *Container) BindFactory(token string, factory ServiceFactory, scope ServiceScope) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.tokenServices[token] = &ServiceRegistration{
		Factory: factory,
		Scope:   scope,
	}
}

// BindTransient binds a service as transient
func (c *Container) BindTransient(token string, prototype any) {
	c.BindFactory(token, func() any {
		prototypeVal := reflect.ValueOf(prototype)
		if prototypeVal.Kind() == reflect.Ptr {
			elemType := prototypeVal.Elem().Type()
			newInstance := reflect.New(elemType)
			return newInstance.Interface()
		}
		newInstance := reflect.New(prototypeVal.Type())
		return newInstance.Interface()
	}, Transient)
}

// BindScoped binds a service as request-scoped
func (c *Container) BindScopedFactory(token string, factory ServiceFactory) {
	c.BindFactory(token, factory, Scoped)
}

// Resolve by type with optional scope key for scoped services
func (c *Container) Resolve(target any) error {
	return c.ResolveWithScope(target, "")
}

// ResolveWithScope resolves with a scope key for request-scoped services
func (c *Container) ResolveWithScope(target any, scopeKey string) error {
	c.lock.RLock()
	defer c.lock.RUnlock()

	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr {
		return errors.New("target must be a pointer")
	}

	elem := val.Elem()
	var targetType reflect.Type

	// Handle **T case (pointer to pointer)
	if elem.Kind() == reflect.Ptr {
		targetType = elem.Type()
	} else if elem.Kind() == reflect.Struct {
		// Handle *T case (pointer to struct)
		targetType = reflect.PointerTo(elem.Type())
	} else {
		return errors.New("target must be a pointer to a struct or pointer to a pointer")
	}

	registration, ok := c.typeServices[targetType]
	if !ok {
		return fmt.Errorf("no registered service for type %s", targetType.String())
	}

	instance, err := c.getInstance(registration, targetType, "", scopeKey)
	if err != nil {
		return err
	}

	return c.assignInstance(instance, elem, targetType)
}

// getInstance gets an instance based on the registration scope
func (c *Container) getInstance(registration *ServiceRegistration, serviceType reflect.Type, token string, scopeKey string) (any, error) {
	switch registration.Scope {
	case Singleton:
		if registration.Instance != nil {
			return registration.Instance, nil
		}
		if registration.Factory != nil {
			// Create singleton instance and store it
			instance := registration.Factory()
			registration.Instance = instance
			return instance, nil
		}
		return nil, errors.New("no instance or factory registered")

	case Transient:
		if registration.Factory != nil {
			return registration.Factory(), nil
		}
		return nil, errors.New("no factory registered for transient service")

	case Scoped:
		if scopeKey == "" {
			return nil, errors.New("scope key required for scoped service")
		}

		// Check if we already have an instance for this scope
		if token != "" {
			// Token-based lookup
			if scopedMap, exists := c.scopedTokens[scopeKey]; exists {
				if instance, exists := scopedMap[token]; exists {
					return instance, nil
				}
			}
		} else {
			// Type-based lookup
			if scopedMap, exists := c.scopedInstances[scopeKey]; exists {
				if instance, exists := scopedMap[serviceType]; exists {
					return instance, nil
				}
			}
		}

		// Create new instance for this scope
		if registration.Factory != nil {
			instance := registration.Factory()

			// Store the instance for this scope
			if token != "" {
				if c.scopedTokens[scopeKey] == nil {
					c.scopedTokens[scopeKey] = make(map[string]any)
				}
				c.scopedTokens[scopeKey][token] = instance
			} else {
				if c.scopedInstances[scopeKey] == nil {
					c.scopedInstances[scopeKey] = make(map[reflect.Type]any)
				}
				c.scopedInstances[scopeKey][serviceType] = instance
			}

			return instance, nil
		}

		if registration.Instance != nil {
			return registration.Instance, nil
		}

		return nil, errors.New("no factory or instance registered for scoped service")

	default:
		return nil, errors.New("unknown service scope")
	}
}

// assignInstance assigns the resolved instance to the target
func (c *Container) assignInstance(instance any, elem reflect.Value, targetType reflect.Type) error {
	instanceVal := reflect.ValueOf(instance)

	if elem.Kind() == reflect.Ptr {
		// **T case - assign pointer directly
		if !instanceVal.Type().AssignableTo(elem.Type()) {
			return fmt.Errorf("type mismatch: expected %s, got %s", elem.Type().String(), instanceVal.Type().String())
		}
		elem.Set(instanceVal)
		return nil
	}

	if elem.Kind() == reflect.Struct {
		// *T case - assign struct value
		if instanceVal.Kind() != reflect.Ptr {
			return fmt.Errorf("registered service for %s is not a pointer", targetType.String())
		}
		elem.Set(instanceVal.Elem())
		return nil
	}

	return errors.New("unsupported target type")
}

// Resolve by token with optional scope key
func (c *Container) ResolveBy(token string, target any) error {
	return c.ResolveByWithScope(token, target, "")
}

// ResolveByWithScope resolves by token with scope key for request-scoped services
func (c *Container) ResolveByWithScope(token string, target any, scopeKey string) error {
	c.lock.RLock()
	defer c.lock.RUnlock()

	registration, ok := c.tokenServices[token]
	if !ok {
		return fmt.Errorf("no registered service for token %s", token)
	}

	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr {
		return errors.New("target must be a pointer")
	}

	elem := val.Elem()

	instance, err := c.getInstance(registration, reflect.TypeOf(target), token, scopeKey)
	if err != nil {
		return err
	}

	instanceVal := reflect.ValueOf(instance)

	// Handle **T case (pointer to pointer)
	if elem.Kind() == reflect.Ptr {
		if !instanceVal.Type().AssignableTo(elem.Type()) {
			return fmt.Errorf("type mismatch: cannot assign %s to %s",
				instanceVal.Type().String(), elem.Type().String())
		}
		elem.Set(instanceVal)
		return nil
	}

	// Handle *T case (pointer to struct)
	if elem.Kind() == reflect.Struct {
		if instanceVal.Kind() == reflect.Ptr {
			if !instanceVal.Elem().Type().AssignableTo(elem.Type()) {
				return fmt.Errorf("type mismatch: cannot assign %s to %s",
					instanceVal.Elem().Type().String(), elem.Type().String())
			}
			elem.Set(instanceVal.Elem())
		} else {
			if !instanceVal.Type().AssignableTo(elem.Type()) {
				return fmt.Errorf("type mismatch: cannot assign %s to %s",
					instanceVal.Type().String(), elem.Type().String())
			}
			elem.Set(instanceVal)
		}
		return nil
	}

	// Handle interface types
	if elem.Kind() == reflect.Interface {
		if !instanceVal.Type().AssignableTo(elem.Type()) {
			return fmt.Errorf("type mismatch: cannot assign %s to %s",
				instanceVal.Type().String(), elem.Type().String())
		}
		elem.Set(instanceVal)
		return nil
	}

	// Handle other pointer types (like *int, *string, etc.)
	if elem.Kind() == reflect.Ptr {
		if !instanceVal.Type().AssignableTo(elem.Type()) {
			return fmt.Errorf("type mismatch: cannot assign %s to %s",
				instanceVal.Type().String(), elem.Type().String())
		}
		elem.Set(instanceVal)
		return nil
	}

	return fmt.Errorf("unsupported target type: %s", elem.Type().String())
}

// Autowire fills fields with `inject:"type"` or `inject:"token"`
func (c *Container) Autowire(target any) error {
	return c.AutowireWithScope(target, "")
}

// AutowireWithScope autowires with scope key for request-scoped services
func (c *Container) AutowireWithScope(target any, scopeKey string) error {
	c.lock.RLock()
	defer c.lock.RUnlock()

	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return errors.New("target must be a pointer to a struct")
	}

	val = val.Elem() // Dereference pointer to struct
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("inject")
		if tag == "" {
			continue
		}

		fieldVal := val.Field(i)
		if !fieldVal.CanSet() {
			continue
		}

		var instance any
		var err error

		if tag == "type" {
			registration, ok := c.typeServices[field.Type]
			if !ok {
				return fmt.Errorf("missing dependency for type %s", field.Type)
			}
			instance, err = c.getInstance(registration, field.Type, "", scopeKey)
		} else {
			registration, ok := c.tokenServices[tag]
			if !ok {
				return fmt.Errorf("missing dependency for token %s", tag)
			}
			instance, err = c.getInstance(registration, field.Type, tag, scopeKey)
		}

		if err != nil {
			return err
		}

		instanceVal := reflect.ValueOf(instance)

		// Handle type compatibility
		if field.Type.Kind() == reflect.Ptr {
			// Field expects a pointer
			if instanceVal.Kind() == reflect.Ptr {
				if !instanceVal.Type().AssignableTo(field.Type) {
					return fmt.Errorf("type mismatch for field %s: cannot assign %s to %s",
						field.Name, instanceVal.Type().String(), field.Type.String())
				}
				fieldVal.Set(instanceVal)
			} else {
				// Need to get pointer to instance
				if instanceVal.CanAddr() {
					fieldVal.Set(instanceVal.Addr())
				} else {
					return fmt.Errorf("cannot get address of instance for field %s", field.Name)
				}
			}
		} else {
			// Field expects a value
			if instanceVal.Kind() == reflect.Ptr {
				if !instanceVal.Elem().Type().AssignableTo(field.Type) {
					return fmt.Errorf("type mismatch for field %s: cannot assign %s to %s",
						field.Name, instanceVal.Elem().Type().String(), field.Type.String())
				}
				fieldVal.Set(instanceVal.Elem())
			} else {
				if !instanceVal.Type().AssignableTo(field.Type) {
					return fmt.Errorf("type mismatch for field %s: cannot assign %s to %s",
						field.Name, instanceVal.Type().String(), field.Type.String())
				}
				fieldVal.Set(instanceVal)
			}
		}
	}

	return nil
}

// ClearScope clears all instances for a given scope (e.g., end of request)
func (c *Container) ClearScope(scopeKey string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.scopedInstances, scopeKey)
	delete(c.scopedTokens, scopeKey)
}

// CreateScope creates a new scoped container for request-scoped services
func (c *Container) CreateScope(scopeKey string) *ScopedContainer {
	return &ScopedContainer{
		container: c,
		scopeKey:  scopeKey,
	}
}

// ScopedContainer wraps the main container with a scope key
type ScopedContainer struct {
	container *Container
	scopeKey  string
}

func (sc *ScopedContainer) Resolve(target any) error {
	return sc.container.ResolveWithScope(target, sc.scopeKey)
}

func (sc *ScopedContainer) ResolveBy(token string, target any) error {
	return sc.container.ResolveByWithScope(token, target, sc.scopeKey)
}

func (sc *ScopedContainer) Autowire(target any) error {
	return sc.container.AutowireWithScope(target, sc.scopeKey)
}

func (sc *ScopedContainer) MustResolve(target any) {
	if err := sc.Resolve(target); err != nil {
		panic(err)
	}
}

func (sc *ScopedContainer) MustResolveBy(token string, target any) {
	if err := sc.ResolveBy(token, target); err != nil {
		panic(err)
	}
}

func (sc *ScopedContainer) MustAutowire(target any) {
	if err := sc.Autowire(target); err != nil {
		panic(err)
	}
}

func (sc *ScopedContainer) ClearScope() {
	sc.container.ClearScope(sc.scopeKey)
}

func (c *Container) MustResolve(target any) {
	if err := c.Resolve(target); err != nil {
		panic(err)
	}
}

func (c *Container) MustResolveBy(token string, target any) {
	if err := c.ResolveBy(token, target); err != nil {
		panic(err)
	}
}

func (c *Container) MustAutowire(target any) {
	if err := c.Autowire(target); err != nil {
		panic(err)
	}
}

/**
 * AddForAutowiring adds a component to the pending autowire list.
 * This is useful for components that need to be autowired after all registrations are done.
 * It allows you to register components that will be autowired later, ensuring they are resolved
 * with the correct dependencies when AutowireAll is called.
 */
func (c *Container) AddForAutowiring(component any) {
	if c.pendingAutowire == nil {
		c.pendingAutowire = make([]any, 0)
	}
	c.pendingAutowire = append(c.pendingAutowire, component)
}

/**
 * AutowireAll resolves all pending components that were added for autowiring
 * This is useful for batch autowiring after all components have been registered
 */
func (c *Container) AutowireAll() error {
	for _, component := range c.pendingAutowire {
		if err := c.Autowire(component); err != nil {
			return fmt.Errorf("failed to autowire %T: %w", component, err)
		}
	}
	c.pendingAutowire = nil // Clear after autowiring
	return nil
}

// Invoke calls a function with dependencies injected into its arguments
func (c *Container) Invoke(fn any) ([]reflect.Value, error) {
	val := reflect.ValueOf(fn)
	if val.Kind() != reflect.Func {
		return nil, errors.New("argument must be a function")
	}

	t := val.Type()
	numArgs := t.NumIn()
	args := make([]reflect.Value, numArgs)

	for i := 0; i < numArgs; i++ {
		argType := t.In(i)
		// Handle pointer vs value types for resolution
		var targetPtr reflect.Value

		// We need to pass a pointer to Resolve.
		// If argType is *Service, we pass **Service.
		// If argType is Service, we pass *Service.
		if argType.Kind() == reflect.Ptr {
			targetPtr = reflect.New(argType)
		} else {
			targetPtr = reflect.New(argType)
		}

		// Resolve the dependency
		if err := c.Resolve(targetPtr.Interface()); err != nil {
			return nil, fmt.Errorf("failed to resolve argument %d (%v): %w", i, argType, err)
		}

		// targetPtr is now populated.
		// If we passed **Service, targetPtr.Elem() is the *Service we need.
		// If we passed *Service, targetPtr.Elem() is the Service struct we need.
		args[i] = targetPtr.Elem()
	}

	results := val.Call(args)

	// Check if the last return value is an error
	if len(results) > 0 {
		last := results[len(results)-1]
		if last.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			if !last.IsNil() {
				return results, last.Interface().(error)
			}
		}
	}

	return results, nil
}
