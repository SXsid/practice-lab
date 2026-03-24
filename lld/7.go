package main

// withouut factory
// main.go
func main() {
	if env == "prod" {
		db = createProdDB()
	}

	// some other file
	if env == "prod" {
		cache = createProdCache()
	}

	// another file
	if env == "prod" {
		logger = createProdLogger()
	}
}

// every file have same if else logic based on the env
// Factory
// factory.go - the ONLY place that knows about environments
func NewApp(env string) *Application {
	switch env {
	case "prod":
		db := createProdDB()
		cache := createProdCache()
		logger := createProdLogger()
		return &Application{db, cache, logger}

	case "test":
		db := createTestDB()
		cache := createMockCache()
		logger := createTestLogger()
		return &Application{db, cache, logger}
	}
}
