package WebSockHub

import (
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func (s *Service) handleBroadcast(context *gin.Context) {
	message, err := ioutil.ReadAll(context.Request.Body)

	if err != nil {
		s.Logging.Debug(1007,
			"could not read response body, reason: %s\n",
			err.Error(),
		)
		context.JSON(makeInternalServerError(err.Error(), nil))
		return
	}

	s.sync <- true
	defer func() {
		<-s.sync
	}()

	for authToken := range s.currentSubscriptions {
		s.currentSubscriptions.publishMessage(authToken, message)
	}

	context.JSON(
		makeStatusOk(
			fmt.Sprintf(
				"message sent to %d recipients\n",
				len(s.currentSubscriptions),
			),
			string(message),
		),
	)
	return
}
