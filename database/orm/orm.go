package orm

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/goravel/framework/contracts/config"
	contractsorm "github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/log"
	"github.com/goravel/framework/database/factory"
	"github.com/goravel/framework/database/gorm"
)

const BindingOrm = "goravel.orm"

type Orm struct {
	ctx             context.Context
	config          config.Config
	connection      string
	log             log.Log
	modelToObserver []contractsorm.ModelToObserver
	mutex           sync.Mutex
	query           contractsorm.Query
	queries         map[string]contractsorm.Query
	refresh         func(key any)
}

func NewOrm(
	ctx context.Context,
	config config.Config,
	connection string,
	query contractsorm.Query,
	queries map[string]contractsorm.Query,
	log log.Log,
	modelToObserver []contractsorm.ModelToObserver,
	refresh func(key any),
) *Orm {
	return &Orm{
		ctx:             ctx,
		config:          config,
		connection:      connection,
		log:             log,
		modelToObserver: modelToObserver,
		query:           query,
		queries:         queries,
		refresh:         refresh,
	}
}

func BuildOrm(ctx context.Context, config config.Config, connection string, log log.Log, refresh func(key any)) (*Orm, error) {
	query, err := gorm.BuildQuery(ctx, config, connection, log, nil)
	if err != nil {
		return nil, err
	}

	queries := map[string]contractsorm.Query{
		connection: query,
	}

	return NewOrm(ctx, config, connection, query, queries, log, nil, refresh), nil
}

func (r *Orm) Connection(name string) contractsorm.Orm {
	if name == "" {
		name = r.config.GetString("database.default")
	}
	if instance, exist := r.queries[name]; exist {
		return NewOrm(r.ctx, r.config, name, instance, r.queries, r.log, r.modelToObserver, r.refresh)
	}

	query, err := gorm.BuildQuery(r.ctx, r.config, name, r.log, r.modelToObserver)
	if err != nil || query == nil {
		r.log.Errorf("[Orm] Init %s connection error: %v", name, err)

		return NewOrm(r.ctx, r.config, name, nil, r.queries, r.log, r.modelToObserver, r.refresh)
	}

	r.queries[name] = query

	return NewOrm(r.ctx, r.config, name, query, r.queries, r.log, r.modelToObserver, r.refresh)
}

func (r *Orm) DB() (*sql.DB, error) {
	return r.query.DB()
}

func (r *Orm) Factory() contractsorm.Factory {
	return factory.NewFactoryImpl(r.Query())
}

func (r *Orm) DatabaseName() string {
	return r.config.GetString(fmt.Sprintf("database.connections.%s.database", r.connection))
}

func (r *Orm) Name() string {
	return r.connection
}

func (r *Orm) Observe(model any, observer contractsorm.Observer) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.modelToObserver = append(r.modelToObserver, contractsorm.ModelToObserver{
		Model:    model,
		Observer: observer,
	})

	for _, query := range r.queries {
		if queryWithObserver, ok := query.(contractsorm.QueryWithObserver); ok {
			queryWithObserver.Observe(model, observer)
		}
	}

	if queryWithObserver, ok := r.query.(contractsorm.QueryWithObserver); ok {
		queryWithObserver.Observe(model, observer)
	}
}

func (r *Orm) Query() contractsorm.Query {
	return r.query
}

func (r *Orm) SetQuery(query contractsorm.Query) {
	r.query = query
}

func (r *Orm) Refresh() {
	r.refresh(BindingOrm)
}

func (r *Orm) Transaction(txFunc func(tx contractsorm.Query) error) error {
	tx, err := r.Query().Begin()
	if err != nil {
		return err
	}

	if err := txFunc(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	} else {
		return tx.Commit()
	}
}

func (r *Orm) WithContext(ctx context.Context) contractsorm.Orm {
	for _, query := range r.queries {
		if queryWithSetContext, ok := query.(contractsorm.QueryWithSetContext); ok {
			queryWithSetContext.SetContext(ctx)
		}
	}

	if queryWithSetContext, ok := r.query.(contractsorm.QueryWithSetContext); ok {
		queryWithSetContext.SetContext(ctx)
	}

	return NewOrm(ctx, r.config, r.connection, r.query, r.queries, r.log, r.modelToObserver, r.refresh)
}
