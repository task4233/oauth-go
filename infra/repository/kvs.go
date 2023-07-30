package repository

type KVS interface {
	Get(key string) (map[string]string, error)
	Set(key string, value map[string]string) error
}
