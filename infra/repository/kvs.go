//go:generate mockgen -source $GOFILE -destination ../mock/mock_$GOFILE -package mock

package repository

type KVS interface {
	Get(key string) (map[string]string, error)
	Set(key string, value map[string]string) error
}
