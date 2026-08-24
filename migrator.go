package smooth

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Migrator struct {
	Logger *zap.SugaredLogger
}

func NewMigrator() (*Migrator, error) {
	logger, err := zap.Config{
		Encoding:          "json",
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel), // Define o nível de log
		OutputPaths:       []string{"stdout"},
		DisableStacktrace: true,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "severity",
			TimeKey:        "time",
			NameKey:        "logger",
			CallerKey:      "caller",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05"),
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}.Build()
	if err != nil {
		return nil, err
	}

	return &Migrator{
		Logger: logger.Sugar(),
	}, nil
}

func (m *Migrator) Init(engine Engine) error {
	//Verificar se a tabela está criada
	var hasTable bool
	rawErr := engine.Raw(context.Background(), &hasTable, Query{
		Raw: &Raw{
			Query: "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'migrations');",
		},
	})

	if rawErr != nil {
		return rawErr
	}

	if !hasTable {
		m.Logger.Infoln("[SmoothORM]: Migration table not exist.")
		//Criar a tabela
		cErr := m.CreateMigrationTable(engine)
		if cErr != nil {
			return cErr
		}
		m.Logger.Infoln("[SmoothORM]: Migration table created.")
	}

	return nil
}

func (m *Migrator) CreateMigrationTable(engine Engine) error {
	return engine.Exec(context.Background(), `
		CREATE SEQUENCE migrations_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1;
		CREATE TABLE "public"."migrations" (
			"id" bigint DEFAULT nextval('migrations_id_seq') NOT NULL,
			"name" character varying(255) NOT NULL,
			CONSTRAINT "migrations_pkey" PRIMARY KEY ("id")
		)
		WITH (oids = false);

		CREATE UNIQUE INDEX uni_migrations_name ON public.migrations USING btree (name);
	`)
}

func (m *Migrator) RunMigration(s Schema, identifier string, engine Engine) error {
	//Verificar se já rodou
	var mig *Migration
	err := engine.First(context.Background(), &mig, Query{
		Where: &[]Where{
			{
				Column: "name",
				Value:  identifier,
			},
		},
	})
	if err != nil {
		// Significa que deu algum error
		return err
	}

	if mig != nil {
		//Significa que já tem. Então já rodou
		return nil
	}

	//Executar Migration e salvar o identificador atomicamente
	txErr := engine.Transaction(context.Background(), func(ctx context.Context) error {
		excErr := engine.Exec(ctx, s.SQL)
		if excErr != nil {
			return excErr
		}

		new := Migration{
			Name: identifier,
		}
		return engine.Create(ctx, &new)
	})
	if txErr != nil {
		return txErr
	}

	m.Logger.Infoln(fmt.Sprint("[SmoothORM]: Migration '", identifier, "' executed successfuly."))
	return nil
}

type Migration struct {
	ID   uint   `json:"id" gorm:"primarykey;autoIncrement"`
	Name string `json:"name"`
}
