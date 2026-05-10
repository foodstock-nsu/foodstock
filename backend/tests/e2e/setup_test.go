///go:build e2e

package e2e

import (
	"backend/internal/testhelpers"
	pkgpostgres "backend/pkg/postgres"
	"context"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

const migrationVersion = 7

type TestApp struct {
	Server *httptest.Server
	Pg     *testhelpers.PostgresContainer
}

var (
	appInstance *TestApp
	once        sync.Once
)

func setupE2E(t *testing.T) *TestApp {
	once.Do(func() {
		ctx := context.Background()

		container, err := testhelpers.StartPostgresContainer(ctx)
		require.NoError(t, err)

		err = container.MigrateUp(migrationVersion)
		require.NoError(t, err)
	})

	appInstance.cleanData(t, context.Background())

	return appInstance
}

// cleanData Clears a whole database between tests
func (a *TestApp) cleanData(t *testing.T, ctx context.Context) {
	query := `
	DO $$ 
	DECLARE 
	    r RECORD;
	BEGIN
	    FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename != 'schema_migrations') LOOP
	        EXECUTE 'TRUNCATE TABLE ' || quote_ident(r.tablename) || ' CASCADE';
	    END LOOP;
	END $$;`

	client, err := pkgpostgres.NewClient(ctx, appInstance.Pg.Config)
	require.NoError(t, err)

	_, err = client.Pool.Exec(ctx, query)
	require.NoError(t, err)
}
