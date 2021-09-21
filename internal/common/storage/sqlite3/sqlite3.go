package sqlite3

import (
	"github.com/yakumioto/alkaid/internal/common/storage"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type sqlite3 struct {
	db *gorm.DB
}

func NewDB(path string) (*sqlite3, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &sqlite3{
		db: db,
	}, nil
}

func (s *sqlite3) Create(value interface{}) error {
	if tx := s.db.Create(value); tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (s *sqlite3) Update(values interface{}, options *storage.UpdateOptions) error {
	if options == nil {
		return storage.ErrNeedUpdateOptions
	}

	if tx := s.db.Model(values).Where(options.Query, options.Args).Updates(values); tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (s *sqlite3) FindByID(dest interface{}, conditions ...interface{}) error {
	if tx := s.db.First(dest, conditions); tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (s *sqlite3) FindByQuery(dest interface{}, options *storage.QueryOptions) error {
	if options == nil {
		options = storage.NewQueryOptions()
	}

	tx := s.db.Order(options.Order()).Limit(options.Limit()).Offset(options.Offset())

	if where := options.Where(); where != nil {
		tx.Where(where.Query, where.Args)
	}

	if or := options.Or(); or != nil {
		tx.Where(or.Query, or.Args)
	}

	if not := options.Not(); not != nil {
		tx.Where(not.Query, not.Args)
	}

	if tx = tx.Find(dest); tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (s *sqlite3) Delete(value interface{}, conditions ...interface{}) error {
	if tx := s.db.Delete(value, conditions); tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (s *sqlite3) Begin() storage.Storage {
	return &sqlite3{db: s.db.Begin()}
}

func (s *sqlite3) Commit() error {
	if tx := s.db.Commit(); tx.Error != nil {
		return tx.Error
	}

	return nil
}
