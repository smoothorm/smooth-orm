package smooth

type Options struct {
	MigrationSystem bool
}

type Smooth struct {
	DB       Engine
	Migrator *Migrator
}

func New(engine Engine, options *Options) (*Smooth, error) {
	var newSmooth Smooth = Smooth{}
	ok, err := engine.Health()
	if !ok {
		return nil, err
	}
	newSmooth.DB = engine

	if options != nil {
		if options.MigrationSystem {
			migrator, err := NewMigrator()
			if err != nil {
				return nil, err
			}
			imErr := migrator.Init(engine)
			if imErr != nil {
				return nil, imErr
			}
			newSmooth.Migrator = migrator
		}
	}

	return &newSmooth, nil
}

func (db *Smooth) Migration(s Schema, identifier string) error {
	if db.Migrator == nil {
		return ErrMigrationSystemDisabled
	}

	return db.Migrator.RunMigration(s, identifier, db.DB)
}

type Schema struct {
	SQL string
}
