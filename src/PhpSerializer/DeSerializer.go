package PhpSerializer

import (
	"github.com/yvasiyarov/php_session_decoder"
)

func Deserialize(sessionData string) (map[string]interface{}, error) {
	decoder := php_session_decoder.NewPhpDecoder(sessionData)
	response := make(map[string]interface{})
	result, err := decoder.Decode()
	if err == nil {
		for key, value := range result {
			response[key] = value
		}
	}
	return response, err
}
